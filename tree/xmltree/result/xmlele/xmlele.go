package xmlele

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/goxpath/pathexpr"
	"github.com/ChrisTrenkamp/goxpath/goxpath/xconst"
	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/result/xmlattr"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/result/xmlchd"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/result/xmlns"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/xmlres"
)

//XMLEle is an implementation of XPRes for XML elements
type XMLEle struct {
	xml.StartElement
	NS       []*xmlns.XMLNS
	Attrs    []*xmlattr.XMLAttr
	Children []xmlres.XMLNode
	Parent   xmlres.XMLElem
	tree.NodePos
}

//GetToken returns the xml.Token representation of the node
func (x *XMLEle) GetToken() xml.Token {
	return x.StartElement
}

//GetParent returns the parent node, or itself if it's the root
func (x *XMLEle) GetParent() tree.Elem {
	return x.Parent
}

//GetChildren returns all child nodes of the element
func (x *XMLEle) GetChildren() []tree.Node {
	ret := make([]tree.Node, len(x.Children))

	for i := range x.Children {
		ret[i] = x.Children[i]
	}

	return ret
}

//GetNSAttrs returns all namespaces and attributes of the element
func (x *XMLEle) GetNSAttrs() []tree.Node {
	ret := make([]tree.Node, 0, len(x.NS)+len(x.Attrs))
	for i := range x.NS {
		ret = append(ret, x.NS[i])
	}
	for i := range x.Attrs {
		ret = append(ret, x.Attrs[i])
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
		return x.checkNameAndSpace(p)
	}

	if p.NodeType == xconst.NodeTypeNode {
		return true
	}

	return false
}

func (x *XMLEle) checkNameAndSpace(p *pathexpr.PathExpr) bool {
	if p.Name.Local == "*" && p.Name.Space == "" {
		return true
	}

	if p.Name.Space != "*" && x.StartElement.Name.Space != p.NS[p.Name.Space] {
		return false
	}

	if p.Name.Local == "*" && p.Axis != xconst.AxisAttribute && p.Axis != xconst.AxisNamespace {
		return true
	}

	if p.Name.Local == x.StartElement.Name.Local {
		return true
	}

	return false
}
