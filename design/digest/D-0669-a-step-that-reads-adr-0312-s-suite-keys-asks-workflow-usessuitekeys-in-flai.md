---
id: D-0669
title: A step that reads ADR-0312's suite keys asks Workflow.UsesSuiteKeys in flai/internal/workflow/suites.go, and reads a workflow file for which it is false as giving the baseline's keys
kind: decision
sources: [S-0806, T-3608, TH-1289]
topics: [workflow, verify, flai, go, code]
paths: [flai/internal/workflow/suites.go, flai/internal/check/check_workflow.go, design/system/workflow.md]
date: 2026-11-18
status: active
absorbed_by: ""
asked_again: []
cited: [S-0807, S-0809, S-0812]
---

# D-0669 A step that reads ADR-0312's suite keys asks Workflow.UsesSuiteKeys in flai/internal/workflow/suites.go, and reads a workflow file for which it is false as giving the baseline's keys

S-0806 built the rule the operator confirmed on TH-1289 into one method: a workflow file that gives none of `verify`, `runs`, and `verify_at` reads as the baseline's (`design/system/workflow.md` § Suites a step needs and runs). Its narrative's Decisions (`wip/archive/agents/S-0806.md`) record where the rule lives:

> T-3608 (TH-1289): the rule lives in Workflow.UsesSuiteKeys, so that S-0807, S-0808, S-0809, and S-0812 can read a file with none of the keys as the baseline's.

`Workflow.UsesSuiteKeys()` in `flai/internal/workflow/suites.go` reports whether the file gives any of the three keys. `flai/internal/check/check_workflow.go` already calls it to skip `workflow.suite-unnamed`. A story that makes the move, the close-out, the acceptance's gates, the batch, a state's runs, the publish, or `flai test --at ci` read the keys calls this method and falls back to the baseline's keys when it is false. It does not write a second test of its own.

Source: S-0806, T-3608, TH-1289.
