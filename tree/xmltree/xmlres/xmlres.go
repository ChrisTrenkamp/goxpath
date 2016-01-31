package xmlres

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/tree"
)

//XMLPrinter implements the XMLPrint method for printing xml.Tokens,
//which is used by the tree/xmltree/result/... packages.
type XMLPrinter interface {
	XMLPrint(e *xml.Encoder) error
}

//XMLNode combines the XPRes interface and XMLPrint for xml-printing
//in the tree/xmltree/result/... packages.
type XMLNode interface {
	tree.Node
	XMLPrinter
}

//XMLElem combines the XPResEle interface and XMLPrint for xml-printing
//in the tree/xmltree/result/xmlele packages.
type XMLElem interface {
	tree.Elem
	XMLPrinter
}
