package tree

import "encoding/xml"

//NodePos is a helper for representing the node's document order
type NodePos int

//Pos returns the node's document order position
func (n NodePos) Pos() int {
	return int(n)
}

//Res is the result of an XPath expression
type Res interface {
	//String prints the node's string value
	String() string
}

//Node is a XPath result that is a node except elements
type Node interface {
	Res
	//Pos returns the node's position in the document order
	Pos() int
	//GetToken returns the xml.Token representation of the node
	GetToken() xml.Token
	//GetParent returns the parent node, which will always be an XML element
	GetParent() Elem
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
	GetNS() []Node
}
