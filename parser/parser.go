package parser

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/ChrisTrenkamp/goxpath/lexer"
	"github.com/ChrisTrenkamp/goxpath/parser/result"
	"github.com/ChrisTrenkamp/goxpath/parser/result/pathexpr"
	"github.com/ChrisTrenkamp/goxpath/parser/result/pathres"
	"github.com/ChrisTrenkamp/goxpath/parser/xmltree"
	"github.com/ChrisTrenkamp/goxpath/xconst"
)

//Parser parses an XML document and generates output from the Lexer
type Parser struct {
	tree   pathres.PathRes
	ctx    pathres.PathRes
	pExpr  pathexpr.PathExpr
	filter []pathres.PathRes
	ns     map[string]string
}

type expTkns []lexer.XItemType
type lexFn func(*Parser, string) (expTkns, error)

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

//CreateParser creates a Parser from a XML reader
func CreateParser(r io.Reader, nsLookup map[string]string) (Parser, error) {
	t, err := xmltree.ParseXML(r)

	return Parser{tree: t, ctx: t, ns: nsLookup}, err
}

//Parse generates output from the Lexer
func (p *Parser) Parse(c chan lexer.XItem) ([]pathres.PathRes, error) {
	var err error
	tok := expTkns{}

	for item := range c {
		if item.Typ == lexer.XItemError {
			return []pathres.PathRes{}, fmt.Errorf(item.Val)
		}

		tok, err = p.eval(item.Typ, item.Val, tok)

		if err != nil {
			return []pathres.PathRes{}, err
		}
	}

	return p.filter, nil
}

func (p *Parser) eval(typ lexer.XItemType, val string, tkns expTkns) (expTkns, error) {
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
		return expTkns{}, fmt.Errorf("Unexpected token: %d", typ)
	}

	if f, ok := parseMap[typ]; ok {
		return f(p, val)
	}

	return expTkns{}, fmt.Errorf("Unsupported token emitted: %d", typ)
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

func absLocPath(p *Parser, val string) (expTkns, error) {
	p.ctx = p.tree
	return pathStartToks(), nil
}

func abbrAbsLocPath(p *Parser, val string) (expTkns, error) {
	p.ctx = p.tree
	p.pExpr = abbrPathExpr()
	p.find()
	return pathStartToks(), nil
}

func relLocPath(p *Parser, val string) (expTkns, error) {
	return pathStartToks(), nil
}

func abbrRelLocPath(p *Parser, val string) (expTkns, error) {
	p.pExpr = abbrPathExpr()
	p.find()
	return pathStartToks(), nil
}

func axis(p *Parser, val string) (expTkns, error) {
	p.pExpr.Axis = val
	return expTkns{lexer.XItemNCName, lexer.XItemQName, lexer.XItemNodeType}, nil
}

func abbrAxis(p *Parser, val string) (expTkns, error) {
	p.pExpr.Axis = xconst.AxisAttribute
	return expTkns{lexer.XItemNCName, lexer.XItemQName}, nil
}

func ncName(p *Parser, val string) (expTkns, error) {
	p.pExpr.Name.Space = val
	return expTkns{lexer.XItemQName}, nil
}

func qName(p *Parser, val string) (expTkns, error) {
	p.pExpr.Name.Local = val
	return expTkns{lexer.XItemPredicate, lexer.XItemEndPath}, nil
}

func nodeType(p *Parser, val string) (expTkns, error) {
	p.pExpr.NodeType = val
	ret := expTkns{lexer.XItemPredicate, lexer.XItemEndPath}
	if val == xconst.NodeTypeProcInst {
		ret = append(ret, lexer.XItemProcLit)
	}
	return ret, nil
}

func procInstLit(p *Parser, val string) (expTkns, error) {
	p.pExpr.ProcInstLit = val
	return expTkns{lexer.XItemPredicate, lexer.XItemEndPath}, nil
}

func endPath(p *Parser, val string) (expTkns, error) {
	p.find()
	return expTkns{lexer.XItemRelLocPath, lexer.XItemAbbrRelLocPath}, nil
}

func (p *Parser) find() {
	vals := []pathres.PathRes{}

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
		p.filter = []pathres.PathRes{p.ctx}
	}

	p.pExpr.NS = p.ns

	for i := range p.filter {
		vals = append(vals, result.Find(p.filter[i], p.pExpr)...)
	}

	p.filter = vals
	p.pExpr = pathexpr.PathExpr{}
}
