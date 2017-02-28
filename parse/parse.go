package rql

import (
	"fmt"
	"strconv"
	"strings"
)

// Tree is the representation of a single parsed statement
type Tree struct {
	Name string         // The name of the statement represented by the tree
	Root *StatementNode // top-level root of the tree
	text string         // The text to be parsed

	// Parsing only; cleared after parse.
	lex       *lexer
	token     [3]item // three-token lookahead for parser
	peekCount int
}

// Copy returns a copy of the Tree. Any parsing state is discarded.
func (t *Tree) Copy() *Tree {
	if t == nil {
		return nil
	}
	return &Tree{
		Name: t.Name,
		Root: t.Root.CopyStatement(),
		text: t.text,
	}
}

// Parse returns a map from statement name to parse.Tree, created by parsing the
// statements described in the argument string. The top-level statement will be
// given the specified name. If an error occurs, parsing stops and an empty map
// is returned with the error.
func Parse(name, text string) (*Tree, error) {
	t := New(name)
	t.text = text
	return t.Parse(text)
}

// next returns the next token.
func (t *Tree) next() item {
	if t.peekCount > 0 {
		t.peekCount--
	} else {
		t.token[0] = t.lex.nextItem()
	}
	return t.token[t.peekCount]
}

// backup backs the input stream up one token.
func (t *Tree) backup() {
	t.peekCount++
}

// peek returns but does not consume the next token
func (t *Tree) peek() item {
	if t.peekCount > 0 {
		return t.token[t.peekCount-1]
	}
	t.peekCount = 1
	t.token[0] = t.lex.nextItem()
	return t.token[0]
}

// nextNonSpace returns the next non-whitespace token
func (t *Tree) nextNonSpace() (token item) {
	for {
		token = t.next()
		if token.typ != itemWhitespace {
			break
		}
	}
	return token
}

// peekNonSpace returns but does not consume the next non-space token
func (t *Tree) peekNonSpace() (token item) {
	for {
		token = t.next()
		if token.typ != itemWhitespace {
			break
		}
	}
	t.backup()
	return token
}

// Parsing

// New allocates a new parse tree with the given name
func New(name string) *Tree {
	return &Tree{
		Name: name,
	}
}

// ErrorContext returns a textual representation of the location of the node in the input text.
func (t *Tree) ErrorContext(n Node) (location, context string) {
	pos := int(n.Position())
	tree := n.tree()
	if tree == nil {
		tree = t
	}
	text := tree.text[:pos]
	byteNum := strings.LastIndex(text, "\n")
	if byteNum == -1 {
		byteNum = pos
	} else {
		byteNum++
		byteNum = pos - byteNum
	}
	lineNum := 1 + strings.Count(text, "\n")
	context = n.String()
	if len(context) > 20 {
		context = fmt.Sprintf("%.20s...", context)
	}
	return fmt.Sprintf("%s:%d:%d", tree.Name, lineNum, byteNum), context
}

// errorf formats the error and terminates processing
func (t *Tree) errorf(format string, args ...interface{}) {
	t.Root = nil
	format = fmt.Sprintf("statement: %s:%d: %s", t.Name, t.token[0].line, format)
	panic(fmt.Errorf(format, args...))
}

// error terminates processing
func (t *Tree) error(err error) {
	t.errorf("%s", err)
}

// expect consumes the next token and guarantees it has the required type
func (t *Tree) expect(expected itemType, context string) item {
	token := t.nextNonSpace()
	if token.typ != expected {
		t.unexpected(token, context)
	}
	return token
}

// expectOneOf consumes the next token and guarantees it has one of the required types
func (t *Tree) expectOneOf(expected []itemType, context string) item {
	token := t.nextNonSpace()
	found := false
	for _, e := range expected {
		if token.typ == e {
			found = true
			break
		}
	}
	if !found {
		t.unexpected(token, context)
	}
	return token
}

// unexpected complains about the token and terminates processing
func (t *Tree) unexpected(token item, context string) {
	t.errorf("unexpected %s in %s", token, context)
}

// recover is the handler that turns panics into returns from the top level of Parse
func (t *Tree) recover(errp *error) {
	e := recover()
	if e != nil {
		if t != nil {
			t.lex.drain()
			t.stopParse()
		}
		*errp = e.(error)
	}
	return
}

// startParse intiializes the parser, using the lexer.
func (t *Tree) startParse(lex *lexer) {
	t.Root = nil
	t.lex = lex
}

// stopParse terminates parsing
func (t *Tree) stopParse() {
	t.lex = nil
}

// Parse parses the statement string to construct a representation of the statement for
// translation.
func (t *Tree) Parse(text string) (tree *Tree, err error) {
	defer t.recover(&err)
	t.startParse(lex(t.Name, text))
	t.text = text
	t.parse()
	t.stopParse()
	return t, nil
}

// IsEmptyTree reports whether this tree (node) is empty of everything but space
func IsEmptyTree(n Node) bool {
	switch n := n.(type) {
	case nil:
		return true
	case *ListNode:
		for _, node := range n.Nodes {
			if !IsEmptyTree(node) {
				return false
			}
		}
	case *IdentifierNode:
		return n.Ident == ""
	case *NullNode:
	case *BoolNode:
	case *NumberNode:
	case *StringNode:
		return false
	case *OperatorNode:
		return IsEmptyTree(n.Operands)
	case *StatementNode:
		if n.Operator == nil {
			return true
		}
		return IsEmptyTree(n.Operator.Operands)
	default:
		panic("unknown node: " + n.String())
	}
	return false
}

// parse is the top-level parser for a statement
// runs to EOF
func (t *Tree) parse() {
	t.Root = t.newStatement(t.peek().pos, nil)
	op := t.peekNonSpace()
	if itemOperatorsStart <= op.typ && op.typ <= itemOperatorsEnd {
		t.Root.Operator = t.operator()
	}
	tok := t.nextNonSpace()
	if tok.typ != itemEOF {
		t.errorf("unexpected token after operator: %q", tok)
	}
}

// operator returns an operator
func (t *Tree) operator() *OperatorNode {
	token := t.expectOneOf(operators, "operator")
	return t.newOperator(token.val, token.pos, t.list())
}

func (t *Tree) list() *ListNode {
	list := t.newList(t.expect(itemLeftParen, "left parentheses").pos)
	expectComma := false
Loop:
	for {
		if expectComma {
			tok := t.expectOneOf([]itemType{itemComma, itemRightParen}, "comma or right parentheses")
			if tok.typ == itemRightParen {
				t.backup()
			}
		}
		switch token := t.nextNonSpace(); {
		case token.typ == itemIdentifier:
			list.append(NewIdentifier(token.val).SetTree(t).SetPos(token.pos))
		case token.typ == itemString:
			s, err := strconv.Unquote(token.val)
			if err != nil {
				t.error(err)
			}
			list.append(t.newString(token.pos, token.val, s))
		case token.typ == itemBool:
			list.append(t.newBool(token.pos, token.val == "true"))
		case token.typ == itemNumber:
			number, err := t.newNumber(token.pos, token.val)
			if err != nil {
				t.error(err)
			}
			list.append(number)
		case token.typ == itemNull:
			list.append(t.newNull(token.pos))
		case itemOperatorsStart <= token.typ && token.typ <= itemOperatorsEnd:
			t.backup()
			op := t.operator()
			list.append(op)
		case token.typ == itemLeftParen:
			t.backup()
			l := t.list()
			list.append(l)
		case token.typ == itemRightParen:
			break Loop
		}
		expectComma = true
	}
	return list
}
