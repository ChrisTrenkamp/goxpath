package findutil

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/parser/pathexpr"
	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/xconst"
)

const (
	wildcard = "*"
)

type findFunc func(tree.Adapter, interface{}, *pathexpr.PathExpr, *[]interface{})

var findMap = map[string]findFunc{
	xconst.AxisAncestor:         findAncestor,
	xconst.AxisAncestorOrSelf:   findAncestorOrSelf,
	xconst.AxisAttribute:        findAttribute,
	xconst.AxisChild:            findChild,
	xconst.AxisDescendent:       findDescendent,
	xconst.AxisDescendentOrSelf: findDescendentOrSelf,
	xconst.AxisFollowing:        findFollowing,
	xconst.AxisFollowingSibling: findFollowingSibling,
	xconst.AxisNamespace:        findNamespace,
	xconst.AxisParent:           findParent,
	xconst.AxisPreceding:        findPreceding,
	xconst.AxisPrecedingSibling: findPrecedingSibling,
	xconst.AxisSelf:             findSelf,
}

//Find finds nodes based on the pathexpr.PathExpr
func Find(a tree.Adapter, x interface{}, p pathexpr.PathExpr) []interface{} {
	ret := []interface{}{}

	if p.Axis == "" {
		findChild(a, x, &p, &ret)
		return ret
	}

	f := findMap[p.Axis]
	f(a, x, &p, &ret)

	return ret
}

func findAncestor(a tree.Adapter, x interface{}, p *pathexpr.PathExpr, ret *[]interface{}) {
	if a.GetNodeType(x) == tree.NtRoot {
		return
	}

	addNode(a, a.GetParent(x), p, ret)
	findAncestor(a, a.GetParent(x), p, ret)
}

func findAncestorOrSelf(a tree.Adapter, x interface{}, p *pathexpr.PathExpr, ret *[]interface{}) {
	findSelf(a, x, p, ret)
	findAncestor(a, x, p, ret)
}

func findAttribute(a tree.Adapter, x interface{}, p *pathexpr.PathExpr, ret *[]interface{}) {
	a.ForEachAttr(x, func(attr xml.Attr, ptr interface{}) {
		if evalAttr(p, attr) {
			*ret = append(*ret, ptr)
		}
	})
}

func findChild(a tree.Adapter, x interface{}, p *pathexpr.PathExpr, ret *[]interface{}) {
	a.ForEachChild(x, func(ptr interface{}) {
		addNode(a, ptr, p, ret)
	})
}

func findDescendent(a tree.Adapter, x interface{}, p *pathexpr.PathExpr, ret *[]interface{}) {
	a.ForEachChild(x, func(ptr interface{}) {
		addNode(a, ptr, p, ret)
		findDescendent(a, ptr, p, ret)
	})
}

func findDescendentOrSelf(a tree.Adapter, x interface{}, p *pathexpr.PathExpr, ret *[]interface{}) {
	findSelf(a, x, p, ret)
	findDescendent(a, x, p, ret)
}

func findFollowing(a tree.Adapter, x interface{}, p *pathexpr.PathExpr, ret *[]interface{}) {
	if a.GetNodeType(x) == tree.NtRoot {
		return
	}
	par := a.GetParent(x)

	seen := false
	a.ForEachChild(par, func(ptr interface{}) {
		if !seen {
			seen = ptr == x
		} else {
			findDescendentOrSelf(a, ptr, p, ret)
		}
	})
	findFollowing(a, par, p, ret)
}

func findFollowingSibling(a tree.Adapter, x interface{}, p *pathexpr.PathExpr, ret *[]interface{}) {
	if a.GetNodeType(x) == tree.NtRoot {
		return
	}
	par := a.GetParent(x)
	seen := false
	a.ForEachChild(par, func(ptr interface{}) {
		if !seen {
			seen = ptr == x
		} else {
			findSelf(a, ptr, p, ret)
		}
	})
}

func findNamespace(a tree.Adapter, x interface{}, p *pathexpr.PathExpr, ret *[]interface{}) {
	ns := a.GetNamespaces(x)
	for _, i := range ns {
		addNode(a, i, p, ret)
	}
}

func findParent(a tree.Adapter, x interface{}, p *pathexpr.PathExpr, ret *[]interface{}) {
	if a.GetNodeType(x) != tree.NtRoot {
		addNode(a, a.GetParent(x), p, ret)
	}
}

