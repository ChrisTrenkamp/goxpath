package directive

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/parser/result/pathexpr"
	"github.com/ChrisTrenkamp/goxpath/parser/result/pathres"
)

//PathResDirective is an implementation of PathRes for XML directives
type PathResDirective struct {
	Value  interface{}
	Parent pathres.PathRes
}

//Interface returns the data representing the directive
func (d *PathResDirective) Interface() interface{} {
	return d.Value
}

//GetParent returns the parent node
func (d *PathResDirective) GetParent() pathres.PathRes {
	return d.Parent
}

//GetChildren returns nothing
func (d *PathResDirective) GetChildren() []pathres.PathRes {
	return []pathres.PathRes{}
}

//GetValue returns the value of the directive
func (d *PathResDirective) GetValue() string {
	//TODO: Make this return the value
	return ""
}

//Print prints the XML directive in string form
func (d *PathResDirective) Print(e *xml.Encoder) error {
	var err error
	if _, ok := d.Value.(xml.Directive); ok {
		val := d.Value.(xml.Directive)
		err = e.EncodeToken(val)
	}
	return err
}

//EvalPath evaluates the XPath path instruction on the element
func (d *PathResDirective) EvalPath(p *pathexpr.PathExpr) bool {
	//TODO: Implement
	return false
}
