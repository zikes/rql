package rql

import (
	"fmt"
	"strings"
	"testing"
)

// Make the types prettyprint
var itemName = map[itemType]string{
	itemError: "error",
	itemEOF:   "EOF",

	itemIdentifier: "identifier",
	itemString:     "string",
	itemBool:       "boolean",
	itemNumber:     "number",
	itemLeftParen:  "(",
	itemRightParen: ")",
	itemWhitespace: "whitespace",

	itemAnd:  "and",
	itemOr:   "or",
	itemEq:   "eq",
	itemNe:   "ne",
	itemLt:   "lt",
	itemGt:   "gt",
	itemLe:   "le",
	itemGe:   "ge",
	itemIn:   "in",
	itemNull: "null",
}

func (i itemType) String() string {
	s := itemName[i]
	if s == "" {
		return fmt.Sprintf("item%d", int(i))
	}
	return s
}

type lexTest struct {
	name  string
	input string
	items []item
}

func mkItem(typ itemType, text string) item {
	return item{
		typ: typ,
		val: text,
	}
}

var (
	tEOF        = mkItem(itemEOF, "")
	tLeftParen  = mkItem(itemLeftParen, "(")
	tRightParen = mkItem(itemRightParen, ")")
	tSpace      = mkItem(itemWhitespace, " ")
	tAnd        = mkItem(itemAnd, "and")
	tOr         = mkItem(itemOr, "or")
	tEq         = mkItem(itemEq, "eq")
	tNe         = mkItem(itemNe, "ne")
	tLt         = mkItem(itemLt, "lt")
	tGt         = mkItem(itemGt, "gt")
	tLe         = mkItem(itemLe, "le")
	tGe         = mkItem(itemGe, "ge")
	tIn         = mkItem(itemIn, "in")
	tNull       = mkItem(itemNull, "null")
	tTrue       = mkItem(itemBool, "true")
	tFalse      = mkItem(itemBool, "false")
	tComma      = mkItem(itemComma, ",")
)

var lexTests = []lexTest{
	{"empty", "", []item{tEOF}},
	{"spaces", " \t\n", []item{mkItem(itemWhitespace, " \t\n"), tEOF}},
	{"parens", "(((3)))", []item{
		tLeftParen,
		tLeftParen,
		tLeftParen,
		mkItem(itemNumber, "3"),
		tRightParen,
		tRightParen,
		tRightParen,
		tEOF,
	}},
	{"empty parens", "()", []item{tLeftParen, tRightParen, tEOF}},
	{"and", "and", []item{tAnd, tEOF}},
	{"or", "or", []item{tOr, tEOF}},
	{"eq", "eq", []item{tEq, tEOF}},
	{"ne", "ne", []item{tNe, tEOF}},
	{"lt", "lt", []item{tLt, tEOF}},
	{"gt", "gt", []item{tGt, tEOF}},
	{"le", "le", []item{tLe, tEOF}},
	{"ge", "ge", []item{tGe, tEOF}},
	{"in", "in", []item{tIn, tEOF}},
	{"null", "null", []item{tNull, tEOF}},
	{"true", "true", []item{tTrue, tEOF}},
	{"false", "false", []item{tFalse, tEOF}},
	{"string", `"test"`, []item{mkItem(itemString, `"test"`), tEOF}},
	{"escaped string", `"this \"is\" a test"`, []item{mkItem(itemString, `"this \"is\" a test"`), tEOF}},
	{"numbers", "1 02 -1 +1.2 .50", []item{
		mkItem(itemNumber, "1"),
		tSpace,
		mkItem(itemNumber, "02"),
		tSpace,
		mkItem(itemNumber, "-1"),
		tSpace,
		mkItem(itemNumber, "+1.2"),
		tSpace,
		mkItem(itemNumber, ".50"),
		tEOF,
	}},
	{"bools", "true false", []item{tTrue, tSpace, tFalse, tEOF}},
	{"identifiers", "and(id,12)", []item{
		tAnd,
		tLeftParen,
		mkItem(itemIdentifier, "id"),
		tComma,
		mkItem(itemNumber, "12"),
		tRightParen,
		tEOF,
	}},

	// errors
	{"badchar", "\x01", []item{
		mkItem(itemError, "unrecognized character in statement: U+0001"),
	}},
	{"EOF in parens", "(", []item{tLeftParen, mkItem(itemError, "unexpected end of statement")}},
	{"unclosed quote", `"`, []item{mkItem(itemError, "unterminated quoted string")}},
	{"bad number", "3k", []item{mkItem(itemError, `bad number syntax: "3k"`)}},
	{"extra right paren", "())", []item{tLeftParen, tRightParen, tRightParen, mkItem(itemError, "unexpected right paren U+0029 ')'")}},
	{"unterminated identifier", "abc123\x01", []item{
		mkItem(itemError, "bad character U+0001"),
	}},
	{"newline in string", `"\` + "\n", []item{
		mkItem(itemError, "unterminated quoted string"),
	}},
}

// collect gathers the emitted items into a slice
func collect(t *lexTest) (items []item) {
	l := lex(t.name, t.input)
	for {
		item := l.nextItem()
		items = append(items, item)
		if item.typ == itemEOF || item.typ == itemError {
			break
		}
	}
	return
}

func stringify(items []item) string {
	buf := []string{}
	for _, i := range items {
		buf = append(buf, i.String())
	}
	return strings.Join(buf, "")
}

func equal(i1, i2 []item, checkPos bool) bool {
	if len(i1) != len(i2) {
		return false
	}
	for k := range i1 {
		if i1[k].typ != i2[k].typ {
			// fmt.Printf("%d - 1: %q, 2: %q", k, i1[k].typ, i2[k].typ)
			return false
		}
		if i1[k].val != i2[k].val {
			// fmt.Printf("%d - 1: %q, 2: %q", k, i1[k].val, i2[k].val)
			return false
		}
		if checkPos && i1[k].pos != i2[k].pos {
			return false
		}
	}
	return true
}

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		items := collect(&test)
		if !equal(items, test.items, false) {
			t.Errorf("%s: got\n\t%+v\nexpected\n\t%+v", test.name, items, test.items)
		}
	}
}

var lexStrings = []struct {
	name  string
	exp   string
	items []item
}{
	{"EOF", "EOF", []item{tEOF}},
	{"error", "test error", []item{mkItem(itemError, "test error")}},
	{"keywords", "<and><or><eq><ne><lt><gt><le><ge><in><null>", []item{
		tAnd,
		tOr,
		tEq,
		tNe,
		tLt,
		tGt,
		tLe,
		tGe,
		tIn,
		tNull,
	}},
	{"long identifiers", `"1234567890"...`, []item{mkItem(itemIdentifier, "1234567890abcdef")}},
	{"identifiers", `"abc123"`, []item{mkItem(itemIdentifier, "abc123")}},
}

func TestString(t *testing.T) {
	for _, test := range lexStrings {
		if stringify(test.items) != test.exp {
			t.Errorf("%s: got\n\t%+v\nexpected\n\t%+v", test.name, stringify(test.items), test.exp)
		}
	}
}
