package execxp

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	"github.com/ChrisTrenkamp/goxpath/internal/execxp/findutil"
	"github.com/ChrisTrenkamp/goxpath/internal/execxp/intfns"
	"github.com/ChrisTrenkamp/goxpath/lexer"
	"github.com/ChrisTrenkamp/goxpath/parser"
	"github.com/ChrisTrenkamp/goxpath/parser/pathexpr"
	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/xconst"
)

type xpFilt struct {
	t         interface{}
	ctx       tree.Result
	expr      pathexpr.PathExpr
	ns        map[string]string
	ctxPos    int
	ctxSize   int
	proxPos   map[int]int
	fns       map[xml.Name]tree.Wrap
	variables map[string]tree.Result
}

type xpExecFn func(tree.Adapter, *xpFilt, string)

var xpFns = map[lexer.XItemType]xpExecFn{
	lexer.XItemAbsLocPath:     xfAbsLocPath,
	lexer.XItemAbbrAbsLocPath: xfAbbrAbsLocPath,
	lexer.XItemRelLocPath:     xfRelLocPath,
	lexer.XItemAbbrRelLocPath: xfAbbrRelLocPath,
	lexer.XItemAxis:           xfAxis,
	lexer.XItemAbbrAxis:       xfAbbrAxis,
	lexer.XItemNCName:         xfNCName,
	lexer.XItemQName:          xfQName,
	lexer.XItemNodeType:       xfNodeType,
	lexer.XItemProcLit:        xfProcInstLit,
	lexer.XItemStrLit:         xfStrLit,
	lexer.XItemNumLit:         xfNumLit,
}

func xfExec(a tree.Adapter, f *xpFilt, n *parser.Node) (err error) {
	for n != nil {
		if fn, ok := xpFns[n.Val.Typ]; ok {
			fn(a, f, n.Val.Val)
			n = n.Left
		} else if n.Val.Typ == lexer.XItemPredicate {
			if err = xfPredicate(a, f, n.Left); err != nil {
				return
			}

			n = n.Right
		} else if n.Val.Typ == lexer.XItemFunction {
			if err = xfFunction(a, f, n); err != nil {
				return
			}

			n = n.Right
		} else if n.Val.Typ == lexer.XItemOperator {
			lf := xpFilt{
				t:         f.t,
				ns:        f.ns,
				ctx:       f.ctx,
				ctxPos:    f.ctxPos,
				ctxSize:   f.ctxSize,
				proxPos:   f.proxPos,
				fns:       f.fns,
				variables: f.variables,
			}
			left, err := exec(a, &lf, n.Left)
			if err != nil {
				return err
			}

			rf := xpFilt{
				t:         f.t,
				ns:        f.ns,
				ctx:       f.ctx,
				fns:       f.fns,
				variables: f.variables,
			}
			right, err := exec(a, &rf, n.Right)
			if err != nil {
				return err
			}

			return xfOperator(a, left, right, f, n.Val.Val)
		} else if n.Val.Typ == lexer.XItemVariable {
			if res, ok := f.variables[n.Val.Val]; ok {
				f.ctx = res
				return nil
			}
			return fmt.Errorf("Invalid variable '%s'", n.Val.Val)
		} else if string(n.Val.Typ) == "" {
			n = n.Left
			//} else {
			//	return fmt.Errorf("Cannot process " + string(n.Val.Typ))
		}
	}

	return
}

func xfPredicate(a tree.Adapter, f *xpFilt, n *parser.Node) (err error) {
	res := f.ctx.(tree.NodeSet)
	nodes := res.GetNodes()
	newRes := make([]interface{}, 0, len(nodes))

	for i := range nodes {
		pf := xpFilt{
			t:         f.t,
			ns:        f.ns,
			ctxPos:    i,
			ctxSize:   f.ctxSize,
			ctx:       a.NewNodeSet([]interface{}{nodes[i]}),
			fns:       f.fns,
			variables: f.variables,
		}

		predRes, err := exec(a, &pf, n)
		if err != nil {
			return err
		}

		ok, err := checkPredRes(a, predRes, f, nodes[i])
		if err != nil {
			return err
		}

		if ok {
			newRes = append(newRes, nodes[i])
		}
	}

	f.proxPos = make(map[int]int)
	for pos, j := range newRes {
		f.proxPos[a.NodePos(j)] = pos + 1
	}

	f.ctx = a.NewNodeSet(newRes)
	f.ctxSize = len(newRes)

	return
}

func checkPredRes(a tree.Adapter, ret tree.Result, f *xpFilt, node interface{}) (bool, error) {
	if num, ok := ret.(tree.Num); ok {
		if float64(f.proxPos[a.NodePos(node)]) == float64(num) {
			return true, nil
		}
		return false, nil
	}

	if b, ok := ret.(tree.IsBool); ok {
		return bool(b.Bool()), nil
	}

	return false, fmt.Errorf("Cannot convert argument to boolean")
}

func xfFunction(a tree.Adapter, f *xpFilt, n *parser.Node) error {
	spl := strings.Split(n.Val.Val, ":")
	var name xml.Name
	if len(spl) == 1 {
		name.Local = spl[0]
	} else {
		name.Space = f.ns[spl[0]]
		name.Local = spl[1]
	}
	fn, ok := intfns.BuiltIn[name]
	if !ok {
		fn, ok = f.fns[name]
	}

	if ok {
		args := []tree.Result{}
		param := n.Left

		for param != nil {
			pf := xpFilt{
				t:         f.t,
				ctx:       f.ctx,
				ns:        f.ns,
				ctxPos:    f.ctxPos,
				ctxSize:   f.ctxSize,
				fns:       f.fns,
				variables: f.variables,
			}
			res, err := exec(a, &pf, param.Left)
			if err != nil {
				return err
			}

			args = append(args, res)
			param = param.Right
		}

		filt, err := fn.Call(a, tree.Ctx{NodeSet: f.ctx.(tree.NodeSet), Size: f.ctxSize, Pos: f.ctxPos + 1}, args...)
		f.ctx = filt
		return err
	}

	return fmt.Errorf("Unknown function: %s", n.Val.Val)
}

