package tree

import "encoding/xml"

//NodePos is a helper for representing the node's document order
type NodePos int

//Pos returns the node's document order position
func (n NodePos) Pos() int {
	return int(n)
}

//NodeType is a safer way to determine a node's type than type assertions.
type NodeType int

//GetNodeType returns the node's type.
func (t NodeType) GetNodeType() NodeType {
	return t
}

//These are all the possible node types
const (
	NtAttr NodeType = iota
	NtChd
	NtComm
	NtEle
	NtNs
	NtRoot
	NtPi
)

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
	//GetNodeType returns the node's type
	GetNodeType() NodeType
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
	GetNS() *NSStruct
}

//NSStruct is a helper implementation of NSElem.
type NSStruct struct {
	NS map[xml.Name]NS
	Elem
	Parent *NSStruct
}

//GetNS returns an elements NSStruct if it exists
func (x *NSStruct) GetNS() *NSStruct {
	return x
}

//BuildNS resolves all the namespace nodes of the element and returns them
func (x *NSStruct) BuildNS() map[xml.Name]NS {
	ret := make(map[xml.Name]NS)

	x.buildNS(ret)

	i := 1

	for k, v := range ret {
		if v.Attr.Name.Local == "xmlns" && v.Attr.Name.Space == "" && v.Attr.Value == "" {
			delete(ret, k)
		} else {
			ret[k] = NS{
				Attr:     v.Attr,
				Parent:   x.Elem,
				NodePos:  NodePos(x.Elem.Pos() + i),
				NodeType: NtNs,
			}
			i++
		}
	}

	return ret
}

func (x *NSStruct) buildNS(ret map[xml.Name]NS) {
	if x == nil {
		return
	}

	x.Parent.buildNS(ret)

	if x.NS != nil {
		for k, v := range x.NS {
			ret[k] = v
		}
	}
}

//NS is a namespace node.
type NS struct {
	xml.Attr
	Parent Elem
	NodePos
	NodeType
}

//GetToken returns the xml.Token representation of the namespace.
func (ns NS) GetToken() xml.Token {
	return ns.Attr
}

//GetParent returns the parent node of the namespace.
func (ns NS) GetParent() Elem {
	return ns.Parent
}

//String returns the string value of the namespace
func (ns NS) String() string {
	return ns.Attr.Value
}
