package goxpath

import (
	"github.com/ChrisTrenkamp/goxpath/goxpath/internal/parser"
	"github.com/ChrisTrenkamp/goxpath/tree"
)

//XPathExec is the XPath executor, compiled from an XPath string
type XPathExec []parser.XPExec

//Exec executes the XPath expression, xp, against the tree, t, with the
//namespace mappings, ns.
func Exec(xp XPathExec, t tree.XPRes, ns map[string]string) []tree.XPRes {
	return parser.Exec(xp, t, ns)
}

//MustParse is like Parse, but panics instead of returning an error.
func MustParse(xp string) XPathExec {
	return parser.MustParse(xp)
}

//Parse parses the XPath expression, xp, returning an XPath executor.
func Parse(xp string) (XPathExec, error) {
	return parser.Parse(xp)
}
