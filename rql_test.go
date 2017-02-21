package rql_test

import "git.nwaonline.com/rune/rql"

var exp_equal_int = &rql.Operator{
	Kind: rql.EQ,
	Operands: rql.ExpressionList{
		&rql.Identifier{Kind: rql.IDENT, Name: "column"},
		&rql.Literal{
			Kind:  rql.NUMERIC,
			Value: "12",
		},
	},
}

var exp_nequal_int = &rql.Operator{
	Kind: rql.NE,
	Operands: rql.ExpressionList{
		&rql.Identifier{Kind: rql.IDENT, Name: "my_col"},
		&rql.Literal{
			Kind:  rql.NUMERIC,
			Value: "-12",
		},
	},
}

var exp_equal_string = &rql.Operator{
	Kind: rql.EQ,
	Operands: rql.ExpressionList{
		&rql.Identifier{Kind: rql.IDENT, Name: "my_col"},
		&rql.Literal{
			Kind:  rql.STRING,
			Value: `"this is a test"`,
		},
	},
}

var exp_equal_boolean = &rql.Operator{
	Kind: rql.EQ,
	Operands: rql.ExpressionList{
		&rql.Identifier{Kind: rql.IDENT, Name: "my_col"},
		&rql.Literal{
			Kind:  rql.BOOLEAN,
			Value: "true",
		},
	},
}

var exp_equal_identifier = &rql.Operator{
	Kind: rql.EQ,
	Operands: rql.ExpressionList{
		&rql.Identifier{Kind: rql.IDENT, Name: "my_col"},
		&rql.Identifier{Kind: rql.IDENT, Name: "my_other_col"},
	},
}

var exp_equal_float = &rql.Operator{
	Kind: rql.EQ,
	Operands: rql.ExpressionList{
		&rql.Identifier{Kind: rql.IDENT, Name: "my_col"},
		&rql.Literal{
			Kind:  rql.NUMERIC,
			Value: "12.123",
		},
	},
}

var exp_equal_negative_float = &rql.Operator{
	Kind: rql.EQ,
	Operands: rql.ExpressionList{
		&rql.Identifier{Kind: rql.IDENT, Name: "my_col"},
		&rql.Literal{
			Kind:  rql.NUMERIC,
			Value: "-12.123",
		},
	},
}

var exp_and = &rql.Operator{
	Kind: rql.AND,
	Operands: rql.ExpressionList{
		exp_equal_int,
		exp_nequal_int,
	},
}

var exp_or = &rql.Operator{
	Kind: rql.OR,
	Operands: rql.ExpressionList{
		exp_equal_int,
		exp_nequal_int,
	},
}

var exp_many_nested = &rql.Operator{
	Kind: rql.AND,
	Operands: rql.ExpressionList{
		exp_equal_int,
		exp_equal_string,
		exp_equal_float,
		exp_and,
		exp_or,
	},
}

var exp_not = &rql.Operator{
	Kind: rql.NOT,
	Operands: rql.ExpressionList{
		exp_and,
	},
}

var exp_in = &rql.Operator{
	Kind: rql.IN,
	Operands: rql.ExpressionList{
		&rql.Identifier{Kind: rql.IDENT, Name: "primes"},
		&rql.ExpressionList{
			&rql.Literal{Kind: rql.NUMERIC, Value: "1"},
			&rql.Literal{Kind: rql.NUMERIC, Value: "2"},
			&rql.Literal{Kind: rql.NUMERIC, Value: "3"},
			&rql.Literal{Kind: rql.NUMERIC, Value: "5"},
			&rql.Literal{Kind: rql.NUMERIC, Value: "7"},
		},
	},
}

var exp_null = &rql.Operator{
	Kind: rql.EQ,
	Operands: rql.ExpressionList{
		&rql.Identifier{Kind: rql.IDENT, Name: "my_col"},
		&rql.Literal{Kind: rql.NULL, Value: "null"},
	},
}

var exp_quoted_string = &rql.Operator{
	Kind: rql.EQ,
	Operands: rql.ExpressionList{
		&rql.Identifier{Kind: rql.IDENT, Name: "col"},
		&rql.Literal{Kind: rql.STRING, Value: `"a string with "quotes" in it"`},
	},
}
