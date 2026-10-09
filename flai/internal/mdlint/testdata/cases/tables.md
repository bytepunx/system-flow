# Tables

I-0077: a row's pipes split it into cells inside a code span too, so a
code span with a pipe in it gives the row more cells than its header,
which markdownlint reports as MD056.

| Column | Other |
|--------|-------|
| a | b |
| `a\|b` | escaped |
| x \| y |
| x \\| y |
| `a|b|c` | x |
| `a | b` | x |
| one |
| a | b | c |
||
|
a | b

| No | Pipes |
| --- | --- |
a line without a pipe is a row of one cell

Text the table interrupts.
| A | B |
|---|---|
| 1 | 2 | 3 |
- an item ends the table | a | b | c |

| A | B |
|---|---|
## A heading ends it | a | b | c |

| Head | Count | Differ |
|------|-------|
| a | b | c | d |

> | A | B |
> |---|---|
> | 1 | 2 | 3 |
> | 1 |

- | A | B |
  | --- | --- |
  | 1 | 2 | 3 |
  | 1 |
lazy | is | not | a | row

| A | B |
|---|---|
    | indented | is | code |

| A | B |
| :-- | --: |
| `x` | `y | z` |
| `` a|b `` | ` c|d ` |
