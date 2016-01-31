package tree

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/goxpath/pathexpr"
)

//Res is the result of an XPath expression
type Res interface {
	//String prints the node's string value
	String() string
}

//Node is a XPath result that is a node except elements
type Node interface {
	Res
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
}
