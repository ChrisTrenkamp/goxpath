package xpath

import (
	"bytes"
	"encoding/xml"
)

//Result interface for the result of an XPath expression
type Result interface {
	XPathString() string
}

//XMLResult is returned when the result of an XPath expression is XML data
type XMLResult struct {
	Value    xml.Token
	Children []XMLResult
}

//XPathString pretty-prints the result of XML data
func (x XMLResult) XPathString() string {
	ret := bytes.NewBufferString("")
	e := xml.NewEncoder(ret)
	x.toString(e)
	err := e.Flush()
	if err != nil {
		return err.Error()
	}
	return ret.String()
}

func (x XMLResult) toString(e *xml.Encoder) {
	e.EncodeToken(x.Value)
	for i := range x.Children {
		x.Children[i].toString(e)
	}

	switch x.Value.(type) {
	case xml.StartElement:
		se := x.Value.(xml.StartElement)
		e.EncodeToken(xml.EndElement{Name: se.Name})
	}
}
