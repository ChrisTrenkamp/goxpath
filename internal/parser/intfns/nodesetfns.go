package intfns

import (
	"encoding/xml"
	"fmt"

	"github.com/ChrisTrenkamp/goxpath/literals/numlit"
	"github.com/ChrisTrenkamp/goxpath/literals/strlit"
	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/xfn"
)

func last(c xfn.Ctx, args ...[]tree.Res) ([]tree.Res, error) {
	return []tree.Res{numlit.NumLit(c.Size)}, nil
}

func position(c xfn.Ctx, args ...[]tree.Res) ([]tree.Res, error) {
	return []tree.Res{numlit.NumLit(c.Pos)}, nil
}

func count(c xfn.Ctx, args ...[]tree.Res) ([]tree.Res, error) {
	arg, err := xfn.GetNode(xfn.GetOptArg(c, args...), nil)

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
	default:
		(*c)++
	}
}

func localName(c xfn.Ctx, args ...[]tree.Res) ([]tree.Res, error) {
	node, err := xfn.GetFirstNode(xfn.GetNode(xfn.GetOptArg(c, args...), nil))

	if err != nil {
		return nil, err
	}

	tok := node.GetToken()

	switch node.GetNodeType() {
	case tree.NtRoot:
		return nil, fmt.Errorf("Cannot run local-name on root node.")
	case tree.NtEle:
		ret := []tree.Res{strlit.StrLit(tok.(xml.StartElement).Name.Local)}
		return ret, nil
	case tree.NtAttr:
		ret := []tree.Res{strlit.StrLit(tok.(xml.Attr).Name.Local)}
		return ret, nil
	case tree.NtPi:
		ret := []tree.Res{strlit.StrLit(tok.(xml.ProcInst).Target)}
		return ret, nil
	}

	return []tree.Res{strlit.StrLit("")}, nil
}

func namespaceURI(c xfn.Ctx, args ...[]tree.Res) ([]tree.Res, error) {
	node, err := xfn.GetFirstNode(xfn.GetNode(xfn.GetOptArg(c, args...), nil))

	if err != nil {
		return nil, err
	}

	tok := node.GetToken()

	switch node.GetNodeType() {
	case tree.NtRoot:
		return nil, fmt.Errorf("Cannot run namespace-uri on root node.")
	case tree.NtEle:
		ret := []tree.Res{strlit.StrLit(tok.(xml.StartElement).Name.Space)}
		return ret, nil
	case tree.NtAttr:
		ret := []tree.Res{strlit.StrLit(tok.(xml.Attr).Name.Space)}
		return ret, nil
	}

	return []tree.Res{strlit.StrLit("")}, nil
}

func name(c xfn.Ctx, args ...[]tree.Res) ([]tree.Res, error) {
	node, err := xfn.GetFirstNode(xfn.GetNode(xfn.GetOptArg(c, args...), nil))

	if err != nil {
		return nil, err
	}

	switch node.GetNodeType() {
	case tree.NtRoot:
		return nil, fmt.Errorf("Cannot run name on root node.")
	case tree.NtEle:
		t := node.GetToken().(xml.StartElement)
		space := ""

		if t.Name.Space != "" {
			space = fmt.Sprintf("{%s}", t.Name.Space)
		}

		res := fmt.Sprintf("%s%s", space, t.Name.Local)
		ret := []tree.Res{strlit.StrLit(res)}

		return ret, nil
	case tree.NtAttr:
		t := node.GetToken().(xml.Attr)
		space := ""

		if t.Name.Space != "" {
			space = fmt.Sprintf("{%s}", t.Name.Space)
		}

		res := fmt.Sprintf("%s%s", space, t.Name.Local)
		ret := []tree.Res{strlit.StrLit(res)}
		return ret, nil
	case tree.NtPi:
		res := fmt.Sprintf("%s", node.GetToken().(xml.ProcInst).Target)
		ret := []tree.Res{strlit.StrLit(res)}
		return ret, nil
	}

	return []tree.Res{strlit.StrLit("")}, nil
}
