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
		{`""`, &rql.Literal{rql.STRING, `""`}},
		{`"test"`, &rql.Literal{rql.STRING, `"test"`}},

		// Literal Numbers
		{"12", &rql.Literal{rql.NUMERIC, "12"}},
		{"-12", &rql.Literal{rql.NUMERIC, `-12`}},

		// Literal Booleans
		{"true", &rql.Literal{rql.BOOLEAN, "true"}},
		{"false", &rql.Literal{rql.BOOLEAN, "false"}},

		// Identifiers
		{"col", &rql.Identifier{Name: "col"}},
		{"col_2", &rql.Identifier{Name: "col_2"}},
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

	exp_equal_int := &rql.Operator{
		Name: "eq",
		Operands: &rql.ParenExpr{
			Expressions: []rql.Expression{
				&rql.Identifier{Name: "column"},
				&rql.Literal{
					Kind:  rql.NUMERIC,
					Value: "12",
				},
			},
		},
	}

	exp_nequal_int := &rql.Operator{
		Name: "ne",
		Operands: &rql.ParenExpr{
			Expressions: []rql.Expression{
				&rql.Identifier{Name: "my_col"},
				&rql.Literal{
					Kind:  rql.NUMERIC,
					Value: "-12",
				},
			},
		},
	}

	exp_equal_string := &rql.Operator{
		Name: "eq",
		Operands: &rql.ParenExpr{
			Expressions: []rql.Expression{
				&rql.Identifier{Name: "my_col"},
				&rql.Literal{
					Kind:  rql.STRING,
					Value: `"this is a test"`,
				},
			},
		},
	}
	exp_equal_boolean := &rql.Operator{
		Name: "eq",
		Operands: &rql.ParenExpr{
			Expressions: []rql.Expression{
				&rql.Identifier{Name: "my_col"},
				&rql.Literal{
					Kind:  rql.BOOLEAN,
					Value: "true",
				},
			},
		},
	}
	exp_equal_identifier := &rql.Operator{
		Name: "eq",
		Operands: &rql.ParenExpr{
			Expressions: []rql.Expression{
				&rql.Identifier{Name: "my_col"},
				&rql.Identifier{Name: "my_other_col"},
			},
		},
	}
	exp_equal_float := &rql.Operator{
		Name: "eq",
		Operands: &rql.ParenExpr{
			Expressions: []rql.Expression{
				&rql.Identifier{Name: "my_col"},
				&rql.Literal{
					Kind:  rql.NUMERIC,
					Value: "12.123",
				},
			},
		},
	}
	exp_equal_negative_float := &rql.Operator{
		Name: "eq",
		Operands: &rql.ParenExpr{
			Expressions: []rql.Expression{
				&rql.Identifier{Name: "my_col"},
				&rql.Literal{
					Kind:  rql.NUMERIC,
					Value: "-12.123",
				},
			},
		},
	}

	exp_and := &rql.Operator{
		Name: "and",
		Operands: &rql.ParenExpr{
			Expressions: []rql.Expression{
				exp_equal_int,
				exp_nequal_int,
			},
		},
	}

	exp_or := &rql.Operator{
		Name: "or",
		Operands: &rql.ParenExpr{
			Expressions: []rql.Expression{
				exp_equal_int,
				exp_nequal_int,
			},
		},
	}

	exp_many_nested := &rql.Operator{
		Name: "and",
		Operands: &rql.ParenExpr{
			[]rql.Expression{
				exp_equal_int,
				exp_equal_string,
				exp_equal_float,
				exp_and,
				exp_or,
			},
		},
	}

	exp_not := &rql.Operator{
		Name: "not",
		Operands: &rql.ParenExpr{
			[]rql.Expression{
				exp_and,
			},
		},
	}

	exp_in := &rql.Operator{
		Name: "in",
		Operands: &rql.ParenExpr{
			[]rql.Expression{
				&rql.Identifier{Name: "primes"},
				&rql.ParenExpr{
					[]rql.Expression{
						&rql.Literal{Kind: rql.NUMERIC, Value: "1"},
						&rql.Literal{Kind: rql.NUMERIC, Value: "2"},
						&rql.Literal{Kind: rql.NUMERIC, Value: "3"},
						&rql.Literal{Kind: rql.NUMERIC, Value: "5"},
						&rql.Literal{Kind: rql.NUMERIC, Value: "7"},
					},
				},
			},
		},
	}

	var tests = []struct {
		s    string
		stmt *rql.Statement
		err  string
	}{
		{
			s:    "eq(column,12)",
			stmt: &rql.Statement{[]rql.Expression{exp_equal_int}},
		},
		{
			s:    "eq(column, 12)",
			stmt: &rql.Statement{[]rql.Expression{exp_equal_int}},
		},
		{
			s:    "ne(my_col,-12)",
			stmt: &rql.Statement{[]rql.Expression{exp_nequal_int}},
		},
		{
			s:    `eq(my_col,"this is a test")`,
			stmt: &rql.Statement{[]rql.Expression{exp_equal_string}},
		},
		{
			s:    `eq(my_col,true)`,
			stmt: &rql.Statement{[]rql.Expression{exp_equal_boolean}},
		},
		{
			s:    `eq(my_col,my_other_col)`,
			stmt: &rql.Statement{[]rql.Expression{exp_equal_identifier}},
		},
		{
			s:    `eq(my_col,12.123)`,
			stmt: &rql.Statement{[]rql.Expression{exp_equal_float}},
		},
		{
			s:    `eq(my_col,-12.123)`,
			stmt: &rql.Statement{[]rql.Expression{exp_equal_negative_float}},
		},
		{
			s:    `and(eq(column,12),ne(my_col,-12))`,
			stmt: &rql.Statement{[]rql.Expression{exp_and}},
		},
		{
			s:    `or(eq(column,12),ne(my_col,-12))`,
			stmt: &rql.Statement{[]rql.Expression{exp_or}},
		},
		{
			s:    `and(eq(column,12),eq(my_col,"this is a test"),eq(my_col,12.123),and(eq(column,12),ne(my_col,-12)),or(eq(column,12),ne(my_col,-12)))`,
			stmt: &rql.Statement{[]rql.Expression{exp_many_nested}},
		},
		{
			s:    `not(and(eq(column,12),ne(my_col,-12)))`,
			stmt: &rql.Statement{[]rql.Expression{exp_not}},
		},
		{
			s:    `in(primes,(1,2,3,5,7))`,
			stmt: &rql.Statement{[]rql.Expression{exp_in}},
		},
	}

	fmt.Printf("Testing Parser\n")

	for i, tt := range tests {
		stmt, err := rql.NewParser(strings.NewReader(tt.s)).Parse()
		fmt.Printf("  %d %s = %s\n", i, stmt.String(), tt.s)
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
