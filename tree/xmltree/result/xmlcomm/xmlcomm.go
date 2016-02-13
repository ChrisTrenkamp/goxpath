package xmlcomm

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/tree"
)

//XMLComm is an implementation of XPRes for XML attributes
type XMLComm struct {
	xml.Comment
	Parent tree.Elem
	tree.NodePos
}

//GetToken returns the xml.Token representation of the node
func (c *XMLComm) GetToken() xml.Token {
	return c.Comment
}

//GetParent returns the parent node
func (c *XMLComm) GetParent() tree.Elem {
	return c.Parent
}

//String returns the value of the comment
func (c *XMLComm) String() string {
	return string(c.Comment)
}
