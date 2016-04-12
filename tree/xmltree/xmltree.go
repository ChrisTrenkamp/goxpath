package xmltree

import (
	"encoding/xml"
	"fmt"
	"io"

	"golang.org/x/net/html/charset"

	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/internal"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/result/xmlele"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/result/xmlnode"
)

//ParseOptions is a set of methods and function pointers that alter
//the way the XML decoder works and the Node types that are created.
//Options that are not set will default to what is set in internal/defoverride.go
type ParseOptions struct {
	Strict    bool
	RootNode  func() tree.Elem
	StartElem func(ele *xmlele.XMLEle, pos tree.Elem, dec *xml.Decoder) tree.Elem
	Node      func(n xmlnode.XMLNode, pos tree.Elem, dec *xml.Decoder)
	EndElem   func(ele xml.EndElement, pos tree.Elem, dec *xml.Decoder) tree.Elem
	Directive func(dir xml.Directive, dec *xml.Decoder)
}

//ParseSettings is a function for setting the ParseOptions you want when
//parsing an XML tree.
type ParseSettings func(s *ParseOptions)

//MustParseXML is like ParseXML, but panics instead of returning an error.
func MustParseXML(r io.Reader, op ...ParseSettings) tree.Node {
	ret, err := ParseXML(r, op...)

	if err != nil {
		panic(err)
	}

	return ret
}

//ParseXML creates an XMLTree structure from an io.Reader.
func ParseXML(r io.Reader, op ...ParseSettings) (tree.Node, error) {
	ov := ParseOptions{
		Strict:    true,
		RootNode:  defoverride.RootNode,
		StartElem: defoverride.StartElem,
		Node:      defoverride.Node,
		EndElem:   defoverride.EndElem,
		Directive: defoverride.Directive,
	}
	for _, i := range op {
		i(&ov)
	}

	dec := xml.NewDecoder(r)
	dec.CharsetReader = charset.NewReaderLabel
	dec.Strict = ov.Strict

	ordrPos := 1
	xmlTree := ov.RootNode()
	pos := xmlTree

	t, err := dec.Token()

	if err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("Premature end of XML file")
		}
		return nil, err
	}

	brokenHeader := false
	switch t := t.(type) {
	case xml.ProcInst:
		if t.Target != "xml" {
			brokenHeader = true
		}
	default:
		brokenHeader = true
	}

	if brokenHeader {
		if ov.Strict {
			return nil, fmt.Errorf("Malformed XML file")
		}
	} else {
		t, err = dec.Token()
	}

	for err == nil {
		switch xt := t.(type) {
		case xml.StartElement:
			ch := createEle(pos, xt.Copy(), &ordrPos)
			pos = ov.StartElem(ch, pos, dec)
		case xml.CharData:
			ch := xmlnode.XMLNode{Token: xt.Copy(), Parent: pos, NodePos: tree.NodePos(ordrPos), NodeType: tree.NtChd}
			ov.Node(ch, pos, dec)
			ordrPos++
		case xml.Comment:
			ch := xmlnode.XMLNode{Token: xt.Copy(), Parent: pos, NodePos: tree.NodePos(ordrPos), NodeType: tree.NtComm}
			ov.Node(ch, pos, dec)
			ordrPos++
		case xml.ProcInst:
			ch := xmlnode.XMLNode{Token: xt.Copy(), Parent: pos, NodePos: tree.NodePos(ordrPos), NodeType: tree.NtPi}
			ov.Node(ch, pos, dec)
			ordrPos++
		case xml.EndElement:
			pos = ov.EndElem(xt, pos, dec)
		case xml.Directive:
			ov.Directive(xt.Copy(), dec)
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
		StartElement: ele,
		NSStruct:     &tree.NSStruct{NS: make(map[xml.Name]tree.NS)},
		Children:     []tree.Node{},
		Parent:       pos,
		NodePos:      tree.NodePos(*ordrPos),
		NodeType:     tree.NtEle,
	}
	*ordrPos++

	ch.NSStruct.Elem = ch
	if nselem, ok := pos.(tree.NSElem); ok {
		ch.NSStruct.Parent = nselem.GetNS()
	}

	if pos.GetNodeType() == tree.NtRoot {
		xns := xml.Name{Space: "", Local: "xml"}
		ch.NSStruct.NS[xns] = tree.NS{Attr: xml.Attr{Name: xns, Value: "http://www.w3.org/XML/1998/namespace"}}
	}

	attrs := make([]xmlnode.XMLNode, 0, len(ele.Attr))

	for i := range ele.Attr {
		attr := ele.Attr[i].Name
		val := ele.Attr[i].Value

		if (attr.Local == "xmlns" && attr.Space == "") || attr.Space == "xmlns" {
			ch.NSStruct.NS[attr] = tree.NS{Attr: xml.Attr{Name: attr, Value: val}}
		} else {
			attrs = append(attrs, xmlnode.XMLNode{Token: &ele.Attr[i], Parent: ch, NodeType: tree.NtAttr})
		}
	}

	ch.Attrs = attrs

	nsLen := len(ch.BuildNS())

	for i := range ch.Attrs {
		attr := &ch.Attrs[i]
		attr.NodePos = tree.NodePos(int(*ordrPos) + nsLen + i + 1)
		*ordrPos++
	}

	return ch
}
