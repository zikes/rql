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

func (s *Scanner) read(skipWhitespace bool) rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	if skipWhitespace && unicode.IsSpace(ch) {
		return s.read(skipWhitespace)
	}
	return ch
}

func (s *Scanner) unread() { s.r.UnreadRune() }

func (s *Scanner) Scan() Expression {
	c := s.read(false)

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
		return nil
	case rune(tokens[COMMA][0]):
		return &Punctuation{COMMA, tokens[COMMA]}
	case rune(tokens[LPAREN][0]):
		return &Punctuation{LPAREN, tokens[LPAREN]}
	case rune(tokens[RPAREN][0]):
		return &Punctuation{RPAREN, tokens[RPAREN]}
	}

	return &Illegal{ILLEGAL, string(c)}
}

func (s *Scanner) scanWhitespace() *Whitespace {
	var buf bytes.Buffer

	for {
		if c := s.read(false); c == eof {
			break
		} else if !unicode.IsSpace(c) {
			s.unread()
			break
		} else {
			buf.WriteRune(c)
		}
	}

	return &Whitespace{WS, buf.String()}
}
func (s *Scanner) scanIdentifier() Expression {
	var buf bytes.Buffer

	for {
		if c := s.read(true); c == eof {
			break
		} else if !unicode.IsLetter(c) && c != '.' && c != '_' && !unicode.IsDigit(c) {
			s.unread()
			break
		} else {
			buf.WriteRune(c)
		}
	}

	str := buf.String()

	opToken := LookupOperator(strings.ToLower(str))
	if opToken > 0 {
		return &Operator{Kind: opToken}
	}

	switch strings.ToLower(str) {
	case "true", "false":
		return &Literal{Kind: BOOLEAN, Value: strings.ToLower(str)}
	case "null":
		return &Literal{Kind: NULL, Value: strings.ToLower(str)}
	}

	return &Identifier{IDENT, str}
}

func (s *Scanner) scanStringLiteral() *Literal {
	var buf bytes.Buffer
	delim := s.read(false)
	buf.WriteRune(delim)

	skip := false

	for {
		if c := s.read(false); c == eof {
			break
		} else if skip {
			skip = false
			buf.WriteRune(c)
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

	d := string(delim)

	return &Literal{STRING, strings.Trim(strings.Replace(buf.String(), `\`+d, d, -1), d)}
}
func (s *Scanner) scanNumericLiteral() *Literal {
	var buf bytes.Buffer

	for {
		if c := s.read(true); c == eof {
			break
		} else if !unicode.IsDigit(c) && c != '-' && c != '.' {
			s.unread()
			break
		} else {
			buf.WriteRune(c)
		}
	}

	return &Literal{NUMERIC, buf.String()}
}
