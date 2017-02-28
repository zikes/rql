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

## Operators

| name  | usage                 | description                        |
|-------|-----------------------|------------------------------------|
| `eq`  | `eq(value, value)`    | Test equality.                     |
| `neq` | `neq(value, value)`   | Test inequality.                   |
| `and` | `and(value, value)`   | Logical and.                       |
| `or`  | `or(value, value)`    | Logical or.                        |
| `not` | `not(value)`          | Logical not.                       |
| `lt`  | `lt(value, value)`    | Less than comparison.              |
| `gt`  | `gt(value, value)`    | Greater than comparison.           |
| `le`  | `le(value, value)`    | Less than or equals comparison.    |
| `ge`  | `ge(value, value)`    | Greater than or equals comparison. |
| `in`  | `in(col, (1,2,3))`    | Check if value is one of a series. |

## Literals

| name       | usage                 | description                        |
|------------|-----------------------|------------------------------------|
| string     | `eq(value, "string")` | String literal data type.          |
| numeric    | `neq(value, -12.123)` | Numeric literal data type.         |
| boolean    | `and(value, true)`    | Boolean literal data type.         |
| null       | `or(value, null)`     | Null literal data type.            |
| identifier | `not(my_col)`         | Identifier literal data type.      |

## Arrays

| name           | usage                 | description                                                    |
|----------------|-----------------------|----------------------------------------------------------------|
| ExpressionList | `(val,val,val,[...])` | Collects a list of values and presents them as a single value. |
