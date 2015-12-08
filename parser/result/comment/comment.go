package comment

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/parser/result/pathexpr"
	"github.com/ChrisTrenkamp/goxpath/parser/result/pathres"
	"github.com/ChrisTrenkamp/goxpath/xconst"
)

//PathResComment is an implementation of PathRes for XML attributes
type PathResComment struct {
	Value  interface{}
	Parent pathres.PathRes
}

//Interface returns the data representing the comment
func (c *PathResComment) Interface() interface{} {
	return c.Value
}

//GetParent returns the parent node
func (c *PathResComment) GetParent() pathres.PathRes {
	return c.Parent
}

//GetChildren returns nothing
func (c *PathResComment) GetChildren() []pathres.PathRes {
	return []pathres.PathRes{}
}

//GetValue returns the value of the comment
func (c *PathResComment) GetValue() string {
	return string(c.Value.(xml.Comment))
}

//Print prints the XML comment in string form
func (c *PathResComment) Print(e *xml.Encoder) error {
	var err error
	if _, ok := c.Value.(xml.Comment); ok {
		val := c.Value.(xml.Comment)
		err = e.EncodeToken(val)
	}
	return err
}

//EvalPath evaluates the XPath path instruction on the element
func (c *PathResComment) EvalPath(p *pathexpr.PathExpr) bool {
	if p.NodeType == xconst.NodeTypeComment || p.NodeType == xconst.NodeTypeNode {
		return true
	}
	return false
}
