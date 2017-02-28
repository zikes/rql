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
	{"boolean", "eq(id,true,false)", noError, "eq(id,true,false)"},
	{"multiple value", `eq(id, -12.3, "test")`, noError, `eq(id,-12.3,"test")`},
	{"nested operators", `and(eq(id,12),gt(age,21))`, noError, `and(eq(id,12),gt(age,21))`},
	{"in - non-empty", `in(first_name, ("Jason","Kevin"))`, noError, `in(first_name,("Jason","Kevin"))`},

	// errors
	{"unexpected token", `12`, hasError, `statement: unexpected token:1: unexpected token after operator: "\"12\""`},
	{"unexpected token 2", `eq(id 12)`, hasError, `statement: unexpected token 2:1: unexpected "12" in comma or right parentheses`},
	{"unexpected token 3", `eq,(id 12)`, hasError, `statement: unexpected token 3:1: unexpected "," in left parentheses`},
	{"unterminated string", `eq(id,"test)`, hasError, `statement: unterminated string:0: unexpected  in comma or right parentheses`},
	{"invalid number", `eq(-12e3)`, hasError, `statement: invalid number:0: unexpected  in comma or right parentheses`},
	{"number", "eq(+2.2.2)", hasError, `statement: number:0: unexpected  in comma or right parentheses`},
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
			if fmt.Sprintf("%s", err) != test.result {
				t.Errorf("%q: error mismatch: expected\n  %s\ngot\n  %s", test.name, test.result, err)
			}
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

func TestError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("error failed to panic")
		}
	}()
	tr, err := New("root").Parse("eq(id,12)")
	if err != nil {
		t.Fatalf("unexpected parse failure: %v", err)
	}
	tr.error(fmt.Errorf("test error"))
}

var errorContextTests = []string{
	"and(eq(id,12),lt(height,500))",
	`and(
		eq(id,12),
		lt(height,500),
	)`,
	"and(eq(id,12),or(lt(height,500),eq(id,12),lt(height,500),eq(id,12),lt(height,500),eq(id,12),lt(height,500),eq(id,12),lt(height,500)))",
}

func TestErrorContextWithTreeCopy(t *testing.T) {
	for _, test := range errorContextTests {
		tree, err := New("root").Parse(test)
		if err != nil {
			t.Fatalf("unexpected parse failure: %v", err)
		}
		treeCopy := tree.Copy()

		last := len(tree.Root.Operator.Operands.Nodes) - 1

		wantLocation, wantContext := tree.ErrorContext(tree.Root.Operator.Operands.Nodes[last])
		gotLocation, gotContext := treeCopy.ErrorContext(treeCopy.Root.Operator.Operands.Nodes[last])
		if wantLocation != gotLocation {
			t.Errorf("wrong error location want %q got %q", wantLocation, gotLocation)
		}
		if wantContext != gotContext {
			t.Errorf("wrong error context want %q got %q", wantContext, gotContext)
		}
	}
}

func TestErrorContextWithDetachedNode(t *testing.T) {
	tree, err := New("root").Parse("eq(id,12)")
	if err != nil {
		t.Fatalf("unexpected parse failure: %v", err)
	}
	op := tree.Root.Operator.Copy()
	op.(*OperatorNode).tr = nil
	wantLocation, wantContext := tree.ErrorContext(tree.Root.Operator)
	gotLocation, gotContext := tree.ErrorContext(op)
	if wantLocation != gotLocation {
		t.Errorf("wrong error location want %q got %q", wantLocation, gotLocation)
	}
	if wantContext != gotContext {
		t.Errorf("wrong error context want %q got %q", wantContext, gotContext)
	}
}

type mysteryNode struct {
	NodeType
	Pos
	tr *Tree
}

func (m *mysteryNode) Type() NodeType {
	return NodeType(-1)
}
func (m *mysteryNode) String() string { return "" }
func (m *mysteryNode) tree() *Tree    { return nil }
func (m *mysteryNode) Copy() Node     { return nil }

var isEmptyTests = []struct {
	name  string
	input string
	empty bool
}{
	{"empty", "", true},
	{"nonempty", "eq()", false},
	{"spaces", "\n\t \n\t ", true},
}

func TestIsEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("error failed to panic")
		}
	}()
	if !IsEmptyTree(nil) {
		t.Errorf("nil tree is not empty")
	}
	for _, test := range isEmptyTests {
		tree, err := New("root").Parse(test.input)
		if err != nil {
			t.Errorf("%q: unexpected error: %v", test.name, err)
			continue
		}
		if empty := IsEmptyTree(tree.Root); empty != test.empty {
			t.Errorf("%q: expected %t got %t", test.name, test.empty, empty)
		}
	}

	tree := New("root")

	listNode := tree.newList(Pos(0))
	listNode.append(tree.newBool(Pos(0), true))
	num, err := tree.newNumber(Pos(0), "0")
	if err != nil {
		t.Errorf("unexpected parse error: %v", err)
	}

	nodes := []Node{
		listNode,
		NewIdentifier("id"),
		tree.newNull(Pos(0)),
		tree.newOperator("eq", Pos(0), listNode),
		tree.newBool(Pos(0), true),
		tree.newString(Pos(0), "", ""),
		num,
		&mysteryNode{},
	}

	for _, node := range nodes {
		if IsEmptyTree(node) != false {
			t.Errorf("%T: expected %t got %t", node, false, true)
		}
	}
}

func TestPeekBackup(t *testing.T) {
	tree := New("test peek")
	tree.startParse(lex(tree.Name, "eq(id, 12)"))
	tree.text = "eq(id, 12)"
	tree.Root = tree.newStatement(tree.peek().pos, nil)
	tok := tree.next()
	tree.backup()
	tok2 := tree.peek()
	if !reflect.DeepEqual(tok, tok2) {
		t.Errorf("token mismatch, %q != %q", tok, tok2)
	}
}

type numberTest struct {
	text    string
	isInt   bool
	isUint  bool
	isFloat bool
	int64
	uint64
	float64
}

var numberTests = []numberTest{
	// basics
	{"0", true, true, true, 0, 0, 0},
	{"-0", true, true, true, 0, 0, 0},
	{"73", true, true, true, 73, 73, 73},
	{"073", true, true, true, 073, 073, 073},
	{"0x73", true, true, true, 0x73, 0x73, 0x73},
	{"-73", true, false, true, -73, 0, -73},
	{"+73", true, false, true, 73, 0, 73},
	{"100", true, true, true, 100, 100, 100},
	{"1e9", true, true, true, 1e9, 1e9, 1e9},
	{"-1e9", true, false, true, -1e9, 0, -1e9},
	{"-1.2", false, false, true, 0, 0, -1.2},
	{"1e19", false, true, true, 0, 1e19, 1e19},
	{"-1e19", false, false, true, 0, 0, -1e19},
	{"18446744073709551615", false, true, true, 0, 18446744073709551615, 18446744073709551615},
	{"18446744073709551616", false, false, false, 0, 0, 0},
	// errors
	{text: "+-2"},
	{text: "0x123."},
	{text: "1e."},
	{text: "0xi."},
	{text: "1+2."},
	{text: "'x"},
	{text: "'xx'"},
	{text: "'433937734937734969526500969526500'"},
}

func TestNumberParse(t *testing.T) {
	for _, test := range numberTests {
		var tree *Tree
		n, err := tree.newNumber(0, test.text)
		ok := test.isInt || test.isUint || test.isFloat
		if ok && err != nil {
			t.Errorf("unexpected error for %q: %s", test.text, err)
		}
		if !ok && err == nil {
			t.Errorf("expected error for %q", test.text)
		}
		if !ok {
			continue
		}
		if test.isInt {
			if !n.IsInt {
				t.Errorf("expected integer for %q", test.text)
			}
			if n.Int64 != test.int64 {
				t.Errorf("int64 for %q should be %d - is %d", test.text, test.int64, n.Int64)
			}
		} else if n.IsInt {
			t.Errorf("did not expect integer for %q", test.text)
		}
		if test.isUint {
			if !n.IsUint {
				t.Errorf("expected unsigned integer for %q", test.text)
			}
			if n.Uint64 != test.uint64 {
				t.Errorf("uint64 for %q should be %d - is %d", test.text, test.uint64, n.Uint64)
			}
		} else if n.IsUint {
			t.Errorf("did not expect unsigned integer for %q", test.text)
		}
		if test.isFloat {
			if !n.IsFloat {
				t.Errorf("unexpected float for %q", test.text)
			}
			if n.Float64 != test.float64 {
				t.Errorf("float64 for %q should be %g - is %g", test.text, test.float64, n.Float64)
			}
		} else if n.IsFloat {
			t.Errorf("did not expect float for %q", test.text)
		}
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
