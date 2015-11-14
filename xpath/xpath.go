package xpath

import (
	"bytes"
	"encoding/xml"
	"io"

	"github.com/ChrisTrenkamp/goxpath/lexer"
	"github.com/ChrisTrenkamp/goxpath/parser"
	"github.com/ChrisTrenkamp/goxpath/parser/result/pathres"
)

//FromStr runs an XPath expression on the XML string
func FromStr(xpath, x string) ([]pathres.PathRes, error) {
	it := lexer.Lex(xpath)
	p, err := parser.CreateParserStr(x)

	if err != nil {
		return nil, err
	}

	return p.Parse(it)
}

//FromReader runs an XPath expression on the XML reader
func FromReader(xpath string, r io.Reader) ([]pathres.PathRes, error) {
	it := lexer.Lex(xpath)
	p, err := parser.CreateParser(r)

	if err != nil {
		return nil, err
	}

	return p.Parse(it)
}

//Print prints out the XPath result
func Print(r pathres.PathRes) (string, error) {
	ret := bytes.NewBufferString("")
	e := xml.NewEncoder(ret)
	err := r.Print(e)
	if err != nil {
		return "", err
	}

	err = e.Flush()
	if err != nil {
		return "", err
	}

	return ret.String(), nil
}
