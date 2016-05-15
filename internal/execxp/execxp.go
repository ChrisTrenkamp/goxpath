package execxp

import (
	"github.com/ChrisTrenkamp/goxpath/internal/parser"
	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/xtypes"
)

//Exec executes the XPath expression, xp, against the tree, t, with the
//namespace mappings, ns.
func Exec(n *parser.Node, t tree.Node, ns map[string]string) (xtypes.Result, error) {
	f := xpFilt{
		t:   t,
		ns:  ns,
		ctx: t,
	}

	return exec(&f, n)
}

func exec(f *xpFilt, n *parser.Node) (xtypes.Result, error) {
	err := xfExec(f, n)
	return f.res, err
}
