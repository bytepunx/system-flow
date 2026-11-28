---
id: TH-1332
title: "Plan for S-0855: hold starts on a failed harness login and tell the operator once"
anchor:
  path: wip/kanban/stories/S-0855-flai-serve-restarts-the-orchestrator-every-minute-while-the-harness-s-login-has.md
  item: S-0855
status: resolved
participants: [planner-S-0855, claude-concierge]
created: 2026-11-28T00:00:58Z
updated: 2026-11-28T00:01:51Z
---

# TH-1332 Plan for S-0855: hold starts on a failed harness login and tell the operator once

On wip/kanban/stories/S-0855-flai-serve-restarts-the-orchestrator-every-minute-while-the-harness-s-login-has.md.

## Entries

### 2026-11-28T00:00:58Z planner-S-0855
S-0855 is planned with six tasks in five layers. The touches, the forecast (4h, raised from flai's 5m), and the cost of delay (2100 USD a week, as `flai cod` gives it) are written. The reasons are in the story's Notes under `### Planning`.

## Tasks and layers

| Layer | Task | What it does |
|-------|------|--------------|
| 1 | T-3643 | ADR: decides from the real failed-run log how a failed login is detected, held, told, and cleared |
| 2 | T-3644 | `usage` reads a claude-code result that failed to authenticate as a login limit, with a fixture from the log |
| 3 | T-3647 | The login limit feeds the account's budget state, so every start waits for its backoff and probe and not `orchestrateRetry`. Includes the test that reproduces I-0412 |
| 4 | T-3649 | One thread to the operator in each project the hold covers, resolved when it clears |
| 4 | T-3650 | The dashboard's budget banner says the login failed and how to renew it |
| 5 | T-3651 | Design and docs, and `flai issue close I-0412` |

T-3649 and T-3650 run together: they share no path. Each task waits for the one before because it reads what that one writes. T-3649 also edits `budget.go` after T-3647.

## Assumptions

- The fix reuses ADR-0209 §5's budget state (one per account, 5m backoff doubling to 1h, a probe, a notice) with a new kind, `login`. This is T-3643's recommended decision, not a settled one: the ADR may change it.
- The failed run's `result` line carries a status or error that marks the failed login (401, `authentication_error`, or the text "Failed to authenticate"). Nobody has read it in the repository yet. T-3643 reads it from `~/.flai/serve/agents/` on the host.
- A failed login also holds story agents' starts and spends no `agent.auto_restarts` (ADR-0108), since one login serves them all.
- The ADR's number is not known yet: S-0852 claims 0322. The story keeps the folder touch `design/adrs` for that reason alone.
- No `after`: S-0855 reads no other open story's output.

## Proposed additions, not drafted

- A separate story could bound consecutive failed runs of any strategic role, whatever the cause. Today only a usage limit or a budget holds them, and an unknown harness error would loop the same way. I recommend it as its own story and not as part of this one, so that this fix stays small.

### 2026-11-28T00:01:51Z claude-concierge
For the operator, who delegated completing drafts and queueing ready work in this session: the plan stands as drafted, assumptions included.

Resolved: plan reviewed and the story finalized on the operator's delegation
