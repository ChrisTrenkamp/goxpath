package execxp

import (
	"fmt"
	"math"

	"github.com/ChrisTrenkamp/goxpath/tree"
)

func bothNodeOperator(a tree.Adapter, left tree.NodeSet, right tree.NodeSet, f *xpFilt, op string) error {
	var err error
	for _, l := range left.GetNodes() {
		for _, r := range right.GetNodes() {
			lStr := a.StringValue(l)
			rStr := a.StringValue(r)

			if eqOps[op] {
				err = equalsOperator(tree.String(lStr), tree.String(rStr), f, op)
				if err == nil && f.ctx.String() == tree.True {
					return nil
				}
			} else {
				err = numberOperator(tree.String(lStr), tree.String(rStr), f, op)
				if err == nil && f.ctx.String() == tree.True {
					return nil
				}
			}
		}
	}

	f.ctx = tree.Bool(false)

	return nil
}

func leftNodeOperator(a tree.Adapter, left tree.NodeSet, right tree.Result, f *xpFilt, op string) error {
	var err error
	for _, l := range left.GetNodes() {
		lStr := a.StringValue(l)

		if eqOps[op] {
			err = equalsOperator(tree.String(lStr), right, f, op)
			if err == nil && f.ctx.String() == tree.True {
				return nil
			}
		} else {
			err = numberOperator(tree.String(lStr), right, f, op)
			if err == nil && f.ctx.String() == tree.True {
				return nil
			}
		}
	}

	f.ctx = tree.Bool(false)

	return nil
}

func rightNodeOperator(a tree.Adapter, left tree.Result, right tree.NodeSet, f *xpFilt, op string) error {
	var err error
	for _, r := range right.GetNodes() {
		rStr := a.StringValue(r)

		if eqOps[op] {
			err = equalsOperator(left, tree.String(rStr), f, op)
			if err == nil && f.ctx.String() == "true" {
				return nil
			}
		} else {
			err = numberOperator(left, tree.String(rStr), f, op)
			if err == nil && f.ctx.String() == "true" {
				return nil
			}
		}
	}

	f.ctx = tree.Bool(false)

	return nil
}

func equalsOperator(left, right tree.Result, f *xpFilt, op string) error {
	_, lOK := left.(tree.Bool)
	_, rOK := right.(tree.Bool)

	if lOK || rOK {
		lTest, lt := left.(tree.IsBool)
		rTest, rt := right.(tree.IsBool)
		if !lt || !rt {
			return fmt.Errorf("Cannot convert argument to boolean")
		}

		if op == "=" {
			f.ctx = tree.Bool(lTest.Bool() == rTest.Bool())
		} else {
			f.ctx = tree.Bool(lTest.Bool() != rTest.Bool())
		}

		return nil
	}

	_, lOK = left.(tree.Num)
	_, rOK = right.(tree.Num)
	if lOK || rOK {
		return numberOperator(left, right, f, op)
	}

	lStr := left.String()
	rStr := right.String()

	if op == "=" {
		f.ctx = tree.Bool(lStr == rStr)
	} else {
		f.ctx = tree.Bool(lStr != rStr)
	}

	return nil
}

func numberOperator(left, right tree.Result, f *xpFilt, op string) error {
	lt, lOK := left.(tree.IsNum)
	rt, rOK := right.(tree.IsNum)
	if !lOK || !rOK {
		return fmt.Errorf("Cannot convert data type to number")
	}

	ln, rn := lt.Num(), rt.Num()

	switch op {
	case "*":
		f.ctx = ln * rn
	case "div":
		if rn != 0 {
			f.ctx = ln / rn
		} else {
			if ln == 0 {
				f.ctx = tree.Num(math.NaN())
			} else {
				if math.Signbit(float64(ln)) == math.Signbit(float64(rn)) {
					f.ctx = tree.Num(math.Inf(1))
				} else {
					f.ctx = tree.Num(math.Inf(-1))
				}
			}
		}
	case "mod":
		f.ctx = tree.Num(int(ln) % int(rn))
	case "+":
		f.ctx = ln + rn
	case "-":
		f.ctx = ln - rn
	case "=":
		f.ctx = tree.Bool(ln == rn)
	case "!=":
		f.ctx = tree.Bool(ln != rn)
	case "<":
		f.ctx = tree.Bool(ln < rn)
	case "<=":
		f.ctx = tree.Bool(ln <= rn)
	case ">":
		f.ctx = tree.Bool(ln > rn)
	case ">=":
		f.ctx = tree.Bool(ln >= rn)
	}

	return nil
}

func andOrOperator(left, right tree.Result, f *xpFilt, op string) error {
	lt, lOK := left.(tree.IsBool)
	rt, rOK := right.(tree.IsBool)

	if !lOK || !rOK {
		return fmt.Errorf("Cannot convert argument to boolean")
	}

	l, r := lt.Bool(), rt.Bool()

	if op == "and" {
		f.ctx = l && r
	} else {
		f.ctx = l || r
	}

	return nil
}

func unionOperator(a tree.Adapter, left, right tree.Result, f *xpFilt, op string) error {
	lNode, lOK := left.(tree.NodeSet)
	rNode, rOK := right.(tree.NodeSet)

	if !lOK || !rOK {
		return fmt.Errorf("Cannot convert data type to node-set")
	}
	nodeset := a.NewNodeSet(append(lNode.GetNodes(), rNode.GetNodes()...))
	nodeset.Unique()
	f.ctx = nodeset

	return nil
}
