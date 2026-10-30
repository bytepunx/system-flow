---
id: D-0431
title: "S-0523 is in progress: which permission lets the orchestrator set cost of delay inputs on an epic or a started story (TH-0694 question 2): Recommended for the operator: option 1. Add a new…"
kind: answer
sources: [TH-0764, S-0523]
topics: [orchestration, planning, flai, cli, go, template, conventions, code, all]
paths: [design/adrs, design/adrs/README.md, design/system/strategic-agents.md, design/system/flai-cli.md, design/system/project-manifest.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, template/CHANGELOG.md, docs/users/flai.md, docs/users/flai-reference.md, docs/users/flaiover.md, docs/operators/settings.md, docs/operators/index.md, flai/internal/manifest/settings.go, flai/internal/guard/guard.go, flai/internal/guard/orchestrate_test.go, flai/internal/guard/guard_test.go, flai/internal/workitem/promotable.go, flai/internal/mcpserver/items_write.go, flai/internal/mcpserver/items_write_test.go, flai/cmd/edit.go, flai/cmd/edit_test.go, flai/cmd/orchestrate_permissions_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/issues/I-0187-the-orchestrator-cannot-set-a-cost-of-delay-input-the-operator-decided-on-an-epic-or-a-story-past-the-backlog.md, design/issues/summary.md]
date: 2026-10-30
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0431 S-0523 is in progress: which permission lets the orchestrator set cost of delay inputs on an epic or a started story (TH-0694 question 2): Recommended for the operator: option 1. Add a new…

The operator decided that the orchestrator's power to add or change cost of delay inputs on any open epic or story is a new permission of its own, `orchestration.permissions.set_cost_of_delay`, off by default. `plan_backlog_stories` is not widened: it keeps ADR-0119 § Decision item 4, first inputs on a backlog story that has none.

What was asked: S-0523's agent asked on TH-0764 which permission lets the orchestrator set cost of delay inputs the operator decided on an epic or a story past the backlog (I-0187). TH-0694 question 2 asked the same. The options were a new permission (option 1) or widening `plan_backlog_stories` (option 2).

The answer: on 2026-10-30T16:47:47Z the operator, alex, confirmed the orchestrator's recommendation of option 1. Its reasons:

- Changing inputs the operator set is a new power over money, and the convention says "never overwrite the operator's inputs". A switch of its own keeps that consent explicit for every project made from the template.
- Off by default, it grants nothing until the operator sets it in `system-flow.yaml`. Option 2 would grant it at once to every project with `plan_backlog_stories` on.
- The operator can turn it off without losing story planning.

Both options keep the write to inputs only: no value, no removal, nothing else changed in the same call, never on a done, cancelled, or archived item, and only to a figure the operator decided on a thread.

Sources: TH-0764, S-0523; the same question is TH-0694 question 2.
