package intfns

import (
	"strconv"

	"github.com/ChrisTrenkamp/goxpath/literals/boollit"
	"github.com/ChrisTrenkamp/goxpath/literals/numlit"
	"github.com/ChrisTrenkamp/goxpath/literals/strlit"
	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/xfn"
)

func _string(c xfn.Ctx, args ...[]tree.Res) ([]tree.Res, error) {
	return []tree.Res{strlit.StrLit(StringFunc(args[0]))}, nil
}

//StringFunc returns the XPath string value of the given argument.
func StringFunc(args []tree.Res) string {
	if len(args) == 0 {
		return ""
	}

	switch t := args[0].(type) {
	case tree.Node:
		return t.ResValue()
	case boollit.BoolLit:
		if t {
			return "true"
		}
		return "false"
	case numlit.NumLit:
		return strconv.FormatFloat(float64(t), 'E', -1, 64)
	case strlit.StrLit:
		return string(t)
	}

	return ""
}
