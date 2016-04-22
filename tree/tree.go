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
	Elem
	GetNS() map[xml.Name]string
}

//NSBuilder is a helper-struct for satisfying the NSElem interface
type NSBuilder struct {
	NS map[xml.Name]string
}

//GetNS returns the namespaces found on the current element.  It should not be
//confused with BuildNS, which actually resolves the namespace nodes.
func (ns NSBuilder) GetNS() map[xml.Name]string {
	return ns.NS
}

//BuildNS resolves all the namespace nodes of the element and returns them
func BuildNS(t Elem) (ret []NS) {
	vals := make(map[xml.Name]string)

	if nselem, ok := t.(NSElem); ok {
		buildNS(nselem, vals)

		ret = make([]NS, 0, len(vals))
		i := 1

		for k, v := range vals {
			if !(k.Local == "xmlns" && k.Space == "" && v == "") {
				ret = append(ret, NS{
					Attr:     xml.Attr{Name: k, Value: v},
					Parent:   t,
					NodePos:  NodePos(t.Pos() + i),
					NodeType: NtNs,
				})
				i++
			}
		}
	}

	return ret
}

func buildNS(x NSElem, ret map[xml.Name]string) {
	if x.GetNodeType() == NtRoot {
		return
	}

	if nselem, ok := x.GetParent().(NSElem); ok {
		buildNS(nselem, ret)
	}

	for k, v := range x.GetNS() {
		ret[k] = v
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
