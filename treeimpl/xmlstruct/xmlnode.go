package xmlstruct

import (
	"encoding/xml"
	"fmt"
	"reflect"

	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/treeimpl/xmltree/xmlnode"
)

type XMLNode struct {
	Val      reflect.Value
	pos      int
	prnt     xmlnode.Elem
	nodeType tree.NodeType
	prntTag  string
}

func (x *XMLNode) ResValue() string {
	if x.Val.Kind() == reflect.Ptr {
		return fmt.Sprintf("%v", x.Val.Elem().Interface())
	}
	return fmt.Sprintf("%v", x.Val.Interface())
}

func (x *XMLNode) Pos() int {
	return x.pos
}

func (x *XMLNode) GetToken() xml.Token {
	switch x.nodeType {
	case tree.NtAttr:
		return xml.Attr{Name: getTagInfo(x.prntTag).name, Value: x.ResValue()}
	case tree.NtChd:
		return xml.CharData(x.ResValue())
	}
	//case tree.NtComm:
	return xml.Comment(x.ResValue())
}

func (x *XMLNode) GetParent() xmlnode.Elem {
	return x.prnt
}

func (x *XMLNode) GetNodeType() tree.NodeType {
	return x.nodeType
}
