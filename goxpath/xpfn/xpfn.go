package xpfn

import "github.com/ChrisTrenkamp/goxpath/tree"

//XPFn interfaces XPath function calls with Go
type XPFn interface {
	Call(ctx []tree.Res, args ...[]tree.Res) ([]tree.Res, error)
}
