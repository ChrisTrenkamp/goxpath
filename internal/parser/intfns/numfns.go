package intfns

import (
	"fmt"
	"math"
	"strconv"

	"github.com/ChrisTrenkamp/goxpath/literals/boollit"
	"github.com/ChrisTrenkamp/goxpath/literals/numlit"
	"github.com/ChrisTrenkamp/goxpath/literals/strlit"
	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/xfn"
)

func number(c xfn.Ctx, args ...[]tree.Res) ([]tree.Res, error) {
	num, err := NumberFunc(args[0])
	if err != nil {
		return nil, err
	}

	return []tree.Res{numlit.NumLit(num)}, nil
}

func sum(c xfn.Ctx, args ...[]tree.Res) ([]tree.Res, error) {
	n, err := xfn.GetNode(args[0], nil)
	if err != nil {
		return nil, err
	}

	ret := 0.0
	var num float64
	for _, i := range n {
		num, err = strconv.ParseFloat(i.ResValue(), 64)
		if err != nil {
			num = math.NaN()
		}
		ret += num
	}

	return []tree.Res{numlit.NumLit(ret)}, err
}

func floor(c xfn.Ctx, args ...[]tree.Res) ([]tree.Res, error) {
	n, err := xfn.GetNumber(args[0], nil)
	if err != nil {
		return nil, err
	}

	return []tree.Res{numlit.NumLit(math.Floor(n))}, nil
}

func ceiling(c xfn.Ctx, args ...[]tree.Res) ([]tree.Res, error) {
	n, err := xfn.GetNumber(args[0], nil)
	if err != nil {
		return nil, err
	}

	return []tree.Res{numlit.NumLit(math.Ceil(n))}, nil
}

func round(c xfn.Ctx, args ...[]tree.Res) ([]tree.Res, error) {
	f, err := xfn.GetNumber(args[0], nil)
	if err != nil {
		return nil, err
	}

	if f < -0.5 {
		f = float64(int(f - 0.5))
	} else if f > 0.5 {
		f = float64(int(f + 0.5))
	} else {
		f = 0
	}

	return []tree.Res{numlit.NumLit(f)}, nil
}

//NumberFunc returns the XPath number value of the argument.
func NumberFunc(arg []tree.Res) (float64, error) {
	if len(arg) == 0 {
		return 0, fmt.Errorf("No objects to convert to a number.")
	}

	switch t := arg[0].(type) {
	case tree.Node:
		nodes, err := xfn.GetNode(arg, nil)
		if err != nil {
			return 0, err
		}
		str := ""
		for _, i := range nodes {
			str += i.ResValue()
		}
		return strconv.ParseFloat(str, 64)
	case boollit.BoolLit:
		if t {
			return 1, nil
		}

		return 0, nil
	case numlit.NumLit:
		return float64(t), nil
	case strlit.StrLit:
		return strconv.ParseFloat(string(t), 64)
	}

	return 0, fmt.Errorf("Unknown object to convert to a number")
}
