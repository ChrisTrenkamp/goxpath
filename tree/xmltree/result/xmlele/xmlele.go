package xmlele

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/result/xmlattr"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/result/xmlchd"
)

//XMLEle is an implementation of XPRes for XML elements
type XMLEle struct {
	xml.StartElement
	*tree.NSStruct
	Attrs    []*xmlattr.XMLAttr
	Children []tree.Node
	Parent   tree.Elem
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

//GetAttrs returns all attributes of the element
func (x *XMLEle) GetAttrs() []tree.Node {
	ret := make([]tree.Node, len(x.Attrs))
	for i := range x.Attrs {
		ret[i] = x.Attrs[i]
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
		case tree.Elem:
			ret += t.String()
		}
	}
	return ret
}
