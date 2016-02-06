package ctxpos

import "github.com/ChrisTrenkamp/goxpath/tree"

//CtxPos contains the result of an XPath expression and contains its context position
type CtxPos struct {
	Pos int
	tree.Res
}

//CreateCtxPos wraps tree in a CtxPos slice with their new positions set
func CreateCtxPos(tree []tree.Res) []CtxPos {
	ret := make([]CtxPos, len(tree))
	for i := range tree {
		ret[i] = CtxPos{Pos: i + 1, Res: tree[i]}
	}
	return ret
}
