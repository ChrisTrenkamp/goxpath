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
	done := false
	tree := &element.PathResElement{
		Value:    xml.StartElement{},
		NS:       make(map[xml.Name]string),
		Children: []pathres.PathRes{},
		Parent:   nil,
	}
	pos := tree

	tree.Parent = tree

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
			attrs := make([]pathres.PathRes, 0, len(ele.Attr))
			ns := make(map[xml.Name]string)

			for k, v := range pos.NS {
				ns[k] = v
			}

			for i := range ele.Attr {
				attr := ele.Attr[i].Name
				val := ele.Attr[i].Value

				if (attr.Local == "xmlns" && attr.Space == "") || attr.Space == "xmlns" {
					if attr.Local == "xmlns" && attr.Space == "" && val == "" {
						delete(ns, attr)
					} else {
						ns[attr] = val
					}
				} else {
					attrs = append(attrs, &attribute.PathResAttribute{Value: &ele.Attr[i], Parent: pos})
				}
			}

			ch := &element.PathResElement{
				Value:    xml.CopyToken(ele),
				NS:       ns,
				Attrs:    attrs,
				Children: []pathres.PathRes{},
				Parent:   pos,
			}

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
