package parser

import (
	"encoding/xml"
	"fmt"
	"sort"
	"strconv"

	"github.com/ChrisTrenkamp/goxpath/goxpath/internal/lexer"
	"github.com/ChrisTrenkamp/goxpath/goxpath/internal/parser/findutil"
	"github.com/ChrisTrenkamp/goxpath/goxpath/internal/parser/intfns"
	"github.com/ChrisTrenkamp/goxpath/goxpath/pathexpr"
	"github.com/ChrisTrenkamp/goxpath/goxpath/xconst"
	"github.com/ChrisTrenkamp/goxpath/goxpath/xpfn"
	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/tree/literals/boollit"
	"github.com/ChrisTrenkamp/goxpath/tree/literals/numlit"
	"github.com/ChrisTrenkamp/goxpath/tree/literals/strlit"
)

//Parser parses an XML document and generates output from the Lexer
type Parser struct {
	xpath   *[]XPExec
	tree    tree.Node
	exNum   int
	ns      map[string]string
	ctx     tree.Node
	ctxPos  int
	ctxSize int
	pExpr   pathexpr.PathExpr
	filter  []tree.Res
	stack   *Parser
	parent  *Parser
	fnName  string
	fnArgs  [][]tree.Res
	predEnd int
}

//XPExec is an instruction that operates on XPath trees
type XPExec func(*Parser) error

type expTkns []lexer.XItemType
type lexFn func(string) (expTkns, XPExec)

var parseMap = map[lexer.XItemType]lexFn{
	lexer.XItemAbsLocPath:     absLocPath,
	lexer.XItemAbbrAbsLocPath: abbrAbsLocPath,
	lexer.XItemRelLocPath:     relLocPath,
	lexer.XItemAbbrRelLocPath: abbrRelLocPath,
	lexer.XItemAxis:           axis,
	lexer.XItemAbbrAxis:       abbrAxis,
	lexer.XItemNCName:         ncName,
	lexer.XItemQName:          qName,
	lexer.XItemNodeType:       nodeType,
	lexer.XItemProcLit:        procInstLit,
	lexer.XItemEndPath:        endPath,
	lexer.XItemFunction:       function,
	lexer.XItemArgument:       argument,
	lexer.XItemEndFunction:    endFunction,
	lexer.XItemPredicate:      predicate,
	lexer.XItemEndPredicate:   endPredicate,
	lexer.XItemStrLit:         strLit,
	lexer.XItemNumLit:         numLit,
}

type nodeSort []tree.Res

func (ns nodeSort) Len() int           { return len(ns) }
func (ns nodeSort) Swap(i, j int)      { ns[i], ns[j] = ns[j], ns[i] }
func (ns nodeSort) Less(i, j int) bool { return ns[i].(tree.Node).Pos() < ns[j].(tree.Node).Pos() }

func (p *Parser) exec() ([]tree.Res, error) {
	for p.exNum < len(*p.xpath) {
		err := (*p.xpath)[p.exNum](p)
		if err != nil {
			return nil, err
		}
		p.exNum++
	}

	return p.filter, nil
}

func (p *Parser) push() {
	p.stack = &Parser{
		xpath:   p.xpath,
		tree:    p.tree,
		exNum:   p.exNum + 1,
		ns:      p.ns,
		ctx:     p.ctx,
		ctxSize: p.ctxSize,
		pExpr:   pathexpr.PathExpr{},
		filter:  nil,
		stack:   nil,
		parent:  p,
		fnName:  "",
		fnArgs:  nil,
	}
}

func (p *Parser) pop() {
	if p.parent != nil {
		p.predEnd = p.exNum + 1
		p.parent.exNum = p.exNum
		p.exNum = len(*p.xpath)
	}
}

//Exec executes the XPath expression, xp, against the tree, t, with the
//namespace mappings, ns.
func Exec(xp []XPExec, t tree.Node, ns map[string]string) ([]tree.Res, error) {
	if ns == nil {
		ns = make(map[string]string)
	}

	p := Parser{
		xpath:  &xp,
		tree:   t,
		exNum:  0,
		ns:     ns,
		ctx:    t,
		pExpr:  pathexpr.PathExpr{},
		filter: nil,
		stack:  nil,
		parent: nil,
		fnName: "",
		fnArgs: nil,
	}

	return p.exec()
}

//Parse parses the XPath expression, xp, returning an XPath executor.
func Parse(xp string) ([]XPExec, error) {
	var err error
	var next XPExec
	tok := beginExprToks()
	var ret []XPExec
	c := lexer.Lex(xp)

	for item := range c {
		if item.Typ == lexer.XItemError {
			return nil, fmt.Errorf(item.Val)
		}

		if err != nil {
			continue
		}

		tok, next, err = eval(item.Typ, item.Val, tok)

		ret = append(ret, next)
	}

	return ret, err
}

