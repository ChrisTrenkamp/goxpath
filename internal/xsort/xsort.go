package xsort

import (
	"sort"

	"github.com/ChrisTrenkamp/goxpath/tree"
)

type nodeSort []tree.Node

func (ns nodeSort) Len() int      { return len(ns) }
func (ns nodeSort) Swap(i, j int) { ns[i], ns[j] = ns[j], ns[i] }
func (ns nodeSort) Less(i, j int) bool {
	return ns[i].Pos() < ns[j].Pos()
}

type nsSort []tree.NS

func (ns nsSort) Len() int { return len(ns) }
func (ns nsSort) Swap(i, j int) {
	ns[i], ns[j] = ns[j], ns[i]
	ns[i].NodePos, ns[j].NodePos = ns[j].NodePos, ns[i].NodePos
}
func (ns nsSort) Less(i, j int) bool {
	return ns[i].Value < ns[j].Value
}

//SortNodes sorts the array by the node document order
func SortNodes(res []tree.Node) {
	sort.Sort(nodeSort(res))
}

//SortNS sorts the NS's returned from tree.BuildNS.  It sorts based on the namespace
//URL and assigns the document position accordingly.
func SortNS(ns []tree.NS) {
	sort.Sort(nsSort(ns))
}
