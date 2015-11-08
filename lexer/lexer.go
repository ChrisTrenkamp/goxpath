package lexer

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	//XItemError is an error with the parser input
	XItemError XItemType = iota
	//XItemEOF is the end of the parser input
	XItemEOF
	//XItemAbsLocPath is an absolute path
	XItemAbsLocPath
	//XItemAbbrAbsLocPath represents an abbreviated absolute path
	XItemAbbrAbsLocPath
	//XItemAbbrRelLocPath marks the start of a path expression
	XItemAbbrRelLocPath
	//XItemRelLocPath represents a relative location path
	XItemRelLocPath
	//XItemEndPath marks the end of a path
	XItemEndPath
	//XItemAxis marks an axis specifier of a path
	XItemAxis
	//XItemAbbrAxis marks an abbreviated axis specifier (just @ at this point)
	XItemAbbrAxis
	//XItemNCName marks a namespace name in a node test
	XItemNCName
	//XItemQName marks the local name in an a node test
	XItemQName
	//XItemNodeType marks a node type in a node test
	XItemNodeType
	//XItemPredicate marks a predicate in an axis
	XItemPredicate
)

const (
	eof = -(iota + 1)
)

//XItemType is the parser token types
type XItemType int

//XItem is the token emitted from the parser
type XItem struct {
	Typ XItemType
	Val string
}

func (i XItem) String() string {
	switch i.Typ {
	case XItemEOF:
		return "EOF"
	case XItemError:
		return fmt.Sprintf("ERROR: %q", i.Val)
	}

	return fmt.Sprintf("%q", i.Val)
}

type stateFn func(*Lexer) stateFn

//Lexer lexes out XPath expressions
type Lexer struct {
	input string
	start int
	pos   int
	width int
	items chan XItem
	inter bool
}

//Lex an XPath expresion on the io.Reader
func Lex(xpath string) chan XItem {
	l := &Lexer{
		input: xpath,
		items: make(chan XItem),
	}
	go l.run()
	return l.items
}

func (l *Lexer) run() {
	for state := startState; state != nil && !l.inter; {
		state = state(l)
	}

	close(l.items)
}

func (l *Lexer) emit(t XItemType) {
	if !l.inter {
		l.items <- XItem{t, l.input[l.start:l.pos]}
	}
	l.start = l.pos
}

func (l *Lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])

	l.pos += l.width

	return r
}

func (l *Lexer) ignore() {
	l.start = l.pos
}

func (l *Lexer) backup() {
	l.pos -= l.width
}

func (l *Lexer) peek() rune {
	r := l.next()

	l.backup()
	return r
}

func (l *Lexer) peekAt(n int) rune {
	if n <= 1 {
		return l.peek()
	}

	width := 0
	var ret rune

	for count := 0; count < n; count++ {
		r, s := utf8.DecodeRuneInString(l.input[l.pos+width:])
		width += s

		if l.pos+width >= len(l.input) {
			return eof
		}

		ret = r
	}

	return ret
}

func (l *Lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}

	l.backup()
	return false
}

func (l *Lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- XItem{
		XItemError,
		fmt.Sprintf(format, args),
	}

	return nil
}

func startState(l *Lexer) stateFn {
	if string(l.next()) == "/" {
		l.ignore()

		if string(l.next()) == "/" {
			l.ignore()
			return abbrAbsLocPathState
		}

		l.backup()
		return absLocPathState
	}

	l.backup()
	return relLocPathState
}
