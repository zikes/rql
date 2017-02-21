package rql_test

import (
	"testing"

	"git.nwaonline.com/rune/rql"
)

var result rql.Statement

func benchmarkParseString(s string, b *testing.B) {
	var r rql.Statement
	for n := 0; n < b.N; n++ {
		r, _ = rql.ParseString(s)
	}
	result = r
}

func BenchmarkParseStringSimple(b *testing.B) { benchmarkParseString("eq()", b) }
func BenchmarkParseStringComplex(b *testing.B) {
	benchmarkParseString(`and(eq(id,12),lt(age,30),gt(height,137),or(eq(id,13),ge(age,21),not(false),and(eq(column,12),eq(my_col,"this is a test"),eq(my_col,12.123),and(eq(column,12),ne(my_col,-12)),or(eq(column,12),ne(my_col,-12)))))`, b)
}
