package xmlele

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/parser/pathexpr"
	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/result/xmlattr"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/result/xmlchd"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/xmlxpres"
	"github.com/ChrisTrenkamp/goxpath/xconst"
)

//XMLEle is an implementation of XPRes for XML elements
type XMLEle struct {
	xml.StartElement
	NS       map[xml.Name]string
	Attrs    []*xmlattr.XMLAttr
	Children []xmlxpres.XMLXPRes
	Parent   xmlxpres.XMLXPResEle
}

//GetParent returns the parent node, or itself if it's the root
func (x *XMLEle) GetParent() tree.XPResEle {
	return x.Parent
}

//GetChildren returns all child nodes of the element
func (x *XMLEle) GetChildren() []tree.XPRes {
	ret := make([]tree.XPRes, len(x.Children))

	for i := range x.Children {
		ret[i] = x.Children[i]
	}

	return ret
}

//String returns the string value of the element and children
func (x *XMLEle) String() string {
	ret := ""
	for i := range x.Children {
		switch t := x.Children[i].(type) {
		case *xmlchd.XMLChd:
			ret += t.String()
		case *XMLEle:
			ret += t.String()
		}
	}
	return ret
}

//XMLPrint prints the XML element and children
func (x *XMLEle) XMLPrint(e *xml.Encoder) error {
	val := x.StartElement

	for i := 0; i < len(val.Attr); i++ {
		if val.Attr[i].Name.Local == "xmlns" || val.Attr[i].Name.Space == "xmlns" {
			val.Attr = append(val.Attr[:i], val.Attr[i+1:]...)
			i--
		}
	}

	err := e.EncodeToken(val)
	if err != nil {
		return err
	}

	for i := range x.Children {
		err = x.Children[i].XMLPrint(e)
		if err != nil {
			return err
		}
	}

	return e.EncodeToken(xml.EndElement{Name: val.Name})
}

//EvalPath evaluates the XPath path instruction on the element
func (x *XMLEle) EvalPath(p *pathexpr.PathExpr) bool {
	if p.NodeType == "" {
		if p.Name.Space != "*" {
			if x.StartElement.Name.Space != p.NS[p.Name.Space] {
				return false
			}
		}

		if p.Name.Local == "*" && p.Axis != xconst.AxisAttribute && p.Axis != xconst.AxisNamespace {
			return true
		}

		if p.Name.Local == x.StartElement.Name.Local {
			return true
		}
	} else {
		if p.NodeType == xconst.NodeTypeNode {
			return true
		}
	}

	return false
}
