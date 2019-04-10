package xmlbuilder

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/treeimpl/xmltree/xmlnode"
)

//BuilderOpts supplies all the information needed to create an XML node.
type BuilderOpts struct {
	Dec          *xml.Decoder
	Tok          xml.Token
	NodeType     tree.NodeType
	NS           map[xml.Name]string
	Attrs        []*xml.Attr
	NodePos      int
	AttrStartPos int
}

//XMLBuilder is an interface for creating XML structures.
type XMLBuilder interface {
	xmlnode.Node
	CreateNode(*BuilderOpts) XMLBuilder
	EndElem() XMLBuilder
}