func eval(typ lexer.XItemType, val string, tkns expTkns) (expTkns, XPExec, error) {
	ok := len(tkns) == 0

	if !ok {
		for i := range tkns {
			if typ == tkns[i] {
				ok = true
				break
			}
		}
	}

	if !ok {
		return expTkns{}, nil, fmt.Errorf("Unexpected token: %s", string(typ))
	}

	if f, ok := parseMap[typ]; ok {
		tkns, next := f(val)
		return tkns, next, nil
	}

	return expTkns{}, nil, fmt.Errorf("Unsupported token emitted: %s", string(typ))
}

func nextPathToks() expTkns {
	return expTkns{lexer.XItemAbsLocPath, lexer.XItemAbbrAbsLocPath, lexer.XItemAbbrRelLocPath, lexer.XItemRelLocPath, lexer.XItemPredicate, lexer.XItemEndPredicate, lexer.XItemEndPath}
}

func beginExprToks() expTkns {
	return expTkns{lexer.XItemAbsLocPath, lexer.XItemAbbrAbsLocPath, lexer.XItemAbbrRelLocPath, lexer.XItemRelLocPath, lexer.XItemStrLit, lexer.XItemNumLit}
}

func pathStartToks() expTkns {
	return expTkns{lexer.XItemAxis, lexer.XItemAbbrAxis, lexer.XItemNCName, lexer.XItemQName, lexer.XItemNodeType, lexer.XItemFunction, lexer.XItemPredicate, lexer.XItemEndPredicate}
}

func abbrPathExpr() pathexpr.PathExpr {
	return pathexpr.PathExpr{
		Name:     xml.Name{},
		Axis:     xconst.AxisDescendentOrSelf,
		NodeType: xconst.NodeTypeNode,
	}
}

func absLocPath(val string) (expTkns, XPExec) {
	ret := func(p *Parser) error {
		p.filter = []tree.Res{p.tree}
		p.ctx = p.tree
		return nil
	}

	return append(pathStartToks(), lexer.XItemEndPath), ret
}

func abbrAbsLocPath(val string) (expTkns, XPExec) {
	ret := func(p *Parser) error {
		p.filter = []tree.Res{p.tree}
		p.ctx = p.tree
		p.pExpr = abbrPathExpr()
		return p.find()
	}

	return pathStartToks(), ret
}

func relLocPath(val string) (expTkns, XPExec) {
	ret := func(p *Parser) error {
		return nil
	}

	return pathStartToks(), ret
}

func abbrRelLocPath(val string) (expTkns, XPExec) {
	ret := func(p *Parser) error {
		p.pExpr = abbrPathExpr()
		return p.find()
	}

	return pathStartToks(), ret
}

func axis(val string) (expTkns, XPExec) {
	ret := func(p *Parser) error {
		p.pExpr.Axis = val
		return nil
	}

	return expTkns{lexer.XItemNCName, lexer.XItemQName, lexer.XItemNodeType}, ret
}

func abbrAxis(val string) (expTkns, XPExec) {
	ret := func(p *Parser) error {
		p.pExpr.Axis = xconst.AxisAttribute
		return nil
	}

	return expTkns{lexer.XItemNCName, lexer.XItemQName}, ret
}

func ncName(val string) (expTkns, XPExec) {
	ret := func(p *Parser) error {
		p.pExpr.Name.Space = val
		return nil
	}

	return expTkns{lexer.XItemQName}, ret
}

func qName(val string) (expTkns, XPExec) {
	ret := func(p *Parser) error {
		p.pExpr.Name.Local = val
		return p.find()
	}

	return nextPathToks(), ret
}

func nodeType(val string) (expTkns, XPExec) {
	retFunc := func(p *Parser) error {
		p.pExpr.NodeType = val
		return p.find()
	}

	ret := nextPathToks()

	if val == xconst.NodeTypeProcInst {
		ret = append(ret, lexer.XItemProcLit)
	}

	return ret, retFunc
}

func procInstLit(val string) (expTkns, XPExec) {
	ret := func(p *Parser) error {
		filt := []tree.Res{}
		for i := range p.filter {
			if tok, tOk := p.filter[i].(tree.Node); tOk {
				if proc, pOk := tok.GetToken().(xml.ProcInst); pOk {
					if proc.Target == val {
						filt = append(filt, p.filter[i])
					}
				}
			}
		}

		p.filter = filt
		return nil
	}

	return nextPathToks(), ret
}

func endPath(val string) (expTkns, XPExec) {
	ret := func(p *Parser) error {
		p.pop()
		return nil
	}

	return nil, ret
}

