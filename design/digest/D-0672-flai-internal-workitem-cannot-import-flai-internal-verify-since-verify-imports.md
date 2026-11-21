---
id: D-0672
title: flai/internal/workitem cannot import flai/internal/verify, since verify imports workitem, so Repo.Move reads a story's verify record through the leaf package flai/internal/verifyrecord, which both read
kind: decision
sources: [S-0807, T-3454, T-3457]
topics: [verify, workflow, go, code, flai]
paths: [flai/internal/verifyrecord, flai/internal/workitem/rules.go, flai/internal/verify]
date: 2026-11-18
status: active
absorbed_by: ""
asked_again: []
cited: [S-0809, TH-1316]
---

# D-0672 flai/internal/workitem cannot import flai/internal/verify, since verify imports workitem, so Repo.Move reads a story's verify record through the leaf package flai/internal/verifyrecord, which both read

S-0807 made a story's move by a transition with `verify` run the rule `verified`, which reads the verify record at the head of the story's branch (`design/system/workflow.md` § What the move and the close-out read (S-0807)). The rule's check lives in `flai/internal/workitem/rules.go` (D-0631), but `workitem` cannot call the `verify` package, which imports `workitem` and `context`. S-0807's narrative's Decisions (`wip/archive/agents/S-0807.md`):

> D-0426: verify imports context and workitem, so Repo.Move reads the record through the new leaf package verifyrecord, never through verify.
>
> Plan layers: (1) T-3454, the per-suite verify record in a new leaf package flai/internal/verifyrecord, which both the rule and the selection read

So the record's per-suite shape (`suite_runs`), `Merge`, `AllPassed`, and `AtHead` live in `flai/internal/verifyrecord`, which imports neither `verify` nor `workitem`. A package below `verify` that needs a story's stored result, such as the acceptance's gates (S-0808) or `flai/internal/context` (D-0426, which read the file itself before this package existed), reads it through `verifyrecord`.

`verifyrecord.AtHead` answers a `Verdict` per suite: `AtHead` (the entry is of the head, in a clean run), `Held` (tiers that passed, were reused, or were skipped for one that was), `Queued` (tiers left to the merge queue), and `Passed`, which is true only when every tier is held, so a queued tier makes it false. T-3457's review fix: "a tier left to the merge queue counts as held for the verified rule on a move, as the close-out passes without it (ADR-0174). verifyrecord.Verdict gained AtHead and Queued, so that S-0808 can still be strict at the acceptance." The move's rule counts `Held` and `Queued` together; a gate that must not count a queued tier reads `Passed`.

Source: S-0807, T-3454, T-3457.
