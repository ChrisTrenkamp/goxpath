package intfns

import "github.com/ChrisTrenkamp/goxpath/xfn"

//BuiltIn contains the list of built-in XPath functions
var BuiltIn = map[string]xfn.Wrap{
	//String functions
	"string":           xfn.Wrap{Fn: _string, NArgs: 1, LastArgOpt: xfn.Optional},
	"concat":           xfn.Wrap{Fn: concat, NArgs: 3, LastArgOpt: xfn.Variadic},
	"starts-with":      xfn.Wrap{Fn: startsWith, NArgs: 2},
	"contains":         xfn.Wrap{Fn: contains, NArgs: 2},
	"substring-before": xfn.Wrap{Fn: substringBefore, NArgs: 2},
	"substring-after":  xfn.Wrap{Fn: substringAfter, NArgs: 2},
	"substring":        xfn.Wrap{Fn: substring, NArgs: 3, LastArgOpt: xfn.Optional},
	"string-length":    xfn.Wrap{Fn: stringLength, NArgs: 1, LastArgOpt: xfn.Optional},
	"normalize-space":  xfn.Wrap{Fn: normalizeSpace, NArgs: 1, LastArgOpt: xfn.Optional},
	"translate":        xfn.Wrap{Fn: translate, NArgs: 3},
	//Node set functions
	"last":          xfn.Wrap{Fn: last},
	"position":      xfn.Wrap{Fn: position},
	"count":         xfn.Wrap{Fn: count, NArgs: 1},
	"local-name":    xfn.Wrap{Fn: localName, NArgs: 1, LastArgOpt: xfn.Optional},
	"namespace-uri": xfn.Wrap{Fn: namespaceURI, NArgs: 1, LastArgOpt: xfn.Optional},
	"name":          xfn.Wrap{Fn: name, NArgs: 1, LastArgOpt: xfn.Optional},
	//boolean functions
	"boolean": xfn.Wrap{Fn: boolean, NArgs: 1},
	"not":     xfn.Wrap{Fn: not, NArgs: 1},
	"true":    xfn.Wrap{Fn: _true},
	"false":   xfn.Wrap{Fn: _false},
	"lang":    xfn.Wrap{Fn: lang, NArgs: 1},
	//number functions
	"number":  xfn.Wrap{Fn: number, NArgs: 1, LastArgOpt: xfn.Optional},
	"sum":     xfn.Wrap{Fn: sum, NArgs: 1},
	"floor":   xfn.Wrap{Fn: floor, NArgs: 1},
	"ceiling": xfn.Wrap{Fn: ceiling, NArgs: 1},
	"round":   xfn.Wrap{Fn: round, NArgs: 1},
}
