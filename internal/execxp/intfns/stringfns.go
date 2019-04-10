package intfns

import (
	"math"
	"regexp"
	"strings"

	"github.com/ChrisTrenkamp/goxpath/tree"
)

func _string(a tree.Adapter, c tree.Ctx, args ...tree.Result) (tree.Result, error) {
	if len(args) == 1 {
		return tree.String(args[0].String()), nil
	}

	return tree.String(c.NodeSet.String()), nil
}

func concat(a tree.Adapter, c tree.Ctx, args ...tree.Result) (tree.Result, error) {
	ret := ""

	for _, i := range args {
		ret += i.String()
	}

	return tree.String(ret), nil
}

func startsWith(a tree.Adapter, c tree.Ctx, args ...tree.Result) (tree.Result, error) {
	return tree.Bool(strings.Index(args[0].String(), args[1].String()) == 0), nil
}

func contains(a tree.Adapter, c tree.Ctx, args ...tree.Result) (tree.Result, error) {
	return tree.Bool(strings.Contains(args[0].String(), args[1].String())), nil
}

func substringBefore(a tree.Adapter, c tree.Ctx, args ...tree.Result) (tree.Result, error) {
	ind := strings.Index(args[0].String(), args[1].String())
	if ind == -1 {
		return tree.String(""), nil
	}

	return tree.String(args[0].String()[:ind]), nil
}

func substringAfter(a tree.Adapter, c tree.Ctx, args ...tree.Result) (tree.Result, error) {
	ind := strings.Index(args[0].String(), args[1].String())
	if ind == -1 {
		return tree.String(""), nil
	}

	return tree.String(args[0].String()[ind+len(args[1].String()):]), nil
}

func substring(a tree.Adapter, c tree.Ctx, args ...tree.Result) (tree.Result, error) {
	str := args[0].String()

	bNum, bErr := round(a, c, args[1])
	if bErr != nil {
		return nil, bErr
	}

	b := bNum.(tree.Num).Num()

	if float64(b-1) >= float64(len(str)) || math.IsNaN(float64(b)) {
		return tree.String(""), nil
	}

	if len(args) == 2 {
		if b <= 1 {
			b = 1
		}

		return tree.String(str[int(b)-1:]), nil
	}

	eNum, eErr := round(a, c, args[2])
	if eErr != nil {
		return nil, eErr
	}

	e := eNum.(tree.Num).Num()

	if e <= 0 || math.IsNaN(float64(e)) || (math.IsInf(float64(b), 0) && math.IsInf(float64(e), 0)) {
		return tree.String(""), nil
	}

	if b <= 1 {
		e = b + e - 1
		b = 1
	}

	if float64(b+e-1) >= float64(len(str)) {
		e = tree.Num(len(str)) - b + 1
	}

	return tree.String(str[int(b)-1 : int(b+e)-1]), nil
}

func stringLength(a tree.Adapter, c tree.Ctx, args ...tree.Result) (tree.Result, error) {
	var str string
	if len(args) == 1 {
		str = args[0].String()
	} else {
		str = c.NodeSet.String()
	}

	return tree.Num(len(str)), nil
}

var spaceTrim = regexp.MustCompile(`\s+`)

func normalizeSpace(a tree.Adapter, c tree.Ctx, args ...tree.Result) (tree.Result, error) {
	var str string
	if len(args) == 1 {
		str = args[0].String()
	} else {
		str = c.NodeSet.String()
	}

	str = strings.TrimSpace(str)

	return tree.String(spaceTrim.ReplaceAllString(str, " ")), nil
}

func translate(a tree.Adapter, c tree.Ctx, args ...tree.Result) (tree.Result, error) {
	ret := args[0].String()
	src := args[1].String()
	repl := args[2].String()

	for i := range src {
		r := ""
		if i < len(repl) {
			r = string(repl[i])
		}

		ret = strings.Replace(ret, string(src[i]), r, -1)
	}

	return tree.String(ret), nil
}
