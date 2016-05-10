package execxp

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/ChrisTrenkamp/goxpath/internal/parser"
	"github.com/ChrisTrenkamp/goxpath/internal/parser/findutil"
	"github.com/ChrisTrenkamp/goxpath/internal/parser/intfns"
	"github.com/ChrisTrenkamp/goxpath/literals/numlit"
	"github.com/ChrisTrenkamp/goxpath/literals/strlit"
	"github.com/ChrisTrenkamp/goxpath/xconst"
	"github.com/ChrisTrenkamp/goxpath/xfn"
	"github.com/ChrisTrenkamp/goxpath/xsort"

	"github.com/ChrisTrenkamp/goxpath/internal/lexer"
	"github.com/ChrisTrenkamp/goxpath/internal/parser/pathexpr"
	"github.com/ChrisTrenkamp/goxpath/tree"
)

type xpFilt struct {
	t       tree.Node
	res     []tree.Res
	ctx     tree.Node
	expr    pathexpr.PathExpr
	ns      map[string]string
	ctxPos  int
	ctxSize int
}

type xpExecFn func(*xpFilt, string) error

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

func xfExec(f *xpFilt, n *parser.Node) (err error) {
	for n != nil {
		if fn, ok := xpFns[n.Val.Typ]; ok {
			if err = fn(f, n.Val.Val); err != nil {
				return
			}

			n = n.Left
		} else if n.Val.Typ == lexer.XItemPredicate {
			if err = xfPredicate(f, n.Left); err != nil {
				return
			}

			n = n.Right
		} else if n.Val.Typ == lexer.XItemFunction {
			if err = xfFunction(f, n); err != nil {
				return
			}

			n = n.Right
		} else if n.Val.Typ == lexer.XItemOperator {
			left, err := Exec(n.Left, f.t, f.ns)
			if err != nil {
				return err
			}

			right, err := Exec(n.Right, f.t, f.ns)
			if err != nil {
				return err
			}

			return xfOperator(left, right, f, n.Val.Val)
		} else if string(n.Val.Typ) == "" {
			n = n.Left
		} else {
			return fmt.Errorf("Unknown operator " + string(n.Val.Typ))
		}
	}

	return
}

func xfPredicate(f *xpFilt, n *parser.Node) (err error) {
	newRes := []tree.Res{}

	for i := range f.res {
		pf := xpFilt{
			t:       f.t,
			ns:      f.ns,
			ctxPos:  i,
			ctxSize: f.ctxSize,
		}

		if n, ok := f.res[i].(tree.Node); ok {
			pf.ctx = n
		} else {
			return fmt.Errorf("Cannot run predicate on primitive data type")
		}

		res, err := exec(&pf, n)
		if err != nil {
			return err
		}

		if checkPredRes(res, i) {
			newRes = append(newRes, f.res[i])
		}
	}

	f.res = newRes

	return
}

func checkPredRes(ret []tree.Res, i int) bool {
	if len(ret) == 1 {
		if num, ok := ret[0].(numlit.NumLit); ok {
			return int(num)-1 == i
		}
	}

	return intfns.BooleanFunc(ret)
}

func xfFunction(f *xpFilt, n *parser.Node) error {
	if fn, ok := intfns.BuiltIn[n.Val.Val]; ok {
		args := [][]tree.Res{}

		param := n.Left
		for param != nil {
			pf := xpFilt{
				t:       f.t,
				ctx:     f.ctx,
				ns:      f.ns,
				ctxPos:  f.ctxPos,
				ctxSize: f.ctxSize,
			}
			res, err := exec(&pf, param.Left)
			if err != nil {
				return err
			}

			args = append(args, res)
			param = param.Right
		}

		filt, err := fn.Call(xfn.Ctx{Node: f.ctx, Filter: f.res, Size: f.ctxSize, Pos: f.ctxPos}, args...)
		f.res = filt
		return err
	}

	return fmt.Errorf("Unknown function: %s", n.Val.Val)
}

