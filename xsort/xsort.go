package xsort

import (
	"sort"

	"github.com/ChrisTrenkamp/goxpath/tree"
)

type resSort []tree.Res

func (ns resSort) Len() int      { return len(ns) }
func (ns resSort) Swap(i, j int) { ns[i], ns[j] = ns[j], ns[i] }
func (ns resSort) Less(i, j int) bool {
	ic, iok := ns[i].(tree.Node)
	jc, jok := ns[j].(tree.Node)

	if iok && jok {
		return ic.Pos() < jc.Pos()
	}

	return true
}

type nodeSort []tree.Node

func (ns nodeSort) Len() int      { return len(ns) }
func (ns nodeSort) Swap(i, j int) { ns[i], ns[j] = ns[j], ns[i] }
func (ns nodeSort) Less(i, j int) bool {
	return ns[i].Pos() < ns[j].Pos()
}

//SortRes sorts the array res by the node document order.  If an element in
//the array is not a node, the result will be pushed to the beginning of the slice.
func SortRes(res []tree.Res) {
	sort.Sort(resSort(res))
}

//SortNodes sorts the array by the node document order
func SortNodes(res []tree.Node) {
	sort.Sort(nodeSort(res))
}
