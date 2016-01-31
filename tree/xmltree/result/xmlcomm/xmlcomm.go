package xmlcomm

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/goxpath/pathexpr"
	"github.com/ChrisTrenkamp/goxpath/goxpath/xconst"
	"github.com/ChrisTrenkamp/goxpath/tree"
)

//XMLComm is an implementation of XPRes for XML attributes
type XMLComm struct {
	xml.Comment
	Parent tree.Elem
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

//XMLPrint prints the comment in XML form
func (c *XMLComm) XMLPrint(e *xml.Encoder) error {
	return e.EncodeToken(c.Comment)
}

//EvalPath evaluates the XPath path instruction on the element
func (c *XMLComm) EvalPath(p *pathexpr.PathExpr) bool {
	if p.NodeType == xconst.NodeTypeComment || p.NodeType == xconst.NodeTypeNode {
		return true
	}

	return false
}
