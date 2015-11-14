package chardata

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/parser/result/pathexpr"
	"github.com/ChrisTrenkamp/goxpath/parser/result/pathres"
)

//PathResCharData is an implementation of PathRes for XML attributes
type PathResCharData struct {
	Value  interface{}
	Parent pathres.PathRes
}

//Interface returns the data representing the character data
func (cd *PathResCharData) Interface() interface{} {
	return cd.Value
}

//GetParent returns the parent node, or itself if it's the root
func (cd *PathResCharData) GetParent() pathres.PathRes {
	return cd.Parent
}

//GetChildren returns nothing
func (cd *PathResCharData) GetChildren() []pathres.PathRes {
	return []pathres.PathRes{}
}

//GetValue returns the value of the element
func (cd *PathResCharData) GetValue() string {
	//TODO: Make this return the value
	return ""
}

//Print prints the XML character data in string form
func (cd *PathResCharData) Print(e *xml.Encoder) error {
	var err error
	if _, ok := cd.Value.(xml.CharData); ok {
		val := cd.Value.(xml.CharData)
		err = e.EncodeToken(val)
	}
	return err
}

//EvalPath evaluates the XPath path instruction on the character data
func (cd *PathResCharData) EvalPath(p *pathexpr.PathExpr) bool {
	//TODO: Implement
	return false
}
