---
id: D-0162
title: flai's wip markdown lint does not flag an ordered list item numbered from other than 1 inside a blockquote, so a thread entry reached main and failed a story's close-out
kind: cause
sources: [I-0070]
topics: [flai, cli, go, code, all]
paths: [flai/internal/mdlint, design/issues, design/system/flai-cli.md]
date: 2026-10-05
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0162 flai's wip markdown lint does not flag an ordered list item numbered from other than 1 inside a blockquote, so a thread entry reached main and failed a story's close-out

I-0070 was closed on 2026-10-05.

The cause:

```text
flai's wip markdown lint does not flag an ordered list item numbered from other than 1 inside a blockquote, so a thread entry reached main and failed a story's close-out
```

The fix:

```text
S-0258 T-0860: flai's lint (flai/internal/mdlint/doc.go) parses a blockquote's content as blocks, with blockquotes as containers beside list items, so MD029 flags a quoted ordered list numbered from other than 1, as markdownlint-cli2 does, and the other rules judge quoted content too. testdata/cases/blockquote.md reproduces TH-0101's '> 3.' and '> 6.' entry, and TestQuotedOrderedListOfI0070 shows a thread entry bringing one is refused.
```

Source: I-0070.
