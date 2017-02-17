package rql

import "fmt"

type Parser struct {
	s   *Scanner
	buf struct {
		lit Expression // last read expression
		n   int        // buffer size (max = 1)
	}
}

func (p *Parser) scan() Expression {
	// If there's a token on the buffer, return it
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.lit
	}

	// Read the next token from the scanner
	p.buf.lit = p.s.Scan()

	return p.buf.lit
}

func (p *Parser) unscan() { p.buf.n = 1 }

func (p *Parser) Parse() (Statement, error) {
	p.skipWhitespace()
	p.expectOperator()

	op, err := p.getOperator()
	if err != nil {
		return nil, err
	}

	return Statement{op}, nil
}

func (p *Parser) getOperator() (*Operator, error) {
	for {
		exp := p.scan()
		if exp.Token().IsEOF() {
			return nil, fmt.Errorf("unexpected EOF")
		}
		if exp.Token().IsWhitespace() {
			continue
		}
		if exp.Token().IsIllegal() {
			return nil, fmt.Errorf("unrecognized expression: %v", exp)
		}
		if !exp.Token().IsOperator() {
			return nil, fmt.Errorf("found %v, expected operator", exp)
		}
		op, ok := exp.(*Operator)
		if !ok {
			return nil, fmt.Errorf("unexpected error parsing operator")
		}
		operands, err := p.getExpressionList()
		if err != nil {
			return nil, fmt.Errorf("error parsing expression list: %s", err)
		}
		op.Operands = *operands
		return op, nil
	}
}

func (p *Parser) getExpressionList() (*ExpressionList, error) {
	p.skipWhitespace()

	if err := p.expectLParen(); err != nil {
		return nil, err
	}

	list := &ExpressionList{}

	for {
		if err := p.expectValueExpression(); err != nil {
			return nil, err
		}
		value, err := p.getValue()
		if err != nil {
			return nil, err
		}
		*list = append(*list, value)
		if err := p.expectCommaOrRParen(); err != nil {
			return nil, err
		}
		commaOrRParen := p.scan()
		if commaOrRParen.Token() == RPAREN {
			break
		}
	}
	return list, nil
}

func (p *Parser) getValue() (Expression, error) {
	val := p.scan()
	if !val.Token().IsValue() {
		return nil, fmt.Errorf("found %v, expected value", val)
	}
	if val.Token() == LPAREN {
		p.unscan()
		return p.getExpressionList()
	}
	if val.Token().IsOperator() {
		p.unscan()
		return p.getOperator()
	}
	if val.Token().IsLiteral() || val.Token().IsIdentifier() {
		return val, nil
	}
	return nil, fmt.Errorf("found %v, expected value", val)
}

func (p *Parser) expectLParen() error {
	p.skipWhitespace()
	exp := p.scan()
	if exp.Token() != LPAREN {
		return fmt.Errorf("found %v, expected open parentheses", exp)
	}
	return nil
}

func (p *Parser) expectValueExpression() error {
	p.skipWhitespace()
	exp := p.scan()
	if !exp.Token().IsValue() {
		return fmt.Errorf("found %v, expected value expression", exp)
	}
	p.unscan()
	return nil
}

func (p *Parser) expectCommaOrRParen() error {
	p.skipWhitespace()
	exp := p.scan()
	if exp.Token() != COMMA && exp.Token() != RPAREN {
		return fmt.Errorf("found %v, expected comma or close parentheses", exp)
	}
	p.unscan()
	return nil
}

func (p *Parser) expectOperator() error {
	p.skipWhitespace()
	exp := p.scan()
	if !exp.Token().IsOperator() {
		return fmt.Errorf("found %v, expected operator", exp)
	}
	p.unscan()
	return nil
}

func (p *Parser) skipWhitespace() {
	for {
		exp := p.scan()
		if !exp.Token().IsWhitespace() {
			p.unscan()
			return
		}
	}
}
