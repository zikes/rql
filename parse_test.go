package rql

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

type parseTest struct {
	name   string
	input  string
	ok     bool
	result string
}

const (
	noError  = true
	hasError = false
)

var parseTests = []parseTest{
	{"empty", "", noError, ``},
	{"spaces", " \n\t", noError, ``},
	{"and - empty", "and()", noError, `and()`},
	{"or - empty", "or()", noError, `or()`},
	{"eq - empty", "eq()", noError, `eq()`},
	{"ne - empty", "ne()", noError, `ne()`},
	{"gt - empty", "gt()", noError, `gt()`},
	{"lt - empty", "lt()", noError, `lt()`},
	{"ge - empty", "ge()", noError, `ge()`},
	{"le - empty", "le()", noError, `le()`},
	{"in - empty", "in()", noError, `in()`},
	{"null", "eq(null)", noError, `eq(null)`},
	{"identifier", "eq(id)", noError, `eq(id)`},
	{"number", "eq(-12.3)", noError, `eq(-12.3)`},
	{"string", `eq("test")`, noError, `eq("test")`},
	{"multiple value", `eq(id, -12.3, "test")`, noError, `eq(id,-12.3,"test")`},
	{"nested operators", `and(eq(id,12),gt(age,21))`, noError, `and(eq(id,12),gt(age,21))`},
	{"in - non-empty", `in(first_name, ("Jason","Kevin"))`, noError, `in(first_name,("Jason","Kevin"))`},
}

func testParse(doCopy bool, t *testing.T) {
	textFormat = "%q"
	defer func() { textFormat = "%s" }()
	for _, test := range parseTests {
		stmt, err := New(test.name).Parse(test.input)
		switch {
		case err == nil && !test.ok:
			t.Errorf("%q: expected error; got none", test.name)
			continue
		case err != nil && test.ok:
			t.Errorf("%q: unexpected error: %v", test.name, err)
			continue
		case err != nil && !test.ok:
			continue
		}
		var result string
		if doCopy {
			result = stmt.Root.Copy().String()
		} else {
			result = stmt.Root.String()
		}
		if result != test.result {
			t.Errorf("%s=(%q): got\n\t%v\nexpected\n\t%v", test.name, test.input, result, test.result)
		}
	}
}

func testParseFromTop(doCopy bool, t *testing.T) {
	textFormat = "%q"
	defer func() { textFormat = "%s" }()
	for _, test := range parseTests {
		stmt, err := Parse(test.name, test.input)
		switch {
		case err == nil && !test.ok:
			t.Errorf("%q: expected error; got none", test.name)
			continue
		case err != nil && test.ok:
			t.Errorf("%q: unexpected error: %v", test.name, err)
			continue
		case err != nil && !test.ok:
			continue
		}
		var result string
		if doCopy {
			result = stmt.Root.Copy().String()
		} else {
			result = stmt.Root.String()
		}
		if result != test.result {
			t.Errorf("%s=(%q): got\n\t%v\nexpected\n\t%v", test.name, test.input, result, test.result)
		}
	}
}

func TestParse(t *testing.T) {
	testParse(false, t)
	testParseFromTop(false, t)
}

func TestParseCopy(t *testing.T) {
	testParse(true, t)
	testParseFromTop(true, t)
}

func TestPeekNonSpace(t *testing.T) {
	tree := New("test peekNonSpace")
	tree.startParse(lex(tree.Name, "  eq()"))
	tree.text = "  eq()"
	tree.Root = tree.newStatement(tree.peek().pos, nil)
	tok := tree.peekNonSpace()
	if tok.typ != itemEq {
		t.Errorf("peekNonSpace failed to return non-space token")
	}
	tree2 := tree.Copy()
	tree.lex = nil
	tree.token = [3]item{}
	tree.peekCount = 0
	if !reflect.DeepEqual(tree2, tree) {
		t.Errorf("Tree.Copy() mismatch")
	}
}

func TestTreeCopy(t *testing.T) {
	tree, err := New("root").Parse("eq(id,12)")
	if err != nil {
		t.Fatalf("unexpected parse failure: %v", err)
	}
	treeCopy := tree.Copy()
	tree.lex = nil
	tree.token = [3]item{}
	tree.peekCount = 0
	if !reflect.DeepEqual(tree, treeCopy) {
		t.Errorf("Tree.Copy() mismatch")
	}
	var n *Tree
	if n.Copy() != nil {
		t.Errorf("nil Tree.Copy() not nil")
	}
}

func TestErrorContextWithTreeCopy(t *testing.T) {
	tree, err := New("root").Parse("and(eq(id,12),lt(height,500))")
	if err != nil {
		t.Fatalf("unexpected parse failure: %v", err)
	}
	treeCopy := tree.Copy()
	wantLocation, wantContext := tree.ErrorContext(tree.Root.Operator)
	gotLocation, gotContext := treeCopy.ErrorContext(treeCopy.Root.Operator)
	if wantLocation != gotLocation {
		t.Errorf("wrong error location want %q got %q", wantLocation, gotLocation)
	}
	if wantContext != gotContext {
		t.Errorf("wrong error context want %q got %q", wantContext, gotContext)
	}
}

// Benchmarks
func BenchmarkParseTiny(b *testing.B) {
	text := `eq(id,12)`
	for i := 0; i < b.N; i++ {
		_, err := New("bench").Parse(text)
		if err != nil {
			b.Fatal(err)
		}
	}
}
func BenchmarkParseSmall(b *testing.B) {
	text := generateStatement(1)
	for i := 0; i < b.N; i++ {
		_, err := New("bench").Parse(text)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseMedium(b *testing.B) {
	text := generateStatement(500)
	for i := 0; i < b.N; i++ {
		_, err := New("bench").Parse(text)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseLarge(b *testing.B) {
	text := generateStatement(5000)
	for i := 0; i < b.N; i++ {
		_, err := New("bench").Parse(text)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func generateStatement(len int) string {
	b := new(bytes.Buffer)
	fmt.Fprint(b, "and(")
	for i := 0; i < len; i++ {
		fmt.Fprint(b, `and(eq(id,12),lt(rating,-60.5),gt(height,173),or(in(first_name,("Jason","Kevin")),ne(last_name,"Costa")))`)
		if i < len-1 {
			fmt.Fprint(b, ",")
		}
	}
	fmt.Fprint(b, ")")
	return b.String()
}
