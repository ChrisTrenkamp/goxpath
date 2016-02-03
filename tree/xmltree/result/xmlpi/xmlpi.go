package xmlpi

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/goxpath/pathexpr"
	"github.com/ChrisTrenkamp/goxpath/goxpath/xconst"
	"github.com/ChrisTrenkamp/goxpath/tree"
)

//XMLPI is an implementation of XPRes for XML processing-instructions
type XMLPI struct {
	xml.ProcInst
	Parent tree.Elem
	tree.NodePos
}

//GetToken returns the xml.Token representation of the node
func (pi *XMLPI) GetToken() xml.Token {
	return pi.ProcInst
}

//GetParent returns the parent node
func (pi *XMLPI) GetParent() tree.Elem {
	return pi.Parent
}

//String returns the value of the processing-instruction
func (pi *XMLPI) String() string {
	return string(pi.ProcInst.Inst)
}

//XMLPrint prints the XML processing-instruction
func (pi *XMLPI) XMLPrint(e *xml.Encoder) error {
	return e.EncodeToken(pi.ProcInst)
}

//EvalPath evaluates the XPath path instruction on the processing-instruction
func (pi *XMLPI) EvalPath(p *pathexpr.PathExpr) bool {
	if p.NodeType == xconst.NodeTypeProcInst {
		return true
	}

	if p.NodeType == xconst.NodeTypeProcInst || p.NodeType == xconst.NodeTypeNode {
		return true
	}

	return false
}
