---
id: D-0269
title: flai's markdown lint does not split a table row on pipes inside a code span, as markdownlint does
kind: cause
sources: [I-0077]
topics: [flai, cli, go, code, all]
paths: [flai/internal/mdlint/mdlint.go, flai/internal/mdlint/rules.go, flai/internal/mdlint/doc.go, flai/internal/mdlint/mdlint_test.go, flai/internal/mdlint/testdata/cases/tables.md, flai/internal/mdlint/testdata/cases/expected.txt, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0077-flai-s-markdown-lint-does-not-split-a-table-row-on-pipes-inside-a-code-span-as-markdownlint-does.md, design/issues/summary.md]
date: 2026-10-09
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0269 flai's markdown lint does not split a table row on pipes inside a code span, as markdownlint does

I-0077 was closed on 2026-10-09.

The cause:

```text
flai's markdown lint does not split a table row on pipes inside a code span, as markdownlint does
```

The fix:

```text
S-0289: flai's mdlint now has MD056, table column count, so a row whose unescaped pipes, those inside a code span among them, give it more cells than its header is reported as markdownlint-cli2 reports it (TestTableColumnCountOfI0077, testdata/cases/tables.md).
```

Source: I-0077.
