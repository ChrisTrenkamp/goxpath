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
)

type nodeStack struct {
	n    *Node
	left bool
}

type parseStack struct {
	stack      []nodeStack
	stateTypes []stateType
	cur        *Node
}

func (p *parseStack) push(t stateType) {
	st := nodeStack{
		n: p.cur.Parent,
	}
	if p.cur.Parent != nil {
		if p.cur.Parent.Left == p.cur {
			st.left = true
		}
	}
	p.stack = append(p.stack, st)
	p.stateTypes = append(p.stateTypes, t)
}

func (p *parseStack) pop() error {
	if len(p.stack) == 0 {
		return fmt.Errorf("Malformed XPath expression.")
	}

	stackPos := len(p.stack) - 1
	st := p.stack[stackPos]

	if st.n == nil {
		for p.cur.Parent != nil {
			p.cur = p.cur.Parent
		}
	} else {
		if st.left {
			p.cur = st.n.Left
		} else {
			p.cur = st.n.Right
		}
	}

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
}

//Parse creates an AST tree for XPath expressions.
func Parse(xp string) (*Node, error) {
	var err error
	c := lexer.Lex(xp)
	p := &parseStack{cur: &Node{}}

	for next := range c {
		if err == nil {
			err = parseMap[next.Typ](p, next)
		}
	}

	n := p.cur
	for n.Parent != nil {
		n = n.Parent
	}

	return n, err
}

func xiError(p *parseStack, i lexer.XItem) error {
	return fmt.Errorf(i.Val)
}

func xiXPath(p *parseStack, i lexer.XItem) error {
	if p.curState() == xpathState {
		p.cur.Push(newNode(i))
		return nil
	}

	next := newNode(i)
	p.cur.Add(next)
	p.push(xpathState)
	if p.cur.Left != nil {
		p.cur = p.cur.Left
	} else if p.cur.Right != nil {
		p.cur = p.cur.Right
	}
	return nil
}

func xiEndPath(p *parseStack, i lexer.XItem) error {
	if err := p.pop(); err != nil {
		return err
	}

	if p.cur.Parent != nil {
		p.cur = p.cur.Parent
	}

	return nil
}

func xiFunc(p *parseStack, i lexer.XItem) error {
	p.cur.Push(newNode(i))
	p.push(funcState)
	return nil
}

func xiFuncArg(p *parseStack, i lexer.XItem) error {
	if p.curState() != paramState {
		p.cur.Push(newNode(i))
		p.push(paramState)
		p.cur.Push(&Node{})
	} else {
		err := p.pop()
		if err != nil {
			return err
		}
		p.cur.Add(newNode(i))
		p.cur = p.cur.Right
		p.cur.Push(&Node{})
	}
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
	p.cur.Push(newNode(i))
	p.push(predState)
	p.cur.Push(&Node{})
	return nil
}

func xiEndPred(p *parseStack, i lexer.XItem) error {
	return p.pop()
}

func xiStrLit(p *parseStack, i lexer.XItem) error {
	p.cur.Add(newNode(i))
	return nil
}

func xiNumLit(p *parseStack, i lexer.XItem) error {
	p.cur.Add(newNode(i))
	return nil
}
