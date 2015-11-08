package eval

import (
	"io"

	"github.com/ChrisTrenkamp/goxpath/lexer"
	"github.com/ChrisTrenkamp/goxpath/parser"
	"github.com/ChrisTrenkamp/goxpath/xpath"
)

//FromStr runs an XPath expression on the XML string
func FromStr(xpath, x string) ([]xpath.Result, error) {
	it := lexer.Lex(xpath)
	p, err := parser.CreateParserStr(x)

	if err != nil {
		return nil, err
	}

	return p.Parse(it)
}

//FromReader runs an XPath expression on the XML reader
func FromReader(xpath string, r io.Reader) ([]xpath.Result, error) {
	it := lexer.Lex(xpath)
	p, err := parser.CreateParser(r)

	if err != nil {
		return nil, err
	}

	return p.Parse(it)
}
