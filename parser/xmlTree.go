package parser

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
)

type xmlTree struct {
	value    xml.Token
	children []*xmlTree
	parent   *xmlTree
}

type rootDoc struct{}

//CreateParserStr creates a Parser from an XML string
func CreateParserStr(x string) (Parser, error) {
	t, err := parseXMLString(x)

	if err != nil {
		return Parser{}, err
	}

	return Parser{tree: t, ctx: t}, err
}

//CreateParser creates a Parser from a XML reader
func CreateParser(r io.Reader) (Parser, error) {
	t, err := parseXML(r)

	if err != nil {
		return Parser{}, err
	}

	return Parser{tree: t, ctx: t}, err
}

func parseXMLString(x string) (*xmlTree, error) {
	return parseXML(bytes.NewBufferString(x))
}

func attachTree(data xml.Token, parent *xmlTree) *xmlTree {
	//Create a copy because the token is reused by the decoder
	//If not copied, there will be corrupt data
	d := xml.CopyToken(data)
	ret := &xmlTree{
		value:    d,
		children: []*xmlTree{},
		parent:   parent,
	}

	if parent != nil {
		parent.children = append(parent.children, ret)
	}

	return ret
}

func parseXML(r io.Reader) (*xmlTree, error) {
	dec := xml.NewDecoder(r)
	tree := attachTree(rootDoc{}, nil)
	pos := tree
	done := false

	for !done {
		t, err := dec.Token()

		if err != nil {
			return nil, err
		}

		if t == nil {
			break
		}

		switch t.(type) {
		case xml.StartElement:
			ch := attachTree(t, pos)
			pos = ch

		case xml.CharData, xml.Comment, xml.Directive, xml.ProcInst:
			attachTree(t, pos)

		case xml.EndElement:
			if pos.parent == nil {
				return nil, fmt.Errorf("Malformed XML found.")
			}

			pos = pos.parent

			if pos.parent == nil {
				done = true
			}
		}
	}

	return tree, nil
}
