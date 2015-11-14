package procinst

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/parser/result/pathexpr"
	"github.com/ChrisTrenkamp/goxpath/parser/result/pathres"
)

//PathResProcInst is an implementation of PathRes for XML processing-instructions
type PathResProcInst struct {
	Value  interface{}
	Parent pathres.PathRes
}

//Interface returns the data representing the processing-instructions
func (pi *PathResProcInst) Interface() interface{} {
	return pi.Value
}

//GetParent returns the parent node
func (pi *PathResProcInst) GetParent() pathres.PathRes {
	return pi.Parent
}

//GetChildren returns nothing
func (pi *PathResProcInst) GetChildren() []pathres.PathRes {
	return []pathres.PathRes{}
}

//GetValue returns the value of the processing-instruction
func (pi *PathResProcInst) GetValue() string {
	//TODO: Make this return the value
	return ""
}

//Print prints the XML processing-instruction in string form
func (pi *PathResProcInst) Print(e *xml.Encoder) error {
	var err error
	if _, ok := pi.Value.(xml.ProcInst); ok {
		val := pi.Value.(xml.ProcInst)
		err = e.EncodeToken(val)
	}
	return err
}

//EvalPath evaluates the XPath path instruction on the processing-instruction
func (pi *PathResProcInst) EvalPath(p *pathexpr.PathExpr) bool {
	//TODO: Implement
	return false
}
