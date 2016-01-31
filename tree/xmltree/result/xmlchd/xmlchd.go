package xmlchd

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/goxpath/pathexpr"
	"github.com/ChrisTrenkamp/goxpath/goxpath/xconst"
	"github.com/ChrisTrenkamp/goxpath/tree"
)

//XMLChd is an implementation of XPRes for XML attributes
type XMLChd struct {
	xml.CharData
	Parent tree.Elem
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

//XMLPrint prints the character data as a processing-instruction.
func (cd *XMLChd) XMLPrint(e *xml.Encoder) error {
	return e.EncodeToken(cd.CharData)
}

//EvalPath evaluates the XPath path instruction on the character data
func (cd *XMLChd) EvalPath(p *pathexpr.PathExpr) bool {
	if p.NodeType == xconst.NodeTypeText || p.NodeType == xconst.NodeTypeNode {
		return true
	}

	return false
}
