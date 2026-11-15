---
id: D-0631
title: A workflow rule's name, types, and default severity live in flai/internal/workflow/rules.go and its check in flai/internal/workitem/rules.go, so the workflow package reads no work item
kind: decision
sources: [S-0738]
topics: [flai, go, code]
paths: [flai/internal/workflow/rules.go, flai/internal/workitem/rules.go]
date: 2026-11-15
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0631 A workflow rule's name, types, and default severity live in flai/internal/workflow/rules.go and its check in flai/internal/workitem/rules.go, so the workflow package reads no work item

S-0738 named the rules `Repo.Move` enforces as a vocabulary flai implements (ADR-0301 §4), and split that vocabulary across two packages. Its narrative's Decisions (`wip/archive/agents/S-0738.md`):

> The vocabulary (names, types, default severity) lives in flai/internal/workflow/rules.go; the checks live in flai/internal/workitem/rules.go, since workitem imports workflow and the checks read items, narratives, and worktrees. ADR-0301 §4 puts the rules in workflow/rules.go; the split keeps workflow free of work items, as its package comment says.

So a story that adds a rule to the vocabulary, or changes what one checks, names it in `flai/internal/workflow/rules.go` and writes its check in `flai/internal/workitem/rules.go`. ADR-0301 §4 and `design/system/workflow.md` name only the first file.

Source: S-0738's narrative, `wip/archive/agents/S-0738.md`, § Decisions.
