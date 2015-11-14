package attribute

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/parser/result/pathexpr"
	"github.com/ChrisTrenkamp/goxpath/parser/result/pathres"
	"github.com/ChrisTrenkamp/goxpath/xconst"
)

//PathResAttribute is an implementation of PathRes for XML attributes
type PathResAttribute struct {
	Value  interface{}
	Parent pathres.PathRes
}

//Interface returns the data representing the attribute
func (a *PathResAttribute) Interface() interface{} {
	return a.Value
}

//GetParent returns the parent node, or itself if it's the root
func (a *PathResAttribute) GetParent() pathres.PathRes {
	return a.Parent
}

//GetChildren returns nothing
func (a *PathResAttribute) GetChildren() []pathres.PathRes {
	return []pathres.PathRes{}
}

//GetValue returns the value of the element
func (a *PathResAttribute) GetValue() string {
	//TODO: Make this return the value
	return ""
}

//Print prints the XML attribute in string form
func (a *PathResAttribute) Print(e *xml.Encoder) error {
	attr := a.Value.(*xml.Attr)
	str := attr.Name.Local + `="` + attr.Value + `"`
	if attr.Name.Space != "" {
		str = attr.Name.Space + ":" + str
	}
	pi := xml.ProcInst{
		Target: "attribute",
		Inst:   ([]byte)(str),
	}
	return e.EncodeToken(pi)
}

//EvalPath evaluates the XPath path instruction on the element
func (a *PathResAttribute) EvalPath(p *pathexpr.PathExpr) bool {
	val := a.Value.(*xml.Attr)

	if p.NodeType == "" {
		if p.Name.Space != "" {
			if p.Name.Space != "*" {
				if p.Name.Space != val.Name.Space {
					return false
				}
			}
		}

		if p.Name.Local == "*" && p.Axis == xconst.AxisAttribute {
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
