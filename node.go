package rql

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

var textFormat = "%s" // changed to "%q" in tests for better error messages

// A Node is an element in the parse tree. The interface
// contains an unexported method so that only types local
// to this package can satisfy it.
type Node interface {
	Type() NodeType
	String() string

	// Copy does a deep copy of the Node and all its
	// components. To avoid type assertions, some Nodes
	// have specialized copy methods.
	Copy() Node

	// byte position of the start of the node in full original input string
	Position() Pos

	// tree returns the containing *Tree.
	// It is unexported so all implementations of Node are in this package
	tree() *Tree
}

// NodeType identifies the type of the parse tree node
type NodeType int

// Pos represents a byte position in the original input text from which
// this statement was parsed
type Pos int

// Position returns itself
func (p Pos) Position() Pos {
	return p
}

// Type returns itself and provides an easy default implementation
// for embedding in a Node. Embedded in all non-trivial Nodes.
func (t NodeType) Type() NodeType {
	return t
}

// NodeType constants
const (
	NodeBool       NodeType = iota // A boolean constant.
	NodeIdentifier                 // An identifier
	NodeNull                       // A null constant
	NodeString                     // A string constant
	NodeNumber                     // A numeric constant
	NodeOperator                   // An operator
	NodeList                       // A list of nodes.
	NodeStatement                  // A statement node.
)

// ListNode holds a sequence of Nodes
type ListNode struct {
	NodeType
	Pos
	tr    *Tree
	Nodes []Node // The element nodes in lexical order
}

func (t *Tree) newList(pos Pos) *ListNode {
	return &ListNode{tr: t, NodeType: NodeList, Pos: pos}
}

func (l *ListNode) append(n Node) {
	l.Nodes = append(l.Nodes, n)
}

func (l *ListNode) tree() *Tree {
	return l.tr
}

// String returns the ListNode as a string
func (l *ListNode) String() string {
	b := new(bytes.Buffer)
	for i, n := range l.Nodes {
		if i > 0 {
			fmt.Fprintf(b, ",")
		}
		fmt.Fprint(b, n)
	}
	return "(" + b.String() + ")"
}

// CopyList will copy and return *ListNode
func (l *ListNode) CopyList() *ListNode {
	if l == nil {
		return l
	}
	n := l.tr.newList(l.Pos)
	for _, elem := range l.Nodes {
		n.append(elem.Copy())
	}
	return n
}

// Copy runs CopyList but returns as Node
func (l *ListNode) Copy() Node {
	return l.CopyList()
}

// StatementNode holds a statement.
type StatementNode struct {
	NodeType
	Pos
	Operator *OperatorNode
	tr       *Tree
}

func (t *Tree) newStatement(pos Pos, op *OperatorNode) *StatementNode {
	return &StatementNode{
		NodeType: NodeStatement,
		Pos:      pos,
		tr:       t,
		Operator: op,
	}
}

// CopyStatement returns a copy of the StatementNode as a *StatementNode
func (s *StatementNode) CopyStatement() *StatementNode {
	return s.tr.newStatement(s.Pos, s.Operator)
}

// Copy runs CopyStatement, returning as a Node
func (s *StatementNode) Copy() Node {
	return s.CopyStatement()
}
func (s *StatementNode) tree() *Tree {
	return s.tr
}

// String returns the StatementNode as a string
func (s *StatementNode) String() string {
	if s.Operator != nil {
		return s.Operator.String()
	}
	return ""
}

// IdentifierNode holds an identifier.
type IdentifierNode struct {
	NodeType
	Pos
	tr    *Tree
	Ident string // The identifier's name.
}

// NewIdentifier creates an IdentifierNode
func NewIdentifier(ident string) *IdentifierNode {
	return &IdentifierNode{NodeType: NodeIdentifier, Ident: ident}
}

// SetPos sets the position. Chained for convenience.
func (i *IdentifierNode) SetPos(pos Pos) *IdentifierNode {
	i.Pos = pos
	return i
}

// SetTree sets the tree. Chained for convenience.
func (i *IdentifierNode) SetTree(t *Tree) *IdentifierNode {
	i.tr = t
	return i
}

func (i *IdentifierNode) String() string {
	return i.Ident
}

func (i *IdentifierNode) tree() *Tree {
	return i.tr
}

// Copy copies the IdentifierNode
func (i *IdentifierNode) Copy() Node {
	return NewIdentifier(i.Ident).SetTree(i.tr).SetPos(i.Pos)
}

// NullNode holds the special identifier 'null'
type NullNode struct {
	NodeType
	Pos
	tr *Tree
}

func (t *Tree) newNull(pos Pos) *NullNode {
	return &NullNode{tr: t, NodeType: NodeNull, Pos: pos}
}