var booleanOps = map[string]int{
	"=":  0,
	"!=": 0,
	"<":  0,
	"<=": 0,
	">":  0,
	">=": 0,
}

var numOps = map[string]int{
	"*":   0,
	"div": 0,
	"mod": 0,
	"+":   0,
	"-":   0,
	"=":   0,
	"!=":  0,
	"<":   0,
	"<=":  0,
	">":   0,
	">=":  0,
}

var andOrOps = map[string]int{
	"and": 0,
	"or":  0,
}

func xfOperator(left, right []tree.Res, f *xpFilt, op string) error {
	if _, ok := booleanOps[op]; ok {
		lNode, lErr := xfn.GetNode(left, nil)
		rNode, rErr := xfn.GetNode(right, nil)
		if lErr == nil && rErr == nil {
			return bothNodeOperator(lNode, rNode, f, op)
		}

		if lErr == nil {
			return leftNodeOperator(lNode, right, f, op)
		}

		if rErr == nil {
			return rightNodeOperator(left, rNode, f, op)
		}

		if op == "=" || op == "!=" {
			return equalsOperator(left, right, f, op)
		}
	}

	if _, ok := numOps[op]; ok {
		return numberOperator(left, right, f, op)
	}

	if _, ok := andOrOps[op]; ok {
		return andOrOperator(left, right, f, op)
	}

	if op == "|" {
		return unionOperator(left, right, f, op)
	}

	return fmt.Errorf("Unknown operator " + op)
}

/*
var numOps = map[string]int{
	"*":   0,
	"div": 0,
	"mod": 0,
	"+":   0,
	"-":   0,
}

var andOrOps = map[string]int{
	"and": 0,
	"or":  0,
}

var equalityOps = map[string]int{
	"=":  0,
	"!=": 0,
	"<":  0,
	"<=": 0,
	">":  0,
	">=": 0,
}

func xfOperator(left, right []tree.Res, f *xpFilt, op string) error {
	if _, ok := numOps[op]; ok {
		return numberOp(left, right, f, op)
	} else if _, ok := andOrOps[op]; ok {
		l := intfns.BooleanFunc(left)
		r := intfns.BooleanFunc(right)

		if op == "and" {
			f.res = []tree.Res{boollit.BoolLit(l && r)}
		} else {
			f.res = []tree.Res{boollit.BoolLit(l || r)}
		}
	} else {
		_, lErr := xfn.GetNode(left, nil)
		_, rErr := xfn.GetNode(right, nil)
		if lErr == nil && rErr == nil {
			return stringOp(left, right, f, op)
		}

		if lErr == nil && rErr != nil || lErr != nil && rErr == nil {
			nonNode := left
			if lErr != nil {
				nonNode = right
			}

			if len(nonNode) == 0 {
				f.res = []tree.Res{boollit.BoolLit(false)}
				return nil
			}

			if len(nonNode) != 1 {
				return fmt.Errorf("More than one primitive result.")
			}

			switch nonNode[0].(type) {
			case numlit.NumLit:
				return stringOp(left, right, f, op)
			case boollit.BoolLit:
			case strlit.StrLit:
			}
		}
	}
}

func numberOp(left, right []tree.Res, f *xpFilt, op string) error {
	ln, err := intfns.NumberFunc(left)
	if err != nil {
		return err
	}

	rn, err := intfns.NumberFunc(right)
	if err != nil {
		return err
	}

	switch op {
	case "*":
		f.res = []tree.Res{numlit.NumLit(ln * rn)}
	case "div":
		f.res = []tree.Res{numlit.NumLit(ln / rn)}
	case "mod":
		f.res = []tree.Res{numlit.NumLit(int(ln) % int(rn))}
	case "+":
		f.res = []tree.Res{numlit.NumLit(ln + rn)}
	case "-":
		f.res = []tree.Res{numlit.NumLit(ln - rn)}
	case "=":
		f.res = []tree.Res{boollit.BoolLit(ln == rn)}
	case "!=":
		f.res = []tree.Res{boollit.BoolLit(ln != rn)}
	case "<":
		f.res = []tree.Res{boollit.BoolLit(ln < rn)}
	case "<=":
		f.res = []tree.Res{boollit.BoolLit(ln <= rn)}
	case ">":
		f.res = []tree.Res{boollit.BoolLit(ln > rn)}
	case ">=":
		f.res = []tree.Res{boollit.BoolLit(ln >= rn)}
	}

	return nil
}

func stringOp(left, right []tree.Res, f *xpFilt, op string) error {
	ln := intfns.StringFunc(left)
	rn := intfns.StringFunc(right)

	switch op {
	case "=":
		f.res = []tree.Res{boollit.BoolLit(ln == rn)}
	case "!=":
		f.res = []tree.Res{boollit.BoolLit(ln != rn)}
	case "<":
		f.res = []tree.Res{boollit.BoolLit(ln < rn)}
	case "<=":
		f.res = []tree.Res{boollit.BoolLit(ln <= rn)}
	case ">":
		f.res = []tree.Res{boollit.BoolLit(ln > rn)}
	case ">=":
		f.res = []tree.Res{boollit.BoolLit(ln >= rn)}
	}

	return nil
}
*/

