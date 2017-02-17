package rql

import (
	"io"
	"strings"
)

func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

func ParseString(s string) (Statement, error) {
	return NewParser(strings.NewReader(s)).Parse()
}
