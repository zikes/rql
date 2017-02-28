# rql
--
    import "git.nwaonline.com/rune/rql"


## Usage

#### func  IsEmptyTree

```go
func IsEmptyTree(n Node) bool
```
IsEmptyTree reports whether this tree (node) is empty of everything but space

#### type BoolNode

```go
type BoolNode struct {
	NodeType
	Pos

	True bool // the value
}
```

BoolNode holds a boolean constant.

#### func (*BoolNode) Copy

```go
func (b *BoolNode) Copy() Node
```
Copy returns a copy of the BoolNode

#### func (*BoolNode) String

```go
func (b *BoolNode) String() string
```
String returns a string representation of the BoolNode

#### type IdentifierNode

```go
type IdentifierNode struct {
	NodeType
	Pos

	Ident string // The identifier's name.
}
```

IdentifierNode holds an identifier.

#### func  NewIdentifier

```go
func NewIdentifier(ident string) *IdentifierNode
```
NewIdentifier creates an IdentifierNode

#### func (*IdentifierNode) Copy

```go
func (i *IdentifierNode) Copy() Node
```
Copy copies the IdentifierNode

#### func (*IdentifierNode) SetPos

```go
func (i *IdentifierNode) SetPos(pos Pos) *IdentifierNode
```
SetPos sets the position. Chained for convenience.

#### func (*IdentifierNode) SetTree

```go
func (i *IdentifierNode) SetTree(t *Tree) *IdentifierNode
```
SetTree sets the tree. Chained for convenience.

#### func (*IdentifierNode) String

```go
func (i *IdentifierNode) String() string
```

#### type ListNode

```go
type ListNode struct {
	NodeType
	Pos

	Nodes []Node // The element nodes in lexical order
}
```

ListNode holds a sequence of Nodes

#### func (*ListNode) Copy

```go
func (l *ListNode) Copy() Node
```
Copy runs CopyList but returns as Node

#### func (*ListNode) CopyList

```go
func (l *ListNode) CopyList() *ListNode
```
CopyList will copy and return *ListNode

#### func (*ListNode) String

```go
func (l *ListNode) String() string
```
String returns the ListNode as a string

#### type Node

```go
type Node interface {
	Type() NodeType
	String() string

	// Copy does a deep copy of the Node and all its
	// components. To avoid type assertions, some Nodes
	// have specialized copy methods.
	Copy() Node

	// byte position of the start of the node in full original input string
	Position() Pos
	// contains filtered or unexported methods
}
```

A Node is an element in the parse tree. The interface contains an unexported
method so that only types local to this package can satisfy it.

#### type NodeType

```go
type NodeType int
```

NodeType identifies the type of the parse tree node

```go
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
```
NodeType constants

#### func (NodeType) Type

```go
func (t NodeType) Type() NodeType
```
Type returns itself and provides an easy default implementation for embedding in
a Node. Embedded in all non-trivial Nodes.

#### type NullNode

```go
type NullNode struct {
	NodeType
	Pos
}
```

NullNode holds the special identifier 'null'

#### func (*NullNode) Copy

```go
func (n *NullNode) Copy() Node
```
Copy returns a copy of the NullNode

#### func (*NullNode) String

```go
func (n *NullNode) String() string
```
String returns the string representation of NullNode

#### func (*NullNode) Type

```go
func (n *NullNode) Type() NodeType
```
Type returns the NodeType value

#### type NumberNode

```go
type NumberNode struct {
	NodeType
	Pos

	IsInt   bool
	IsUint  bool
	IsFloat bool
	Int64   int64
	Uint64  uint64
	Float64 float64
	Text    string
}
```

NumberNode holds a number: signed or unsigned int or float. The value is parsed
and stored under all the types that can represent the value.

#### func (*NumberNode) Copy

```go
func (n *NumberNode) Copy() Node
```
Copy returns a copy of the NumberNode

#### func (*NumberNode) String

```go
func (n *NumberNode) String() string
```
String returns a string representation of the NumberNode

#### type OperatorNode

```go
type OperatorNode struct {
	NodeType
	Pos
	Operator string

	Operands *ListNode
}
```

OperatorNode contains an operator and its operands.

#### func (*OperatorNode) Copy

```go
func (o *OperatorNode) Copy() Node
```
Copy returns a copy of the OperatorNode

#### func (*OperatorNode) String

```go
func (o *OperatorNode) String() string
```
String returns the string representation of the OperatorNode

#### type Pos

```go
type Pos int
```

Pos represents a byte position in the original input text from which this
statement was parsed

#### func (Pos) Position

```go
func (p Pos) Position() Pos
```
Position returns itself

#### type StatementNode

```go
type StatementNode struct {
	NodeType
	Pos
	Operator *OperatorNode
}
```

StatementNode holds a statement.

#### func (*StatementNode) Copy

```go
func (s *StatementNode) Copy() Node
```
Copy runs CopyStatement, returning as a Node

#### func (*StatementNode) CopyStatement

```go
func (s *StatementNode) CopyStatement() *StatementNode
```
CopyStatement returns a copy of the StatementNode as a *StatementNode

#### func (*StatementNode) String

```go
func (s *StatementNode) String() string
```
String returns the StatementNode as a string

#### type StringNode

```go
type StringNode struct {
	NodeType
	Pos

	Quoted string // Original text of the string, with quotes
	Text   string // Unquoted string, after processing
}
```

StringNode holds a string constant. The value has been "unquoted".

#### func (*StringNode) Copy

```go
func (s *StringNode) Copy() Node
```
Copy returns a copy of the StringNode

#### func (*StringNode) String

```go
func (s *StringNode) String() string
```
String returns the original quoted version of the StringNode

#### type Tree

```go
type Tree struct {
	Name string         // The name of the statement represented by the tree
	Root *StatementNode // top-level root of the tree
}
```

Tree is the representation of a single parsed statement

#### func  New

```go
func New(name string) *Tree
```
New allocates a new parse tree with the given name

#### func  Parse

```go
func Parse(name, text string) (*Tree, error)
```
Parse returns a map from statement name to parse.Tree, created by parsing the
statements described in the argument string. The top-level statement will be
given the specified name. If an error occurs, parsing stops and an empty map is
returned with the error.

#### func (*Tree) Copy

```go
func (t *Tree) Copy() *Tree
```
Copy returns a copy of the Tree. Any parsing state is discarded.

#### func (*Tree) ErrorContext

```go
func (t *Tree) ErrorContext(n Node) (location, context string)
```
ErrorContext returns a textual representation of the location of the node in the
input text.

#### func (*Tree) Parse

```go
func (t *Tree) Parse(text string) (tree *Tree, err error)
```
Parse parses the statement string to construct a representation of the statement
for translation.
