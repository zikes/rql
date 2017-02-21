package rql

import "strconv"

type Token int

// List of tokens
const (
	ILLEGAL Token = iota
	EOF
	WS

	literal_beg
	IDENT   // main
	NUMERIC // int, float
	STRING  // "string"
	BOOLEAN // true, false
	NULL    // null
	literal_end

	operator_beg
	AND // &&
	OR  // ||
	NOT // !
	LT  // <
	GT  // >
	LE  // <=
	GE  // >=
	EQ  // ==
	NE  // !=
	IN  // in()
	operator_end

	punctuation_beg
	LPAREN // (
	RPAREN // )
	COMMA  // ,
	punctuation_end

	array_beg
	EXPR_LIST
	array_end
)

// [...] syntax creates an array of specific length
var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	WS:      "WS",

	IDENT:   "IDENT",
	NUMERIC: "NUMERIC",
	STRING:  "STRING",
	BOOLEAN: "BOOLEAN",
	NULL:    "NULL",

	AND: "and",
	OR:  "or",
	NOT: "not",
	LT:  "lt",
	GT:  "gt",
	LE:  "le",
	GE:  "ge",
	EQ:  "eq",
	NE:  "ne",
	IN:  "in",

	LPAREN: "(",
	RPAREN: ")",
	COMMA:  ",",

	EXPR_LIST: "EXPR_LIST",
}

func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

func (tok Token) IsLiteral() bool     { return literal_beg < tok && tok < literal_end }
func (tok Token) IsOperator() bool    { return operator_beg < tok && tok < operator_end }
func (tok Token) IsPunctuation() bool { return punctuation_beg < tok && tok < punctuation_end }
func (tok Token) IsArray() bool       { return array_beg < tok && tok < array_end }
func (tok Token) IsWhitespace() bool  { return tok == WS }
func (tok Token) IsEOF() bool         { return tok == EOF }
func (tok Token) IsIllegal() bool     { return tok == ILLEGAL }
func (tok Token) IsIdentifier() bool  { return tok == IDENT }
func (tok Token) IsValue() bool {
	return tok.IsIdentifier() || tok.IsLiteral() || tok.IsOperator() || tok.IsArray() || tok == LPAREN
}

func LookupOperator(s string) Token {
	for i := operator_beg + 1; i < operator_end; i++ {
		if s == tokens[i] {
			return i
		}
	}
	return 0
}
