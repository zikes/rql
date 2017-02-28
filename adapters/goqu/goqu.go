package goquadapter

import (
	rql "git.nwaonline.com/rune/rql/parse"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gopkg.in/doug-martin/goqu.v3"
)

func ToGoqu(n rql.Node) goqu.Expression {
	switch n := n.(type) {
	case *rql.StatementNode:
		return ToGoqu(n.Operator)
	case *rql.OperatorNode:
		left := n.Operands.Nodes[0]
		right := n.Operands.Nodes[1]
		switch n.Operator {
		case "eq":
			return goqu.I(left.(*rql.IdentifierNode).Ident).Eq(value(right))
		}
	}
	return goqu.Ex{}
}

func ToSQL(n rql.Node) string {
	driver, _, _ := sqlmock.New()
	db := goqu.New("default", driver)
	e := ToGoqu(n)
	sql, _, _ := db.From("test").Where(e).ToSql()
	return sql
}

func value(n rql.Node) interface{} {
	switch n := n.(type) {
	case *rql.BoolNode:
		return n.True
	case *rql.NullNode:
		return nil
	case *rql.StringNode:
		return n.Text
	case *rql.NumberNode:
		switch {
		case n.IsInt:
			return n.Int64
		case n.IsUint:
			return n.Uint64
		case n.IsFloat:
			return n.Float64
		}
	case *rql.ListNode:
		vals := []interface{}{}
		for _, v := range n.Nodes {
			vals = append(vals, value(v))
		}
		return vals
	}
	return nil
}
