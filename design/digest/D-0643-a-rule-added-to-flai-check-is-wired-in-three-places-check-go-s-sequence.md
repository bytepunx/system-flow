---
id: D-0643
title: A rule added to flai check is wired in three places, check.go's sequence, concurrent_test.go's oldOrder, and a file in everyRuleProject that makes it find something, and an item rule is called from the per-item pass as well
kind: decision
sources: [S-0741]
topics: [flai, go, code]
paths: [flai/internal/check/check.go, flai/internal/check/concurrent_test.go, flai/internal/check/check_items.go, flai/internal/check/check_workflow.go]
date: 2026-11-15
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0643 A rule added to flai check is wired in three places, check.go's sequence, concurrent_test.go's oldOrder, and a file in everyRuleProject that makes it find something, and an item rule is called from the per-item pass as well

S-0741 added `checker.workflow` to `flai check` (ADR-0301 §2, §3, §6). Its task sub-agent wrote the rule but stopped short of wiring it, and the story's agent finished the job. From its narrative's Decisions (`wip/archive/agents/S-0741.md`):

> T-3204's sub-agent stopped short of wiring the check: the story's agent added checker.workflow to sequence and to concurrent_test.go's oldOrder and every-rule project, and itemStates to the per-item pass, so workflow.item covers the written item on a checked write while the workflow's own findings skip it.

So a story that adds a rule to `flai check` does four things:

- Add it to `sequence` in `flai/internal/check/check.go`, the list of rules `Run` runs.
- Add it to `oldOrder` in `flai/internal/check/concurrent_test.go`. That list is written out by hand on purpose, so a rule dropped from `sequence`, or moved in it, fails the test.
- Add a file to `everyRuleProject` in the same test so that the new rule finds something.
- When the rule judges one item, as `workflow.item`'s `itemStates` does, call it from the per-item pass in `check_items.go` as well. A checked write then reports it for the written item (ADR-0264). The whole-project findings do not run on an item write.

A task sub-agent given the rule alone can miss these steps. Name them in the task's `## Done when`.

Source: S-0741's narrative, `wip/archive/agents/S-0741.md`, § Decisions.
