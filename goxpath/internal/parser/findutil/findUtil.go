package findutil

import (
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
	if x.GetParent().EvalPath(p) {
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
		for _, i := range ele.GetNSAttrs() {
			if i.EvalPath(p) {
				*ret = append(*ret, i)
			}
		}
	}
}

func findChild(x tree.Node, p *pathexpr.PathExpr, ret *[]tree.Node) {
	if ele, ok := x.(tree.Elem); ok {
		ch := ele.GetChildren()
		for i := range ch {
			if ch[i].EvalPath(p) {
				*ret = append(*ret, ch[i])
			}
		}
	}
}

func findDescendent(x tree.Node, p *pathexpr.PathExpr, ret *[]tree.Node) {
	if ele, ok := x.(tree.Elem); ok {
		ch := ele.GetChildren()
		for i := range ch {
			if ch[i].EvalPath(p) {
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
	findAttribute(x, p, ret)
}

func findParent(x tree.Node, p *pathexpr.PathExpr, ret *[]tree.Node) {
	if x.GetParent() != x && x.GetParent().EvalPath(p) {
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
	if x.EvalPath(p) {
		*ret = append(*ret, x)
	}
}
