package execxp

import (
	"github.com/ChrisTrenkamp/goxpath/internal/parser"
	"github.com/ChrisTrenkamp/goxpath/tree"
)

type nodeExec func(n *parser.Node, t tree.Node, ns map[string]string) ([]tree.Res, error)

//Exec executes the XPath expression, xp, against the tree, t, with the
//namespace mappings, ns.
func Exec(n *parser.Node, t tree.Node, ns map[string]string) ([]tree.Res, error) {
	f := xpFilt{
		t:   t,
		ns:  ns,
		ctx: t,
	}

	return exec(&f, n)
}

func exec(f *xpFilt, n *parser.Node) ([]tree.Res, error) {
	err := xfExec(f, n)

	ret := make([]tree.Res, len(f.res))
	for i := range ret {
		ret[i] = f.res[i]
	}

	return ret, err
}
