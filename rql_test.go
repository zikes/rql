package rql_test

import "git.nwaonline.com/rune/rql"

var exp_equal_int = &rql.Operator{
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

var exp_nequal_int = &rql.Operator{
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

var exp_equal_string = &rql.Operator{
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

var exp_equal_boolean = &rql.Operator{
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

var exp_equal_identifier = &rql.Operator{
	Name: "eq",
	Operands: &rql.ParenExpr{
		Expressions: []rql.Expression{
			&rql.Identifier{Name: "my_col"},
			&rql.Identifier{Name: "my_other_col"},
		},
	},
}

var exp_equal_float = &rql.Operator{
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

var exp_equal_negative_float = &rql.Operator{
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

var exp_and = &rql.Operator{
	Name: "and",
	Operands: &rql.ParenExpr{
		Expressions: []rql.Expression{
			exp_equal_int,
			exp_nequal_int,
		},
	},
}

var exp_or = &rql.Operator{
	Name: "or",
	Operands: &rql.ParenExpr{
		Expressions: []rql.Expression{
			exp_equal_int,
			exp_nequal_int,
		},
	},
}

var exp_many_nested = &rql.Operator{
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

var exp_not = &rql.Operator{
	Name: "not",
	Operands: &rql.ParenExpr{
		[]rql.Expression{
			exp_and,
		},
	},
}

var exp_in = &rql.Operator{
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

var exp_null = &rql.Operator{
	Name: "eq",
	Operands: &rql.ParenExpr{
		[]rql.Expression{
			&rql.Identifier{Name: "my_col"},
			&rql.Literal{Kind: rql.NULL, Value: "null"},
		},
	},
}
