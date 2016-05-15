package intfns

import (
	"fmt"

	"github.com/ChrisTrenkamp/goxpath/xfn"
	"github.com/ChrisTrenkamp/goxpath/xtypes"
)

func _string(c xfn.Ctx, args ...xtypes.Result) (xtypes.Result, error) {
	if b, ok := args[0].(xtypes.IsString); ok {
		return xtypes.String(b.String()), nil
	}

	return nil, fmt.Errorf("Cannot convert object to a boolean")
}
