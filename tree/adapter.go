package tree

import (
	"encoding/xml"
)

//XMLSpace is the W3C XML namespace
const XMLSpace = "http://www.w3.org/XML/1998/namespace"

//NodeType is a safer way to determine a node's type than type assertions.
type NodeType int

//These are all the possible node types
const (
	NtAttr NodeType = iota
	NtChd
	NtComm
	NtElem
	NtNs
	NtRoot
	NtPi
)

// Adapter adapts an XML tree representation to xpath. The
// requirements for the underlying XML implementation is minimal:
//
//   - The XML implementation must have consistent pointers for its
//   nodes, that is, if two interface{} values are passed to the
//   functions of the adapter interface, those two must point to the
//   same backend node if and only if the interface{} values are
//   equal.
//   - The underlying XML implementation must have a way of
//     identifying positions of nodes in a unique way. There must be
//     only one node for any node position. Node positions are not required
//     to be consecutive.
//
type Adapter interface {
	// Returns the type of the node
	GetNodeType(interface{}) NodeType
	// Returns the parent node of the given node
	GetParent(interface{}) interface{}
	// Returns the position of the node within the document. This is
	// used to order nodes within nodesets
	NodePos(interface{}) int
	// The input is a pointer to an element node. The call should build
	// the namespace of the nodes, and return an array of namespace
	// objects (of type NtNs)
	GetNamespaces(interface{}) []interface{}

	// GetElementName should return the name of the element
	GetElementName(interface{}) xml.Name

	// The input is an attribute. This must return the attribute token
	GetAttrTok(interface{}) xml.Attr
	// The input in a namespace object. This must return a namespace attribute
	GetNamespaceTok(interface{}) xml.Attr
	// The input is an element. It should return an element token
	GetElemTok(interface{}) xml.StartElement
	// The input is a processing instruction. It should return a ProcInst token
	GetProcInstTok(interface{}) xml.ProcInst
	// The input is chardata. It should return a chardata token
	GetCharDataTok(interface{}) xml.CharData
	// The input is a comment. It should return a comment token
	GetCommentTok(interface{}) xml.Comment

	// If the node is an element, iterate attributes of element and call
	// the func with the attribute, and a pointer to the attribute node
	ForEachAttr(interface{}, func(xml.Attr, interface{}))
	// If the node is an element, iterate direct children of the node,
	// and call func with a pointer to each child
	ForEachChild(interface{}, func(interface{}))

	// NewNodeSet must return a new nodeset implementation containing the given nodes
	NewNodeSet([]interface{}) NodeSet
	// StringValue must return the string value of the nodeq
	StringValue(interface{}) string
}

// FindAttributeForElement iuses ForEachAttr to find the requested attribute
func FindAttributeForElement(a Adapter, node interface{}, lname, space string) (xml.Attr, bool) {
	var found xml.Attr
	ok := false
	a.ForEachAttr(node, func(attr xml.Attr, a interface{}) {
		if attr.Name.Space == space && attr.Name.Local == lname {
			found = attr
			ok = true
		}
	})
	return found, ok
}
