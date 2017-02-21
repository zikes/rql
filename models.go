// https://github.com/jlyonsmith/Rql/wiki/Overview
// https://github.com/persvr/rql
// http://jeremy.marzhillstudios.com/entries/Using-bufioScanner-to-build-a-Tokenizer/
package rql

import (
	"fmt"
	"strconv"
	"strings"
)

type Expression interface {
	exprNode()
	Token() Token
	fmt.Stringer
}
type ExpressionList []Expression

func (e ExpressionList) exprNode() {}
func (e ExpressionList) String() string {
	out := []string{}
	for _, v := range e {
		if _, ok := v.(Whitespace); ok {
			continue
		}
		out = append(out, v.String())
	}
	return "(" + strings.Join(out, ",") + ")"
}
func (e ExpressionList) Token() Token {
	return EXPR_LIST
}

type Operator struct {
	Kind     Token // AND, OR, NOT, LT, GT, LE, GE, EQ, NE, IN
	Operands ExpressionList
}

func (Operator) exprNode() {}
func (o Operator) String() string {
	return o.Kind.String() + o.Operands.String()
}
func (o Operator) Token() Token {
	return o.Kind
}

type Identifier struct {
	Kind Token // IDENT
	Name string
}

func (Identifier) exprNode() {}
func (i Identifier) String() string {
	return i.Name
}
func (i Identifier) Token() Token {
	return i.Kind
}

type Literal struct {
	Kind  Token // NUMERIC, STRING, BOOLEAN, NULL
	Value string
}

func (Literal) exprNode() {}
func (l Literal) String() string {
	if l.Kind == STRING {
		return `"` + strings.Replace(l.Value, `"`, `\"`, -1) + `"`
	}
	if l.Kind == NULL {
		return "null"
	}
	return l.Value
}
func (l Literal) Token() Token {
	return l.Kind
}
func (l Literal) IntValue() int {
	i, err := strconv.Atoi(l.Value)
	if err != nil {
		return 0
	}
	return i
}
func (l Literal) StringValue() string {
	return l.Value
}
func (l Literal) BoolValue() bool {
	return l.Value == "true"
}
func (l Literal) FloatValue() float64 {
	f, err := strconv.ParseFloat(l.Value, 64)
	if err != nil {
		return 0
	}
	return f
}

type Punctuation struct {
	Kind  Token  // COMMA, LPAREN, RPAREN
	Value string // ",", "(", ")"
}

func (Punctuation) exprNode() {}
func (p Punctuation) String() string {
	return p.Value
}
func (p Punctuation) Token() Token {
	return p.Kind
}

type Whitespace struct {
	Kind  Token // WS
	Value string
}

func (Whitespace) exprNode() {}
func (w Whitespace) String() string {
	return w.Value
}
func (w Whitespace) Token() Token {
	return w.Kind
}

type Statement ExpressionList

func (s Statement) String() string {
	out := []string{}
	for _, v := range s {
		if _, ok := v.(Whitespace); ok {
			continue
		}
		out = append(out, v.String())
	}
	return strings.Join(out, ",")
}
