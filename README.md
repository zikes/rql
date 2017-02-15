# rql

RQL is a simple, composable query language that's easy to read and write, without being
ambiguous.

Operators, comparators, and functions all take the form `expressionName(operand, ...)`.
Operands may also be expressions, allowing you to nest as deeply as necessary.

```rql
or(eq(my_column,12),and(gt(my_column,20),lt(my_column,30))
```

In SQL, the above would be

```sql
[...] WHERE my_column=12 OR (my_column > 20 AND my_column < 30)
```
