package parser

import (
	"fmt"

	"github.com/ChrisTrenkamp/goxpath/internal/lexer"
)

type stateType int

const (
	defState stateType = iota
	xpathState
	funcState
	paramState
	predState
	parenState
)

type parseStack struct {
	stack      []*Node
	stateTypes []stateType
	cur        *Node
}

func (p *parseStack) push(t stateType) {
	p.stack = append(p.stack, p.cur)
	p.stateTypes = append(p.stateTypes, t)
}

func (p *parseStack) pop() error {
	if len(p.stack) == 0 {
		return fmt.Errorf("Malformed XPath expression.")
	}

	stackPos := len(p.stack) - 1

	p.cur = p.stack[stackPos]
	p.stack = p.stack[:stackPos]
	p.stateTypes = p.stateTypes[:stackPos]
	return nil
}

func (p *parseStack) curState() stateType {
	if len(p.stateTypes) == 0 {
		return defState
	}
	return p.stateTypes[len(p.stateTypes)-1]
}

type lexFn func(*parseStack, lexer.XItem) error

var parseMap = map[lexer.XItemType]lexFn{
	lexer.XItemError:          xiError,
	lexer.XItemAbsLocPath:     xiXPath,
	lexer.XItemAbbrAbsLocPath: xiXPath,
	lexer.XItemAbbrRelLocPath: xiXPath,
	lexer.XItemRelLocPath:     xiXPath,
	lexer.XItemEndPath:        xiEndPath,
	lexer.XItemAxis:           xiXPath,
	lexer.XItemAbbrAxis:       xiXPath,
	lexer.XItemNCName:         xiXPath,
	lexer.XItemQName:          xiXPath,
	lexer.XItemNodeType:       xiXPath,
	lexer.XItemProcLit:        xiXPath,
	lexer.XItemFunction:       xiFunc,
	lexer.XItemArgument:       xiFuncArg,
	lexer.XItemEndFunction:    xiEndFunc,
	lexer.XItemPredicate:      xiPred,
	lexer.XItemEndPredicate:   xiEndPred,
	lexer.XItemStrLit:         xiStrLit,
	lexer.XItemNumLit:         xiNumLit,
	lexer.XItemOperator:       xiOp,
}

var opPrecedence = map[string]int{
	"|":   1,
	"*":   2,
	"div": 2,
	"mod": 2,
	"+":   3,
	"-":   3,
	"=":   4,
	"!=":  4,
	"<":   4,
	"<=":  4,
	">":   4,
	">=":  4,
	"and": 5,
	"or":  6,
}

//Parse creates an AST tree for XPath expressions.
func Parse(xp string) (*Node, error) {
	var err error
	c := lexer.Lex(xp)
	n := &Node{}
	p := &parseStack{cur: n}

	for next := range c {
		if err == nil {
			err = parseMap[next.Typ](p, next)
			/*
				fmt.Println(next, "Parent:", p.cur.Parent, "Left:", p.cur.Left, "Right:", p.cur.Right)

				for i := 0; i < 7; i++ {
					n.prettyPrint(i, 16)
				}
			*/
		}
	}

	return n, err
}

func xiError(p *parseStack, i lexer.XItem) error {
	return fmt.Errorf(i.Val)
}

func xiXPath(p *parseStack, i lexer.XItem) error {
	if p.curState() == xpathState {
		p.cur.push(i)
		p.cur = p.cur.next
		return nil
	}

	p.cur.pushNotEmpty(i)
	p.push(xpathState)
	p.cur = p.cur.next
	return nil
}

func xiEndPath(p *parseStack, i lexer.XItem) error {
	return p.pop()
}

func xiFunc(p *parseStack, i lexer.XItem) error {
	p.cur.push(i)
	p.cur = p.cur.next
	p.push(funcState)
	return nil
}

func xiFuncArg(p *parseStack, i lexer.XItem) error {
	if p.curState() != funcState {
		err := p.pop()
		if err != nil {
			return err
		}
	}

	p.cur.push(i)
	p.cur = p.cur.next
	p.push(paramState)
	p.cur.push(lexer.XItem{Typ: Empty, Val: ""})
	p.cur = p.cur.next
	return nil
}

func xiEndFunc(p *parseStack, i lexer.XItem) error {
	if p.curState() == paramState {
		err := p.pop()
		if err != nil {
			return err
		}
	}
	return p.pop()
}

func xiPred(p *parseStack, i lexer.XItem) error {
	p.cur.push(i)
	p.cur = p.cur.next
	p.push(predState)
	p.cur.push(lexer.XItem{Typ: Empty, Val: ""})
	p.cur = p.cur.next
	return nil
}

func xiEndPred(p *parseStack, i lexer.XItem) error {
	return p.pop()
}

func xiStrLit(p *parseStack, i lexer.XItem) error {
	p.cur.add(i)
	return nil
}

func xiNumLit(p *parseStack, i lexer.XItem) error {
	p.cur.add(i)
	return nil
}

func xiOp(p *parseStack, i lexer.XItem) error {
	if i.Val == "(" {
		p.cur.push(lexer.XItem{Typ: Empty, Val: ""})
		p.push(parenState)
		p.cur = p.cur.next
		return nil
	}

	if i.Val == ")" {
		return p.pop()
	}

	if p.cur.Val.Typ == lexer.XItemOperator {
		if opPrecedence[p.cur.Val.Val] <= opPrecedence[i.Val] {
			p.cur.add(i)
		} else {
			p.cur.push(i)
		}
	} else {
		p.cur.add(i)
	}
	p.cur = p.cur.next

	return nil
}
