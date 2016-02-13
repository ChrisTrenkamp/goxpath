package xmlchd

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/tree"
)

//XMLChd is an implementation of XPRes for XML attributes
type XMLChd struct {
	xml.CharData
	Parent tree.Elem
	tree.NodePos
}

//GetToken returns the xml.Token representation of the node
func (cd *XMLChd) GetToken() xml.Token {
	return cd.CharData
}

//GetParent returns the parent node
func (cd *XMLChd) GetParent() tree.Elem {
	return cd.Parent
}

//String returns the value of the character data
func (cd *XMLChd) String() string {
	return string(cd.CharData)
}
