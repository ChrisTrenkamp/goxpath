package xpfn

import (
	"fmt"

	"github.com/ChrisTrenkamp/goxpath/tree"
)

//Ctx represents the current context position, size, node, and the current filtered result
type Ctx struct {
	tree.Node
	Filter []tree.Res
	Pos    int
	Size   int
}

//Fn is a XPath function, written in Go
type Fn func(c Ctx, args ...[]tree.Res) ([]tree.Res, error)

//Wrap interfaces XPath function calls with Go
type Wrap struct {
	Fn Fn
	//NArgs represents the number of arguments to the XPath function.  -1 represents a single optional argument
	NArgs int
	//Variadic makes the last argument variadic
	Variadic bool
}

//Call checks the arguments and calls Fn if they are valid
func (w Wrap) Call(c Ctx, args ...[]tree.Res) ([]tree.Res, error) {
	if w.NArgs == -1 {
		if len(args) != 0 && len(args) != 1 {
			return nil, fmt.Errorf("Too many arguments.")
		}

		return w.Fn(c, args...)
	}

	if len(args) < w.NArgs {
		return nil, fmt.Errorf("Not enough arguments")
	}

	if len(args) > w.NArgs && !w.Variadic {
		return nil, fmt.Errorf("Too many arguments")
	}

	return w.Fn(c, args...)
}
