package goquadapter

import (
	"reflect"

	rql "git.nwaonline.com/rune/rql/parse"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gopkg.in/doug-martin/goqu.v3"
)

func ToGoqu(n rql.Node) goqu.Expression {
	switch n := n.(type) {
	case *rql.StatementNode:
		return ToGoqu(n.Operator)
	case *rql.OperatorNode:
		if n == nil || len(n.Operands.Nodes) == 0 {
			return goqu.Ex{}
		}
		left := n.Operands.Nodes[0]
		right := n.Operands.Nodes[1]
		switch n.Operator {
		case "eq":
			return goqu.I(left.(*rql.IdentifierNode).Ident).Eq(value(right))
		case "ne":
			return goqu.I(left.(*rql.IdentifierNode).Ident).Neq(value(right))
		case "lt":
			return goqu.I(left.(*rql.IdentifierNode).Ident).Lt(value(right))
		case "gt":
			return goqu.I(left.(*rql.IdentifierNode).Ident).Gt(value(right))
		case "le":
			return goqu.I(left.(*rql.IdentifierNode).Ident).Lte(value(right))
		case "ge":
			return goqu.I(left.(*rql.IdentifierNode).Ident).Gte(value(right))
		case "in":
			values := []interface{}{}
			for _, v := range n.Operands.Nodes[1:] {
				values = append(values, value(v))
			}
			return goqu.I(left.(*rql.IdentifierNode).Ident).In(values)
		case "or":
			exOr, err := goqu.ExOr{}.ToExpressions()
			if err != nil {
				panic(err)
			}
			for _, v := range n.Operands.Nodes {
				exOr = exOr.Append(ToGoqu(v))
			}
			return exOr
		case "and":
			ex, err := goqu.Ex{}.ToExpressions()
			if err != nil {
				panic(err)
			}
			for _, v := range n.Operands.Nodes {
				ex = ex.Append(ToGoqu(v))
			}
			return ex
		}
	}
	return goqu.Ex{}
}

func ToSQL(n rql.Node) string {
	driver, _, _ := sqlmock.New()
	db := goqu.New("default", driver)
	e := ToGoqu(n)
	if reflect.DeepEqual(e, goqu.Ex{}) {
		sql, _, _ := db.From("test").ToSql()
		return sql
	}
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
