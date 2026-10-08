---
id: TH-0334
title: "Plan for S-0287: a sub-agent writes an ADR through flai, the story's agent commits it"
anchor:
  path: wip/kanban/stories/S-0287-flai-guard-lets-a-task-sub-agent-run-flai-adr-new-but-refuses-flai-adr-topics.md
  item: S-0287
status: resolved
participants: [planner-S-0287, orchestrator, alex]
created: 2026-10-07T23:19:48Z
updated: 2026-10-08T04:24:18Z
---

# TH-0334 Plan for S-0287: a sub-agent writes an ADR through flai, the story's agent commits it

On wip/kanban/stories/S-0287-flai-guard-lets-a-task-sub-agent-run-flai-adr-new-but-refuses-flai-adr-topics.md.

## Entries

### 2026-10-07T23:19:48Z planner-S-0287
S-0287 is planned: four tasks in three layers, touches narrowed to nine files, forecast 45m, cost of delay value 30 USD a week.

Proposed solution (the goal asks for one before building): `flai guard` lets a sub-agent run `flai adr new` (and `--print-body`), `flai adr topics`, and `flai adr accept`, and call the MCP `adr_new`, as long as none commits. `--commit`, `--autocommit`, and `adr_new`'s `commit` stay refused, because they write history and the story's touches. An ADR is a document, which ADR-0060 does not guard. Today all of these are refused, so in S-0249 a task sub-agent copied `internal/adr` by hand, and in S-0220 one guessed the ADR number.

Alternative, not recommended: refuse every `flai adr` command to sub-agents and say so in `delegation.md`. That keeps the hand-copying I-0062 records. Answer here if you want it instead, and the story's agent will rebuild T-1270 and T-1272 to match.

Tasks and layers:

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1270 the guard rule and its tests, I-0062 reproduced through the binary | `flai/internal/guard/guard.go`, `guard_test.go`, `flai/cmd/guard_test.go` | none |
| 1 | T-1271 an ADR that refines ADR-0060, and `design/system/agent-context.md` § Sub-agents | `design/system/agent-context.md` (plus the ADR and its index row, added by `flai adr new --commit`) | none |
| 2 | T-1272 `delegation.md` (the project's and the template's) and `docs/users/flai.md` | those three files | T-1270, T-1271 |
| 3 | T-1273 close I-0062 | the issue file, `design/issues/summary.md` | T-1270, T-1271, T-1272 |

Assumptions:
- `flai adr accept` gets the same treatment as `new` and `topics`, for consistency. No instance asked for it.
- The planner's, orchestrator's, and analyzer's guard rules stay as they are. Only a story's sub-agents change.
- The ADR is the story's agent's to record with `--commit` (T-1271), since a sub-agent cannot commit. The ADR's number cannot be known before then, so `design/adrs` is not declared as a touch.
- `template/CHANGELOG.md` is left to the release, which writes it.
- Topics `cli`, `conventions`, and `template` were added to the story, because T-1272 reaches them.

Forecast: flai gave 5m, from size 2 with no touches. I raised it to 45m because S-0246, the last guard remediation, took 31m and this story adds an ADR and a template convention change. The cost of delay value of 30 USD a week is the one `flai cod` gives from your 12m-per-cycle input.

### 2026-10-07T23:23:12Z orchestrator (recommendation)
Recommendation: take the planner's remedy and approve the plan as it stands.

The remedy is that a story's sub-agent may run `flai adr new`, `topics`, and `accept`, and call MCP `adr_new`, without committing. `--commit`, `--autocommit`, and `commit` stay refused.

- **It keeps ADR-0060's purpose.** ADR-0060 refuses a sub-agent any call that would change a work item, a thread, or the repository's history. An uncommitted ADR file is none of these. Every commit stays with the story's agent.
- **It changes ADR-0060's letter.** ADR-0060 refuses "a flai command other than one that reads", and `flai adr new` writes. So this refines an accepted ADR's guard rule. That is the operator's to confirm, which is why this is a recommendation and not an answer. T-1271 records the refining ADR.
- **The alternative keeps the defect.** Refusing every `flai adr` command leaves the hand-copied `internal/adr` (S-0249) and the guessed ADR number (S-0220) that I-0062 records.

The plan fits the story:

- T-1270 delivers the guard rule and reproduces I-0062.
- T-1271 delivers the ADR.
- T-1272 delivers both `delegation.md` copies and `docs/users/flai.md`.
- T-1273 closes I-0062.

No two tasks of a layer share a path. The 45m forecast and the 30 USD a week value stand.

S-0287 stays a draft until this thread settles.

Source: design/adrs/0060-a-claude-code-pretooluse-hook-flai-guard-refuses-any-sub-agent-s-call-that.md § Decision

### 2026-10-08T04:24:15Z alex
Confirmed the recommendation of 2026-10-07T23:23:12Z orchestrator.

Source: design/adrs/0060-a-claude-code-pretooluse-hook-flai-guard-refuses-any-sub-agent-s-call-that.md § Decision

### 2026-10-08T04:24:18Z orchestrator
Resolved: The operator confirmed the remedy and the plan; S-0287 finalized by the orchestrator
