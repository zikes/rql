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

func Array(e ...interface{}) (ExpressionList, error) {
	out := ExpressionList{}
	for _, v := range e {
		switch v := v.(type) {
		case nil:
			out = append(out, Literal{Kind: NULL, Value: "null"})
		case string:
			out = append(out, Literal{Kind: STRING, Value: v})
		case int:
			out = append(out, Literal{Kind: NUMERIC, Value: strconv.Itoa(v)})
		case float64:
			out = append(out, Literal{Kind: NUMERIC, Value: strconv.FormatFloat(v, 'f', -1, 64)})
		case Expression:
			out = append(out, v)
		default:
			// unknown type
			return out, fmt.Errorf("Unknown type: %v", v)
		}
	}
	return out, nil
}
