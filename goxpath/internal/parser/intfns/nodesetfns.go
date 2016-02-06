package intfns

import (
	"encoding/xml"
	"fmt"

	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/tree/literals/numlit"
	"github.com/ChrisTrenkamp/goxpath/tree/literals/strlit"
)

func last(arg []tree.Res) ([]tree.Res, error) {
	if len(arg) == 0 {
		return nil, nil
	}

	return []tree.Res{arg[len(arg)-1]}, nil
}

func count(arg []tree.Res) ([]tree.Res, error) {
	return []tree.Res{numlit.NumLit(len(arg))}, nil
}

func localName(arg []tree.Res) ([]tree.Res, error) {
	if len(arg) == 0 {
		return []tree.Res{strlit.StrLit("")}, nil
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

func namespaceURI(arg []tree.Res) ([]tree.Res, error) {
	if len(arg) == 0 {
		return []tree.Res{strlit.StrLit("")}, nil
	}

	node, ok := arg[0].(tree.Node)

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

func name(arg []tree.Res) ([]tree.Res, error) {
	if len(arg) == 0 {
		return []tree.Res{strlit.StrLit("")}, nil
	}

	node, ok := arg[0].(tree.Node)

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