func findPreceding(a tree.Adapter, x interface{}, p *pathexpr.PathExpr, ret *[]interface{}) {
	if a.GetNodeType(x) == tree.NtRoot {
		return
	}
	par := a.GetParent(x)
	seen := false
	a.ForEachChild(par, func(ptr interface{}) {
		if !seen {
			seen = ptr == x
			if !seen {
				findDescendentOrSelf(a, ptr, p, ret)
			}
		}
	})
	findPreceding(a, par, p, ret)
}

func findPrecedingSibling(a tree.Adapter, x interface{}, p *pathexpr.PathExpr, ret *[]interface{}) {
	if a.GetNodeType(x) == tree.NtRoot {
		return
	}
	par := a.GetParent(x)

	seen := false
	a.ForEachChild(par, func(ptr interface{}) {
		if !seen {
			seen = ptr == x
			if !seen {
				findSelf(a, ptr, p, ret)
			}
		}
	})
}

func findSelf(a tree.Adapter, x interface{}, p *pathexpr.PathExpr, ret *[]interface{}) {
	addNode(a, x, p, ret)
}

func addNode(a tree.Adapter, x interface{}, p *pathexpr.PathExpr, ret *[]interface{}) {
	add := false

	switch a.GetNodeType(x) {
	case tree.NtAttr:
		add = evalAttr(p, a.GetAttrTok(x))
	case tree.NtChd:
		add = evalChd(p)
	case tree.NtComm:
		add = evalComm(p)
	case tree.NtElem, tree.NtRoot:
		add = evalEle(p, a.GetElementName(x))
	case tree.NtNs:
		add = evalNS(p, a.GetNamespaceTok(x))
	case tree.NtPi:
		add = evalPI(p)
	}

	if add {
		*ret = append(*ret, x)
	}
}

func evalAttr(p *pathexpr.PathExpr, a xml.Attr) bool {
	if p.NodeType == "" {
		if p.Name.Space != wildcard {
			if a.Name.Space != p.NS[p.Name.Space] {
				return false
			}
		}

		if p.Name.Local == wildcard && p.Axis == xconst.AxisAttribute {
			return true
		}

		if p.Name.Local == a.Name.Local {
			return true
		}
	} else {
		if p.NodeType == xconst.NodeTypeNode {
			return true
		}
	}

	return false
}

func evalChd(p *pathexpr.PathExpr) bool {
	if p.NodeType == xconst.NodeTypeText || p.NodeType == xconst.NodeTypeNode {
		return true
	}

	return false
}

func evalComm(p *pathexpr.PathExpr) bool {
	if p.NodeType == xconst.NodeTypeComment || p.NodeType == xconst.NodeTypeNode {
		return true
	}

	return false
}

func evalEle(p *pathexpr.PathExpr, ele xml.Name) bool {
	if p.NodeType == "" {
		return checkNameAndSpace(p, ele)
	}

	if p.NodeType == xconst.NodeTypeNode {
		return true
	}

	return false
}

func checkNameAndSpace(p *pathexpr.PathExpr, ele xml.Name) bool {
	if p.Name.Local == wildcard && p.Name.Space == "" {
		return true
	}

	if p.Name.Space != wildcard && ele.Space != p.NS[p.Name.Space] {
		return false
	}

	if p.Name.Local == wildcard && p.Axis != xconst.AxisAttribute && p.Axis != xconst.AxisNamespace {
		return true
	}

	if p.Name.Local == ele.Local {
		return true
	}

	return false
}

func evalNS(p *pathexpr.PathExpr, ns xml.Attr) bool {
	if p.NodeType == "" {
		if p.Name.Space != "" && p.Name.Space != wildcard {
			return false
		}

		if p.Name.Local == wildcard && p.Axis == xconst.AxisNamespace {
			return true
		}

		if p.Name.Local == ns.Name.Local {
			return true
		}
	} else {
		if p.NodeType == xconst.NodeTypeNode {
			return true
		}
	}

	return false
}

func evalPI(p *pathexpr.PathExpr) bool {
	if p.NodeType == xconst.NodeTypeProcInst || p.NodeType == xconst.NodeTypeNode {
		return true
	}

	return false
}