func function(val string) (expTkns, XPExec) {
	ret := func(p *Parser) error {
		p.fnName = val
		p.fnArgs = make([][]tree.Res, 0)
		return nil
	}

	return expTkns{lexer.XItemArgument, lexer.XItemEndFunction}, ret
}

func argument(val string) (expTkns, XPExec) {
	ret := func(p *Parser) error {
		p.push()
		filt, err := p.stack.exec()

		if err != nil {
			return err
		}

		p.fnArgs = append(p.fnArgs, filt)

		return nil
	}

	return append(beginExprToks(), lexer.XItemEndFunction), ret
}

func endFunction(val string) (expTkns, XPExec) {
	ret := func(p *Parser) error {
		if fn, ok := intfns.BuiltIn[p.fnName]; ok {
			filt, err := fn.Call(xpfn.Ctx{Node: p.ctx, Filter: p.filter, Size: p.ctxSize, Pos: p.ctxPos}, p.fnArgs...)

			if err != nil {
				return err
			}

			if filt == nil {
				filt = []tree.Res{}
			}

			p.filter = filt
			return nil
		}

		return fmt.Errorf("Unknown function: '%s'", p.fnName)
	}
	return nil, ret
}

func predicate(val string) (expTkns, XPExec) {
	ret := func(p *Parser) error {
		exNum := p.exNum
		filt := make([]tree.Res, 0, len(p.filter))

		if len(p.filter) == 0 {
			p.push()
			_, err := p.stack.exec()
			if err != nil {
				return err
			}
			p.exNum = p.stack.predEnd
			p.pop()
			return nil
		}

		for i := range p.filter {
			p.exNum = exNum
			p.push()
			if n, ok := p.filter[i].(tree.Node); ok {
				p.stack.ctx = n
			}
			p.stack.filter = []tree.Res{p.filter[i]}
			p.ctxPos = i

			ret, err := p.stack.exec()
			if err != nil {
				return err
			}

			if len(ret) == 1 {
				if num, ok := ret[0].(numlit.NumLit); ok {
					if int(num)-1 == i {
						filt = append(filt, p.filter[i])
					}
					p.pop()
					continue
				}
			}

			boolFn := intfns.BuiltIn["boolean"]
			res, err := boolFn.Fn(xpfn.Ctx{}, ret)

			if err != nil {
				return err
			}

			if res[0].(boollit.BoolLit) {
				filt = append(filt, p.filter[i])
			}

			p.pop()
		}

		p.filter = filt
		p.exNum = p.stack.predEnd

		return nil
	}

	return beginExprToks(), ret
}

func endPredicate(val string) (expTkns, XPExec) {
	ret := func(p *Parser) error {
		return nil
	}

	return nextPathToks(), ret
}

func strLit(val string) (expTkns, XPExec) {
	ret := func(p *Parser) error {
		p.filter = []tree.Res{strlit.StrLit(val)}
		p.pop()
		return nil
	}
	return nil, ret
}

func numLit(val string) (expTkns, XPExec) {
	ret := func(p *Parser) error {
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		p.filter = []tree.Res{numlit.NumLit(f)}
		p.pop()
		return nil
	}
	return nil, ret
}

func (p *Parser) find() error {
	vals := []tree.Res{}

	if p.pExpr.Axis == "" && p.pExpr.NodeType == "" && p.pExpr.Name.Space == "" {
		if p.pExpr.Name.Local == "." {
			p.pExpr = pathexpr.PathExpr{
				Name:     xml.Name{},
				Axis:     xconst.AxisSelf,
				NodeType: xconst.NodeTypeNode,
			}
		}

		if p.pExpr.Name.Local == ".." {
			p.pExpr = pathexpr.PathExpr{
				Name:     xml.Name{},
				Axis:     xconst.AxisParent,
				NodeType: xconst.NodeTypeNode,
			}
		}
	}

	if p.filter == nil {
		p.filter = []tree.Res{p.ctx}
	}

	p.pExpr.NS = p.ns

	for _, i := range p.filter {
		if node, ok := i.(tree.Node); ok {
			for _, j := range findutil.Find(node, p.pExpr) {
				vals = append(vals, j)
			}
		} else {
			return fmt.Errorf("Cannot run path expression on primitive data type.")
		}
	}

	p.filter = remDupsAndSort(vals)
	p.pExpr = pathexpr.PathExpr{}
	p.ctxSize = len(p.filter)

	return nil
}

func remDupsAndSort(filt []tree.Res) []tree.Res {
	dupFilt := make(map[tree.Res]int)

	for _, i := range filt {
		dupFilt[i] = 0
	}

	filt = make([]tree.Res, 0, len(dupFilt))
	for i := range dupFilt {
		filt = append(filt, i)
	}

	sort.Sort(nodeSort(filt))
	return filt
}
