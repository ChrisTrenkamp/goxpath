package execxp

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/internal/parser"
	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/xfn"
	"github.com/ChrisTrenkamp/goxpath/xtypes"
)

//Exec executes the XPath expression, xp, against the tree, t, with the
//namespace mappings, ns.
func Exec(n *parser.Node, t tree.Node, ns map[string]string, fns map[xml.Name]xfn.Wrap) (xtypes.Result, error) {
	f := xpFilt{
		t:   t,
		ns:  ns,
		ctx: xtypes.NodeSet{t},
		fns: fns,
	}

	return exec(&f, n)
}

func exec(f *xpFilt, n *parser.Node) (xtypes.Result, error) {
	err := xfExec(f, n)
	return f.ctx, err
}
