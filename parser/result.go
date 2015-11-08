package parser

import "github.com/ChrisTrenkamp/goxpath/xpath"

func (p *Parser) createRes() ([]xpath.Result, error) {
	res := []xpath.Result{}

	for i := range p.filter {
		res = append(res, createXMLRes(p.filter[i]))
	}

	return res, nil
}

func createXMLRes(x *xmlTree) xpath.XMLResult {
	ret := xpath.XMLResult{
		Value:    x.value,
		Children: make([]xpath.XMLResult, len(x.children)),
	}

	for i := range x.children {
		ret.Children[i] = createXMLRes(x.children[i])
	}

	return ret
}