func xfAbsLocPath(f *xpFilt, val string) error {
	f.res = []tree.Res{f.t}
	f.ctx = f.t
	return nil
}

func xfAbbrAbsLocPath(f *xpFilt, val string) error {
	f.res = []tree.Res{f.t}
	f.ctx = f.t
	f.expr = abbrPathExpr()
	return find(f)
}

func xfRelLocPath(f *xpFilt, val string) error {
	return nil
}

func xfAbbrRelLocPath(f *xpFilt, val string) error {
	f.expr = abbrPathExpr()
	return find(f)
}

func xfAxis(f *xpFilt, val string) error {
	f.expr.Axis = val
	return nil
}

func xfAbbrAxis(f *xpFilt, val string) error {
	f.expr.Axis = xconst.AxisAttribute
	return nil
}

func xfNCName(f *xpFilt, val string) error {
	f.expr.Name.Space = val
	return nil
}

func xfQName(f *xpFilt, val string) error {
	f.expr.Name.Local = val
	return find(f)
}

func xfNodeType(f *xpFilt, val string) error {
	f.expr.NodeType = val
	return find(f)
}

func xfProcInstLit(f *xpFilt, val string) error {
	filt := []tree.Res{}
	for i := range f.res {
		if tok, tOk := f.res[i].(tree.Node); tOk {
			if proc, pOk := tok.GetToken().(xml.ProcInst); pOk {
				if proc.Target == val {
					filt = append(filt, f.res[i])
				}
			}
		}
	}

	f.res = filt
	return nil
}

func xfStrLit(f *xpFilt, val string) error {
	f.res = []tree.Res{strlit.StrLit(val)}
	return nil
}

func xfNumLit(f *xpFilt, val string) error {
	num, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return err
	}

	f.res = []tree.Res{numlit.NumLit(num)}
	return nil
}

func abbrPathExpr() pathexpr.PathExpr {
	return pathexpr.PathExpr{
		Name:     xml.Name{},
		Axis:     xconst.AxisDescendentOrSelf,
		NodeType: xconst.NodeTypeNode,
	}
}

func find(f *xpFilt) error {
	dupFilt := make(map[int]tree.Res)

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

	if f.res == nil {
		f.res = []tree.Res{f.ctx}
	}

	f.expr.NS = f.ns

	for _, i := range f.res {
		if node, ok := i.(tree.Node); ok {
			for _, j := range findutil.Find(node, f.expr) {
				dupFilt[j.Pos()] = j
			}
		} else {
			return fmt.Errorf("Cannot run path expression on primitive data type.")
		}
	}

	f.res = make([]tree.Res, 0, len(dupFilt))
	for _, i := range dupFilt {
		f.res = append(f.res, i)
	}

	xsort.SortRes(f.res)

	f.expr = pathexpr.PathExpr{}
	f.ctxSize = len(f.res)

	return nil
}
