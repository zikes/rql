package rql_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"git.nwaonline.com/rune/rql"
)

func TestExpression_Constructors(t *testing.T) {
	var tests = []struct {
		exp rql.Expression
		got rql.Expression
	}{
		// Literals
		{rql.Literal{rql.STRING, ``}, rql.Lit("")},
		{rql.Literal{rql.STRING, `"testing"`}, rql.Lit(`"testing"`)},
		{rql.Literal{rql.STRING, `test "testing" test`}, rql.Lit(`test "testing" test`)},

		{rql.Literal{rql.NUMERIC, "0"}, rql.Lit(0)},
		{rql.Literal{rql.NUMERIC, "1"}, rql.Lit(1)},
		{rql.Literal{rql.NUMERIC, "-1"}, rql.Lit(-1)},
		{rql.Literal{rql.NUMERIC, "0.1"}, rql.Lit(0.1)},
		{rql.Literal{rql.NUMERIC, "-0.1"}, rql.Lit(-0.1)},

		{rql.Literal{rql.NULL, ""}, rql.Lit(nil)},

		{rql.Literal{rql.BOOLEAN, "true"}, rql.Lit(true)},
		{rql.Literal{rql.BOOLEAN, "false"}, rql.Lit(false)},

		// Identifiers
		{rql.Identifier{rql.IDENT, "id"}, rql.Ident("id")},

		// Arrays
		{rql.ExpressionList{rql.Lit(1), rql.Lit(2)}, rql.Array(1, 2)},

		// Operators
		{rql.Operator{rql.AND, rql.ExpressionList{}}, rql.And()},
		{rql.Operator{rql.AND, rql.ExpressionList{rql.Lit(true), rql.Lit(false)}}, rql.And(true, false)},

		{rql.Operator{rql.OR, rql.ExpressionList{}}, rql.Or()},
		{rql.Operator{rql.OR, rql.ExpressionList{rql.Lit(true), rql.Lit(false)}}, rql.Or(true, false)},
		{rql.Operator{rql.NOT, rql.ExpressionList{}}, rql.Not()},
		{rql.Operator{rql.NOT, rql.ExpressionList{rql.Lit(true), rql.Lit(false)}}, rql.Not(true, false)},
		{rql.Operator{rql.LT, rql.ExpressionList{}}, rql.Lt()},
		{rql.Operator{rql.LT, rql.ExpressionList{rql.Identifier{rql.IDENT, "id"}, rql.Literal{rql.NUMERIC, "12"}}}, rql.Lt(rql.Ident("id"), 12)},
		{rql.Operator{rql.GT, rql.ExpressionList{}}, rql.Gt()},
		{rql.Operator{rql.GT, rql.ExpressionList{rql.Identifier{rql.IDENT, "id"}, rql.Literal{rql.NUMERIC, "12"}}}, rql.Gt(rql.Ident("id"), 12)},
		{rql.Operator{rql.LE, rql.ExpressionList{}}, rql.Le()},
		{rql.Operator{rql.LE, rql.ExpressionList{rql.Identifier{rql.IDENT, "id"}, rql.Literal{rql.NUMERIC, "12"}}}, rql.Le(rql.Ident("id"), 12)},
		{rql.Operator{rql.GE, rql.ExpressionList{}}, rql.Ge()},
		{rql.Operator{rql.GE, rql.ExpressionList{rql.Identifier{rql.IDENT, "id"}, rql.Literal{rql.NUMERIC, "12"}}}, rql.Ge(rql.Ident("id"), 12)},
		{rql.Operator{rql.EQ, rql.ExpressionList{}}, rql.Eq()},
		{rql.Operator{rql.EQ, rql.ExpressionList{rql.Identifier{rql.IDENT, "id"}, rql.Literal{rql.NUMERIC, "12"}}}, rql.Eq(rql.Ident("id"), 12)},
		{rql.Operator{rql.NE, rql.ExpressionList{}}, rql.Ne()},
		{rql.Operator{rql.NE, rql.ExpressionList{rql.Identifier{rql.IDENT, "id"}, rql.Literal{rql.NUMERIC, "12"}}}, rql.Ne(rql.Ident("id"), 12)},
		{rql.Operator{rql.IN, rql.ExpressionList{}}, rql.In()},
		{rql.Operator{rql.IN, rql.ExpressionList{rql.Identifier{rql.IDENT, "id"}, rql.ExpressionList{rql.Literal{rql.NUMERIC, "1"}, rql.Literal{rql.NUMERIC, "2"}, rql.Literal{rql.NUMERIC, "3"}, rql.Literal{rql.NUMERIC, "4"}, rql.Literal{rql.NUMERIC, "5"}}}}, rql.In(rql.Ident("id"), rql.Array(1, 2, 3, 4, 5))},
	}

	fmt.Printf("Testing Constructors\n")
	for i, tt := range tests {
		if !reflect.DeepEqual(tt.exp, tt.got) {
			t.Errorf("  %d %q\nmismatch:\n    exp=%s\n    got=%s\n\n", i, tt.got, tt.exp, tt.got)
		}
	}
}

