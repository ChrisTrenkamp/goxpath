package findutil

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/goxpath/pathexpr"
	"github.com/ChrisTrenkamp/goxpath/goxpath/xconst"
	"github.com/ChrisTrenkamp/goxpath/tree"
)

type findFunc func(tree.Node, *pathexpr.PathExpr, *[]tree.Node)

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
func Find(x tree.Node, p pathexpr.PathExpr) []tree.Node {
	ret := []tree.Node{}

	if p.Axis == "" {
		findChild(x, &p, &ret)
		return ret
	}

	f := findMap[p.Axis]
	f(x, &p, &ret)

	return ret
}

func findAncestor(x tree.Node, p *pathexpr.PathExpr, ret *[]tree.Node) {
	if checkNode(x.GetParent(), p) {
		*ret = append(*ret, x.GetParent())
	}
	if x.GetParent() != x {
		findAncestor(x.GetParent(), p, ret)
	}
}

func findAncestorOrSelf(x tree.Node, p *pathexpr.PathExpr, ret *[]tree.Node) {
	findSelf(x, p, ret)
	findAncestor(x, p, ret)
}

func findAttribute(x tree.Node, p *pathexpr.PathExpr, ret *[]tree.Node) {
	if ele, ok := x.(tree.Elem); ok {
		for _, i := range ele.GetAttrs() {
			if checkNode(i, p) {
				*ret = append(*ret, i)
			}
		}
	}
}

func findChild(x tree.Node, p *pathexpr.PathExpr, ret *[]tree.Node) {
	if ele, ok := x.(tree.Elem); ok {
		ch := ele.GetChildren()
		for i := range ch {
			if checkNode(ch[i], p) {
				*ret = append(*ret, ch[i])
			}
		}
	}
}

func findDescendent(x tree.Node, p *pathexpr.PathExpr, ret *[]tree.Node) {
	if ele, ok := x.(tree.Elem); ok {
		ch := ele.GetChildren()
		for i := range ch {
			if checkNode(ch[i], p) {
				*ret = append(*ret, ch[i])
			}
			findDescendent(ch[i], p, ret)
		}
	}
}

func findDescendentOrSelf(x tree.Node, p *pathexpr.PathExpr, ret *[]tree.Node) {
	findSelf(x, p, ret)
	findDescendent(x, p, ret)
}

func findFollowing(x tree.Node, p *pathexpr.PathExpr, ret *[]tree.Node) {
	if x == x.GetParent() {
		return
	}
	par := x.GetParent()
	ch := par.GetChildren()
	i := 0
	for x != ch[i] {
		i++
	}
	i++
	for i < len(ch) {
		findDescendentOrSelf(ch[i], p, ret)
		i++
	}
	findFollowing(par, p, ret)
}

func findFollowingSibling(x tree.Node, p *pathexpr.PathExpr, ret *[]tree.Node) {
	if x == x.GetParent() {
		return
	}
	par := x.GetParent()
	ch := par.GetChildren()
	i := 0
	for x != ch[i] {
		i++
	}
	i++
	for i < len(ch) {
		findSelf(ch[i], p, ret)
		i++
	}
}

func findNamespace(x tree.Node, p *pathexpr.PathExpr, ret *[]tree.Node) {
	if ele, ok := x.(tree.NSElem); ok {
		for _, i := range ele.GetNS() {
			attr := i.GetToken().(xml.Attr)
			if evalNS(p, attr) {
				*ret = append(*ret, i)
			}
		}
	}
}

func findParent(x tree.Node, p *pathexpr.PathExpr, ret *[]tree.Node) {
	if x.GetParent() != x && checkNode(x.GetParent(), p) {
		*ret = append(*ret, x.GetParent())
	}
}

func findPreceding(x tree.Node, p *pathexpr.PathExpr, ret *[]tree.Node) {
	if x == x.GetParent() {
		return
	}
	par := x.GetParent()
	ch := par.GetChildren()
	i := len(ch) - 1
	for x != ch[i] {
		i--
	}
	i--
	for i >= 0 {
		findDescendentOrSelf(ch[i], p, ret)
		i--
	}
	findPreceding(par, p, ret)
}

func findPrecedingSibling(x tree.Node, p *pathexpr.PathExpr, ret *[]tree.Node) {
	if x == x.GetParent() {
		return
	}
	par := x.GetParent()
	ch := par.GetChildren()
	i := len(ch) - 1
	for x != ch[i] {
		i--
	}
	i--
	for i >= 0 {
		findSelf(ch[i], p, ret)
		i--
	}
}

func findSelf(x tree.Node, p *pathexpr.PathExpr, ret *[]tree.Node) {
	if checkNode(x, p) {
		*ret = append(*ret, x)
	}
}

func checkNode(x tree.Node, p *pathexpr.PathExpr) bool {
	tok := x.GetToken()
	switch t := tok.(type) {
	case xml.Attr:
		return evalAttr(p, t)
	case xml.CharData:
		return evalChd(p)
	case xml.Comment:
		return evalComm(p)
	case xml.StartElement:
		return evalEle(p, t)
	case xml.ProcInst:
		return evalPI(p)
	}
	return false
}

func evalAttr(p *pathexpr.PathExpr, a xml.Attr) bool {
	if p.NodeType == "" {
		if p.Name.Space != "*" {
			if a.Name.Space != p.NS[p.Name.Space] {
				return false
			}
		}

		if p.Name.Local == "*" && p.Axis == xconst.AxisAttribute {
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

func evalEle(p *pathexpr.PathExpr, ele xml.StartElement) bool {
	if p.NodeType == "" {
		return checkNameAndSpace(p, ele)
	}

	if p.NodeType == xconst.NodeTypeNode {
		return true
	}

	return false
}

func checkNameAndSpace(p *pathexpr.PathExpr, ele xml.StartElement) bool {
	if p.Name.Local == "*" && p.Name.Space == "" {
		return true
	}

	if p.Name.Space != "*" && ele.Name.Space != p.NS[p.Name.Space] {
		return false
	}

	if p.Name.Local == "*" && p.Axis != xconst.AxisAttribute && p.Axis != xconst.AxisNamespace {
		return true
	}

	if p.Name.Local == ele.Name.Local {
		return true
	}

	return false
}

func evalNS(p *pathexpr.PathExpr, ns xml.Attr) bool {
	if p.NodeType == "" {
		if p.Name.Space != "" && p.Name.Space != "*" {
			return false
		}

		if p.Name.Local == "*" && p.Axis == xconst.AxisNamespace {
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
	if p.NodeType == xconst.NodeTypeProcInst {
		return true
	}

	if p.NodeType == xconst.NodeTypeProcInst || p.NodeType == xconst.NodeTypeNode {
		return true
	}

	return false
}
