package rql_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"git.nwaonline.com/rune/rql"
)

func TestExpressions_String(t *testing.T) {
	var tests = []struct {
		s string
		e rql.Expression
	}{
		// Literal Strings
		{`""`, rql.Literal{rql.STRING, `""`}},
		{`"test"`, rql.Literal{rql.STRING, `"test"`}},

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
		{"or()", rql.Operator{Kind: rql.OR}},
		{"not()", rql.Operator{Kind: rql.NOT}},
		{"lt()", rql.Operator{Kind: rql.LT}},
		{"gt()", rql.Operator{Kind: rql.GT}},
		{"le()", rql.Operator{Kind: rql.LE}},
		{"ge()", rql.Operator{Kind: rql.GE}},
		{"eq()", rql.Operator{Kind: rql.EQ}},
		{"ne()", rql.Operator{Kind: rql.NE}},
		{"in()", rql.Operator{Kind: rql.IN}},
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
