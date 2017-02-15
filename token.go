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
)

// [...] syntax creates an array of specific length
var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",

	IDENT:   "IDENT",
	NUMERIC: "NUMERIC",
	STRING:  "STRING",
	BOOLEAN: "BOOLEAN",

	AND: "and",
	OR:  "or",
	NOT: "not",
	LT:  "le",
	GT:  "gt",
	LE:  "le",
	GE:  "ge",
	EQ:  "eq",
	NE:  "ne",
	IN:  "in",

	LPAREN: "(",
	RPAREN: ")",
	COMMA:  ",",
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

func LookupOperator(s string) Token {
	for i := operator_beg + 1; i < operator_end; i++ {
		if s == tokens[i] {
			return i
		}
	}
	return 0
}