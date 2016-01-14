package xmlpi

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/parser/pathexpr"
	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/xconst"
)

//XMLPI is an implementation of XPRes for XML processing-instructions
type XMLPI struct {
	xml.ProcInst
	Parent tree.XPResEle
}

//GetParent returns the parent node
func (pi *XMLPI) GetParent() tree.XPResEle {
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
	if p.ProcInstLit != "" && p.NodeType == xconst.NodeTypeProcInst {
		return p.ProcInstLit == pi.ProcInst.Target
	}

	if p.NodeType == xconst.NodeTypeProcInst || p.NodeType == xconst.NodeTypeNode {
		return true
	}

	return false
}
