package intfns

import (
	"github.com/ChrisTrenkamp/goxpath/goxpath/xpfn"
	"github.com/ChrisTrenkamp/goxpath/goxpath/xpfn/arg"
	"github.com/ChrisTrenkamp/goxpath/goxpath/xpfn/noarg"
	"github.com/ChrisTrenkamp/goxpath/goxpath/xpfn/optarg"
)

//BuiltIn contains the list of built-in XPath functions
var BuiltIn = map[string]xpfn.XPFn{
	//Node set functions
	"last":          noarg.Wrap{Fn: last},
	"count":         noarg.Wrap{Fn: count},
	"local-name":    optarg.Wrap{Fn: localName},
	"namespace-uri": optarg.Wrap{Fn: namespaceURI},
	"name":          optarg.Wrap{Fn: name},
	//boolean functions
	"boolean": arg.Wrap{Fn: boolean, NArgs: 1},
	"not":     arg.Wrap{Fn: not, NArgs: 1},
	"true":    noarg.Wrap{Fn: _true},
	"false":   noarg.Wrap{Fn: _false},
}
