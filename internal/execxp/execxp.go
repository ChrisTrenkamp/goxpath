package execxp

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/parser"
	"github.com/ChrisTrenkamp/goxpath/tree"
)

//Exec executes the XPath expression, xp, against the tree, t, with the
//namespace mappings, ns.
func Exec(n *parser.Node, a tree.Adapter, t interface{}, ns map[string]string, fns map[xml.Name]tree.Wrap, v map[string]tree.Result) (tree.Result, error) {
	f := xpFilt{
		t:         t,
		ns:        ns,
		ctx:       a.NewNodeSet([]interface{}{t}),
		fns:       fns,
		variables: v,
	}

	return exec(a, &f, n)
}

func exec(a tree.Adapter, f *xpFilt, n *parser.Node) (tree.Result, error) {
	err := xfExec(a, f, n)
	return f.ctx, err
}
