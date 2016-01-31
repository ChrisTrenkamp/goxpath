package fns

import (
	"encoding/xml"
	"fmt"

	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/tree/literals/strlit"
)

type xpFn func(ctx []tree.Res, args ...[]tree.Res) ([]tree.Res, error)

//BuiltIn contains the list of built-in XPath functions
var BuiltIn = map[string]xpFn{
	"local-name": name,
}

func name(ctx []tree.Res, args ...[]tree.Res) ([]tree.Res, error) {
	var arg tree.Res
	var argSet []tree.Res
	if len(args) == 0 {
		if len(ctx) > 1 {
			return nil, fmt.Errorf("Cannot run local-name() on multiple contexts.")
		}

		arg = ctx[0]
	} else if len(args) > 1 {
		return nil, fmt.Errorf("Cannot run local-name() on multiple contexts.")
	} else {
		argSet = args[0]

		if len(argSet) > 1 {
			return nil, fmt.Errorf("Cannot run local-name() on multiple contexts.")
		}

		arg = argSet[0]
	}

	node, ok := arg.(tree.Node)

	if !ok {
		return nil, fmt.Errorf("Cannot run local-name on non-nodes.")
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

	return nil, fmt.Errorf("Cannot run local-name on result.")
}
