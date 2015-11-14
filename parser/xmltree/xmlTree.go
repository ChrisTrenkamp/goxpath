package xmltree

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/ChrisTrenkamp/goxpath/parser/result/attribute"
	"github.com/ChrisTrenkamp/goxpath/parser/result/chardata"
	"github.com/ChrisTrenkamp/goxpath/parser/result/comment"
	"github.com/ChrisTrenkamp/goxpath/parser/result/element"
	"github.com/ChrisTrenkamp/goxpath/parser/result/pathres"
	"github.com/ChrisTrenkamp/goxpath/parser/result/procinst"
)

//ParseXMLStr creates an XMLTree structure from an XML string
func ParseXMLStr(x string) (pathres.PathRes, error) {
	return ParseXML(bytes.NewBufferString(x))
}

//ParseXML creates an XMLTree structure from an io.Reader
func ParseXML(r io.Reader) (pathres.PathRes, error) {
	dec := xml.NewDecoder(r)
	tree := &element.PathResElement{Value: xml.StartElement{}, Children: []pathres.PathRes{}, Parent: nil}
	tree.Parent = tree
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
			ele := t.(xml.StartElement)
			attrs := make([]pathres.PathRes, len(ele.Attr))
			for i := range attrs {
				attrs[i] = &attribute.PathResAttribute{Value: &ele.Attr[i], Parent: pos}
			}
			ch := &element.PathResElement{Value: xml.CopyToken(ele), Attrs: attrs, Children: []pathres.PathRes{}, Parent: pos}
			pos.Children = append(pos.Children, ch)
			pos = ch

		case xml.CharData:
			ch := &chardata.PathResCharData{Value: xml.CopyToken(t), Parent: pos}
			pos.Children = append(pos.Children, ch)
		case xml.Comment:
			ch := &comment.PathResComment{Value: xml.CopyToken(t), Parent: pos}
			pos.Children = append(pos.Children, ch)
		case xml.ProcInst:
			if pos.Parent != pos {
				ch := &procinst.PathResProcInst{Value: xml.CopyToken(t), Parent: pos}
				pos.Children = append(pos.Children, ch)
			}
		case xml.EndElement:
			if pos.Parent == pos {
				return nil, fmt.Errorf("Malformed XML found.")
			}

			pos = pos.Parent.(*element.PathResElement)

			if pos.Parent == pos {
				done = true
			}
		}
	}

	return tree, nil
}
