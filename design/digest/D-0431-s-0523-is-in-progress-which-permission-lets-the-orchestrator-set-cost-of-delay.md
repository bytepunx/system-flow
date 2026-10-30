---
id: D-0431
title: "S-0523 is in progress: which permission lets the orchestrator set cost of delay inputs on an epic or a started story (TH-0694 question 2): Recommended for the operator: option 1. Add a new…"
kind: answer
sources: [TH-0764]
topics: [orchestration, planning, flai, cli, go, template, conventions, code, all]
paths: [design/adrs, design/adrs/README.md, design/system/strategic-agents.md, design/system/flai-cli.md, design/system/project-manifest.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, template/CHANGELOG.md, docs/users/flai.md, docs/users/flai-reference.md, docs/users/flaiover.md, docs/operators/settings.md, docs/operators/index.md, flai/internal/manifest/settings.go, flai/internal/guard/guard.go, flai/internal/guard/orchestrate_test.go, flai/internal/guard/guard_test.go, flai/internal/workitem/promotable.go, flai/internal/mcpserver/items_write.go, flai/internal/mcpserver/items_write_test.go, flai/cmd/edit.go, flai/cmd/edit_test.go, flai/cmd/orchestrate_permissions_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/issues/I-0187-the-orchestrator-cannot-set-a-cost-of-delay-input-the-operator-decided-on-an-epic-or-a-story-past-the-backlog.md, design/issues/summary.md]
date: 2026-10-30
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0431 S-0523 is in progress: which permission lets the orchestrator set cost of delay inputs on an epic or a started story (TH-0694 question 2): Recommended for the operator: option 1. Add a new…

The operator, alex, answered TH-0764 at 2026-10-30T16:47:47Z.

What was asked:

```text
S-0523 has started. Most of its work does not depend on your answer, and I am doing that now. One choice is yours. TH-0694 asked it two days ago, and it is still open there.
```

The answer:

```text
**Recommended for the operator:** option 1. Add a new permission, `orchestration.permissions.set_cost_of_delay`, off by default. Leave `plan_backlog_stories` as ADR-0119 § Decision item 4 has it: inputs on a backlog story only.

## Why this is yours to decide

- It is a decision not yet recorded. ADR-0119 limits my cost of delay writes to a backlog story. Either option widens that.
- It widens my own powers over money figures. I should not settle that myself.

## Why option 1

- **You turn it on yourself.** The new key is off by default. Confirming it grants me nothing until you set it in `system-flow.yaml`. Option 2 grants the power at once to every project that already has `plan_backlog_stories` on.
- **It keeps consent separate.** The convention says "never overwrite the operator's inputs". Option 1 lets me change inputs you set. A separate switch keeps that a choice you make on its own. You can also turn it off without losing story planning.
- **The cost is small.** It adds one manifest key, one settings toggle, and docs rows. The story's agent already builds the gate as one place.

## The same question on TH-0694

TH-0694 question 2 asks the same thing. My recommendation there, from 2026-10-28T14:34:37Z, is the same answer and still awaits you. One answer on either thread settles both.

The thread stays open for your answer.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision
```

Source: TH-0764.
