package intfns

import (
	"encoding/xml"
	"fmt"

	"github.com/ChrisTrenkamp/goxpath/goxpath/xpfn"
	"github.com/ChrisTrenkamp/goxpath/goxpath/xpfn/noarg"
	"github.com/ChrisTrenkamp/goxpath/goxpath/xpfn/optarg"
	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/tree/literals/numlit"
	"github.com/ChrisTrenkamp/goxpath/tree/literals/strlit"
)

//BuiltIn contains the list of built-in XPath functions
var BuiltIn = map[string]xpfn.XPFn{
	"last":       noarg.Wrap{Fn: last},
	"count":      noarg.Wrap{Fn: count},
	"local-name": optarg.Wrap{Fn: name},
}

func last(arg []tree.Res) ([]tree.Res, error) {
	if len(arg) == 0 {
		return nil, nil
	}

	return []tree.Res{arg[len(arg)-1]}, nil
}

func count(arg []tree.Res) ([]tree.Res, error) {
	return []tree.Res{numlit.NumLit(len(arg))}, nil
}

func name(arg []tree.Res) ([]tree.Res, error) {
	if len(arg) == 0 {
		return nil, fmt.Errorf("No node in argument.")
	}

	node, ok := arg[0].(tree.Node)

	if !ok {
		return nil, fmt.Errorf("Argument is not a node")
	}

	tok := node.GetToken()

	switch t := tok.(type) {
	case xml.StartElement:
		if t.Name.Local == "" {
			return nil, fmt.Errorf("Cannot run local-name on root node.")
		}
		ret := []tree.Res{strlit.StrLit(t.Name.Local)}
		return ret, nil
	case xml.Attr:
		ret := []tree.Res{strlit.StrLit(t.Name.Local)}
		return ret, nil
	case xml.ProcInst:
		ret := []tree.Res{strlit.StrLit(t.Target)}
		return ret, nil
	}

	return []tree.Res{strlit.StrLit("")}, nil
}
