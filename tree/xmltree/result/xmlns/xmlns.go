package xmlns

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/tree"
)

//XMLNS is an implementation of XPRes for XML attributes
type XMLNS struct {
	xml.Attr
	Parent tree.Elem
	tree.NodePos
}

//GetToken returns the xml.Token representation of the node
func (ns *XMLNS) GetToken() xml.Token {
	return ns.Attr
}

//GetParent returns the parent node
func (ns *XMLNS) GetParent() tree.Elem {
	return ns.Parent
}

//String returns the string value of the namespace
func (ns *XMLNS) String() string {
	return ns.Attr.Value
}
