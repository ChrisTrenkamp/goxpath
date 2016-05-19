package parser

import "github.com/ChrisTrenkamp/goxpath/internal/lexer"

//NodeType enumerations
const (
	Empty lexer.XItemType = ""
)

//Node builds an AST tree for operating on XPath expressions
type Node struct {
	Val    lexer.XItem
	Left   *Node
	Right  *Node
	Parent *Node
	next   *Node
}

var beginPathType = map[lexer.XItemType]bool{
	lexer.XItemAbsLocPath:     true,
	lexer.XItemAbbrAbsLocPath: true,
	lexer.XItemAbbrRelLocPath: true,
	lexer.XItemRelLocPath:     true,
}

func (n *Node) add(i lexer.XItem) {
	if n.Val.Typ == Empty {
		n.Val = i
	} else if n.Left == nil {
		n.Left = &Node{Val: n.Val, Parent: n}
		n.Val = i
	} else if beginPathType[n.Val.Typ] {
		next := &Node{Val: n.Val, Left: n.Left, Parent: n}
		n.Left = next
		n.Val = i
	} else if n.Right == nil {
		n.Right = &Node{Val: i, Parent: n}
	} else {
		next := &Node{Val: n.Val, Left: n.Left, Right: n.Right, Parent: n}
		n.Left, n.Right = next, nil
		n.Val = i
	}
	n.next = n
}

func (n *Node) push(i lexer.XItem) {
	if n.Left == nil {
		n.Left = &Node{Val: i, Parent: n}
		n.next = n.Left
	} else if n.Right == nil {
		n.Right = &Node{Val: i, Parent: n}
		n.next = n.Right
	} else {
		next := &Node{Val: i, Left: n.Right, Parent: n}
		n.Right = next
		n.next = n.Right
	}
}

func (n *Node) pushNotEmpty(i lexer.XItem) {
	if n.Val.Typ == Empty {
		n.add(i)
	} else {
		n.push(i)
	}
}
