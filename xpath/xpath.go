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
func FromStr(xpath, x string, ns map[string]string) ([]pathres.PathRes, error) {
	return FromReader(xpath, bytes.NewBufferString(x), ns)
}

//FromReader runs an XPath expression on the XML reader
func FromReader(xpath string, r io.Reader, ns map[string]string) ([]pathres.PathRes, error) {
	it := lexer.Lex(xpath)
	p, err := parser.CreateParser(r, ns)

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
