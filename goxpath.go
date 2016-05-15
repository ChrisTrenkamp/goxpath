package goxpath

import (
	"github.com/ChrisTrenkamp/goxpath/internal/execxp"
	"github.com/ChrisTrenkamp/goxpath/internal/parser"
	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/xtypes"
)

//XPathExec is the XPath executor, compiled from an XPath string
type XPathExec struct {
	n *parser.Node
}

//Exec executes the XPath expression, xp, against the tree, t, with the
//namespace mappings, ns.
func Exec(xp XPathExec, t tree.Node, ns map[string]string) (xtypes.Result, error) {
	return execxp.Exec(xp.n, t, ns)
}

//ExecStr combines Parse and Exec into one method.
func ExecStr(xpstr string, t tree.Node, ns map[string]string) (xtypes.Result, error) {
	xp, err := Parse(xpstr)
	if err != nil {
		return nil, err
	}
	return execxp.Exec(xp.n, t, ns)
}

//MustExec is like Exec, but panics instead of returning an error.
func MustExec(xp XPathExec, t tree.Node, ns map[string]string) xtypes.Result {
	res, err := Exec(xp, t, ns)
	if err != nil {
		panic(err)
	}
	return res
}

//Parse parses the XPath expression, xp, returning an XPath executor.
func Parse(xp string) (XPathExec, error) {
	n, err := parser.Parse(xp)
	return XPathExec{n: n}, err
}

//MustParse is like Parse, but panics instead of returning an error.
func MustParse(xp string) XPathExec {
	ret, err := Parse(xp)
	if err != nil {
		panic(err)
	}
	return ret
}
