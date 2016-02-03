package tree

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/goxpath/pathexpr"
)

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
	//EvalPath evaluates this node against the XPath step, p, and returns true
	//if it passes, and false otherwise.
	EvalPath(p *pathexpr.PathExpr) bool
}

//Elem is a XPath result that is an element node
type Elem interface {
	Node
	//GetChildren returns the elements children.
	GetChildren() []Node
	//GetNSAttrs returns the namespaces and attributes of the element
	GetNSAttrs() []Node
}
