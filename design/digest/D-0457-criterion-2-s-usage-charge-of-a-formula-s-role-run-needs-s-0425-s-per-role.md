---
id: D-0457
title: "Criterion 2's usage charge of a formula's role run needs S-0425's per-role usage: move that clause to S-0425: Recommended: take agent-S-0422's recommendation and move the usage-charge clause of…"
kind: answer
sources: [TH-0800, S-0422, S-0425]
topics: [cli, orchestration, conventions, template, planning, flai, go, code, all]
paths: [design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, design/system/strategic-agents.md, design/system/flai-cli.md, docs/operators/index.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flai-reference.md, template/CHANGELOG.md, flai/cmd/accept.go, flai/cmd/accept_formula_test.go, flai/cmd/formula.go, flai/cmd/formula_test.go, flai/cmd/role.go, flai/cmd/role_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/check/check.go, flai/internal/check/formulas_test.go, flai/internal/formulas/formulas.go, flai/internal/formulas/formulas_test.go, flai/internal/guard/roles.go, flai/internal/guard/roles_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/mcpserver/role.go, flai/internal/mcpserver/role_test.go, flai/internal/preview/orchestrator.go, flai/internal/preview/accept_test.go, flai/internal/serve/activity.go, flai/internal/serve/activity_test.go, flai/internal/serve/agents.go, flai/internal/serve/formula_schedule.go, flai/internal/serve/formula_schedule_test.go, flai/internal/serve/role.go, flai/internal/serve/role_test.go, flai/internal/serve/serve.go, flai/internal/roles/roles.go, design/adrs/0239-a-formula-s-instance-is-worked-by-its-story-s-agent-step-by-step-a-role-step.md, design/adrs/README.md, design/workflows/formulas/README.md, template/root/design/workflows/formulas/README.md, design/system/metrics.md]
date: 2026-11-02
status: absorbed
absorbed_by: "design/adrs/0249-strategic-usage-is-per-role-a-role-s-activity-is-charged-by-its-role-file-s.md#decision"
asked_again: []
cited: [S-0425, TH-0801]
---

# D-0457 Criterion 2's usage charge of a formula's role run needs S-0425's per-role usage: move that clause to S-0425: Recommended: take agent-S-0422's recommendation and move the usage-charge clause of…

The operator, alex, answered TH-0800 at 2026-11-02T03:42:04Z: a formula's role run is not charged to its instance under `usage.strategic` until S-0425's ADR on `design/system/metrics.md` lets `usage.strategic` take a defined role's name. S-0422 dropped that clause from its criterion 2, and S-0425's criterion 2 carries it.

What was asked:

```text
**Recommended:** move the clause of S-0422's criterion 2 that charges a role step's usage to the instance under `usage.strategic` to S-0425, and tick criterion 2 here for the rest (the role_run lift, the trigger, and the wake, all done in T-1775).
```

The answer:

```text
**Recommended:** take agent-S-0422's recommendation and move the usage-charge clause of S-0422's criterion 2 to S-0425.

**What changes, and who does it:**

- **S-0422, criterion 2.** Drop the clause "and its usage is charged to the instance under `usage.strategic` as a planner's is to the item it planned". Keep the rest: the `role_run` lift, the trigger, and the wake. The agent then ticks it and records the move under `### Planning` in its Notes. The criterion's words change, so you approve that here.
- **S-0425, criterion 2.** Add "including a formula's role step, charged to the instance its run was started for". S-0425 is finalized, so that edit is yours.

**Why:**

- Only an ADR can let `usage.strategic` take a defined role's name, because `design/system/metrics.md` is the contract between flai and the dashboard. That ADR is S-0425's first task, T-1790.
- S-0425's criterion 2 already charges a defined role's run, under the `item` rule, to the item it was started for. T-1775 makes the instance that item, so the charge follows with no further code.
- The alternative is to write that ADR and change `workitem/usage.go`, `metrics.md`, and `serve/activity.go` in S-0422. That duplicates most of T-1790 and T-1791, and those files are in S-0425's touches.
- Nothing is lost while S-0425 waits. S-0425 is held behind S-0423 and S-0424, which are both held behind S-0422. Until then, a formula's role run is not charged, the same as any defined role's run today.

**What happens next:** with this answer, S-0422's agent finishes criterion 2 without the charge and goes on to its other tasks.

Source: design/conventions/decisions.md § Project additions
```

S-0422's narrative, `## Decisions`, records the outcome:

```text
Criterion 2's usage charge moved to S-0425 on TH-0800, confirmed by the operator: `usage.strategic` takes a defined role only through S-0425's ADR on `metrics.md`. The criterion's words and the Notes say so.
```

S-0425 carried the clause out. Its narrative, `## Decisions`, records:

```text
Criterion 2 gained "including a formula's role step, charged to the instance its run was started for", as the operator confirmed on TH-0800 (D-0457); answered TH-0801.
```

[ADR-0249](../adrs/0249-strategic-usage-is-per-role-a-role-s-activity-is-charged-by-its-role-file-s.md) § Decision, item 3, now states the charge: a role run started for a formula's instance has the instance as its item, so a role whose rule is `item` is charged to the instance and the items above it, and the analyzer started for a role step keeps its own rule, `issues`.

Sources: TH-0800, S-0422 (wip/archive/agents/S-0422.md), S-0425 (wip/archive/agents/S-0425.md).
