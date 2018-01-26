package goquadapter

import (
	"testing"

	rql "github.com/zikes/rql/parse"
)

type parseTest struct {
	name   string
	input  string
	result string
}

var parseTests = []parseTest{
	{"empty", "", `SELECT * FROM "test"`},
	{"nested", "and(eq(id,12),or(lt(age,21),gt(height,156.2)))", `SELECT * FROM "test" WHERE (("id" = 12) AND (("age" < 21) OR ("height" > 156.2)))`},

	// operators
	{"equals", "eq(id,12)", `SELECT * FROM "test" WHERE ("id" = 12)`},
	{"not equals", "ne(id,12)", `SELECT * FROM "test" WHERE ("id" != 12)`},
	{"less than", "lt(id,12)", `SELECT * FROM "test" WHERE ("id" < 12)`},
	{"greater than", "gt(id,12)", `SELECT * FROM "test" WHERE ("id" > 12)`},
	{"less than equals", "le(id,12)", `SELECT * FROM "test" WHERE ("id" <= 12)`},
	{"greater than equals", "ge(id,12)", `SELECT * FROM "test" WHERE ("id" >= 12)`},
	{"and", "and(eq(id,12),lt(age,21))", `SELECT * FROM "test" WHERE (("id" = 12) AND ("age" < 21))`},
	{"or", "or(eq(id,12),lt(age,21))", `SELECT * FROM "test" WHERE (("id" = 12) OR ("age" < 21))`},
	{"in", "in(id,(12,13,14))", `SELECT * FROM "test" WHERE ("id" IN ((12, 13, 14)))`},

	{"null", "eq(id,null)", `SELECT * FROM "test" WHERE ("id" IS NULL)`},
	{"bool", "eq(id,true)", `SELECT * FROM "test" WHERE ("id" IS TRUE)`},
	{"string", `eq(id,"test")`, `SELECT * FROM "test" WHERE ("id" = 'test')`},
}

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		stmt, err := rql.New(test.name).Parse(test.input)
		if err != nil {
			t.Fatalf("unexpected parse failure: %v", err)
		}
		got := ToSQL(stmt.Root)
		if got != test.result {
			t.Errorf("%s: SQL mismatch\n\texpected:\n\t\t%s\n\tgot:\n\t\t%s", test.name, test.result, got)
		}
	}
}
