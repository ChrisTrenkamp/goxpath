package intfns

import "github.com/ChrisTrenkamp/goxpath/goxpath/xpfn"

//BuiltIn contains the list of built-in XPath functions
var BuiltIn = map[string]xpfn.Wrap{
	//Node set functions
	"last":          xpfn.Wrap{Fn: last},
	"position":      xpfn.Wrap{Fn: position},
	"count":         xpfn.Wrap{Fn: count, NArgs: 1},
	"local-name":    xpfn.Wrap{Fn: localName, NArgs: -1},
	"namespace-uri": xpfn.Wrap{Fn: namespaceURI, NArgs: -1},
	"name":          xpfn.Wrap{Fn: name, NArgs: -1},
	//boolean functions
	"boolean": xpfn.Wrap{Fn: boolean, NArgs: 1},
	"not":     xpfn.Wrap{Fn: not, NArgs: 1},
	"true":    xpfn.Wrap{Fn: _true},
	"false":   xpfn.Wrap{Fn: _false},
}
