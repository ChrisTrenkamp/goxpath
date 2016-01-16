package parser

import (
	"encoding/xml"
	"fmt"

	"github.com/ChrisTrenkamp/goxpath/goxpath/internal/lexer"
	"github.com/ChrisTrenkamp/goxpath/goxpath/internal/parser/findutil"
	"github.com/ChrisTrenkamp/goxpath/goxpath/pathexpr"
	"github.com/ChrisTrenkamp/goxpath/goxpath/xconst"
	"github.com/ChrisTrenkamp/goxpath/tree"
)

//Parser parses an XML document and generates output from the Lexer
type Parser struct {
	tree   tree.XPRes
	ns     map[string]string
	ctx    tree.XPRes
	pExpr  pathexpr.PathExpr
	filter []tree.XPRes
}

//XPExec is an instruction that operates on XPath trees
type XPExec func(*Parser)

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
}

//Exec executes the XPath expression, xp, against the tree, t, with the
//namespace mappings, ns.
func Exec(xp []XPExec, t tree.XPRes, ns map[string]string) []tree.XPRes {
	if ns == nil {
		ns = make(map[string]string)
	}

	p := Parser{tree: t, ns: ns, ctx: t}

	for _, i := range xp {
		i(&p)
	}

	return p.filter
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
	tok := expTkns{}
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
		return expTkns{}, nil, fmt.Errorf("Unexpected token: %d", typ)
	}

	if f, ok := parseMap[typ]; ok {
		tkns, next := f(val)
		return tkns, next, nil
	}

	return expTkns{}, nil, fmt.Errorf("Unsupported token emitted: %d", typ)
}

func pathStartToks() expTkns {
	return expTkns{lexer.XItemAxis, lexer.XItemAbbrAxis, lexer.XItemNCName, lexer.XItemQName, lexer.XItemNodeType}
}

func abbrPathExpr() pathexpr.PathExpr {
	return pathexpr.PathExpr{
		Name:     xml.Name{},
		Axis:     xconst.AxisDescendentOrSelf,
		NodeType: xconst.NodeTypeNode,
	}
}

func absLocPath(val string) (expTkns, XPExec) {
	ret := func(p *Parser) {
		p.ctx = p.tree
	}

	return pathStartToks(), ret
}

func abbrAbsLocPath(val string) (expTkns, XPExec) {
	ret := func(p *Parser) {
		p.ctx = p.tree
		p.pExpr = abbrPathExpr()
		p.find()
	}

	return pathStartToks(), ret
}

func relLocPath(val string) (expTkns, XPExec) {
	ret := func(p *Parser) {
	}

	return pathStartToks(), ret
}

func abbrRelLocPath(val string) (expTkns, XPExec) {
	ret := func(p *Parser) {
		p.pExpr = abbrPathExpr()
		p.find()
	}

	return pathStartToks(), ret
}

func axis(val string) (expTkns, XPExec) {
	ret := func(p *Parser) {
		p.pExpr.Axis = val
	}

	return expTkns{lexer.XItemNCName, lexer.XItemQName, lexer.XItemNodeType}, ret
}

func abbrAxis(val string) (expTkns, XPExec) {
	ret := func(p *Parser) {
		p.pExpr.Axis = xconst.AxisAttribute
	}

	return expTkns{lexer.XItemNCName, lexer.XItemQName}, ret
}

func ncName(val string) (expTkns, XPExec) {
	ret := func(p *Parser) {
		p.pExpr.Name.Space = val
	}

	return expTkns{lexer.XItemQName}, ret
}

func qName(val string) (expTkns, XPExec) {
	ret := func(p *Parser) {
		p.pExpr.Name.Local = val
	}

	return expTkns{lexer.XItemPredicate, lexer.XItemEndPath}, ret
}

func nodeType(val string) (expTkns, XPExec) {
	retFunc := func(p *Parser) {
		p.pExpr.NodeType = val
	}

	ret := expTkns{lexer.XItemPredicate, lexer.XItemEndPath}

	if val == xconst.NodeTypeProcInst {
		ret = append(ret, lexer.XItemProcLit)
	}

	return ret, retFunc
}

func procInstLit(val string) (expTkns, XPExec) {
	ret := func(p *Parser) {
		p.pExpr.ProcInstLit = val
	}

	return expTkns{lexer.XItemPredicate, lexer.XItemEndPath}, ret
}

func endPath(val string) (expTkns, XPExec) {
	ret := func(p *Parser) {
		p.find()
	}

	return expTkns{lexer.XItemRelLocPath, lexer.XItemAbbrRelLocPath}, ret
}

func (p *Parser) find() {
	vals := []tree.XPRes{}

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
		p.filter = []tree.XPRes{p.ctx}
	}

	p.pExpr.NS = p.ns

	for i := range p.filter {
		vals = append(vals, findutil.Find(p.filter[i], p.pExpr)...)
	}

	p.filter = vals
	p.pExpr = pathexpr.PathExpr{}
}
