package findutil

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/parser/pathexpr"
	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/result/xmlele"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/result/xmlns"
	"github.com/ChrisTrenkamp/goxpath/xconst"
)

type findFunc func(tree.XPRes, *pathexpr.PathExpr, *[]tree.XPRes)

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
func Find(x tree.XPRes, p pathexpr.PathExpr) []tree.XPRes {
	ret := []tree.XPRes{}

	if p.Axis == "" {
		findChild(x, &p, &ret)
		return ret
	}

	f := findMap[p.Axis]
	f(x, &p, &ret)

	return ret
}

func findAncestor(x tree.XPRes, p *pathexpr.PathExpr, ret *[]tree.XPRes) {
	if x.GetParent().EvalPath(p) {
		*ret = append(*ret, x.GetParent())
	}
	if x.GetParent() != x {
		findAncestor(x.GetParent(), p, ret)
	}
}

func findAncestorOrSelf(x tree.XPRes, p *pathexpr.PathExpr, ret *[]tree.XPRes) {
	findSelf(x, p, ret)
	findAncestor(x, p, ret)
}

func findAttribute(x tree.XPRes, p *pathexpr.PathExpr, ret *[]tree.XPRes) {
	if ele, ok := x.(*xmlele.XMLEle); ok {
		for i := range ele.Attrs {
			if ele.Attrs[i].EvalPath(p) {
				*ret = append(*ret, ele.Attrs[i])
			}
		}
	}
}

func findChild(x tree.XPRes, p *pathexpr.PathExpr, ret *[]tree.XPRes) {
	if ele, ok := x.(tree.XPResEle); ok {
		ch := ele.GetChildren()
		for i := range ch {
			if ch[i].EvalPath(p) {
				*ret = append(*ret, ch[i])
			}
		}
	}
}

func findDescendent(x tree.XPRes, p *pathexpr.PathExpr, ret *[]tree.XPRes) {
	if ele, ok := x.(tree.XPResEle); ok {
		ch := ele.GetChildren()
		for i := range ch {
			if ch[i].EvalPath(p) {
				*ret = append(*ret, ch[i])
			}
			findDescendent(ch[i], p, ret)
		}
	}
}

func findDescendentOrSelf(x tree.XPRes, p *pathexpr.PathExpr, ret *[]tree.XPRes) {
	findSelf(x, p, ret)
	findDescendent(x, p, ret)
}

func findFollowing(x tree.XPRes, p *pathexpr.PathExpr, ret *[]tree.XPRes) {
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

func findFollowingSibling(x tree.XPRes, p *pathexpr.PathExpr, ret *[]tree.XPRes) {
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

func findNamespace(x tree.XPRes, p *pathexpr.PathExpr, ret *[]tree.XPRes) {
	if ele, ok := x.(*xmlele.XMLEle); ok {
		for k, v := range ele.NS {
			ns := &xmlns.XMLNS{
				Attr:   xml.Attr{Name: k, Value: v},
				Parent: ele,
			}
			if ns.EvalPath(p) {
				*ret = append(*ret, ns)
			}
		}
	}
}

func findParent(x tree.XPRes, p *pathexpr.PathExpr, ret *[]tree.XPRes) {
	if x.GetParent() != x && x.GetParent().EvalPath(p) {
		*ret = append(*ret, x.GetParent())
	}
}

func findPreceding(x tree.XPRes, p *pathexpr.PathExpr, ret *[]tree.XPRes) {
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

func findPrecedingSibling(x tree.XPRes, p *pathexpr.PathExpr, ret *[]tree.XPRes) {
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

func findSelf(x tree.XPRes, p *pathexpr.PathExpr, ret *[]tree.XPRes) {
	if x.EvalPath(p) {
		*ret = append(*ret, x)
	}
}
