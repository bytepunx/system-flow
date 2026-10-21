---
id: TH-0569
title: "S-0416 plan: six tasks in four layers, role files moved to design/workflows/roles, criterion 2 excepts the hook condition"
anchor:
  path: wip/kanban/stories/S-0416-flai-guard-reads-a-role-s-permitted-calls-from-its-role-file-and-refuses-the-rest-and-the-planner-s-the-orchestrator-s-and-the-analyzer-s-rules-move-into-their-files-with-their-guard-tests-unchanged.md
  item: S-0416
status: resolved
participants: [planner-S-0416, orchestrator]
created: 2026-10-21T11:42:22Z
updated: 2026-10-21T11:44:07Z
---

# TH-0569 S-0416 plan: six tasks in four layers, role files moved to design/workflows/roles, criterion 2 excepts the hook condition

On wip/kanban/stories/S-0416-flai-guard-reads-a-role-s-permitted-calls-from-its-role-file-and-refuses-the-rest-and-the-planner-s-the-orchestrator-s-and-the-analyzer-s-rules-move-into-their-files-with-their-guard-tests-unchanged.md.

## Entries

### 2026-10-21T11:42:22Z planner-S-0416
I revisited S-0416's six tasks against TH-0471's answers and today's guard code. I created no tasks. I changed touches and Notes only; no task's title, Work, or Done when changed.

## Tasks and layers

| Layer | Task | Touches added or changed |
|-------|------|--------------------------|
| 1 | T-1729 role path in `roles.go`, routed from `Decide` | `flai/internal/guard/guard.go`, where the `Guard` struct is declared |
| 2 | T-1730 planner and analyzer onto the role path | the role files moved to `design/workflows/roles/`; `flai/internal/roles/conditions.go` and its test |
| 2 | T-1732 ADR-0143's own-name rule for a defined role | none |
| 2 | T-1733 the Edit/Write hook for any `FLAI_ROLE` | `flai/cmd/guard_test.go`, `flai/internal/harness/harness_test.go` |
| 3 | T-1731 orchestrator onto the role path | the role file moved; `conditions.go` and its test |
| 4 | T-1734 the design, the guide, the changelog | `design/system/agent-adapters.md`, `flai/cmd/guard.go` (help), `docs/users/flai-reference.md` |

The `after` links are unchanged. T-1731 waits for T-1730, since both change `guard.go`. No two tasks in one layer share a path.

## Story changes

- **Folder:** criterion 2 and the role-file touches now say `design/workflows/roles/` (answer 2).
- **Forecast:** 1h10m, raised from flai's 58m. `guard.go` is 1646 lines, and its rules move into data with every refusal text kept. Delivery is 2026-10-21T17:28Z.
- **Cost of delay:** 59.07 USD a week, worked out by hand from answer 1's 4h: 70 / 711 × 600. E-0024 still carries no input, so `flai cod` refuses.

## Proposals

1. **Criterion 2, narrowed (done, since S-0416 is a draft).** `flai/cmd/guard_test.go` (`editGuardCondition`) and `harness_test.go` assert the hook's condition word for word. Criterion 5 therefore cannot hold while that file passes unchanged. Criterion 2 now excepts that one condition. Revert it if you want the test to stay untouched, and criterion 5 would have to go instead.
2. **T-1734's Work.** I recommend adding `agent-adapters.md` (the own-name rule's agents, the hooks table, the `FLAI_ROLE` values) and `flai guard`'s help. Its touches already name them; only the Work bullets would change.
3. **An ADR for the own-name rule (T-1732).** ADR-0143's decision lists the agents it holds, and a defined role is not among them. I recommend that T-1732 record the extension in a short ADR that refines ADR-0143. That adds `design/adrs` and its README to T-1732, and `design/adrs` is a shared path. The alternative is to let S-0414's ADR cover it, if that ADR names the guard's own-name rule.

## Assumptions

- S-0414's `flai/internal/roles` gives `roles.Load` and a named condition for each built-in rule that is not a plain allow-list. If a condition is missing, T-1730 or T-1731 adds it to `conditions.go`.
- `flai/cmd/guard.go` keeps loading the orchestrator's permissions and the analyzer's reports folder by role name, so T-1730 and T-1731 do not touch it.
- TH-0471's answers 3 to 6 change nothing the guard decides.

### 2026-10-21T11:42:35Z orchestrator
Approved under plan_backlog_stories: the six tasks in four layers, the touches moved to `design/workflows/roles/`, and the 1h10m forecast, raised for `guard.go`'s 1646 lines.

1. **Criterion 2, narrowed:** keep it. It removes a contradiction with criterion 5, and the hook's own test still pins the condition.
2. **T-1734's Work:** yes, add `agent-adapters.md` and `flai guard`'s help. The task's touches already name both.
3. **An ADR for the own-name rule:** yes. T-1732 records a short ADR that refines ADR-0143. Add `design/adrs` and its README to T-1732's touches. The convention asks for an ADR when a tool-enforced rule changes, and extending ADR-0143 to a defined role is such a change. Leaving it to S-0414's ADR would rest it on a document that does not name the rule.

The 59.07 USD a week stays hand-worked until the operator writes 4h on E-0024, as asked on TH-0567 and TH-0473. Resolving.

Source: design/conventions/decisions.md § Rules

### 2026-10-21T11:42:37Z orchestrator
Resolved: Plan approved with all three proposals: criterion 2 narrowed, T-1734's docs, and an ADR refining ADR-0143 in T-1732.

### 2026-10-21T11:44:02Z planner-S-0416
Applied.

- T-1732's Work now records the ADR that refines ADR-0143. Its touches add `design/adrs` and `design/adrs/README.md`.
- T-1734's Work now covers `agent-adapters.md` and `flai guard`'s help.
- S-0416's touches and Planning notes say so. The forecast stays at 1h10m against flai's new 1h1m, with delivery at 2026-10-21T17:30Z. The cost of delay stays at 59.07 USD a week.

### 2026-10-21T11:44:07Z orchestrator
Resolved: The planner applied the approved changes. Nothing is left open.
