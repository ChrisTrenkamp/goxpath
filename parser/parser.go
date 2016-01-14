package parser

import (
	"encoding/xml"
	"fmt"

	"github.com/ChrisTrenkamp/goxpath/lexer"
	"github.com/ChrisTrenkamp/goxpath/parser/findutil"
	"github.com/ChrisTrenkamp/goxpath/parser/pathexpr"
	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/xconst"
)

//Parser parses an XML document and generates output from the Lexer
type Parser struct {
	Tree   tree.XPRes
	NS     map[string]string
	ctx    tree.XPRes
	pExpr  pathexpr.PathExpr
	filter []tree.XPRes
}

type xpExec func(*Parser)

//XPathExec is the XPath executor, compiled from an XPath string
type XPathExec []xpExec

type expTkns []lexer.XItemType
type lexFn func(string) (expTkns, xpExec)

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
func Exec(xp XPathExec, t tree.XPRes, ns map[string]string) []tree.XPRes {
	if ns == nil {
		ns = make(map[string]string)
	}

	p := Parser{Tree: t, NS: ns, ctx: t}

	for _, i := range xp {
		i(&p)
	}

	return p.filter
}

//MustParse is like Parse, but panics instead of returning an error.
func MustParse(xp string) XPathExec {
	ret, err := Parse(xp)

	if err != nil {
		panic(err)
	}

	return ret
}

//Parse parses the XPath expression, xp, returning an XPath executor.
func Parse(xp string) (XPathExec, error) {
	var err error
	var next xpExec
	tok := expTkns{}
	ret := make(XPathExec, 0)
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

func eval(typ lexer.XItemType, val string, tkns expTkns) (expTkns, xpExec, error) {
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

func absLocPath(val string) (expTkns, xpExec) {
	ret := func(p *Parser) {
		p.ctx = p.Tree
	}

	return pathStartToks(), ret
}

func abbrAbsLocPath(val string) (expTkns, xpExec) {
	ret := func(p *Parser) {
		p.ctx = p.Tree
		p.pExpr = abbrPathExpr()
		p.find()
	}

	return pathStartToks(), ret
}

func relLocPath(val string) (expTkns, xpExec) {
	ret := func(p *Parser) {
	}

	return pathStartToks(), ret
}

func abbrRelLocPath(val string) (expTkns, xpExec) {
	ret := func(p *Parser) {
		p.pExpr = abbrPathExpr()
		p.find()
	}

	return pathStartToks(), ret
}

func axis(val string) (expTkns, xpExec) {
	ret := func(p *Parser) {
		p.pExpr.Axis = val
	}

	return expTkns{lexer.XItemNCName, lexer.XItemQName, lexer.XItemNodeType}, ret
}

func abbrAxis(val string) (expTkns, xpExec) {
	ret := func(p *Parser) {
		p.pExpr.Axis = xconst.AxisAttribute
	}

	return expTkns{lexer.XItemNCName, lexer.XItemQName}, ret
}

func ncName(val string) (expTkns, xpExec) {
	ret := func(p *Parser) {
		p.pExpr.Name.Space = val
	}

	return expTkns{lexer.XItemQName}, ret
}

func qName(val string) (expTkns, xpExec) {
	ret := func(p *Parser) {
		p.pExpr.Name.Local = val
	}

	return expTkns{lexer.XItemPredicate, lexer.XItemEndPath}, ret
}

func nodeType(val string) (expTkns, xpExec) {
	retFunc := func(p *Parser) {
		p.pExpr.NodeType = val
	}

	ret := expTkns{lexer.XItemPredicate, lexer.XItemEndPath}

	if val == xconst.NodeTypeProcInst {
		ret = append(ret, lexer.XItemProcLit)
	}

	return ret, retFunc
}

func procInstLit(val string) (expTkns, xpExec) {
	ret := func(p *Parser) {
		p.pExpr.ProcInstLit = val
	}

	return expTkns{lexer.XItemPredicate, lexer.XItemEndPath}, ret
}

func endPath(val string) (expTkns, xpExec) {
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

	p.pExpr.NS = p.NS

	for i := range p.filter {
		vals = append(vals, findutil.Find(p.filter[i], p.pExpr)...)
	}

	p.filter = vals
	p.pExpr = pathexpr.PathExpr{}
}
