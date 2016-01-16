package tree

import "github.com/ChrisTrenkamp/goxpath/goxpath/pathexpr"

//XPRes is a XPath node within a tree
type XPRes interface {
	//String prints the node's string value
	String() string
	//GetParent returns the parent node, which will always be an XML element
	GetParent() XPResEle
	//EvalPath evaluates this node against the XPath step, p, and returns true
	//if it passes, and false otherwise.
	EvalPath(p *pathexpr.PathExpr) bool
}

//XPResEle represents an XML element node, which implements the XPRes interface,
//and adds another method for getting the node's children.
type XPResEle interface {
	XPRes
	//GetChildren returns the elements children.
	GetChildren() []XPRes
}
