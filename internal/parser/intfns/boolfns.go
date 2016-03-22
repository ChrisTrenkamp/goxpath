package intfns

import (
	"github.com/ChrisTrenkamp/goxpath/literals/boollit"
	"github.com/ChrisTrenkamp/goxpath/literals/numlit"
	"github.com/ChrisTrenkamp/goxpath/literals/strlit"
	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/xfn"
)

func boolean(c xfn.Ctx, args ...[]tree.Res) ([]tree.Res, error) {
	return []tree.Res{boollit.BoolLit(BooleanFunc(args[0]))}, nil
}

//BooleanFunc returns the XPath boolean value of the argument.  This is used
//for the function itself as well as predicates (as defined by the spec).
func BooleanFunc(arg []tree.Res) bool {
	if len(arg) == 0 {
		return false
	}

	switch t := arg[0].(type) {
	case tree.Node:
		return true
	case boollit.BoolLit:
		return bool(t)
	case numlit.NumLit:
		return float64(t) > 0
	case strlit.StrLit:
		return len(string(t)) > 0
	}

	return false
}

func not(c xfn.Ctx, arg ...[]tree.Res) ([]tree.Res, error) {
	ret, err := boolean(c, arg...)
	return []tree.Res{boollit.BoolLit(!(ret[0].(boollit.BoolLit)))}, err
}

func _true(c xfn.Ctx, arg ...[]tree.Res) ([]tree.Res, error) {
	return []tree.Res{boollit.BoolLit(true)}, nil
}

func _false(c xfn.Ctx, arg ...[]tree.Res) ([]tree.Res, error) {
	return []tree.Res{boollit.BoolLit(false)}, nil
}
