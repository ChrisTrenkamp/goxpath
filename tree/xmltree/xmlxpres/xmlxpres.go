package xmlxpres

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/tree"
)

//XMLPrint implements the XMLPrint method for printing xml.Tokens,
//which is used by the tree/xmltree/result/... packages.
type XMLPrint interface {
	XMLPrint(e *xml.Encoder) error
}

//XMLXPRes combines the XPRes interface and XMLPrint for xml-printing
//in the tree/xmltree/result/... packages.
type XMLXPRes interface {
	tree.XPRes
	XMLPrint
}

//XMLXPResEle combines the XPResEle interface and XMLPrint for xml-printing
//in the tree/xmltree/result/xmlele packages.
type XMLXPResEle interface {
	tree.XPResEle
	XMLPrint
}
