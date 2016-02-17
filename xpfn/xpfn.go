package xpfn

import (
	"fmt"

	"github.com/ChrisTrenkamp/goxpath/literals/boollit"
	"github.com/ChrisTrenkamp/goxpath/literals/strlit"
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

//GetOptArg returns the context filter if there is no argument in args.  Otherwise,
//it returns the first element in args.
func GetOptArg(c Ctx, args ...[]tree.Res) []tree.Res {
	if len(args) == 0 {
		return c.Filter
	}

	return args[0]
}

//GetNode casts res into an array of nodes.  If one of the elements is not a
//node, it will return an error.
func GetNode(res []tree.Res, e error) ([]tree.Node, error) {
	if e != nil {
		return nil, e
	}

	ret := make([]tree.Node, len(res))
	for i := range res {
		if n, ok := res[i].(tree.Node); ok {
			ret[i] = n
		} else {
			return nil, fmt.Errorf("One or more arguments are not nodes")
		}
	}

	return ret, nil
}

//GetFirstNode returns the first element in nodes, which will also be the
//first node in document order.
func GetFirstNode(nodes []tree.Node, e error) (tree.Node, error) {
	if e != nil {
		return nil, e
	}

	if len(nodes) > 0 {
		return nodes[0], nil
	}

	return nil, fmt.Errorf("No nodes in set")
}

//GetBool casts the first argument in res to a bool.  If the length of res
//is not 1, then it returns an error.
func GetBool(res []tree.Res, e error) (bool, error) {
	if e != nil {
		return false, e
	}

	if len(res) != 1 {
		return false, fmt.Errorf("Result set is not a single boolean")
	}

	if b, ok := res[0].(boollit.BoolLit); ok {
		return b.String() == "true", nil
	}

	return false, fmt.Errorf("Result set is not a single boolean")
}

//GetString casts the first argument in res to a string.  If the length of res
//is not 1, then it returns an error.
func GetString(res []tree.Res, e error) (string, error) {
	if e != nil {
		return "", e
	}

	if len(res) != 1 {
		return "", fmt.Errorf("Result set is not a single string")
	}

	if s, ok := res[0].(strlit.StrLit); ok {
		return s.String(), nil
	}

	return "", fmt.Errorf("Result set is not a single string")
}
