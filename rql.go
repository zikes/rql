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
func Lit(i interface{}) (Expression, error) {
	switch i := i.(type) {
	case int:
		return Literal{Kind: NUMERIC, Value: strconv.Itoa(i)}, nil
	case float64:
		return Literal{Kind: NUMERIC, Value: strconv.FormatFloat(i, 'f', -1, 64)}, nil
	case string:
		return Literal{Kind: STRING, Value: i}, nil
	case nil:
		return Literal{Kind: NULL, Value: "null"}, nil
	case bool:
		s := "false"
		if i {
			s = "true"
		}
		return Literal{Kind: BOOLEAN, Value: s}, nil
	default:
		return nil, fmt.Errorf("unknown type")
	}
}

// Operators

func And(e ...Expression) (Expression, error) {
	return Operator{
		Kind:     AND,
		Operands: ExpressionList(e),
	}, nil
}
func Or(e ...Expression) (Expression, error) {
	return Operator{
		Kind:     OR,
		Operands: ExpressionList(e),
	}, nil
}
func Not(e ...Expression) (Expression, error) {
	return Operator{
		Kind:     NOT,
		Operands: ExpressionList(e),
	}, nil
}
func Lt(e ...Expression) (Expression, error) {
	return Operator{
		Kind:     LT,
		Operands: ExpressionList(e),
	}, nil
}
func Gt(e ...Expression) (Expression, error) {
	return Operator{
		Kind:     GT,
		Operands: ExpressionList(e),
	}, nil
}
func Le(e ...Expression) (Expression, error) {
	return Operator{
		Kind:     LE,
		Operands: ExpressionList(e),
	}, nil
}
func Ge(e ...Expression) (Expression, error) {
	return Operator{
		Kind:     GE,
		Operands: ExpressionList(e),
	}, nil
}
func Eq(e ...Expression) (Expression, error) {
	return Operator{
		Kind:     EQ,
		Operands: ExpressionList(e),
	}, nil
}
func Ne(e ...Expression) (Expression, error) {
	return Operator{
		Kind:     NE,
		Operands: ExpressionList(e),
	}, nil
}
func In(e ...Expression) (Expression, error) {
	return Operator{
		Kind:     IN,
		Operands: ExpressionList(e),
	}, nil
}

// Arrays

func Array(e ...interface{}) (ExpressionList, error) {
	out := ExpressionList{}

	for _, v := range e {
		lit, err := Lit(v)
		if err == nil {
			out = append(out, lit)
			continue
		}
		switch v := v.(type) {
		case Expression:
			out = append(out, v)
		default:
			// unknown type
			return out, fmt.Errorf("Unknown type: %v", v)
		}
	}
	return out, nil
}
