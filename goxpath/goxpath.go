package goxpath

import (
	"github.com/ChrisTrenkamp/goxpath/goxpath/internal/parser"
	"github.com/ChrisTrenkamp/goxpath/tree"
)

//XPathExec is the XPath executor, compiled from an XPath string
type XPathExec []parser.XPExec

//Exec executes the XPath expression, xp, against the tree, t, with the
//namespace mappings, ns.
func Exec(xp XPathExec, t tree.Node, ns map[string]string) ([]tree.Res, error) {
	return parser.Exec(xp, t, ns)
}

//MustExec is like Exec, but panics instead of returning an error.
func MustExec(xp XPathExec, t tree.Node, ns map[string]string) []tree.Res {
	res, err := parser.Exec(xp, t, ns)
	if err != nil {
		panic(err)
	}
	return res
}

//MustParse is like Parse, but panics instead of returning an error.
func MustParse(xp string) XPathExec {
	return parser.MustParse(xp)
}

//Parse parses the XPath expression, xp, returning an XPath executor.
func Parse(xp string) (XPathExec, error) {
	return parser.Parse(xp)
}
