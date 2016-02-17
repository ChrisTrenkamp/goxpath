package intfns

import (
	"encoding/xml"
	"fmt"

	"github.com/ChrisTrenkamp/goxpath/literals/numlit"
	"github.com/ChrisTrenkamp/goxpath/literals/strlit"
	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/xpfn"
)

func last(c xpfn.Ctx, args ...[]tree.Res) ([]tree.Res, error) {
	return []tree.Res{numlit.NumLit(c.Size)}, nil
}

func position(c xpfn.Ctx, args ...[]tree.Res) ([]tree.Res, error) {
	return []tree.Res{numlit.NumLit(c.Pos)}, nil
}

func count(c xpfn.Ctx, args ...[]tree.Res) ([]tree.Res, error) {
	arg, err := xpfn.GetNode(xpfn.GetOptArg(c, args...), nil)

	if err != nil {
		return nil, err
	}

	ret := 0

	for i := range arg {
		countArg(arg[i], &ret)
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

func localName(c xpfn.Ctx, args ...[]tree.Res) ([]tree.Res, error) {
	node, err := xpfn.GetFirstNode(xpfn.GetNode(xpfn.GetOptArg(c, args...), nil))

	if err != nil {
		return nil, err
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

func namespaceURI(c xpfn.Ctx, args ...[]tree.Res) ([]tree.Res, error) {
	node, err := xpfn.GetFirstNode(xpfn.GetNode(xpfn.GetOptArg(c, args...), nil))

	if err != nil {
		return nil, err
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

func name(c xpfn.Ctx, args ...[]tree.Res) ([]tree.Res, error) {
	node, err := xpfn.GetFirstNode(xpfn.GetNode(xpfn.GetOptArg(c, args...), nil))

	if err != nil {
		return nil, err
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
