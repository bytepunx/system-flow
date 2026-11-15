---
id: D-0629
title: A refusal of a move lists the allowed moves in the workflow file's transition order, so a task in in-progress lists review, done, ready, cancelled
kind: decision
sources: [S-0737]
topics: [flai, go, code]
paths: [flai/internal/workitem/rules.go, flai/internal/workflow/baseline/board.md]
date: 2026-11-14
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0629 A refusal of a move lists the allowed moves in the workflow file's transition order, so a task in in-progress lists review, done, ready, cancelled

S-0737 made `Repo.Move` take its transition map from the workflow file, and one refusal's words changed with it. Its narrative's Decisions:

> One refusal's words change. A task in in-progress lists its allowed moves in the workflow file's order, review, done, ready, cancelled; the old map gave review, ready, cancelled, done. No test depends on the order. Reordering the baseline file to keep the old order would have changed ADR-0301 §7's text.

So the order of the moves a refusal names follows the order of `transitions` in `flai/internal/workflow/baseline/board.md`, or a project's own file, not a list in code. A test that asserts a refusal's words follows the file, and the baseline's order stays as ADR-0301 §7 gives it.

Source: S-0737's narrative, `wip/archive/agents/S-0737.md`, § Decisions.
