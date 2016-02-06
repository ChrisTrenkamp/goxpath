package arg

import (
	"fmt"

	"github.com/ChrisTrenkamp/goxpath/tree"
)

//Fn defines an XPath function that has one or more arguments
type Fn func(arg ...[]tree.Res) ([]tree.Res, error)

//Wrap is wraps the OptArgFn method with XPFn
type Wrap struct {
	Fn       Fn
	NArgs    int
	Variadic bool
}

//Call satisfies the XPFn interface for optarg.Wrap
func (fn Wrap) Call(ctx []tree.Res, args ...[]tree.Res) ([]tree.Res, error) {
	if len(args) < fn.NArgs {
		return nil, fmt.Errorf("Not enough arguments")
	}

	if len(args) > fn.NArgs && !fn.Variadic {
		return nil, fmt.Errorf("Too many arguments")
	}

	return fn.Fn(args...)
}
