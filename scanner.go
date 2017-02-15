package rql

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"unicode"
)

var eof = rune(0)

type Scanner struct {
	r *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) unread() { s.r.UnreadRune() }

func (s *Scanner) Scan() (Token, Expression) {
	c := s.read()

	if unicode.IsSpace(c) {
		s.unread()
		return s.scanWhitespace()
	}

	if unicode.IsLetter(c) || c == '_' {
		s.unread()
		return s.scanIdentifier()
	}

	if unicode.In(c, unicode.Quotation_Mark) {
		s.unread()
		return s.scanStringLiteral()
	}

	if unicode.IsDigit(c) || c == '-' {
		s.unread()
		return s.scanNumericLiteral()
	}

	switch c {
	case eof:
		return EOF, nil
	case rune(tokens[COMMA][0]):
		return COMMA, &Punctuation{COMMA, tokens[COMMA]}
	case rune(tokens[LPAREN][0]):
		return LPAREN, &Punctuation{LPAREN, tokens[LPAREN]}
	case rune(tokens[RPAREN][0]):
		return RPAREN, &Punctuation{RPAREN, tokens[RPAREN]}
	}

	return ILLEGAL, &Identifier{string(c)}
}

func (s *Scanner) scanWhitespace() (Token, *Whitespace) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if c := s.read(); c == eof {
			break
		} else if !unicode.IsSpace(c) {
			s.unread()
			break
		} else {
			buf.WriteRune(c)
		}
	}

	return WS, &Whitespace{buf.String()}
}
func (s *Scanner) scanIdentifier() (Token, Expression) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if c := s.read(); c == eof {
			break
		} else if !unicode.IsLetter(c) && c != '_' && !unicode.IsDigit(c) {
			s.unread()
			break
		} else {
			buf.WriteRune(c)
		}
	}

	str := buf.String()

	opToken := LookupOperator(strings.ToLower(str))
	if opToken > 0 {
		if ch := s.read(); ch != rune(tokens[LPAREN][0]) {
			s.unread()
			return ILLEGAL, &Identifier{string(ch)}
		}
		p := s.scanParenExpression()
		o := &Operator{
			Name:     str,
			Operands: p,
		}
		return opToken, o
	}

	switch strings.ToLower(str) {
	case "true", "false":
		return BOOLEAN, &Literal{Kind: BOOLEAN, Value: strings.ToLower(str)}
	}

	return IDENT, &Identifier{str}
}
func (s *Scanner) scanParenExpression() *ParenExpr {
	p := &ParenExpr{}
	for {
		tok, exp := s.Scan()
		if tok == COMMA || tok == WS {
			continue
		}
		if tok == RPAREN {
			break
		}
		p.Expressions = append(p.Expressions, exp)
	}
	return p
}

func (s *Scanner) scanStringLiteral() (Token, *Literal) {
	var buf bytes.Buffer
	delim := s.read()
	buf.WriteRune(delim)

	skip := false

	for {
		if c := s.read(); c == eof {
			break
		} else if skip {
			skip = false
			continue
		} else if c == '\\' {
			skip = true
			continue
		} else if c == delim {
			buf.WriteRune(c)
			break
		} else {
			buf.WriteRune(c)
		}
	}

	return STRING, &Literal{STRING, buf.String()}
}
func (s *Scanner) scanNumericLiteral() (Token, *Literal) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if c := s.read(); c == eof {
			break
		} else if !unicode.IsDigit(c) && c != '-' && c != '.' {
			s.unread()
			break
		} else {
			buf.WriteRune(c)
		}
	}

	return NUMERIC, &Literal{NUMERIC, buf.String()}
}
