package xmlele

import (
	"encoding/xml"
	"sort"

	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/treeimpl/xmltree/xmlbuilder"
	"github.com/ChrisTrenkamp/goxpath/treeimpl/xmltree/xmlnode"
)

//NSBuilder is a helper-struct for satisfying the NSElem interface
type NSBuilder struct {
	NS map[xml.Name]string
}

//GetNS returns the namespaces found on the current element.  It should not be
//confused with BuildNS, which actually resolves the namespace nodes.
func (ns NSBuilder) GetNS() map[xml.Name]string {
	return ns.NS
}

//XMLEle is an implementation of XPRes for XML elements
type XMLEle struct {
	xml.StartElement
	NSBuilder
	Attrs    []xmlnode.Node
	Children []xmlnode.Node
	Parent   xmlnode.Elem
	tree.NodeType
	Position int
}

func (x XMLEle) GetNodeType() tree.NodeType { return x.NodeType }

//Root is the default root node builder for xmltree.ParseXML
func Root() xmlbuilder.XMLBuilder {
	return &XMLEle{NodeType: tree.NtRoot}
}

func (x *XMLEle) Pos() int { return x.Position }

//CreateNode is an implementation of xmlbuilder.XMLBuilder.  It appends the node
//specified in opts and returns the child if it is an element.  Otherwise, it returns x.
func (x *XMLEle) CreateNode(opts *xmlbuilder.BuilderOpts) xmlbuilder.XMLBuilder {
	if opts.NodeType == tree.NtElem {
		ele := &XMLEle{
			StartElement: opts.Tok.(xml.StartElement),
			NSBuilder:    NSBuilder{NS: opts.NS},
			Attrs:        make([]xmlnode.Node, len(opts.Attrs)),
			Parent:       x,
			Position:     opts.NodePos,
			NodeType:     opts.NodeType,
		}
		for i := range opts.Attrs {
			ele.Attrs[i] = xmlnode.XMLNode{
				Token:    opts.Attrs[i],
				Position: opts.AttrStartPos + i,
				NodeType: tree.NtAttr,
				Parent:   ele,
			}
		}
		x.Children = append(x.Children, ele)
		return ele
	}

	node := xmlnode.XMLNode{
		Token:    opts.Tok,
		Position: opts.NodePos,
		NodeType: opts.NodeType,
		Parent:   x,
	}
	x.Children = append(x.Children, node)
	return x
}

//EndElem is an implementation of xmlbuilder.XMLBuilder.  It returns x's parent.
func (x *XMLEle) EndElem() xmlbuilder.XMLBuilder {
	return x.Parent.(*XMLEle)
}

//GetToken returns the xml.Token representation of the node
func (x *XMLEle) GetToken() xml.Token {
	return x.StartElement
}

//GetParent returns the parent node, or itself if it's the root
func (x *XMLEle) GetParent() xmlnode.Elem {
	return x.Parent
}

//GetChildren returns all child nodes of the element
func (x *XMLEle) GetChildren() []xmlnode.Node {
	ret := make([]xmlnode.Node, len(x.Children))

	for i := range x.Children {
		ret[i] = x.Children[i]
	}

	return ret
}

//GetAttrs returns all attributes of the element
func (x *XMLEle) GetAttrs() []xmlnode.Node {
	ret := make([]xmlnode.Node, len(x.Attrs))
	for i := range x.Attrs {
		ret[i] = x.Attrs[i]
	}
	return ret
}

//ResValue returns the string value of the element and children
func (x *XMLEle) ResValue() string {
	ret := ""
	for i := range x.Children {
		switch x.Children[i].GetNodeType() {
		case tree.NtChd, tree.NtElem, tree.NtRoot:
			ret += x.Children[i].ResValue()
		}
	}
	return ret
}

type nsValueSort []NS

func (ns nsValueSort) Len() int { return len(ns) }
func (ns nsValueSort) Swap(i, j int) {
	ns[i], ns[j] = ns[j], ns[i]
}
func (ns nsValueSort) Less(i, j int) bool {
	return ns[i].Value < ns[j].Value
}

//BuildNS resolves all the namespace nodes of the element and returns them
func BuildNS(t xmlnode.Elem) (ret []NS) {
	vals := make(map[xml.Name]string)

	if nselem, ok := t.(xmlnode.NSElem); ok {
		buildNS(nselem, vals)

		ret = make([]NS, 0, len(vals))
		i := 1

		for k, v := range vals {
			if !(k.Local == "xmlns" && k.Space == "" && v == "") {
				ret = append(ret, NS{
					Attr:     xml.Attr{Name: k, Value: v},
					Parent:   t,
					NodeType: tree.NtNs,
				})
				i++
			}
		}

		sort.Sort(nsValueSort(ret))
		for i := range ret {
			ret[i].Position = t.Pos() + i + 1
		}
	}

	return ret
}

func buildNS(x xmlnode.NSElem, ret map[xml.Name]string) {
	if x.GetNodeType() == tree.NtRoot {
		return
	}

	if nselem, ok := x.GetParent().(xmlnode.NSElem); ok {
		buildNS(nselem, ret)
	}

	for k, v := range x.GetNS() {
		ret[k] = v
	}
}

//NS is a namespace node.
type NS struct {
	xml.Attr
	Parent   xmlnode.Elem
	Position int
	tree.NodeType
}

func (ns NS) GetNodeType() tree.NodeType { return ns.NodeType }

func (ns NS) Pos() int { return ns.Position }

//GetToken returns the xml.Token representation of the namespace.
func (ns NS) GetToken() xml.Token {
	return ns.Attr
}

//GetParent returns the parent node of the namespace.
func (ns NS) GetParent() xmlnode.Elem {
	return ns.Parent
}

//ResValue returns the string value of the namespace
func (ns NS) ResValue() string {
	return ns.Attr.Value
}
