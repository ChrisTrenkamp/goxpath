package parser

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/ChrisTrenkamp/goxpath/goxpath/internal/lexer"
	"github.com/ChrisTrenkamp/goxpath/goxpath/internal/parser/findutil"
	"github.com/ChrisTrenkamp/goxpath/goxpath/internal/parser/fns"
	"github.com/ChrisTrenkamp/goxpath/goxpath/pathexpr"
	"github.com/ChrisTrenkamp/goxpath/goxpath/xconst"
	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/tree/literals/numlit"
	"github.com/ChrisTrenkamp/goxpath/tree/literals/strlit"
)

//Parser parses an XML document and generates output from the Lexer
type Parser struct {
	xpath  *[]XPExec
	tree   tree.Node
	exNum  int
	ns     map[string]string
	ctx    tree.Node
	pExpr  pathexpr.PathExpr
	filter []tree.Res
	stack  *Parser
	parent *Parser
	fnName string
	fnArgs [][]tree.Res
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
	lexer.XItemStrLit:         strLit,
	lexer.XItemNumLit:         numLit,
}

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
		xpath:  p.xpath,
		tree:   p.tree,
		exNum:  p.exNum + 1,
		ns:     p.ns,
		ctx:    p.ctx,
		pExpr:  pathexpr.PathExpr{},
		filter: nil,
		stack:  nil,
		parent: p,
		fnName: "",
		fnArgs: nil,
	}
}

func (p *Parser) pop() {
	if p.parent != nil {
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

//MustParse is like Parse, but panics instead of returning an error.
func MustParse(xp string) []XPExec {
	ret, err := Parse(xp)

	if err != nil {
		panic(err)
	}

	return ret
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

		tok, next, err = eval(item.Typ, item.Val, tok)

		if err != nil {
			return nil, err
		}

		ret = append(ret, next)
	}

	return ret, nil
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

func beginExprToks() expTkns {
	return expTkns{lexer.XItemAbsLocPath, lexer.XItemAbbrAbsLocPath, lexer.XItemAbbrRelLocPath, lexer.XItemRelLocPath, lexer.XItemStrLit, lexer.XItemNumLit}
}

func pathStartToks() expTkns {
	return expTkns{lexer.XItemAxis, lexer.XItemAbbrAxis, lexer.XItemNCName, lexer.XItemQName, lexer.XItemNodeType, lexer.XItemFunction}
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
		p.ctx = p.tree
		return nil
	}

	return pathStartToks(), ret
}

func abbrAbsLocPath(val string) (expTkns, XPExec) {
	ret := func(p *Parser) error {
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

	return expTkns{lexer.XItemPredicate, lexer.XItemEndPath}, ret
}

func nodeType(val string) (expTkns, XPExec) {
	retFunc := func(p *Parser) error {
		p.pExpr.NodeType = val
		return p.find()
	}

	ret := expTkns{lexer.XItemPredicate, lexer.XItemEndPath}

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

	return expTkns{lexer.XItemPredicate, lexer.XItemEndPath}, ret
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
		if fn, ok := fns.BuiltIn[p.fnName]; ok {
			filt, err := fn(p.filter, p.fnArgs...)

			if err != nil {
				return err
			}

			p.filter = filt
			return nil
		}

		return fmt.Errorf("Unknown function: '%s'\n", p.fnName)
	}
	return nil, ret
}

func strLit(val string) (expTkns, XPExec) {
	ret := func(p *Parser) error {
		p.filter = []tree.Res{strlit.StrLit(val)}
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

	p.filter = vals
	p.pExpr = pathexpr.PathExpr{}

	return nil
}
