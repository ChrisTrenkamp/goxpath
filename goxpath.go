package goxpath

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/internal/execxp"
	"github.com/ChrisTrenkamp/goxpath/internal/parser"
	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/xfn"
	"github.com/ChrisTrenkamp/goxpath/xtypes"
)

//Opts defines namespace mappings and custom functions for XPath expressions.
type Opts struct {
	NS    map[string]string
	Funcs map[xml.Name]xfn.Wrap
}

//FuncOpts is a function wrapper for Opts.
type FuncOpts func(*Opts)

//XPathExec is the XPath executor, compiled from an XPath string
type XPathExec struct {
	n *parser.Node
}

//Exec executes the XPath expression, xp, against the tree, t, with the
//namespace mappings, ns.
func Exec(xp XPathExec, t tree.Node, opts ...FuncOpts) (xtypes.Result, error) {
	o := &Opts{
		NS:    make(map[string]string),
		Funcs: make(map[xml.Name]xfn.Wrap),
	}
	for _, i := range opts {
		i(o)
	}
	return execxp.Exec(xp.n, t, o.NS, o.Funcs)
}

//MustExec is like Exec, but panics instead of returning an error.
func MustExec(xp XPathExec, t tree.Node, opts ...FuncOpts) xtypes.Result {
	res, err := Exec(xp, t, opts...)
	if err != nil {
		panic(err)
	}
	return res
}

//ExecStr combines Parse and Exec into one method.
func ExecStr(xpstr string, t tree.Node, opts ...FuncOpts) (xtypes.Result, error) {
	xp, err := Parse(xpstr)
	if err != nil {
		return nil, err
	}
	return Exec(xp, t, opts...)
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
