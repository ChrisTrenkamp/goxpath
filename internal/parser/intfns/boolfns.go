package intfns

import (
	"fmt"

	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/xfn"
	"github.com/ChrisTrenkamp/goxpath/xtypes"
	"golang.org/x/text/language"
)

func boolean(c xfn.Ctx, args ...xtypes.Result) (xtypes.Result, error) {
	if b, ok := args[0].(xtypes.IsBool); ok {
		return b.Bool(), nil
	}

	return nil, fmt.Errorf("Cannot convert object to a boolean")
}

func not(c xfn.Ctx, args ...xtypes.Result) (xtypes.Result, error) {
	b, ok := args[0].(xtypes.Bool)
	if !ok {
		return nil, fmt.Errorf("Object is not a boolean")
	}
	return !b, nil
}

func _true(c xfn.Ctx, args ...xtypes.Result) (xtypes.Result, error) {
	return xtypes.Bool(true), nil
}

func _false(c xfn.Ctx, args ...xtypes.Result) (xtypes.Result, error) {
	return xtypes.Bool(false), nil
}

func lang(c xfn.Ctx, args ...xtypes.Result) (xtypes.Result, error) {
	lStr, ok := args[0].(xtypes.String)
	if !ok {
		return nil, fmt.Errorf("Argument is not a string")
	}

	srcLang := language.Make(string(lStr))
	srcRegion, srcRegionConf := srcLang.Region()
	match := language.NewMatcher([]language.Tag{srcLang})

	if elem, ok := c.Node.(tree.Elem); ok {
		if attr, ok := tree.GetAttribute(elem, "lang", tree.XMLSpace); ok {
			targLang := language.Make(attr.Value)
			targRegion, targRegionConf := targLang.Region()
			if srcRegionConf == language.Exact && targRegionConf != language.Exact {
				return xtypes.Bool(false), nil
			}
			if srcRegion != targRegion && srcRegionConf == language.Exact && targRegionConf == language.Exact {
				return xtypes.Bool(false), nil
			}
			_, _, conf := match.Match(language.Make(attr.Value))
			return xtypes.Bool(conf >= language.High), nil
		}
	}

	n := c.Node.GetParent()
	for n.GetNodeType() != tree.NtRoot {
		if attr, ok := tree.GetAttribute(n, "lang", tree.XMLSpace); ok {
			targLang := language.Make(attr.Value)
			targRegion, targRegionConf := targLang.Region()
			if srcRegionConf == language.Exact && targRegionConf != language.Exact {
				return xtypes.Bool(false), nil
			}
			if srcRegion != targRegion && srcRegionConf == language.Exact && targRegionConf == language.Exact {
				return xtypes.Bool(false), nil
			}
			_, _, conf := match.Match(language.Make(attr.Value))
			return xtypes.Bool(conf >= language.High), nil
		}
	}

	return xtypes.Bool(false), nil
}
