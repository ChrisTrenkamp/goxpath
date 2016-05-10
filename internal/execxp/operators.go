package execxp

import (
	"math"

	"github.com/ChrisTrenkamp/goxpath/internal/parser/intfns"
	"github.com/ChrisTrenkamp/goxpath/literals/boollit"
	"github.com/ChrisTrenkamp/goxpath/literals/numlit"
	"github.com/ChrisTrenkamp/goxpath/literals/strlit"
	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/xfn"
)

func bothNodeOperator(left []tree.Node, right []tree.Node, f *xpFilt, op string) error {
	var err error
	for _, l := range left {
		for _, r := range right {
			lStr := l.ResValue()
			rStr := r.ResValue()

			if op == "=" || op == "!=" {
				err = equalsOperator([]tree.Res{strlit.StrLit(lStr)}, []tree.Res{strlit.StrLit(rStr)}, f, op)
				if err == nil && f.res[0].ResValue() == "true" {
					return nil
				}
			} else {
				err = numberOperator([]tree.Res{strlit.StrLit(lStr)}, []tree.Res{strlit.StrLit(rStr)}, f, op)
				if err == nil && f.res[0].ResValue() == "true" {
					return nil
				}
			}
		}
	}

	return nil
}

func leftNodeOperator(left []tree.Node, right []tree.Res, f *xpFilt, op string) error {
	var err error
	for _, l := range left {
		lStr := l.ResValue()

		if op == "=" || op == "!=" {
			err = equalsOperator([]tree.Res{strlit.StrLit(lStr)}, right, f, op)
			if err == nil && f.res[0].ResValue() == "true" {
				return nil
			}
		} else {
			err = numberOperator([]tree.Res{strlit.StrLit(lStr)}, right, f, op)
			if err == nil && f.res[0].ResValue() == "true" {
				return nil
			}
		}
	}

	return nil
}

func rightNodeOperator(left []tree.Res, right []tree.Node, f *xpFilt, op string) error {
	var err error
	for _, r := range right {
		rStr := r.ResValue()

		if op == "=" || op == "!=" {
			err = equalsOperator(left, []tree.Res{strlit.StrLit(rStr)}, f, op)
			if err == nil && f.res[0].ResValue() == "true" {
				return nil
			}
		} else {
			err = numberOperator(left, []tree.Res{strlit.StrLit(rStr)}, f, op)
			if err == nil && f.res[0].ResValue() == "true" {
				return nil
			}
		}
	}

	return nil
}

func equalsOperator(left, right []tree.Res, f *xpFilt, op string) error {
	lBool, lErr := xfn.GetBool(left, nil)
	rBool, rErr := xfn.GetBool(right, nil)
	if lErr == nil || rErr == nil {
		lBool = intfns.BooleanFunc(left)
		rBool = intfns.BooleanFunc(right)
		if op == "=" {
			f.res = []tree.Res{boollit.BoolLit(lBool == rBool)}
		} else {
			f.res = []tree.Res{boollit.BoolLit(lBool != rBool)}
		}

		return nil
	}

	_, lErr = xfn.GetNumber(left, nil)
	_, rErr = xfn.GetNumber(right, nil)
	if lErr == nil || rErr == nil {
		return numberOperator(left, right, f, op)
	}

	lStr := intfns.StringFunc(left)
	rStr := intfns.StringFunc(right)

	if op == "=" {
		f.res = []tree.Res{boollit.BoolLit(lStr == rStr)}
	} else {
		f.res = []tree.Res{boollit.BoolLit(lStr != rStr)}
	}

	return nil
}

func numberOperator(left, right []tree.Res, f *xpFilt, op string) error {
	ln, err := intfns.NumberFunc(left)
	if err != nil {
		return err
	}

	rn, err := intfns.NumberFunc(right)
	if err != nil {
		return err
	}

	switch op {
	case "*":
		f.res = []tree.Res{numlit.NumLit(ln * rn)}
	case "div":
		if rn != 0 {
			f.res = []tree.Res{numlit.NumLit(ln / rn)}
		} else {
			f.res = []tree.Res{numlit.NumLit(math.NaN())}
		}
	case "mod":
		f.res = []tree.Res{numlit.NumLit(int(ln) % int(rn))}
	case "+":
		f.res = []tree.Res{numlit.NumLit(ln + rn)}
	case "-":
		f.res = []tree.Res{numlit.NumLit(ln - rn)}
	case "=":
		f.res = []tree.Res{boollit.BoolLit(ln == rn)}
	case "!=":
		f.res = []tree.Res{boollit.BoolLit(ln != rn)}
	case "<":
		f.res = []tree.Res{boollit.BoolLit(ln < rn)}
	case "<=":
		f.res = []tree.Res{boollit.BoolLit(ln <= rn)}
	case ">":
		f.res = []tree.Res{boollit.BoolLit(ln > rn)}
	case ">=":
		f.res = []tree.Res{boollit.BoolLit(ln >= rn)}
	}

	return nil
}

func andOrOperator(left, right []tree.Res, f *xpFilt, op string) error {
	l := intfns.BooleanFunc(left)
	r := intfns.BooleanFunc(right)

	if op == "and" {
		f.res = []tree.Res{boollit.BoolLit(l && r)}
	} else {
		f.res = []tree.Res{boollit.BoolLit(l || r)}
	}

	return nil
}

func unionOperator(left, right []tree.Res, f *xpFilt, op string) error {
	lNode, err := xfn.GetNode(left, nil)
	if err != nil {
		return err
	}

	rNode, err := xfn.GetNode(right, nil)
	if err != nil {
		return err
	}

	uniq := make(map[int]tree.Node)
	for _, i := range lNode {
		uniq[i.Pos()] = i
	}
	for _, i := range rNode {
		uniq[i.Pos()] = i
	}

	f.res = make([]tree.Res, 0, len(uniq))
	for _, v := range uniq {
		f.res = append(f.res, v)
	}

	return nil
}
