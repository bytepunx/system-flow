---
id: D-0539
title: The fenced quote helpers in digest/seed.go and mcpserver/permission.go stay apart from mdlint.Fenced, and folding them in is a cleanup story of its own
kind: decision
sources: [S-0614, TH-0923]
topics: [markdown-lint, flai, code]
paths: [flai/internal/mdlint/quote.go, flai/internal/digest/seed.go, flai/internal/mcpserver/permission.go]
date: 2026-11-09
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0539 The fenced quote helpers in digest/seed.go and mcpserver/permission.go stay apart from mdlint.Fenced, and folding them in is a cleanup story of its own

S-0614 added `mdlint.Fenced` (`flai/internal/mdlint/quote.go`) so the merge queue's and the pull request's send-back threads quote tool output lint-safe: tabs as spaces, control characters dropped, blank runs folded. Its plan, on TH-0923, left the two other `fenced` helpers where they are:

- `fenced` in `flai/internal/digest/seed.go` is already lint-safe.
- `fenced` in `flai/internal/mcpserver/permission.go` turns MD010 off around its block on purpose, so the permission prompt quotes a write's text with its tabs as they are. `mdlint.Fenced` would make them spaces.

The orchestrator approved the plan on TH-0923: "A queue that holds itself after repeated refusals, or one shared quote helper, would each be a story of its own." A story that wants one shared quote helper is that cleanup story. It must keep the permission prompt's tabs verbatim.
