package pathres

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/parser/result/pathexpr"
)

const (
	//PathResElement is an element result from an XPath expression
	PathResElement int = iota
	//PathResAttribute is an attribute result from an XPath expression
	PathResAttribute
	//PathResRoot is a root document result from an XPath expression
	PathResRoot
	//PathResNamespace is a namespace node result from an XPath expression
	PathResNamespace
	//PathResProcInstr is a processing-instruction result from an XPath expression
	PathResProcInstr
	//PathResComment is a comment result from an XPath expression
	PathResComment
	//PathResText is a text result of an XPath expression
	PathResText
)

//PathRes is an interface for XPath results
type PathRes interface {
	Interface() interface{}
	GetParent() PathRes
	GetChildren() []PathRes
	GetValue() string
	Print(*xml.Encoder) error
	EvalPath(*pathexpr.PathExpr) bool
}
