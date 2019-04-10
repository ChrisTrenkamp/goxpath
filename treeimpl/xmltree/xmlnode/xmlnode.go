package xmlnode

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/tree"
)

//Node is a XPath result that is a node except elements
type Node interface {
	//ResValue prints the node's string value
	ResValue() string
	//Pos returns the node's position in the document order
	Pos() int
	//GetToken returns the xml.Token representation of the node
	GetToken() xml.Token
	//GetParent returns the parent node, which will always be an XML element
	GetParent() Elem
	//GetNodeType returns the node's type
	GetNodeType() tree.NodeType
}

//Elem is a XPath result that is an element node
type Elem interface {
	Node
	//GetChildren returns the elements children.
	GetChildren() []Node
	//GetAttrs returns the attributes of the element
	GetAttrs() []Node
}

//NSElem is a node that keeps track of namespaces.
type NSElem interface {
	Elem
	GetNS() map[xml.Name]string
}

//XMLNode will represent an attribute, character data, comment, or processing instruction node
type XMLNode struct {
	xml.Token
	tree.NodeType
	Parent   Elem
	Position int
}

func (a XMLNode) Pos() int { return a.Position }

func (a XMLNode) GetNodeType() tree.NodeType { return a.NodeType }

//GetToken returns the xml.Token representation of the node
func (a XMLNode) GetToken() xml.Token {
	if a.NodeType == tree.NtAttr {
		ret := a.Token.(*xml.Attr)
		return *ret
	}
	return a.Token
}

//GetParent returns the parent node
func (a XMLNode) GetParent() Elem {
	return a.Parent
}

//ResValue returns the string value of the attribute
func (a XMLNode) ResValue() string {
	switch a.NodeType {
	case tree.NtAttr:
		return a.Token.(*xml.Attr).Value
	case tree.NtChd:
		return string(a.Token.(xml.CharData))
	case tree.NtComm:
		return string(a.Token.(xml.Comment))
	}
	//case tree.NtPi:
	return string(a.Token.(xml.ProcInst).Inst)
}