var eqOps = map[string]bool{
	"=":  true,
	"!=": true,
}

var booleanOps = map[string]bool{
	"=":  true,
	"!=": true,
	"<":  true,
	"<=": true,
	">":  true,
	">=": true,
}

var numOps = map[string]bool{
	"*":   true,
	"div": true,
	"mod": true,
	"+":   true,
	"-":   true,
	"=":   true,
	"!=":  true,
	"<":   true,
	"<=":  true,
	">":   true,
	">=":  true,
}

var andOrOps = map[string]bool{
	"and": true,
	"or":  true,
}

func xfOperator(a tree.Adapter, left, right tree.Result, f *xpFilt, op string) error {
	if booleanOps[op] {
		lNode, lOK := left.(tree.NodeSet)
		rNode, rOK := right.(tree.NodeSet)
		if lOK && rOK {
			return bothNodeOperator(a, lNode, rNode, f, op)
		}

		if lOK {
			return leftNodeOperator(a, lNode, right, f, op)
		}

		if rOK {
			return rightNodeOperator(a, left, rNode, f, op)
		}

		if eqOps[op] {
			return equalsOperator(left, right, f, op)
		}
	}

	if numOps[op] {
		return numberOperator(left, right, f, op)
	}

	if andOrOps[op] {
		return andOrOperator(left, right, f, op)
	}

	//if op == "|" {
	return unionOperator(a, left, right, f, op)
	//}

	//return fmt.Errorf("Unknown operator " + op)
}

func xfAbsLocPath(a tree.Adapter, f *xpFilt, val string) {
	i := f.t
	for a.GetNodeType(i) != tree.NtRoot {
		i = a.GetParent(i)
	}
	f.ctx = a.NewNodeSet([]interface{}{i})
}

func xfAbbrAbsLocPath(a tree.Adapter, f *xpFilt, val string) {
	i := f.t
	for a.GetNodeType(i) != tree.NtRoot {
		i = a.GetParent(i)
	}
	f.ctx = a.NewNodeSet([]interface{}{i})
	f.expr = abbrPathExpr()
	find(a, f)
}

func xfRelLocPath(a tree.Adapter, f *xpFilt, val string) {
}

func xfAbbrRelLocPath(a tree.Adapter, f *xpFilt, val string) {
	f.expr = abbrPathExpr()
	find(a, f)
}

func xfAxis(a tree.Adapter, f *xpFilt, val string) {
	f.expr.Axis = val
}

func xfAbbrAxis(a tree.Adapter, f *xpFilt, val string) {
	f.expr.Axis = xconst.AxisAttribute
}

func xfNCName(a tree.Adapter, f *xpFilt, val string) {
	f.expr.Name.Space = val
}

func xfQName(a tree.Adapter, f *xpFilt, val string) {
	f.expr.Name.Local = val
	find(a, f)
}

func xfNodeType(a tree.Adapter, f *xpFilt, val string) {
	f.expr.NodeType = val
	find(a, f)
}

func xfProcInstLit(a tree.Adapter, f *xpFilt, val string) {
	filt := make([]interface{}, 0)
	for _, i := range f.ctx.(tree.NodeSet).GetNodes() {
		if a.GetNodeType(i) == tree.NtPi {
			if a.GetProcInstTok(i).Target == val {
				filt = append(filt, i)
			}
		}
	}
	f.ctx = a.NewNodeSet(filt)
}

func xfStrLit(a tree.Adapter, f *xpFilt, val string) {
	f.ctx = tree.String(val)
}

func xfNumLit(a tree.Adapter, f *xpFilt, val string) {
	num, _ := strconv.ParseFloat(val, 64)
	f.ctx = tree.Num(num)
}

func abbrPathExpr() pathexpr.PathExpr {
	return pathexpr.PathExpr{
		Name:     xml.Name{},
		Axis:     xconst.AxisDescendentOrSelf,
		NodeType: xconst.NodeTypeNode,
	}
}

func find(a tree.Adapter, f *xpFilt) {
	dupFilt := make(map[int]interface{})
	f.proxPos = make(map[int]int)

	if f.expr.Axis == "" && f.expr.NodeType == "" && f.expr.Name.Space == "" {
		if f.expr.Name.Local == "." {
			f.expr = pathexpr.PathExpr{
				Name:     xml.Name{},
				Axis:     xconst.AxisSelf,
				NodeType: xconst.NodeTypeNode,
			}
		}

		if f.expr.Name.Local == ".." {
			f.expr = pathexpr.PathExpr{
				Name:     xml.Name{},
				Axis:     xconst.AxisParent,
				NodeType: xconst.NodeTypeNode,
			}
		}
	}

	f.expr.NS = f.ns

	for _, i := range f.ctx.(tree.NodeSet).GetNodes() {
		for pos, j := range findutil.Find(a, i, f.expr) {
			ps := a.NodePos(j)
			dupFilt[ps] = j
			f.proxPos[ps] = pos + 1
		}
	}

	res := make([]interface{}, 0, len(dupFilt))
	for _, i := range dupFilt {
		res = append(res, i)
	}

	f.expr = pathexpr.PathExpr{}
	f.ctxSize = len(res)
	nodeset := a.NewNodeSet(res)
	nodeset.Sort()
	f.ctx = nodeset
}
