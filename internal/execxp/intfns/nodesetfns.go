package intfns

import (
	"fmt"

	"github.com/ChrisTrenkamp/goxpath/tree"
)

func last(a tree.Adapter, c tree.Ctx, args ...tree.Result) (tree.Result, error) {
	return tree.Num(c.Size), nil
}

func position(a tree.Adapter, c tree.Ctx, args ...tree.Result) (tree.Result, error) {
	return tree.Num(c.Pos), nil
}

func count(a tree.Adapter, c tree.Ctx, args ...tree.Result) (tree.Result, error) {
	n, ok := args[0].(tree.NodeSet)
	if !ok {
		return nil, fmt.Errorf("Cannot convert object to a node-set")
	}

	return tree.Num(len(n.GetNodes())), nil
}

func localName(a tree.Adapter, c tree.Ctx, args ...tree.Result) (tree.Result, error) {
	var n tree.NodeSet
	ok := true
	if len(args) == 1 {
		n, ok = args[0].(tree.NodeSet)
	} else {
		n = c.NodeSet
	}
	if !ok {
		return nil, fmt.Errorf("Cannot convert object to a node-set")
	}

	ret := ""
	nodes := n.GetNodes()
	if len(nodes) == 0 {
		return tree.String(ret), nil
	}
	node := nodes[0]
	switch a.GetNodeType(node) {
	case tree.NtElem:
		ret = a.GetElementName(node).Local
	case tree.NtAttr:
		ret = a.GetAttrTok(node).Name.Local
	case tree.NtPi:
		ret = a.GetProcInstTok(node).Target
	}

	return tree.String(ret), nil
}

func namespaceURI(a tree.Adapter, c tree.Ctx, args ...tree.Result) (tree.Result, error) {
	var n tree.NodeSet
	ok := true
	if len(args) == 1 {
		n, ok = args[0].(tree.NodeSet)
	} else {
		n = c.NodeSet
	}
	if !ok {
		return nil, fmt.Errorf("Cannot convert object to a node-set")
	}

	ret := ""
	nodes := n.GetNodes()
	if len(nodes) == 0 {
		return tree.String(ret), nil
	}
	node := nodes[0]

	switch a.GetNodeType(node) {
	case tree.NtElem:
		ret = a.GetElementName(node).Space
	case tree.NtAttr:
		ret = a.GetAttrTok(node).Name.Space
	}

	return tree.String(ret), nil
}

func name(a tree.Adapter, c tree.Ctx, args ...tree.Result) (tree.Result, error) {
	var n tree.NodeSet
	ok := true
	if len(args) == 1 {
		n, ok = args[0].(tree.NodeSet)
	} else {
		n = c.NodeSet
	}
	if !ok {
		return nil, fmt.Errorf("Cannot convert object to a node-set")
	}

	ret := ""
	nodes := n.GetNodes()
	if len(nodes) == 0 {
		return tree.String(ret), nil
	}
	node := nodes[0]

	switch a.GetNodeType(node) {
	case tree.NtElem:
		t := a.GetElementName(node)
		space := ""

		if t.Space != "" {
			space = fmt.Sprintf("{%s}", t.Space)
		}

		ret = fmt.Sprintf("%s%s", space, t.Local)
	case tree.NtAttr:
		t := a.GetAttrTok(node)
		space := ""

		if t.Name.Space != "" {
			space = fmt.Sprintf("{%s}", t.Name.Space)
		}

		ret = fmt.Sprintf("%s%s", space, t.Name.Local)
	case tree.NtPi:
		t := a.GetProcInstTok(node)
		ret = fmt.Sprintf("%s", t.Target)
	}

	return tree.String(ret), nil
}
