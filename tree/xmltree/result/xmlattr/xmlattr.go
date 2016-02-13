package xmlattr

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/tree"
)

//XMLAttr is an implementation of XPRes for XML attributes
type XMLAttr struct {
	xml.Attr
	Parent tree.Elem
	tree.NodePos
}

//GetToken returns the xml.Token representation of the node
func (a *XMLAttr) GetToken() xml.Token {
	return a.Attr
}

//GetParent returns the parent node
func (a *XMLAttr) GetParent() tree.Elem {
	return a.Parent
}

//String returns the string value of the attribute
func (a *XMLAttr) String() string {
	return a.Attr.Value
}
