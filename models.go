// https://github.com/jlyonsmith/Rql/wiki/Overview
// https://github.com/persvr/rql
// http://jeremy.marzhillstudios.com/entries/Using-bufioScanner-to-build-a-Tokenizer/
package rql

import "strings"

type Expression interface {
	exprNode()
	String() string
}
type ExpressionList []Expression

func (e *ExpressionList) String() []string {
	out := []string{}
	for _, v := range *e {
		out = append(out, v.String())
	}
	return out
}

type Operator struct {
	Name     string
	Operands *ParenExpr
}

func (*Operator) exprNode() {}
func (o *Operator) String() string {
	return o.Name + o.Operands.String()
}

type ParenExpr struct {
	Expressions ExpressionList // parenthesized expression
}

func (*ParenExpr) exprNode() {}
func (p *ParenExpr) String() string {
	return "(" + strings.Join(p.Expressions.String(), ",") + ")"
}

type Identifier struct {
	Name string
}

func (*Identifier) exprNode() {}
func (i *Identifier) String() string {
	return i.Name
}

type Literal struct {
	Kind  Token // NUMERIC, STRING, BOOLEAN
	Value string
}

func (*Literal) exprNode() {}
func (l *Literal) String() string {
	return l.Value
}

type Punctuation struct {
	Kind  Token  // COMMA, LPAREN, RPAREN
	Value string // ",", "(", ")"
}

func (*Punctuation) exprNode() {}
func (p *Punctuation) String() string {
	return p.Value
}

type Whitespace struct {
	Value string
}

func (*Whitespace) exprNode() {}
func (w *Whitespace) String() string {
	return w.Value
}

type Statement struct {
	Expressions ExpressionList
}

func (s *Statement) String() string {
	return strings.Join(s.Expressions.String(), "")
}
