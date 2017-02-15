package rql

import (
	"fmt"
	"io"
)

type Parser struct {
	s   *Scanner
	buf struct {
		tok Token      // last read token
		lit Expression // last read expression
		n   int        // buffer size (max = 1)
	}
}

func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

func (p *Parser) scan() (Token, Expression) {
	// If there's a token on the buffer, return it
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Read the next token from the scanner
	p.buf.tok, p.buf.lit = p.s.Scan()

	return p.buf.tok, p.buf.lit
}

func (p *Parser) unscan() { p.buf.n = 1 }

func (p *Parser) scanIgnoreWhitespace() (Token, Expression) {
	tok, lit := p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return tok, lit
}

func (p *Parser) Parse() (*Statement, error) {
	stmt := &Statement{}
	for {
		tok, exp := p.scanIgnoreWhitespace()
		if tok == EOF {
			break
		}

		if !tok.IsOperator() {
			return nil, fmt.Errorf("found %s:%q, expected operator", tok, exp)
		}

		stmt.Expressions = append(stmt.Expressions, exp)
	}
	return stmt, nil
}
