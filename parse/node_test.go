package rql

import "testing"

func TestNodeType_Type(t *testing.T) {
	if NodeBool.Type() != NodeBool {
		t.Errorf("NodeType.Type() mismatch")
	}
}

func TestListNode_Tree(t *testing.T) {
	node := &ListNode{tr: nil, NodeType: NodeList, Pos: Pos(0)}
	if node.tree() != nil {
		t.Errorf("unexpected ListNode.tree() value")
	}
}

func TestListNode_Copy(t *testing.T) {
	var node *ListNode
	if node.Copy() != node {
		t.Errorf("unexpected ListNode.CopyList() value")
	}
}

func TestStatementNode_Tree(t *testing.T) {
	node := &StatementNode{tr: nil}
	if node.tree() != nil {
		t.Errorf("unexpected StatementNode.tree() value")
	}
}

func TestIdentifierNode_Tree(t *testing.T) {
	node := &IdentifierNode{tr: nil}
	if node.tree() != nil {
		t.Errorf("unexpected IdentifierNode.tree() value")
	}
}

func TestNullNode_Type(t *testing.T) {
	var node *NullNode
	if node.Type() != NodeNull {
		t.Errorf("unexpected NullNode.Type() value")
	}
}
func TestNullNode_Tree(t *testing.T) {
	node := &NullNode{tr: nil}
	if node.tree() != nil {
		t.Errorf("unexpected NullNode.tree() value")
	}
}
func TestNullNode_Copy(t *testing.T) {
	node := &NullNode{tr: nil, Pos: Pos(0)}
	if node.Copy().tree() != nil {
		t.Errorf("unexpected NullNode.Copy() value")
	}
}

func TestBoolNode_Tree(t *testing.T) {
	node := &BoolNode{tr: nil}
	if node.tree() != nil {
		t.Errorf("unexpected BoolNode.tree() value")
	}
}
func TestBoolNode_Copy(t *testing.T) {
	node := &BoolNode{tr: nil, True: false}
	if node.Copy().String() != node.String() {
		t.Errorf("unexpected BoolNode.Copy() value")
	}
}

func TestNumberNode_Tree(t *testing.T) {
	node := &NumberNode{tr: nil}
	if node.tree() != nil {
		t.Errorf("unexpected NumberNode.tree() value")
	}
}

func TestStringNode_Tree(t *testing.T) {
	node := &StringNode{tr: nil}
	if node.tree() != nil {
		t.Errorf("unexpected StringNode.tree() value")
	}
}
func TestStringNode_Copy(t *testing.T) {
	node := &StringNode{tr: nil, Quoted: "", Text: ""}
	if node.Copy().String() != node.String() {
		t.Errorf("unexpected StringNode.Copy() value")
	}
}
