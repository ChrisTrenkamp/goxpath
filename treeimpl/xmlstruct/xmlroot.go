package xmlstruct

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/treeimpl/xmltree/xmlnode"
)

type XMLRoot struct {
	Ele *XMLEle
}

func (x *XMLRoot) ResValue() string {
	return x.Ele.ResValue()
}

func (x *XMLRoot) Pos() int {
	return 0
}

func (x *XMLRoot) GetToken() xml.Token {
	return xml.StartElement{}
}

func (x *XMLRoot) GetParent() xmlnode.Elem {
	return x
}

func (x *XMLRoot) GetNodeType() tree.NodeType {
	return tree.NtRoot
}

func (x *XMLRoot) GetChildren() []xmlnode.Node {
	return []xmlnode.Node{x.Ele}
}

func (x *XMLRoot) GetAttrs() []xmlnode.Node {
	return nil
}
