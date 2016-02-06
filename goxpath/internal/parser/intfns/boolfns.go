package intfns

import (
	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/tree/literals/boollit"
	"github.com/ChrisTrenkamp/goxpath/tree/literals/numlit"
	"github.com/ChrisTrenkamp/goxpath/tree/literals/strlit"
)

func boolean(args ...[]tree.Res) ([]tree.Res, error) {
	arg := args[0]

	if len(arg) == 0 {
		return []tree.Res{boollit.BoolLit(false)}, nil
	}

	switch t := arg[0].(type) {
	case tree.Node:
		return []tree.Res{boollit.BoolLit(true)}, nil
	case boollit.BoolLit:
		return []tree.Res{boollit.BoolLit(t)}, nil
	case numlit.NumLit:
		return []tree.Res{boollit.BoolLit(float64(t) > 0)}, nil
	case strlit.StrLit:
		return []tree.Res{boollit.BoolLit(len(string(t)) > 0)}, nil
	}

	return []tree.Res{boollit.BoolLit(false)}, nil
}

func not(arg ...[]tree.Res) ([]tree.Res, error) {
	ret, err := boolean(arg...)
	return []tree.Res{boollit.BoolLit(!(ret[0].(boollit.BoolLit)))}, err
}

func _true(arg []tree.Res) ([]tree.Res, error) {
	return []tree.Res{boollit.BoolLit(true)}, nil
}

func _false(arg []tree.Res) ([]tree.Res, error) {
	return []tree.Res{boollit.BoolLit(false)}, nil
}
