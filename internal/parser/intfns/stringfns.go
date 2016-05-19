package intfns

import (
	"regexp"
	"strings"

	"github.com/ChrisTrenkamp/goxpath/xfn"
	"github.com/ChrisTrenkamp/goxpath/xtypes"
)

func _string(c xfn.Ctx, args ...xtypes.Result) (xtypes.Result, error) {
	if len(args) == 1 {
		return xtypes.String(args[0].String()), nil
	}

	return xtypes.String(c.Node.ResValue()), nil
}

func concat(c xfn.Ctx, args ...xtypes.Result) (xtypes.Result, error) {
	ret := ""

	for _, i := range args {
		ret += i.String()
	}

	return xtypes.String(ret), nil
}

func startsWith(c xfn.Ctx, args ...xtypes.Result) (xtypes.Result, error) {
	return xtypes.Bool(strings.Index(args[0].String(), args[1].String()) == 0), nil
}

func contains(c xfn.Ctx, args ...xtypes.Result) (xtypes.Result, error) {
	return xtypes.Bool(strings.Index(args[0].String(), args[1].String()) != -1), nil
}

func substringBefore(c xfn.Ctx, args ...xtypes.Result) (xtypes.Result, error) {
	ind := strings.Index(args[0].String(), args[1].String())
	if ind == -1 {
		return xtypes.String(""), nil
	}

	return xtypes.String(args[0].String()[:ind]), nil
}

func substringAfter(c xfn.Ctx, args ...xtypes.Result) (xtypes.Result, error) {
	ind := strings.Index(args[0].String(), args[1].String())
	if ind == -1 {
		return xtypes.String(""), nil
	}

	return xtypes.String(args[0].String()[ind+len(args[1].String()):]), nil
}

func substring(c xfn.Ctx, args ...xtypes.Result) (xtypes.Result, error) {
	str := args[0].String()

	bNum, bErr := round(c, args[1])
	if bErr != nil {
		return nil, bErr
	}

	b := int(bNum.(xtypes.Num).Num()) - 1

	if b < 0 {
		b = 0
	}

	if b >= len(str) {
		return xtypes.String(""), nil
	}

	if len(args) == 3 {
		eNum, eErr := round(c, args[1])
		if eErr != nil {
			return nil, eErr
		}

		e := int(eNum.(xtypes.Num).Num()) + 1
		if e <= 0 {
			return xtypes.String(""), nil
		}

		if b+e >= len(str) {
			e = len(str) - b + 1
		}

		return xtypes.String(str[b : b+e]), nil
	}

	return xtypes.String(str[b:]), nil
}

func stringLength(c xfn.Ctx, args ...xtypes.Result) (xtypes.Result, error) {
	var str string
	if len(args) == 1 {
		str = args[0].String()
	} else {
		str = c.Node.ResValue()
	}

	return xtypes.Num(len(str)), nil
}

var spaceTrim = regexp.MustCompile("\\s+")

func normalizeSpace(c xfn.Ctx, args ...xtypes.Result) (xtypes.Result, error) {
	var str string
	if len(args) == 1 {
		str = args[0].String()
	} else {
		str = c.Node.ResValue()
	}

	str = strings.TrimSpace(str)

	return xtypes.String(spaceTrim.ReplaceAllString(str, " ")), nil
}

func translate(c xfn.Ctx, args ...xtypes.Result) (xtypes.Result, error) {
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

	return xtypes.String(ret), nil
}