// Type returns the NodeType value
func (n *NullNode) Type() NodeType {
	return NodeNull
}

// String returns the string representation of NullNode
func (n *NullNode) String() string {
	return "null"
}

func (n *NullNode) tree() *Tree {
	return n.tr
}

// Copy returns a copy of the NullNode
func (n *NullNode) Copy() Node {
	return n.tr.newNull(n.Pos)
}

// BoolNode holds a boolean constant.
type BoolNode struct {
	NodeType
	Pos
	tr   *Tree
	True bool // the value
}

func (t *Tree) newBool(pos Pos, true bool) *BoolNode {
	return &BoolNode{tr: t, NodeType: NodeBool, Pos: pos, True: true}
}

// String returns a string representation of the BoolNode
func (b *BoolNode) String() string {
	if b.True {
		return "true"
	}
	return "false"
}

func (b *BoolNode) tree() *Tree {
	return b.tr
}

// Copy returns a copy of the BoolNode
func (b *BoolNode) Copy() Node {
	return b.tr.newBool(b.Pos, b.True)
}

// NumberNode holds a number: signed or unsigned int or float.
// The value is parsed and stored under all the types that can represent the value.
type NumberNode struct {
	NodeType
	Pos
	tr      *Tree
	IsInt   bool
	IsUint  bool
	IsFloat bool
	Int64   int64
	Uint64  uint64
	Float64 float64
	Text    string
}

func (t *Tree) newNumber(pos Pos, text string) (*NumberNode, error) {
	n := &NumberNode{tr: t, NodeType: NodeNumber, Pos: pos, Text: text}

	u, err := strconv.ParseUint(text, 0, 64) // fails for -0; fixed below
	if err == nil {
		n.IsUint = true
		n.Uint64 = u
	}
	i, err := strconv.ParseInt(text, 0, 64)
	if err == nil {
		n.IsInt = true
		n.Int64 = i
		if i == 0 {
			n.IsUint = true
			n.Uint64 = u
		}
	}
	if n.IsInt {
		n.IsFloat = true
		n.Float64 = float64(n.Int64)
	} else if n.IsUint {
		n.IsFloat = true
		n.Float64 = float64(n.Uint64)
	} else {
		f, err := strconv.ParseFloat(text, 64)
		if err == nil {
			if !strings.ContainsAny(text, ".eE") {
				return nil, fmt.Errorf("integer overflow: %q", text)
			}
			n.IsFloat = true
			n.Float64 = f
			if !n.IsInt && float64(int64(f)) == f {
				n.IsInt = true
				n.Int64 = int64(f)
			}
			if !n.IsUint && float64(uint64(f)) == f {
				n.IsUint = true
				n.Uint64 = uint64(f)
			}
		}
	}
	if !n.IsInt && !n.IsUint && !n.IsFloat {
		return nil, fmt.Errorf("illegal number syntax: %q", text)
	}
	return n, nil
}

// String returns a string representation of the NumberNode
func (n *NumberNode) String() string {
	return n.Text
}

func (n *NumberNode) tree() *Tree {
	return n.tr
}

// Copy returns a copy of the NumberNode
func (n *NumberNode) Copy() Node {
	nn := new(NumberNode)
	*nn = *n
	return nn
}

// StringNode holds a string constant. The value has been "unquoted".
type StringNode struct {
	NodeType
	Pos
	tr     *Tree
	Quoted string // Original text of the string, with quotes
	Text   string // Unquoted string, after processing
}

func (t *Tree) newString(pos Pos, orig, text string) *StringNode {
	return &StringNode{tr: t, NodeType: NodeString, Pos: pos, Quoted: orig, Text: text}
}

// String returns the original quoted version of the StringNode
func (s *StringNode) String() string {
	return s.Quoted
}

func (s *StringNode) tree() *Tree {
	return s.tr
}

// Copy returns a copy of the StringNode
func (s *StringNode) Copy() Node {
	return s.tr.newString(s.Pos, s.Quoted, s.Text)
}

// OperatorNode contains an operator and its operands.
type OperatorNode struct {
	NodeType
	Pos
	Operator string
	tr       *Tree
	Operands *ListNode
}

func (t *Tree) newOperator(op string, pos Pos, list *ListNode) *OperatorNode {
	return &OperatorNode{tr: t, Operator: op, Pos: pos, Operands: list}
}

// String returns the string representation of the OperatorNode
func (o *OperatorNode) String() string {
	return fmt.Sprintf("%s%s", o.Operator, o.Operands)
}

func (o *OperatorNode) tree() *Tree {
	return o.tr
}

// Copy returns a copy of the OperatorNode
func (o *OperatorNode) Copy() Node {
	return o.tr.newOperator(o.Operator, o.Pos, o.Operands.CopyList())
}
