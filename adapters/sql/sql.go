package sqladapter

import (
	"fmt"
	"strings"

	rql "github.com/zikes/rql/parse"
)

func ToSQL(n rql.Node) string {
	switch n := n.(type) {
	case *rql.StatementNode:
		return ToSQL(n.Operator)
	case *rql.BoolNode:
		return fmt.Sprintf("%v", n.True)
	case *rql.NullNode:
		return "NULL"
	case *rql.IdentifierNode:
		return n.Ident
	case *rql.StringNode:
		return n.Quoted
	case *rql.NumberNode:
		switch {
		case n.IsInt:
			return fmt.Sprintf("%d", n.Int64)
		case n.IsUint:
			return fmt.Sprintf("%d", n.Uint64)
		case n.IsFloat:
			return fmt.Sprintf("%g", n.Float64)
		}
	case *rql.ListNode:
		str := []string{}
		for _, v := range n.Nodes {
			str = append(str, ToSQL(v))
		}
		return "(" + strings.Join(str, ", ") + ")"
	case *rql.OperatorNode:
		left := n.Operands.Nodes[0]
		right := n.Operands.Nodes[1]
		switch n.Operator {
		case "eq":
			if right.Type() == rql.NodeNull {
				return ToSQL(left) + " IS " + ToSQL(right)
			}
			return ToSQL(left) + " = " + ToSQL(right)
		case "ne":
			if right.Type() == rql.NodeNull {
				return ToSQL(left) + " IS NOT " + ToSQL(right)
			}
			return ToSQL(left) + " != " + ToSQL(right)
		case "gt":
			return ToSQL(left) + " > " + ToSQL(right)
		case "lt":
			return ToSQL(left) + " < " + ToSQL(right)
		case "ge":
			return ToSQL(left) + " >= " + ToSQL(right)
		case "le":
			return ToSQL(left) + " <= " + ToSQL(right)
		case "in":
			return ToSQL(left) + " IN " + ToSQL(right)
		case "and":
			str := []string{}
			for _, v := range n.Operands.Nodes {
				str = append(str, ToSQL(v))
			}
			return "(" + strings.Join(str, " AND ") + ")"
		case "or":
			str := []string{}
			for _, v := range n.Operands.Nodes {
				str = append(str, ToSQL(v))
			}
			return "(" + strings.Join(str, " OR ") + ")"
		}
	}
	return ""
}
