package xmltree

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/result/xmlattr"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/result/xmlchd"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/result/xmlcomm"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/result/xmlele"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/result/xmlpi"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/xmlres"
)

//MustParseXML is like ParseXML, but panics instead of returning an error.
func MustParseXML(r io.Reader) xmlres.XMLNode {
	ret, err := ParseXML(r)

	if err != nil {
		panic(err)
	}

	return ret
}

//ParseXML creates an XMLTree structure from an io.Reader.
func ParseXML(r io.Reader) (xmlres.XMLNode, error) {
	dec := xml.NewDecoder(r)
	done := false
	tree := &xmlele.XMLEle{
		StartElement: xml.StartElement{},
		NS:           make(map[xml.Name]string),
		Children:     []xmlres.XMLNode{},
		Parent:       nil,
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
			ns := make(map[xml.Name]string)

			for k, v := range pos.NS {
				ns[k] = v
			}

			ns[xml.Name{Space: "", Local: "xml"}] = "http://www.w3.org/XML/1998/namespace"

			ch := &xmlele.XMLEle{
				StartElement: xml.CopyToken(ele).(xml.StartElement),
				NS:           ns,
				Children:     []xmlres.XMLNode{},
				Parent:       pos,
			}

			attrs := make([]*xmlattr.XMLAttr, 0, len(ele.Attr))

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
					attrs = append(attrs, &xmlattr.XMLAttr{Attr: ele.Attr[i], Parent: ch})
				}
			}

			ch.Attrs = attrs
			pos.Children = append(pos.Children, ch)
			pos = ch

		case xml.CharData:
			ch := &xmlchd.XMLChd{CharData: xml.CopyToken(t).(xml.CharData), Parent: pos}
			pos.Children = append(pos.Children, ch)
		case xml.Comment:
			ch := &xmlcomm.XMLComm{Comment: xml.CopyToken(t).(xml.Comment), Parent: pos}
			pos.Children = append(pos.Children, ch)
		case xml.ProcInst:
			if pos.Parent != pos {
				ch := &xmlpi.XMLPI{ProcInst: xml.CopyToken(t).(xml.ProcInst), Parent: pos}
				pos.Children = append(pos.Children, ch)
			}
		case xml.EndElement:
			if pos.Parent == pos {
				return nil, fmt.Errorf("Malformed XML found.")
			}

			pos = pos.Parent.(*xmlele.XMLEle)

			if pos.Parent == pos {
				done = true
			}
		}
	}

	return tree, nil
}

//Marshal prints the result tree, r, in XML form to w.
func Marshal(r xmlres.XMLPrinter, w io.Writer) error {
	return marshal(r, w)
}

//MarshalStr is like Marhal, but returns a string.
func MarshalStr(r xmlres.XMLPrinter) (string, error) {
	ret := bytes.NewBufferString("")
	err := marshal(r, ret)

	return ret.String(), err
}

func marshal(r xmlres.XMLPrinter, w io.Writer) error {
	e := xml.NewEncoder(w)
	err := r.XMLPrint(e)
	if err != nil {
		return err
	}

	err = e.Flush()
	if err != nil {
		return err
	}

	return nil
}
