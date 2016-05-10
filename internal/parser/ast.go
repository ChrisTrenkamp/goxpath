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
}

func newNode(val lexer.XItem) *Node {
	return &Node{Val: val}
}

func (n *Node) swapVals(i *Node) {
	n.Val, i.Val = i.Val, n.Val
}

var beginPathType = map[lexer.XItemType]bool{
	lexer.XItemAbsLocPath:     true,
	lexer.XItemAbbrAbsLocPath: true,
	lexer.XItemAbbrRelLocPath: true,
	lexer.XItemRelLocPath:     true,
}

//Add appends i into n.
//If n is empty, i's values are copied over
//If n's left child is empty, i is appended to the left and the value's are swapped
//If n's right child is empty, i is appended to the right
//Otherwise, the value's are swapped, n's children is moved to i's, and i is set to n's Left
func (n *Node) Add(i *Node) {
	if n.Val.Typ == Empty {
		n.Val = i.Val
	} else if n.Left == nil {
		n.Left = i
		i.Parent = n
		n.swapVals(i)
	} else if beginPathType[n.Val.Typ] {
		i.Left = n.Left
		n.Left = i
		i.Parent = n
		n.swapVals(i)
	} else if n.Right == nil {
		n.Right = i
		i.Parent = n
	} else {
		n.swapVals(i)
		i.Left, i.Right = n.Left, n.Right
		i.Left.Parent, i.Right.Parent = i, i
		n.Left, n.Right = i, nil
	}
}

//Push appends i to n's Left if it's empty.  Otherwise, it adds n to i's right.
func (n *Node) Push(i *Node) {
	if n.Parent != nil {
		i.Parent = n.Parent
		if n.Parent.Left == n {
			n.Parent.Left = i
		} else {
			n.Parent.Right = i
		}
		i.Left, i.Right = n.Left, n.Right
		n.Left, n.Right = nil, nil
	}
	n.Parent = i
	n.swapVals(i)

	if i.Left == nil {
		i.Left = n
	} else if i.Right == nil {
		i.Left.Parent = i
		i.Right = n
	} else {
		i.Left.Parent = i
		n.Left = i.Right
		i.Right = n
		i.Right.Parent = i
		n.Left.Parent = n
	}
}

//PushNotEmpty add's i if n has an empty value.  Othwise, i is pushed.
func (n *Node) PushNotEmpty(i *Node) {
	if n.Val.Typ == Empty {
		n.Add(i)
	} else {
		n.Push(i)
	}
}
