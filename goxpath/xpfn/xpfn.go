package xpfn

import (
	"github.com/ChrisTrenkamp/goxpath/goxpath/ctxpos"
	"github.com/ChrisTrenkamp/goxpath/tree"
)

//XPFn interfaces XPath function calls with Go
type XPFn interface {
	Call(ctx []ctxpos.CtxPos, args ...[]ctxpos.CtxPos) ([]tree.Res, error)
}
