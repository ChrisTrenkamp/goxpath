package defoverride

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/result/xmlele"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/result/xmlnode"
)

//RootNode creates the root document node.  Return your own custom data type here.
func RootNode() tree.Elem {
	root := &xmlele.XMLEle{
		Children: []tree.Node{},
		NodeType: tree.NtRoot,
	}

	return root
}

//StartElem appends ele to pos's children.  The child element will be returned.
//ele will already have its fields filled out, they just need to be copied over
//to the new data type.
func StartElem(ele *xmlele.XMLEle, pos tree.Elem, dec *xml.Decoder) tree.Elem {
	//When implementing your own StartElem method, pos will be the type of the last call
	//to StartElem, or Root.
	curPos := pos.(*xmlele.XMLEle)

	//Here is where you create your own data type and append to pos's children.
	curPos.Children = append(curPos.Children, ele)
	return ele
}

//Node appends n to pos's children.  Customize the type as needed.
func Node(n xmlnode.XMLNode, pos tree.Elem, dec *xml.Decoder) {
	curPos := pos.(*xmlele.XMLEle)
	curPos.Children = append(curPos.Children, n)
}

//EndElem marks the end of an element.  The return value is pos's parent.
func EndElem(ele xml.EndElement, pos tree.Elem, dec *xml.Decoder) tree.Elem {
	curPos := pos.(*xmlele.XMLEle)
	return curPos.Parent
}

//Directive handles directives that may occur in an XML document.  The default functionality
//does nothing.
func Directive(dir xml.Directive, dec *xml.Decoder) {}