func TestIntLiteral_Value(t *testing.T) {
	var tests = []struct {
		exp int
		lit rql.Literal
	}{
		{0, rql.Lit(0)},
		{1, rql.Lit(1)},
		{-1, rql.Lit(-1)},

		// Overflows
		{0, rql.Literal{rql.NUMERIC, "9223372036854775808"}},
	}

	fmt.Printf("Testing Integer Literal Valuation\n")
	for i, tt := range tests {
		if tt.lit.IntValue() != tt.exp {
			t.Errorf("  %d %q\n\nmismatch:\n    exp=%s\n    got=%s\n\n", i, tt.lit, tt.exp, tt.lit.IntValue())
		}
	}
}

func TestFloatLiteral_Value(t *testing.T) {
	var tests = []struct {
		exp float64
		lit rql.Literal
	}{
		{float64(0), rql.Lit(float64(0))},
		{1.1, rql.Lit(1.1)},
		{-1.1, rql.Lit(-1.1)},
		{float64(9.223372036854776e+26), rql.Literal{rql.NUMERIC, "922337203685477580812345968"}},

		{0, rql.Literal{rql.NUMERIC, "abc"}},
	}

	fmt.Printf("Testing Float Literal Valuation\n")
	for i, tt := range tests {
		if !reflect.DeepEqual(tt.lit.FloatValue(), tt.exp) {
			t.Errorf("  %d %q\n\nmismatch:\n    exp=%s\n    got=%s\n\n", i, tt.lit, tt.exp, tt.lit.FloatValue())
		}
	}
}

func TestStringLiteral_Value(t *testing.T) {
	var tests = []struct {
		exp string
		lit rql.Literal
	}{
		{"", rql.Lit("")},
		{`""`, rql.Lit(`""`)},
		{`this is a "test"`, rql.Lit(`this is a "test"`)},
	}

	fmt.Printf("Testing String Literal Valuation\n")
	for i, tt := range tests {
		if tt.lit.StringValue() != tt.exp {
			t.Errorf("  %d %q\n\nmismatch:\n    exp=%s\n    got=%s\n\n", i, tt.lit, tt.exp, tt.lit.IntValue())
		}
	}
}

func TestBooleanLiteral_Value(t *testing.T) {
	var tests = []struct {
		exp bool
		lit rql.Literal
	}{
		{true, rql.Lit(true)},
		{false, rql.Lit(false)},
	}

	fmt.Printf("Testing Boolean Literal Valuation\n")
	for i, tt := range tests {
		if tt.lit.BoolValue() != tt.exp {
			t.Errorf("  %d %q\n\nmismatch:\n    exp=%s\n    got=%s\n\n", i, tt.lit, tt.exp, tt.lit.IntValue())
		}
	}
}

