# Code spans

I-0072: a task body named a code span with a space before its closing
backtick, which markdownlint reports as MD038.

- Run `- Trigger: ` here.
- Padded ` a ` is fine.
- Double padded `  a  ` is not, at either end.
- Lead ` a` and trail `a ` are not.
- Backticks `` `x` `` are fine, ``  `x` `` is not, `` `x`` is fine.
- Back ``a` `` one, `` ` `` two, ``  a` `` three.
- Only spaces `   ` and ` ` are fine, and so is `` ` `` and `a`.
- Tab `	x` is not.
- Tab trailing `x	` is not, nor is a no-break space ` x` here.

Across `a
b ` lines, and ` a
b` lines, and ` a
b ` lines.

Across `  a
b` lines and `a
  b` lines and `a
  b  ` lines.

Ends `foo
` here and `
foo ` here and `  
foo` here.

Start `
  b` keeps the indentation of a continuation line.

Para `
      e` six spaces.

Para `x
  ` closing alone.

Para `
  y ` padded.

- Item `
  c` here, the container takes the indentation.

- Item `
    c` deeper.

- Lazy item text `
c` lazy.

- Lazy two `
 d` lazy.

1. Ordered `o ` item
   continued `c ` text

## Heading ` x` here

## Heading `x` clean

| A | B |
| --- | --- |
| `x ` | ok `y` |

```text
`x ` in a fence
```

<!-- `x ` in a comment -->

Text <!-- `y ` inline comment --> more.

[link ` t`](https://example.com) and ![img `u `](x.png)

Emphasis *`e `* and **` f`** five.

Unclosed ``` `` x ` y` six.

Hard `x  
y` seven.

Escaped \` opens no code span, and `x``y ` holds a run of two.

I-0363: escaped backticks in a code span, from D-0592.

`storygit`'s sync result carries JSON tags (`RegeneratedAtHead []string \`json:"regenerated_at_head"\`` in `flai/internal/storygit/sync.go`). But`flai stream sync --json` does not marshal that struct. `flai/cmd/stream.go` builds its own `map[string]any`, one for a stopped sync and one for a finished one, naming each key. A field added to the result alone never reaches the command's output.

As the curator first wrote it, which markdownlint reports too:

`storygit`'s sync result carries JSON tags (`RegeneratedAtHead []string \`json:"regenerated_at_head"\`` in `flai/internal/storygit/sync.go`). But `flai stream sync --json` does not marshal that struct. `flai/cmd/stream.go` builds its own `map[string]any`, one for a stopped sync and one for a finished one, naming each key. A field added to the result alone never reaches the command's output.
