---
id: D-0275
title: "S-0363 plan: task_done refuses on an outdated flai (ADR-0136), four tasks in three layers: Recommendation: take the planner's softer form"
kind: answer
sources: [TH-0401]
topics: [flai, cli, go, code, all]
paths: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/internal/mcpserver/task_test.go, flai/cmd/task_done.go, docs/users/flai-reference.md, design/adrs/README.md, design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md, design/issues/I-0127-an-outdated-mcp-flai-s-task-done-commits-every-changed-path-under-the-closing-task-so-a-layer-s-tasks-cannot-close-apart.md, design/issues/summary.md, design/adrs/0142-flai-task-done-refuses-to-close-a-task-while-the-running-flai-s-major-or-minor.md]
date: 2026-10-12
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0275 S-0363 plan: task_done refuses on an outdated flai (ADR-0136), four tasks in three layers: Recommendation: take the planner's softer form

The operator, alex, answered TH-0401 at 2026-10-12T15:29:56Z.

What was asked:

```text
S-0363 is planned. It has 4 draft tasks, 12 file touches, a 20m forecast, and a cost of delay of 12.50 USD a week. The story's `### Planning` note says where each touch and figure came from.
```

The answer:

```text
**Recommendation:** take the planner's softer form. `task_done` refuses only while the running flai's major or minor version is behind the project's history, and warns on a patch behind. Keep the rest of the plan as drafted.

Why this is yours: ADR-0136 would be a new decision. It departs from ADR-0117 § Decision 7, under which an outdated flai warns and is never refused. A recorded decision is not mine to make.

Why the softer form:

- Under the planner's form, a host one patch behind stops every story's agent at its next task close. That lasts until `flai host upgrade` runs.
- Every release puts the host a patch behind. Today the host runs 1.40.1, and I have published 1.40.3 and 1.40.4 since. So under the strict form, each publish would stop all running stories until you upgrade.
- I-0127's instance ran 1.39.13 against 1.40.0 in the history, so a minor difference. The softer form still refuses it.
- The cost: a patch behind can still lack a fix, as 1.39.13 lacked S-0322's in 1.39.14. The inbox's `flai_outdated` warning keeps covering that case.

If you take this, the planner re-plans T-1439 and T-1440 to the minor-version threshold. The four tasks and three layers stay the same. If you prefer the strict form, confirm the plan as drafted.

Either way, I finalize S-0363 once you answer. The figures are flai's (20m, 12.50 USD a week), and its 12 file touches fit the criteria.

This is a recommendation, not an answer: confirm it, or answer otherwise.

Source: design/adrs/0117-the-host-api-installs-a-release-the-dashboard-names-only-when-flai-lists-it-as.md § Decision
```

Source: TH-0401.
