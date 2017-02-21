package rql

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

func ParseString(s string) (Statement, error) {
	return NewParser(strings.NewReader(s)).Parse()
}

// Literals

func Ident(s string) (Expression, error) {
	if s == "" {
		return nil, fmt.Errorf("Identifier must be non-empty")
	}
	return Identifier{
		Kind: IDENT,
		Name: s,
	}, nil
}
func Lit(i interface{}) Expression {
	switch i := i.(type) {
	case int:
		return Literal{Kind: NUMERIC, Value: strconv.Itoa(i)}
	case float64:
		return Literal{Kind: NUMERIC, Value: strconv.FormatFloat(i, 'f', -1, 64)}
	case string:
		return Literal{Kind: STRING, Value: i}
	case nil:
		return Literal{Kind: NULL, Value: ""}
	case bool:
		s := "false"
		if i {
			s = "true"
		}
		return Literal{Kind: BOOLEAN, Value: s}
	default:
		return Illegal{Kind: ILLEGAL, Value: ""}
	}
}

// Operators

func And(e ...Expression) Expression {
	return Operator{
		Kind:     AND,
		Operands: ExpressionList(e),
	}
}
func Or(e ...Expression) Expression {
	return Operator{
		Kind:     OR,
		Operands: ExpressionList(e),
	}
}
func Not(e ...Expression) Expression {
	return Operator{
		Kind:     NOT,
		Operands: ExpressionList(e),
	}
}
func Lt(e ...Expression) Expression {
	return Operator{
		Kind:     LT,
		Operands: ExpressionList(e),
	}
}
func Gt(e ...Expression) Expression {
	return Operator{
		Kind:     GT,
		Operands: ExpressionList(e),
	}
}
func Le(e ...Expression) Expression {
	return Operator{
		Kind:     LE,
		Operands: ExpressionList(e),
	}
}
func Ge(e ...Expression) Expression {
	return Operator{
		Kind:     GE,
		Operands: ExpressionList(e),
	}
}
func Eq(e ...Expression) Expression {
	return Operator{
		Kind:     EQ,
		Operands: ExpressionList(e),
	}
}
func Ne(e ...Expression) Expression {
	return Operator{
		Kind:     NE,
		Operands: ExpressionList(e),
	}
}
func In(e ...Expression) Expression {
	return Operator{
		Kind:     IN,
		Operands: ExpressionList(e),
	}
}

// Arrays

func Array(e ...interface{}) ExpressionList {
	out := ExpressionList{}

	for _, v := range e {
		lit := Lit(v)
		if lit.Token() == ILLEGAL {
			out = append(out, lit)
			continue
		}
		switch v := v.(type) {
		case Expression:
			out = append(out, v)
		default:
			// unknown type
			out = append(out, lit)
		}
	}
	return out
}
