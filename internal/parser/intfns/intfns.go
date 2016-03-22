package intfns

import "github.com/ChrisTrenkamp/goxpath/xfn"

//BuiltIn contains the list of built-in XPath functions
var BuiltIn = map[string]xfn.Wrap{
	//Node set functions
	"last":          xfn.Wrap{Fn: last},
	"position":      xfn.Wrap{Fn: position},
	"count":         xfn.Wrap{Fn: count, NArgs: 1},
	"local-name":    xfn.Wrap{Fn: localName, NArgs: -1},
	"namespace-uri": xfn.Wrap{Fn: namespaceURI, NArgs: -1},
	"name":          xfn.Wrap{Fn: name, NArgs: -1},
	//boolean functions
	"boolean": xfn.Wrap{Fn: boolean, NArgs: 1},
	"not":     xfn.Wrap{Fn: not, NArgs: 1},
	"true":    xfn.Wrap{Fn: _true},
	"false":   xfn.Wrap{Fn: _false},
}
