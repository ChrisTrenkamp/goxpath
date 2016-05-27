package intfns

import "github.com/ChrisTrenkamp/goxpath/xfn"

//BuiltIn contains the list of built-in XPath functions
var BuiltIn = map[string]xfn.Wrap{
	//String functions
	"string":           {Fn: _string, NArgs: 1, LastArgOpt: xfn.Optional},
	"concat":           {Fn: concat, NArgs: 3, LastArgOpt: xfn.Variadic},
	"starts-with":      {Fn: startsWith, NArgs: 2},
	"contains":         {Fn: contains, NArgs: 2},
	"substring-before": {Fn: substringBefore, NArgs: 2},
	"substring-after":  {Fn: substringAfter, NArgs: 2},
	"substring":        {Fn: substring, NArgs: 3, LastArgOpt: xfn.Optional},
	"string-length":    {Fn: stringLength, NArgs: 1, LastArgOpt: xfn.Optional},
	"normalize-space":  {Fn: normalizeSpace, NArgs: 1, LastArgOpt: xfn.Optional},
	"translate":        {Fn: translate, NArgs: 3},
	//Node set functions
	"last":          {Fn: last},
	"position":      {Fn: position},
	"count":         {Fn: count, NArgs: 1},
	"local-name":    {Fn: localName, NArgs: 1, LastArgOpt: xfn.Optional},
	"namespace-uri": {Fn: namespaceURI, NArgs: 1, LastArgOpt: xfn.Optional},
	"name":          {Fn: name, NArgs: 1, LastArgOpt: xfn.Optional},
	//boolean functions
	"boolean": {Fn: boolean, NArgs: 1},
	"not":     {Fn: not, NArgs: 1},
	"true":    {Fn: _true},
	"false":   {Fn: _false},
	"lang":    {Fn: lang, NArgs: 1},
	//number functions
	"number":  {Fn: number, NArgs: 1, LastArgOpt: xfn.Optional},
	"sum":     {Fn: sum, NArgs: 1},
	"floor":   {Fn: floor, NArgs: 1},
	"ceiling": {Fn: ceiling, NArgs: 1},
	"round":   {Fn: round, NArgs: 1},
}
