package intfns

import (
	"encoding/xml"
	"fmt"

	"github.com/ChrisTrenkamp/goxpath/goxpath/ctxpos"
	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/tree/literals/numlit"
	"github.com/ChrisTrenkamp/goxpath/tree/literals/strlit"
)

func last(arg []ctxpos.CtxPos) ([]tree.Res, error) {
	return []tree.Res{numlit.NumLit(len(arg))}, nil
}

func count(args ...[]ctxpos.CtxPos) ([]tree.Res, error) {
	arg := args[0]

	if len(arg) == 0 {
		return nil, nil
	}

	if _, ok := arg[0].Res.(tree.Node); !ok {
		return nil, fmt.Errorf("Argument is not a node-set")
	}

	ret := 0

	for i := range arg {
		countArg(arg[i].Res, &ret)
	}
	return []tree.Res{numlit.NumLit(ret)}, nil
}

func countArg(r tree.Res, c *int) {
	switch t := r.(type) {
	case tree.Elem:
		for _, i := range t.GetChildren() {
			countArg(i, c)
		}
		(*c)++
	case tree.Node:
		(*c)++
	}
}

func localName(arg []ctxpos.CtxPos) ([]tree.Res, error) {
	if len(arg) == 0 {
		return []tree.Res{strlit.StrLit("")}, nil
	}

	node, ok := arg[0].Res.(tree.Node)

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

func namespaceURI(arg []ctxpos.CtxPos) ([]tree.Res, error) {
	if len(arg) == 0 {
		return []tree.Res{strlit.StrLit("")}, nil
	}

	node, ok := arg[0].Res.(tree.Node)

	if !ok {
		return nil, fmt.Errorf("Argument is not a node")
	}

	tok := node.GetToken()

	switch t := tok.(type) {
	case xml.StartElement:
		if t.Name.Local == "" {
			return nil, fmt.Errorf("Cannot run namespace-uri on root node.")
		}
		ret := []tree.Res{strlit.StrLit(t.Name.Space)}
		return ret, nil
	case xml.Attr:
		ret := []tree.Res{strlit.StrLit(t.Name.Space)}
		return ret, nil
	}

	return []tree.Res{strlit.StrLit("")}, nil
}

func name(arg []ctxpos.CtxPos) ([]tree.Res, error) {
	if len(arg) == 0 {
		return []tree.Res{strlit.StrLit("")}, nil
	}

	node, ok := arg[0].Res.(tree.Node)

	if !ok {
		return nil, fmt.Errorf("Argument is not a node")
	}

	tok := node.GetToken()

	switch t := tok.(type) {
	case xml.StartElement:
		if t.Name.Local == "" {
			return nil, fmt.Errorf("Cannot run name on root node.")
		}
		space := ""
		if t.Name.Space != "" {
			space = fmt.Sprintf("{%s}", t.Name.Space)
		}
		res := fmt.Sprintf("%s%s", space, t.Name.Local)
		ret := []tree.Res{strlit.StrLit(res)}
		return ret, nil
	case xml.Attr:
		space := ""
		if t.Name.Space != "" {
			space = fmt.Sprintf("{%s}", t.Name.Space)
		}
		res := fmt.Sprintf("%s%s", space, t.Name.Local)
		ret := []tree.Res{strlit.StrLit(res)}
		return ret, nil
	case xml.ProcInst:
		res := fmt.Sprintf("%s", t.Target)
		ret := []tree.Res{strlit.StrLit(res)}
		return ret, nil
	}

	return []tree.Res{strlit.StrLit("")}, nil
}
