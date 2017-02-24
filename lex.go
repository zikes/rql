package rql

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type item struct {
	typ  itemType
	pos  Pos
	val  string
	line int
}

func (i item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.val
	case i.typ > itemKeyword:
		return fmt.Sprintf("<%s>", i.val)
	case len(i.val) > 10:
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

type itemType int

const (
	itemError itemType = iota // error occurred, value is text of error

	itemEOF

	itemIdentifier // alphanumeric identifier
	itemString     // raw quoted string (includes quotes)
	itemBool       // boolean constant
	itemNumber     // numeric value
	itemLeftParen  // '('
	itemRightParen // ')'
	itemComma      // ','
	itemWhitespace // white space separating arguments

	itemKeyword // used only to delimit keywords

	itemOperatorsStart // used only to delimit operators
	itemAnd            // and keyword
	itemOr             // or keyword
	itemEq             // eq keyword
	itemNe             // ne keyword
	itemLt             // lt keyword
	itemGt             // gt keyword
	itemLe             // le keyword
	itemGe             // ge keyword
	itemIn             // in keyword
	itemOperatorsEnd   // used only to delimit operators

	itemNull // null keyword
)

var key = map[string]itemType{
	"and":  itemAnd,
	"or":   itemOr,
	"eq":   itemEq,
	"ne":   itemNe,
	"lt":   itemLt,
	"gt":   itemGt,
	"le":   itemLe,
	"ge":   itemGe,
	"in":   itemIn,
	"null": itemNull,
}

var operators = []itemType{
	itemAnd,
	itemOr,
	itemEq,
	itemNe,
	itemLt,
	itemGt,
	itemLe,
	itemGe,
	itemIn,
}

const eof = -1

type stateFn func(*lexer) stateFn

type lexer struct {
	name    string
	input   string
	state   stateFn
	pos     Pos
	start   Pos
	width   Pos
	lastPos Pos
	items   chan item
	depth   int
	line    int
}

func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width
	if r == '\n' {
		l.line++
	}
	return r
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
	if l.width == 1 && l.input[l.pos] == '\n' {
		l.line--
	}
}

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.start, l.input[l.start:l.pos], l.line}
	switch t {
	case itemWhitespace, itemString:
		l.line += strings.Count(l.input[l.start:l.pos], "\n")
	}
	l.start = l.pos
}

func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, l.start, fmt.Sprintf(format, args...), l.line}
	return nil
}

func (l *lexer) nextItem() item {
	item := <-l.items
	l.lastPos = item.pos
	return item
}

func (l *lexer) drain() {
	for range l.items {
	}
}

func lex(name, input string) *lexer {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item),
		line:  1,
	}
	go l.run()
	return l
}

func (l *lexer) run() {
	for l.state = lexStatement; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.items)
}

// lexStatement scans until it finds an identifier
func lexStatement(l *lexer) stateFn {
	switch r := l.next(); {
	case r == eof:
		if l.depth > 0 {
			return l.errorf("unexpected end of statement")
		}
		l.emit(itemEOF)
		return nil
	case isWhitespace(r):
		return lexWhitespace
	case r == '"':
		return lexString
	case r == '.' || r == '+' || r == '-' || ('0' <= r && r <= '9'):
		return lexNumber
	case r == '(':
		l.emit(itemLeftParen)
		l.depth++
	case r == ')':
		l.emit(itemRightParen)
		l.depth--
		if l.depth < 0 {
			return l.errorf("unexpected right paren %#U", r)
		}
	case r == ',':
		l.emit(itemComma)
	case isAlphaNumeric(r):
		l.backup()
		return lexIdentifier
	default:
		return l.errorf("unrecognized character in statement: %#U", r)
	}
	return lexStatement
}

// lexWhitespace scans until it runs out of whitespace
func lexWhitespace(l *lexer) stateFn {
	for isWhitespace(l.peek()) {
		l.next()
	}
	l.emit(itemWhitespace)
	return lexStatement
}

// lexIdentifier scans an identifier
func lexIdentifier(l *lexer) stateFn {
Loop:
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
			// absorb
		default:
			l.backup()
			word := l.input[l.start:l.pos]
			if !l.atTerminator() {
				return l.errorf("bad character %#U", r)
			}

			switch {
			case key[word] > itemKeyword:
				l.emit(key[word])
			case word == "true", word == "false":
				l.emit(itemBool)
			default:
				l.emit(itemIdentifier)
			}
			break Loop
		}
	}
	return lexStatement
}

// lexNumber scans a numeric value
func lexNumber(l *lexer) stateFn {
	if !l.scanNumber() {
		return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
	}
	l.emit(itemNumber)
	return lexStatement
}

func (l *lexer) scanNumber() bool {
	// optional leading sign
	l.accept("+-")
	digits := "0123456789"
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	// next thing mustn't be alphanumeric
	if isAlphaNumeric(l.peek()) {
		l.next()
		return false
	}
	return true
}

// lexString scans a quoted string
func lexString(l *lexer) stateFn {
Loop:
	for {
		switch l.next() {
		case '\\':
			if r := l.next(); r != eof && r != '\n' {
				break
			}
			fallthrough
		case eof, '\n':
			return l.errorf("unterminated quoted string")
		case '"':
			break Loop
		}
	}
	l.emit(itemString)
	return lexStatement
}

func (l *lexer) atTerminator() bool {
	r := l.peek()
	if isWhitespace(r) {
		return true
	}
	switch r {
	case eof, '.', ',', ')', '(':
		return true
	}
	return false
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n'
}

func isAlphaNumeric(r rune) bool {
	return r == '_' || r == '.' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
