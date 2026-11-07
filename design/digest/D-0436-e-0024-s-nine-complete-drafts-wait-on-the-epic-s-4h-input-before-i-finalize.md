---
id: D-0436
title: "E-0024's nine complete drafts wait on the epic's 4h input before I finalize them: S-0414 to S-0422: finalize as they are"
kind: answer
sources: [TH-0587]
topics: [cli, orchestration, planning, template, conventions, flai, go, flaiover, dashboard, sveltekit, code, all]
paths: [flai/internal/serve, flai/internal/guard, flai/internal/context, flai/internal/hostapi, flai/cmd, flaiover/src, design/system/strategic-agents.md, design/system/metrics.md, design/conventions/strategic-agents.md, ".claude/agents", docs/operators/index.md, docs/users/flai.md, template]
date: 2026-10-31
status: active
absorbed_by: ""
asked_again: []
cited: [TH-0587, TH-0725, TH-0845]
---

# D-0436 E-0024's nine complete drafts wait on the epic's 4h input before I finalize them: S-0414 to S-0422: finalize as they are

The operator let complete drafts be finalized with cost of delay values their planners worked out by hand. Setting the epic's input first, so that `flai cod` could value each story, was not required.

What was asked: on TH-0587 the orchestrator held E-0024's nine complete drafts, S-0414 to S-0422, under `finalize_drafts`. Their criteria, touches, and forecasts fit the approved plans. Their values did not come from `flai cod`:

- `flai cod` refuses while the epic has no input. The planners derived each value by hand from the 4h per cycle the operator had confirmed on TH-0567 and TH-0473.
- flai guard refused the orchestrator the epic's input (I-0187, since remedied by S-0523 and `orchestration.permissions.set_cost_of_delay`, ADR-0224).
- The order was by throughput, so the values changed no pull order. They fed only the release figures.

The orchestrator recommended running `flai edit E-0024 --time-lost-per-cycle 4h` first, then finalizing after the planner's next run. The alternative was to finalize as they were.

The answer: on 2026-10-31T07:52:39Z the operator, alex, answered "finalize as they are".

So when a complete draft's only gap is a value worked out by hand from an input the operator already decided, the operator does not hold its finalization for `flai cod`.

Source: TH-0587.
