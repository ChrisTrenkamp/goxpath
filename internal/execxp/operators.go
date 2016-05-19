package execxp

import (
	"fmt"
	"math"

	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/xtypes"
)

func bothNodeOperator(left xtypes.NodeSet, right xtypes.NodeSet, f *xpFilt, op string) error {
	var err error
	for _, l := range left {
		for _, r := range right {
			lStr := l.ResValue()
			rStr := r.ResValue()

			if eqOps[op] {
				err = equalsOperator(xtypes.String(lStr), xtypes.String(rStr), f, op)
				if err == nil && f.res.String() == xtypes.True {
					return nil
				}
			} else {
				err = numberOperator(xtypes.String(lStr), xtypes.String(rStr), f, op)
				if err == nil && f.res.String() == xtypes.True {
					return nil
				}
			}
		}
	}

	f.res = xtypes.Bool(false)

	return nil
}

func leftNodeOperator(left xtypes.NodeSet, right xtypes.Result, f *xpFilt, op string) error {
	var err error
	for _, l := range left {
		lStr := l.ResValue()

		if eqOps[op] {
			err = equalsOperator(xtypes.String(lStr), right, f, op)
			if err == nil && f.res.String() == xtypes.True {
				return nil
			}
		} else {
			err = numberOperator(xtypes.String(lStr), right, f, op)
			if err == nil && f.res.String() == xtypes.True {
				return nil
			}
		}
	}

	f.res = xtypes.Bool(false)

	return nil
}

func rightNodeOperator(left xtypes.Result, right xtypes.NodeSet, f *xpFilt, op string) error {
	var err error
	for _, r := range right {
		rStr := r.ResValue()

		if eqOps[op] {
			err = equalsOperator(left, xtypes.String(rStr), f, op)
			if err == nil && f.res.String() == "true" {
				return nil
			}
		} else {
			err = numberOperator(left, xtypes.String(rStr), f, op)
			if err == nil && f.res.String() == "true" {
				return nil
			}
		}
	}

	f.res = xtypes.Bool(false)

	return nil
}

func equalsOperator(left, right xtypes.Result, f *xpFilt, op string) error {
	_, lOK := left.(xtypes.Bool)
	_, rOK := right.(xtypes.Bool)
	if lOK || rOK {
		lTest, lt := left.(xtypes.IsBool)
		rTest, rt := right.(xtypes.IsBool)
		if !lt || !rt {
			return fmt.Errorf("Cannot convert argument to boolean")
		}
		if op == "=" {
			f.res = xtypes.Bool(lTest.Bool() == rTest.Bool())
		} else {
			f.res = xtypes.Bool(lTest.Bool() != rTest.Bool())
		}

		return nil
	}

	_, lOK = left.(xtypes.Num)
	_, rOK = right.(xtypes.Num)
	if lOK || rOK {
		return numberOperator(left, right, f, op)
	}

	lStr := left.String()
	rStr := right.String()

	if op == "=" {
		f.res = xtypes.Bool(lStr == rStr)
	} else {
		f.res = xtypes.Bool(lStr != rStr)
	}

	return nil
}

func numberOperator(left, right xtypes.Result, f *xpFilt, op string) error {
	lt, lOK := left.(xtypes.IsNum)
	rt, rOK := right.(xtypes.IsNum)
	if !lOK || !rOK {
		return fmt.Errorf("Cannot convert data type to number")
	}

	ln, rn := lt.Num(), rt.Num()

	switch op {
	case "*":
		f.res = xtypes.Num(ln * rn)
	case "div":
		if rn != 0 {
			f.res = xtypes.Num(ln / rn)
		} else {
			if ln == 0 {
				f.res = xtypes.Num(math.NaN())
			} else if ln >= 0 {
				f.res = xtypes.Num(math.Inf(1))
			} else {
				f.res = xtypes.Num(math.Inf(-1))
			}
		}
	case "mod":
		f.res = xtypes.Num(int(ln) % int(rn))
	case "+":
		f.res = xtypes.Num(ln + rn)
	case "-":
		f.res = xtypes.Num(ln - rn)
	case "=":
		f.res = xtypes.Bool(ln == rn)
	case "!=":
		f.res = xtypes.Bool(ln != rn)
	case "<":
		f.res = xtypes.Bool(ln < rn)
	case "<=":
		f.res = xtypes.Bool(ln <= rn)
	case ">":
		f.res = xtypes.Bool(ln > rn)
	case ">=":
		f.res = xtypes.Bool(ln >= rn)
	}

	return nil
}

func andOrOperator(left, right xtypes.Result, f *xpFilt, op string) error {
	lt, lOK := left.(xtypes.IsBool)
	rt, rOK := right.(xtypes.IsBool)
	if !lOK || !rOK {
		return fmt.Errorf("Cannot convert data type to number")
	}

	l, r := lt.Bool(), rt.Bool()

	if op == "and" {
		f.res = xtypes.Bool(l && r)
	} else {
		f.res = xtypes.Bool(l || r)
	}

	return nil
}

func unionOperator(left, right xtypes.Result, f *xpFilt, op string) error {
	lNode, lOK := left.(xtypes.NodeSet)
	rNode, rOK := right.(xtypes.NodeSet)
	if !lOK || !rOK {
		return fmt.Errorf("Cannot convert data type to node-set")
	}

	uniq := make(map[int]tree.Node)
	for _, i := range lNode {
		uniq[i.Pos()] = i
	}
	for _, i := range rNode {
		uniq[i.Pos()] = i
	}

	res := make(xtypes.NodeSet, 0, len(uniq))
	for _, v := range uniq {
		res = append(res, v)
	}

	f.res = res

	return nil
}
