---
id: T-1433
type: task
nature: remediation
title: flai's markdown lint reports MD056, table column count, on a row whose unescaped pipes, those inside a code span among them, give it more or fewer cells than its header
status: done
parent: S-0289
owner: alex
created: 2026-10-09T18:22:28Z
updated: 2026-10-09T18:28:59Z
transitions:
  - to: ready
    at: 2026-10-09T18:23:54Z
    by: agent-S-0289
  - to: in-progress
    at: 2026-10-09T18:23:54Z
    by: agent-S-0289
  - to: done
    at: 2026-10-09T18:28:59Z
    by: agent-S-0289
stream: S-0289
tags: [flai]
touches: [flai/internal/mdlint/mdlint.go, flai/internal/mdlint/rules.go, flai/internal/mdlint/doc.go, flai/internal/mdlint/mdlint_test.go, flai/internal/mdlint/testdata/cases/tables.md, flai/internal/mdlint/testdata/cases/expected.txt, design/system/flai-cli.md, docs/users/flai.md]
usage:
  source: log
  seconds: 305
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 61
      output: 20628
      cache_read: 2689368
      cache_write: 144343
      cost: 1.9346
---
# T-1433 flai's markdown lint reports MD056, table column count, on a row whose unescaped pipes, those inside a code span among them, give it more or fewer cells than its header

## Work

I-0077's row, in `design/system/flai-cli.md` before commit adf27c49, wrote the merge-base marker as a code span of seven unescaped pipes in a table cell. markdownlint-cli2 0.20.0 splits a row on every unescaped pipe, inside a code span too, as GFM does. It reported MD056, table column count, and MD038 on the code span the split left. flai's `mdlint` already splits cells that way in `doc.go`'s `cells`, and reproduces the MD038 on the row today. It has no MD056, so the row passed. A reproduction of the row on a two-column table gives:

```text
markdownlint-cli2: 5 MD038 (Context "`, or `"), 5 MD056 (Expected: 2; Actual: 9; Too many cells, extra data will be missing)
flai mdlint:       5 MD038 (Context "`, or `")
```

Add MD056 (`table-column-count`, tags `table`) to `mdlint`'s rules in `mdlint.go` and implement it in `rules.go`:

- Take each table's column count from its header row, and count each row's cells as `cells` splits them: an unescaped pipe splits, a backslash-escaped one does not, and a leading and a trailing pipe are left out.
- Report each body row with more or fewer cells, on its line, with markdownlint's detail: `Expected: n; Actual: m; Too many cells, extra data will be missing` or `Too few cells, row will be missing data`.
- Record which rows belong to which table in `doc.go`, where the table is recognised.

Add a `tables.md` fixture under `testdata/cases` with I-0077's row, a row with an escaped pipe inside a code span, a row with too few cells, a table without leading and trailing pipes, and a clean table. Regenerate `expected.txt` with `make mdlint-fixtures`, so that `TestFixturesMatchMarkdownlint` compares flai with markdownlint-cli2 on it. Add a case for I-0077's row to `mdlint_test.go`.

Name MD056 where the rules flai lints are listed: the `flai check` row and the `mdlint/` line of the layout in `design/system/flai-cli.md`, and the paragraph on markdown lint in `docs/users/flai.md`, which names what flai checks ("tables" joins the list). ADR-0061 is accepted and stays as it is, as for MD007 and MD038.

## Done when

- `mdlint` reports MD056 on I-0077's row, with the detail markdownlint-cli2 gives, and `TestFixturesMatchMarkdownlint` passes with `tables.md` and the regenerated `expected.txt`.
- An escaped pipe in a code span splits no cell, and a table that agrees with its header brings no finding.
- `flai check --strict` on the main checkout reports no `markdown.MD056` that markdownlint-cli2 does not, and `TestRepositoryLintsClean` passes.
- `design/system/flai-cli.md` and `docs/users/flai.md` name MD056 among the rules flai lints.
- `flai test` on the changed paths passes.

## Notes

Docs ship in this task, not a later one, because `documentation.md` asks for them in the same commit as the behaviour. MD055 (table pipe style) and MD058 (blanks around tables) are left out: I-0077 names neither.
