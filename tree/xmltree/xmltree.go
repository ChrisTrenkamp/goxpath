package xmltree

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/internal"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/result/xmlattr"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/result/xmlchd"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/result/xmlcomm"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/result/xmlele"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/result/xmlpi"
)

//TreeOverride gives the option to override the tree with a customized structure
type TreeOverride interface {
	Root() tree.Elem
	StartElem(ele *xmlele.XMLEle, pos tree.Elem, dec *xml.Decoder) tree.Elem
	AddNode(n tree.Node, pos tree.Elem, dec *xml.Decoder)
	EndElem(ele xml.EndElement, pos tree.Elem, dec *xml.Decoder) tree.Elem
}

//MustParseXML is like ParseXML, but panics instead of returning an error.
func MustParseXML(r io.Reader, ov TreeOverride) tree.Node {
	ret, err := ParseXML(r, ov)

	if err != nil {
		panic(err)
	}

	return ret
}

//ParseXML creates an XMLTree structure from an io.Reader.
func ParseXML(r io.Reader, ov TreeOverride) (tree.Node, error) {
	if ov == nil {
		ov = defoverride.DefOverride{}
	}

	dec := xml.NewDecoder(r)
	ordrPos := 1
	xmlTree := ov.Root()
	pos := xmlTree

	t, err := dec.Token()

	if err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("Premature end of XML file")
		}
		return nil, err
	}

	switch t := t.(type) {
	case xml.ProcInst:
		if t.Target != "xml" {
			return nil, fmt.Errorf("Malformed XML file")
		}
	default:
		return nil, fmt.Errorf("Malformed XML file")
	}

	t, err = dec.Token()

	for err == nil {
		switch xt := t.(type) {
		case xml.StartElement:
			ch := createEle(pos, xt, &ordrPos)
			pos = ov.StartElem(ch, pos, dec)
		case xml.CharData:
			ch := &xmlchd.XMLChd{CharData: xml.CopyToken(t).(xml.CharData), Parent: pos, NodePos: tree.NodePos(ordrPos)}
			ov.AddNode(ch, pos, dec)
			ordrPos++
		case xml.Comment:
			ch := &xmlcomm.XMLComm{Comment: xml.CopyToken(t).(xml.Comment), Parent: pos, NodePos: tree.NodePos(ordrPos)}
			ov.AddNode(ch, pos, dec)
			ordrPos++
		case xml.ProcInst:
			ch := &xmlpi.XMLPI{ProcInst: xml.CopyToken(t).(xml.ProcInst), Parent: pos, NodePos: tree.NodePos(ordrPos)}
			ov.AddNode(ch, pos, dec)
			ordrPos++
		case xml.EndElement:
			pos = ov.EndElem(xt, pos, dec)
		}

		t, err = dec.Token()
	}

	if err == io.EOF {
		err = nil
	}

	return xmlTree, err
}

func createEle(pos tree.Elem, ele xml.StartElement, ordrPos *int) *xmlele.XMLEle {
	ch := &xmlele.XMLEle{
		StartElement: xml.CopyToken(ele).(xml.StartElement),
		NSStruct:     tree.NSStruct{NS: make(map[xml.Name]tree.NS)},
		Children:     []tree.Node{},
		Parent:       pos,
		NodePos:      tree.NodePos(*ordrPos),
	}
	*ordrPos++

	ns := make(map[xml.Name]tree.NS)

	if nselem, ok := pos.(tree.NSElem); ok {
		ns = createNS(nselem)
	}

	attrs := make([]*xmlattr.XMLAttr, 0, len(ele.Attr))

	for i := range ele.Attr {
		attr := ele.Attr[i].Name
		val := ele.Attr[i].Value

		if (attr.Local == "xmlns" && attr.Space == "") || attr.Space == "xmlns" {
			if attr.Local == "xmlns" && attr.Space == "" && val == "" {
				delete(ns, attr)
			} else {
				ns[attr] = tree.NS{Attr: xml.Attr{Name: attr, Value: val}}
			}
		} else {
			attrs = append(attrs, &xmlattr.XMLAttr{Attr: &ele.Attr[i], Parent: ch})
		}
	}

	for k, v := range ns {
		ch.NSStruct.NS[k] = tree.NS{Attr: v.Attr, Parent: ch, NodePos: tree.NodePos(*ordrPos)}
		*ordrPos++
	}

	ch.Attrs = attrs

	for _, i := range ch.Attrs {
		i.NodePos = tree.NodePos(*ordrPos)
		*ordrPos++
	}

	return ch
}

func createNS(pos tree.NSElem) map[xml.Name]tree.NS {
	ns := pos.GetNS()

	xns := xml.Name{Space: "", Local: "xml"}
	ns[xns] = tree.NS{Attr: xml.Attr{Name: xns, Value: "http://www.w3.org/XML/1998/namespace"}}
	return ns
}
