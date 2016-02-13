package xmlpi

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/tree"
)

//XMLPI is an implementation of XPRes for XML processing-instructions
type XMLPI struct {
	xml.ProcInst
	Parent tree.Elem
	tree.NodePos
}

//GetToken returns the xml.Token representation of the node
func (pi *XMLPI) GetToken() xml.Token {
	return pi.ProcInst
}

//GetParent returns the parent node
func (pi *XMLPI) GetParent() tree.Elem {
	return pi.Parent
}

//String returns the value of the processing-instruction
func (pi *XMLPI) String() string {
	return string(pi.ProcInst.Inst)
}
