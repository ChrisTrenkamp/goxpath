package noarg

import (
	"github.com/ChrisTrenkamp/goxpath/goxpath/ctxpos"
	"github.com/ChrisTrenkamp/goxpath/tree"
)

//Fn defines an XPath function that has one optional argument
type Fn func(arg []ctxpos.CtxPos) ([]tree.Res, error)

//Wrap is wraps the OptArgFn method with XPFn
type Wrap struct {
	Fn Fn
}

//Call satisfies the XPFn interface for optarg.Wrap
func (fn Wrap) Call(ctx []ctxpos.CtxPos, args ...[]ctxpos.CtxPos) ([]tree.Res, error) {
	return fn.Fn(ctx)
}
