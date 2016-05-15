package intfns

import (
	"fmt"
	"math"

	"github.com/ChrisTrenkamp/goxpath/xfn"
	"github.com/ChrisTrenkamp/goxpath/xtypes"
)

func number(c xfn.Ctx, args ...xtypes.Result) (xtypes.Result, error) {
	if b, ok := args[0].(xtypes.IsNum); ok {
		return b.Num(), nil
	}

	return nil, fmt.Errorf("Cannot convert object to a number")
}

func sum(c xfn.Ctx, args ...xtypes.Result) (xtypes.Result, error) {
	n, ok := args[0].(xtypes.NodeSet)
	if !ok {
		return nil, fmt.Errorf("Argument is not a node-set")
	}

	ret := 0.0
	for _, i := range n {
		ret += float64(xtypes.GetNodeNum(i))
	}

	return xtypes.Num(ret), nil
}

func floor(c xfn.Ctx, args ...xtypes.Result) (xtypes.Result, error) {
	n, ok := args[0].(xtypes.Num)
	if !ok {
		return nil, fmt.Errorf("Cannot convert object to a number")
	}

	return xtypes.Num(math.Floor(float64(n))), nil
}

func ceiling(c xfn.Ctx, args ...xtypes.Result) (xtypes.Result, error) {
	n, ok := args[0].(xtypes.Num)
	if !ok {
		return nil, fmt.Errorf("Cannot convert object to a number")
	}

	return xtypes.Num(math.Ceil(float64(n))), nil
}

func round(c xfn.Ctx, args ...xtypes.Result) (xtypes.Result, error) {
	n, ok := args[0].(xtypes.Num)
	if !ok {
		return nil, fmt.Errorf("Cannot convert object to a number")
	}

	if n < -0.5 {
		n = xtypes.Num(int(n - 0.5))
	} else if n > 0.5 {
		n = xtypes.Num(int(n + 0.5))
	} else {
		n = 0
	}

	return n, nil
}