func TestExpressions_String(t *testing.T) {
	var tests = []struct {
		s string
		e rql.Expression
	}{
		// Literal Strings
		{`""`, rql.Literal{rql.STRING, ``}},
		{`"test"`, rql.Literal{rql.STRING, `test`}},

		// Literal Numbers
		{"12", rql.Literal{rql.NUMERIC, "12"}},
		{"-12", rql.Literal{rql.NUMERIC, `-12`}},
		{"12.123", rql.Literal{rql.NUMERIC, "12.123"}},
		{"-12.123", rql.Literal{rql.NUMERIC, `-12.123`}},

		// Literal Booleans
		{"true", rql.Literal{rql.BOOLEAN, "true"}},
		{"false", rql.Literal{rql.BOOLEAN, "false"}},
		{"null", rql.Literal{rql.NULL, "null"}},

		// Identifiers
		{"col", rql.Identifier{Kind: rql.IDENT, Name: "col"}},
		{"col_2", rql.Identifier{Kind: rql.IDENT, Name: "col_2"}},

		// Operators
		{"and()", rql.Operator{Kind: rql.AND}},
		{"and(id,eq(col,12))", rql.Operator{
			Kind: rql.AND,
			Operands: rql.ExpressionList{
				rql.Identifier{Kind: rql.IDENT, Name: "id"},
				rql.Operator{
					Kind: rql.EQ,
					Operands: rql.ExpressionList{
						rql.Identifier{Kind: rql.IDENT, Name: "col"},
						rql.Literal{rql.NUMERIC, "12"},
					},
				},
			},
		}},
		{"or()", rql.Operator{Kind: rql.OR}},
		{"not()", rql.Operator{Kind: rql.NOT}},
		{"lt()", rql.Operator{Kind: rql.LT}},
		{"gt()", rql.Operator{Kind: rql.GT}},
		{"le()", rql.Operator{Kind: rql.LE}},
		{"ge()", rql.Operator{Kind: rql.GE}},
		{"eq()", rql.Operator{Kind: rql.EQ}},
		{"ne()", rql.Operator{Kind: rql.NE}},
		{"in()", rql.Operator{Kind: rql.IN}},

		// Arrays
		{"()", rql.ExpressionList{}},
		{"(id)", rql.ExpressionList{
			rql.Identifier{Kind: rql.IDENT, Name: "id"},
		}},
		{"(id1,id2)", rql.ExpressionList{
			rql.Identifier{Kind: rql.IDENT, Name: "id1"},
			rql.Identifier{Kind: rql.IDENT, Name: "id2"},
		}},
		{"(id1,id2)", rql.ExpressionList{
			rql.Identifier{Kind: rql.IDENT, Name: "id1"},
			rql.Whitespace{Kind: rql.WS, Value: " "},
			rql.Identifier{Kind: rql.IDENT, Name: "id2"},
		}},

		// Whitespace
		{" ", rql.Whitespace{Kind: rql.WS, Value: " "}},
	}

	fmt.Printf("Testing Expression.String()\n")

	for i, tt := range tests {
		got := tt.e.String()
		if tt.s != got {
			t.Errorf("  %d. string mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, got)
		}
	}
}

func TestStatement_String(t *testing.T) {
	var tests = []struct {
		s string
		e rql.Statement
	}{
		// Statement
		{"eq()", rql.Statement{rql.Operator{Kind: rql.EQ}, rql.Whitespace{rql.WS, " "}}},
	}

	fmt.Printf("Testing Statement.String()\n")

	for i, tt := range tests {
		got := tt.e.String()
		if tt.s != got {
			t.Errorf("  %d. string mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, got)
		}
	}
}

func TestParser_ParseStatement(t *testing.T) {
	var tests = []struct {
		s    string
		stmt rql.Statement
		err  string
	}{
		{
			s:    "\teq (\n\tmy_col,\n\tnull\n)",
			stmt: rql.Statement{exp_null},
		},
		{
			s:    "eq(column,12)",
			stmt: rql.Statement{exp_equal_int},
		},
		{
			s:    "ne(my_col,-12)",
			stmt: rql.Statement{exp_nequal_int},
		},
		{
			s:    `eq(my_col,"this is a test")`,
			stmt: rql.Statement{exp_equal_string},
		},
		{
			s:    `eq(my_col,true)`,
			stmt: rql.Statement{exp_equal_boolean},
		},
		{
			s:    `eq(my_col,my_other_col)`,
			stmt: rql.Statement{exp_equal_identifier},
		},
		{
			s:    `eq(my_col,12.123)`,
			stmt: rql.Statement{exp_equal_float},
		},
		{
			s:    `eq(my_col,-12.123)`,
			stmt: rql.Statement{exp_equal_negative_float},
		},
		{
			s:    `and(eq(column,12),ne(my_col,-12))`,
			stmt: rql.Statement{exp_and},
		},
		{
			s:    `or(eq(column,12),ne(my_col,-12))`,
			stmt: rql.Statement{exp_or},
		},
		{
			s:    `and(eq(column,12),eq(my_col,"this is a test"),eq(my_col,12.123),and(eq(column,12),ne(my_col,-12)),or(eq(column,12),ne(my_col,-12)))`,
			stmt: rql.Statement{exp_many_nested},
		},
		{
			s:    `not(and(eq(column,12),ne(my_col,-12)))`,
			stmt: rql.Statement{exp_not},
		},
		{
			s:    `in(primes,(1,2,3,5,7))`,
			stmt: rql.Statement{exp_in},
		},
		{
			s:    `eq(my_col,null)`,
			stmt: rql.Statement{exp_null},
		},
		{
			s:    `eq(col,"a string with \"quotes\" in it")`,
			stmt: rql.Statement{exp_quoted_string},
		},
	}

	fmt.Printf("Testing Parser\n")

	for i, tt := range tests {
		stmt, err := rql.NewParser(strings.NewReader(tt.s)).Parse()
		if !reflect.DeepEqual(tt.err, errstring(err)) {
			t.Errorf("  %d. %q: error mismatch:\n    exp=%s\n    got=%s\n\n", i, tt.s, tt.err, err)
		} else if tt.err == "" && !reflect.DeepEqual(tt.stmt, stmt) {
			t.Errorf("  %d. %q\n\nstmt mismatch:\n    exp=%s\n    got=%s\n\n", i, tt.s, tt.stmt, stmt)
		}
	}
}

func TestInvalidSyntax(t *testing.T) {
	var tests = []struct {
		s   string
		exp error
	}{
		{``, fmt.Errorf("found EOF, expected operator")},
		{`and(*)`, fmt.Errorf("error parsing expression list: found *, expected value expression")},
		{`or`, fmt.Errorf("error parsing expression list: found EOF, expected open parentheses")},
		{`or)`, fmt.Errorf("error parsing expression list: found ), expected open parentheses")},
		{`or{`, fmt.Errorf("error parsing expression list: found {, expected open parentheses")},
		{`or("unclosed string)`, fmt.Errorf("error parsing expression list: found EOF, expected comma or close parentheses")},
	}

	fmt.Printf("Testing Errors\n")
	for i, tt := range tests {
		stmt, err := rql.ParseString(tt.s)
		if !reflect.DeepEqual(tt.exp, err) {
			t.Errorf("  %d. %v - %v", i, stmt, err)
		}
	}
}

func TestExpressionTokenValues(t *testing.T) {
	var tests = []struct {
		e rql.Expression
		t rql.Token
	}{
		// Arrays
		{rql.Array(), rql.EXPR_LIST},

		// Operators
		{rql.And(), rql.AND},
		{rql.Or(), rql.OR},
		{rql.Not(), rql.NOT},
		{rql.Lt(), rql.LT},
		{rql.Gt(), rql.GT},
		{rql.Le(), rql.LE},
		{rql.Ge(), rql.GE},
		{rql.Eq(), rql.EQ},
		{rql.Ne(), rql.NE},
		{rql.In(), rql.IN},

		// Literals
		{rql.Lit(""), rql.STRING},
		{rql.Lit(1), rql.NUMERIC},
		{rql.Lit(1.23), rql.NUMERIC},
		{rql.Lit(nil), rql.NULL},

		// Punctuation
		{rql.Punctuation{rql.COMMA, ","}, rql.COMMA},
		{rql.Punctuation{rql.LPAREN, "("}, rql.LPAREN},
		{rql.Punctuation{rql.RPAREN, ")"}, rql.RPAREN},

		// Whitespace
		{rql.Whitespace{rql.WS, " "}, rql.WS},
	}

	fmt.Printf("Testing Expression Token Values\n")
	for i, tt := range tests {
		testToken := tt.e.Token()
		if !reflect.DeepEqual(testToken, tt.t) {
			t.Errorf("  %d. token mismatch:\n   exp=%s\n   got=%s\n", i, tt.t, testToken)
		}
	}
}

func TestTokens(t *testing.T) {
	var tests = []struct {
		token         rql.Token
		str           string
		isLiteral     bool
		isOperator    bool
		isPunctuation bool
		isArray       bool
		isWhitespace  bool
		isEOF         bool
		isIllegal     bool
		isIdentifier  bool
		isValue       bool
	}{
		// token,       str,         lit,   op,    punct, arr,   ws,    eof,   ill,   ident, val
		{rql.ILLEGAL, `ILLEGAL`, false, false, false, false, false, false, true, false, false},
		{rql.EOF, `EOF`, false, false, false, false, false, true, false, false, false},
		{rql.WS, `WS`, false, false, false, false, true, false, false, false, false},

		{rql.IDENT, `IDENT`, true, false, false, false, false, false, false, true, true},
		{rql.NUMERIC, `NUMERIC`, true, false, false, false, false, false, false, false, true},
		{rql.STRING, `STRING`, true, false, false, false, false, false, false, false, true},
		{rql.BOOLEAN, `BOOLEAN`, true, false, false, false, false, false, false, false, true},
		{rql.NULL, `NULL`, true, false, false, false, false, false, false, false, true},

		{rql.AND, `and`, false, true, false, false, false, false, false, false, true},
		{rql.OR, `or`, false, true, false, false, false, false, false, false, true},
		{rql.NOT, `not`, false, true, false, false, false, false, false, false, true},
		{rql.LT, `lt`, false, true, false, false, false, false, false, false, true},
		{rql.GT, `gt`, false, true, false, false, false, false, false, false, true},
		{rql.LE, `le`, false, true, false, false, false, false, false, false, true},
		{rql.GE, `ge`, false, true, false, false, false, false, false, false, true},
		{rql.EQ, `eq`, false, true, false, false, false, false, false, false, true},
		{rql.NE, `ne`, false, true, false, false, false, false, false, false, true},
		{rql.IN, `in`, false, true, false, false, false, false, false, false, true},

		{rql.EXPR_LIST, `EXPR_LIST`, false, false, false, true, false, false, false, false, true},
		{rql.Token(12345), `token(12345)`, false, false, false, false, false, false, true, false, false},
	}

	fmt.Println("Testing Tokens")

	for i, tt := range tests {
		tok := tt.token
		if tok.String() != tt.str {
			t.Errorf("  %d. string mismatch:\n   exp=%s\n   got=%s\n", i, tt.str, tok.String())
		}
		if tok.IsLiteral() != tt.isLiteral {
			t.Errorf("  %d. type mismatch:\n   exp=%s\n   got=%s\n", i, tt.isLiteral, tok.IsLiteral())
		}
		if tok.IsOperator() != tt.isOperator {
			t.Errorf("  %d. type mismatch:\n   exp=%s\n   got=%s\n", i, tt.isOperator, tok.IsOperator())
		}
		if tok.IsPunctuation() != tt.isPunctuation {
			t.Errorf("  %d. type mismatch:\n   exp=%s\n   got=%s\n", i, tt.isPunctuation, tok.IsPunctuation())
		}
		if tok.IsEOF() != tt.isEOF {
			t.Errorf("  %d. type mismatch:\n   exp=%s\n   got=%s\n", i, tt.isEOF, tok.IsEOF())
		}
		if tok.IsIllegal() != tt.isIllegal {
			t.Errorf("  %d. type mismatch:\n   exp=%s\n   got=%s\n", i, tt.isIllegal, tok.IsIllegal())
		}
		if tok.IsIdentifier() != tt.isIdentifier {
			t.Errorf("  %d. type mismatch:\n   exp=%s\n   got=%s\n", i, tt.isIdentifier, tok.IsIdentifier())
		}
		if tok.IsArray() != tt.isArray {
			t.Errorf("  %d. type mismatch:\n   exp=%s\n   got=%s\n", i, tt.isArray, tok.IsArray())
		}
		if tok.IsWhitespace() != tt.isWhitespace {
			t.Errorf("  %d. type mismatch:\n   exp=%s\n   got=%s\n", i, tt.isWhitespace, tok.IsWhitespace())
		}
		if tok.IsValue() != tt.isValue {
			t.Errorf("  %d. type mismatch:\n   exp=%s\n   got=%s\n", i, tt.isValue, tok.IsValue())
		}
	}
}

func errstring(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
