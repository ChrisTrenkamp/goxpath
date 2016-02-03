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

//XMLPrint prints the namespace as a processing-instruction.
func (ns *XMLNS) XMLPrint(e *xml.Encoder) error {
	pi := xml.ProcInst{
		Target: "namespace",
		Inst:   ([]byte)(ns.Attr.Value),
	}
	return e.EncodeToken(pi)
}

//EvalPath evaluates the XPath path instruction on the element
func (ns *XMLNS) EvalPath(p *pathexpr.PathExpr) bool {
	if p.NodeType == "" {
		if p.Name.Space != "" && p.Name.Space != "*" {
			return false
		}

		if p.Name.Local == "*" && p.Axis == xconst.AxisNamespace {
			return true
		}

		if p.Name.Local == ns.Attr.Name.Local {
			return true
		}
	} else {
		if p.NodeType == xconst.NodeTypeNode {
			return true
		}
	}

	return false
}
