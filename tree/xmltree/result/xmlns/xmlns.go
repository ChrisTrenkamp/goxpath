package xmlns

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/goxpath/pathexpr"
	"github.com/ChrisTrenkamp/goxpath/goxpath/xconst"
	"github.com/ChrisTrenkamp/goxpath/tree"
)

//XMLNS is an implementation of XPRes for XML attributes
type XMLNS struct {
	xml.Attr
	Parent tree.XPResEle
}

//GetParent returns the parent node
func (a *XMLNS) GetParent() tree.XPResEle {
	return a.Parent
}

//String returns the string value of the namespace
func (a *XMLNS) String() string {
	return a.Attr.Value
}

//XMLPrint prints the namespace as a processing-instruction.
func (a *XMLNS) XMLPrint(e *xml.Encoder) error {
	pi := xml.ProcInst{
		Target: "namespace",
		Inst:   ([]byte)(a.Attr.Value),
	}
	return e.EncodeToken(pi)
}

//EvalPath evaluates the XPath path instruction on the element
func (a *XMLNS) EvalPath(p *pathexpr.PathExpr) bool {
	if p.NodeType == "" {
		if p.Name.Space != "" && p.Name.Space != "*" {
			return false
		}

		if p.Name.Local == "*" && p.Axis == xconst.AxisNamespace {
			return true
		}

		if p.Name.Local == a.Attr.Name.Local {
			return true
		}
	} else {
		if p.NodeType == xconst.NodeTypeNode {
			return true
		}
	}

	return false
}
