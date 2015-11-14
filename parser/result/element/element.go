package element

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/parser/result/pathexpr"
	"github.com/ChrisTrenkamp/goxpath/parser/result/pathres"
	"github.com/ChrisTrenkamp/goxpath/xconst"
)

//PathResElement is an implementation of PathRes for XML elements
type PathResElement struct {
	Value    interface{}
	Attrs    []pathres.PathRes
	Children []pathres.PathRes
	Parent   pathres.PathRes
}

//Interface returns the data representing the element
func (x *PathResElement) Interface() interface{} {
	return x.Value
}

//GetParent returns the parent node, or itself if it's the root
func (x *PathResElement) GetParent() pathres.PathRes {
	return x.Parent
}

//GetChildren returns all child nodes of the element
func (x *PathResElement) GetChildren() []pathres.PathRes {
	return x.Children
}

//GetValue returns all text nodes in the element, as per the spec
func (x *PathResElement) GetValue() string {
	//TODO: Replace with a text-node query.  Return all nodes if root
	return ""
}

//Print prints the XML element in string form
func (x *PathResElement) Print(e *xml.Encoder) error {
	var err error
	if _, ok := x.Value.(xml.StartElement); ok {
		val := x.Value.(xml.StartElement)
		err = e.EncodeToken(val)
	}

	if err != nil {
		return err
	}

	for i := range x.Children {
		err = x.Children[i].Print(e)
		if err != nil {
			return err
		}
	}

	if _, ok := x.Value.(xml.StartElement); ok {
		err = e.EncodeToken(xml.EndElement{Name: x.Value.(xml.StartElement).Name})
	}

	return err
}

//EvalPath evaluates the XPath path instruction on the element
func (x *PathResElement) EvalPath(p *pathexpr.PathExpr) bool {
	val := x.Value.(xml.StartElement)

	if p.NodeType == "" {
		if p.Name.Space != "" {
			if p.Name.Space != "*" {
				if p.Name.Space != val.Name.Space {
					return false
				}
			}
		}

		if p.Name.Local == "*" && p.Axis != xconst.AxisAttribute && p.Axis != xconst.AxisNamespace {
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
