package parser

import (
	"encoding/xml"
	"fmt"

	"github.com/ChrisTrenkamp/goxpath/lexer"
	"github.com/ChrisTrenkamp/goxpath/xpath"
)

//Parser parses an XML document and generates output from the Lexer
type Parser struct {
	tree   *xmlTree
	ctx    *xmlTree
	pExpr  pathExpr
	filter []*xmlTree
}

type pathExpr struct {
	name     xml.Name
	axis     string
	abbrAxis string
	abbr     bool
	nodeType string
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
	lexer.XItemEndPath:        endPath,
}

//Parse generates output from the Lexer
func (p *Parser) Parse(c chan lexer.XItem) ([]xpath.Result, error) {
	var err error
	tok := expTkns{}

	for item := range c {
		if item.Typ == lexer.XItemError {
			return []xpath.Result{}, fmt.Errorf(item.Val)
		}

		tok, err = p.eval(item.Typ, item.Val, tok...)

		if err != nil {
			return []xpath.Result{}, err
		}
	}

	return p.createRes()
}

func (p *Parser) eval(typ lexer.XItemType, val string, tkns ...lexer.XItemType) (expTkns, error) {
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
		fmt.Println("INVALID TOKEN FOUND")
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

func absLocPath(p *Parser, val string) (expTkns, error) {
	p.ctx = p.tree
	p.pExpr = pathExpr{abbr: false}
	return pathStartToks(), nil
}

func abbrAbsLocPath(p *Parser, val string) (expTkns, error) {
	p.ctx = p.tree
	p.pExpr = pathExpr{abbr: true}
	return pathStartToks(), nil
}

func relLocPath(p *Parser, val string) (expTkns, error) {
	p.pExpr = pathExpr{abbr: false}
	return pathStartToks(), nil
}

func abbrRelLocPath(p *Parser, val string) (expTkns, error) {
	p.pExpr = pathExpr{abbr: true}
	return pathStartToks(), nil
}

func axis(p *Parser, val string) (expTkns, error) {
	p.pExpr.axis = val
	return expTkns{lexer.XItemNCName, lexer.XItemQName, lexer.XItemNodeType}, nil
}

func abbrAxis(p *Parser, val string) (expTkns, error) {
	p.pExpr.abbrAxis = val
	return expTkns{lexer.XItemNCName, lexer.XItemQName}, nil
}

func ncName(p *Parser, val string) (expTkns, error) {
	p.pExpr.name.Space = val
	return expTkns{lexer.XItemQName}, nil
}

func qName(p *Parser, val string) (expTkns, error) {
	p.pExpr.name.Local = val
	return expTkns{lexer.XItemPredicate, lexer.XItemEndPath}, nil
}

func nodeType(p *Parser, val string) (expTkns, error) {
	p.pExpr.nodeType = val
	return expTkns{lexer.XItemPredicate, lexer.XItemEndPath}, nil
}

func endPath(p *Parser, val string) (expTkns, error) {
	vals := []*xmlTree{}
	if p.filter == nil {
		p.filter = []*xmlTree{p.ctx}
	}

	for i := range p.filter {
		vals = append(vals, p.filter[i].findTag(p.pExpr)...)
	}

	p.filter = vals
	p.pExpr = pathExpr{}
	return expTkns{lexer.XItemRelLocPath, lexer.XItemAbbrRelLocPath}, nil
}
