---
id: D-0657
title: "S-0817 is ready to accept, but the same uncommitted digest entry on main that holds S-0742 fails lint (MD028): I fixed it"
kind: answer
sources: [TH-1269, TH-1270]
topics: [merge, sync, acceptance, digest, lint, markdown, flai, cli, go, code, all]
paths: [design/digest, design/adrs/0248-every-acceptance-commits-the-main-checkout-s-uncommitted-files-under-wip-design.md, flai/internal/storygit/mdmerge.go, design/system/workflow.md]
date: 2026-11-17
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0657 S-0817 is ready to accept, but the same uncommitted digest entry on main that holds S-0742 fails lint (MD028): I fixed it

A digest entry whose body holds two blockquotes with only a blank line between them fails markdownlint's MD028 (no blank line inside a blockquote), which the markdown tier's fix cannot fix. Since every acceptance commits the main checkout's uncommitted files under `design/` ([ADR-0248](../adrs/0248-every-acceptance-commits-the-main-checkout-s-uncommitted-files-under-wip-design.md)), such an entry stops every acceptance until someone edits it in the main checkout. The agents that meet it cannot edit files, so the operator must.

It happened with D-0653, written by the curator for S-0739. The orchestrator asked on TH-1269 (S-0742), 2026-11-17:

```text
It holds two blockquotes with a blank line between them, which markdownlint flags and its `--fix` cannot fix:
design/digest/D-0653-…md:20 error MD028/no-blanks-blockquote Blank line inside blockquote
S-0742 cannot be accepted: the main checkout's uncommitted files hold lint the fix of the markdown tier could not fix … correct them in the main checkout and run flai accept S-0742 again
```

The same file then held S-0817 on TH-1270. The operator, alex, answered there at 2026-11-17T04:40:43Z: "I fixed it", and S-0742 and S-0817 were accepted.

So an entry that quotes two passages puts text between them, such as "It adds:", or joins them into one quote whose blank line holds only `>`.

Sources: TH-1269, TH-1270.
