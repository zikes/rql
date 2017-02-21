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
	}

	fmt.Printf("Testing Constructors\n")
	for i, tt := range tests {
		if !reflect.DeepEqual(tt.exp, tt.got) {
			t.Errorf("  %d %q\n\nmismatch:\n    exp=%s\n    got=%s\n\n", i, tt.got, tt.exp, tt.got)
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
	}

	fmt.Printf("Testing Float Literal Valuation\n")
	for i, tt := range tests {
		if tt.lit.FloatValue() != tt.exp {
			t.Errorf("  %d %q\n\nmismatch:\n    exp=%s\n    got=%s\n\n", i, tt.lit, tt.exp, tt.lit.IntValue())
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
	}

	fmt.Printf("Testing Expression.String()\n")

	for i, tt := range tests {
		got := tt.e.String()
		fmt.Printf("  %d %s = %s\n", i, tt.s, got)
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
		fmt.Printf("  %d %q\n", i, tt.s)
		if !reflect.DeepEqual(tt.err, errstring(err)) {
			t.Errorf("  %d. %q: error mismatch:\n    exp=%s\n    got=%s\n\n", i, tt.s, tt.err, err)
		} else if tt.err == "" && !reflect.DeepEqual(tt.stmt, stmt) {
			t.Errorf("  %d. %q\n\nstmt mismatch:\n    exp=%s\n    got=%s\n\n", i, tt.s, tt.stmt, stmt)
		}
	}
}

func errstring(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
