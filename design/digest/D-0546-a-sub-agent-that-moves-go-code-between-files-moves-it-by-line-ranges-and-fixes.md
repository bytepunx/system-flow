---
id: D-0546
title: A sub-agent that moves Go code between files moves it by line ranges and fixes the imports with golangci-lint fmt, since no goimports is installed and the project's golangci config runs goimports as a formatter
kind: decision
sources: [S-0631]
topics: [go, code, delegation]
paths: [flai/.golangci.yaml]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0546 A sub-agent that moves Go code between files moves it by line ranges and fixes the imports with golangci-lint fmt, since no goimports is installed and the project's golangci config runs goimports as a formatter

S-0631 split six Go files under `flai/` that were over 2,000 lines, through task sub-agents. From its Decisions (`wip/archive/agents/S-0631.md`):

```text
Sub-agents move code with line ranges and `golangci-lint fmt` for imports, not by retyping it: no `goimports` is installed, and the project's golangci config runs goimports as a formatter.
```

Moving by line ranges kept each move reviewable: comparing each split's code lines before and after, outside the import blocks, showed the moves changed nothing else. So a later story that moves code between files briefs its sub-agents the same way, and reviews the move by that comparison.

From the same story: each split ran in its own task worktree, since `cmd`, `serve`, and `hostapi` import the packages the other splits moved, and a half-done move in one would fail a sibling's build (`delegation.md`).

Sources: S-0631.
