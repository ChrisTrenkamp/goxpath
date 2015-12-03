package namespace

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/parser/result/pathexpr"
	"github.com/ChrisTrenkamp/goxpath/parser/result/pathres"
	"github.com/ChrisTrenkamp/goxpath/xconst"
)

//PathResNamespace is an implementation of PathRes for XML attributes
type PathResNamespace struct {
	Value  xml.Attr
	Parent pathres.PathRes
}

//Interface returns the data representing the attribute
func (a *PathResNamespace) Interface() interface{} {
	return a.Value
}

//GetParent returns the parent node
func (a *PathResNamespace) GetParent() pathres.PathRes {
	return a.Parent
}

//GetChildren returns nothing
func (a *PathResNamespace) GetChildren() []pathres.PathRes {
	return []pathres.PathRes{}
}

//GetValue returns the namespace value
func (a *PathResNamespace) GetValue() string {
	return a.Value.Value
}

//Print prints the XML attribute in string form
func (a *PathResNamespace) Print(e *xml.Encoder) error {
	pi := xml.ProcInst{
		Target: "namespace",
		Inst:   ([]byte)(a.Value.Value),
	}
	return e.EncodeToken(pi)
}

//EvalPath evaluates the XPath path instruction on the element
func (a *PathResNamespace) EvalPath(p *pathexpr.PathExpr) bool {
	val := a.Value

	if p.NodeType == "" {
		if p.Name.Space != "" {
			if p.Name.Space != "*" {
				if p.Name.Space != val.Name.Space {
					return false
				}
			}
		}

		if p.Name.Local == "*" && p.Axis == xconst.AxisNamespace {
			return true
		}

		if p.Name.Local == val.Name.Local {
			return true
		}
	} else {
		if p.NodeType == xconst.NodeTypeNode {
			return true
		}
	}

	return false
}
