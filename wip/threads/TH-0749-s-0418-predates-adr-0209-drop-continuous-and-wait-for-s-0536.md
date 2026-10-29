---
id: TH-0749
title: "S-0418 predates ADR-0209: drop continuous and wait for S-0536?"
anchor:
  path: wip/kanban/stories/S-0418-a-role-a-project-defines-starts-on-an-event-or-on-a-schedule-as-its-role-file-s.md
  item: S-0418
status: resolved
participants: [agent-S-0418, orchestrator, alex]
created: 2026-10-29T14:43:32Z
updated: 2026-10-29T17:20:43Z
---

# TH-0749 S-0418 predates ADR-0209: drop continuous and wait for S-0536?

On wip/kanban/stories/S-0418-a-role-a-project-defines-starts-on-an-event-or-on-a-schedule-as-its-role-file-s.md.

## Entries

### 2026-10-29T14:43:32Z agent-S-0418
**Recommendation:** answer A. I rewrite S-0418's criteria and tasks to fit ADR-0209 and ADR-0211, move S-0418 back to ready with `after: S-0536`, and work it once S-0536 is accepted.

## The conflict

- **Criterion 3 contradicts an accepted ADR.** It keeps a `continuous` role running and adds `flai role stop|start` and `role.stop|start`. ADR-0209 §1 says "`flai check` rejects `continuous` on any role". S-0536, now in review, makes `flai check` report it (`flai/internal/roles/roles.go` on `story/S-0536`).
- **S-0536 already holds the machinery criterion 1 needs.** Its reason queue (`flai/internal/serve/orchestrate_reasons.go`) queues events, coalesces them, and starts one run per batch. A project role's event trigger should reuse it, not build a second queue.
- **T-1741 is already done.** S-0414 reads and checks triggers, events, schedules, and conditions, and S-0536 adds the `continuous` finding. The one gap left: nothing checks the value a `roles.<name>.schedule` setting holds.
- **The files overlap.** S-0536 changes `serve.go`, `roles.go`, `cmd/role.go`, `check/roles_test.go`, and every doc S-0418 touches. Writing S-0418 against today's main means a heavy rebase onto S-0536.

## Options

- **A (recommended): re-scope S-0418 and wait for S-0536.**
  - Drop criterion 3. Cancel T-1746 and T-1741.
  - Criterion 1: events and conditions from ADR-0207 and ADR-0211, queued and coalesced through S-0536's reason queue.
  - Criterion 2: the schedule comes from a new `roles.<name>.schedule` manifest key, a cron expression or `daily`, run through `cadence`.
  - Criterion 4: no `continuous` value.
  - Criterion 5: an unknown event or condition, and a `roles.<name>.schedule` value that does not parse or never comes round.
  - I move S-0418 back to ready with `after: S-0536`, and flai serve starts it again once S-0536 is accepted.
- **B: keep the criteria as written.** This needs a new ADR that supersedes part of ADR-0209 §1. It also undoes S-0536's check.
- **C: cancel S-0418.** A later E-0030 story generalises S-0536's reason queue to project roles.

I have changed nothing and committed nothing. Nothing else in S-0418 can go ahead until you answer.

Source: design/adrs/0209-the-orchestrator-runs-for-a-reason-and-ends-started-by-flai-serve-on-the-events.md § Decision

### 2026-10-29T14:43:47Z orchestrator (recommendation)
**Recommended for the operator:** answer A. Re-scope S-0418 to ADR-0209 and ADR-0211, move it back to ready with `after: S-0536`, and let it be worked once S-0536 is accepted.

## Why A

- **The accepted ADR settles criterion 3.** ADR-0209 § Decision point 1 says "`flai check` rejects `continuous` on any role". A criterion that keeps a `continuous` role running, with `flai role stop|start`, contradicts it, and S-0536's check would fail on it. Option B needs a superseding ADR and undoes S-0536's check, for no use the story names.
- **One queue, not two.** Point 2 makes starts coalesced through a reason queue. S-0536 builds that queue in `orchestrate_reasons.go`, so a project role's event trigger reuses it rather than building a second.
- **Waiting is cheaper than rebasing.** S-0536 changes `serve.go`, `roles.go`, `cmd/role.go`, and the docs S-0418 touches. Working S-0418 against today's main means a heavy rebase once S-0536 lands, and S-0536 is in review awaiting only your acceptance (TH-0748).
- **C loses work already planned.** The `roles.<name>.schedule` key and its check are still needed, and A keeps them.

## Why it is yours

Rewriting a finalized story's criteria and cancelling two of its tasks is a change of scope. I promoted S-0418 at 14:42Z after S-0536's move to review freed its hold, before this conflict was known.

Source: design/adrs/0209-the-orchestrator-runs-for-a-reason-and-ends-started-by-flai-serve-on-the-events.md § Decision

Source: design/adrs/0209-the-orchestrator-runs-for-a-reason-and-ends-started-by-flai-serve-on-the-events.md § Decision

### 2026-10-29T16:35:26Z alex
Confirmed the recommendation of 2026-10-29T14:43:47Z orchestrator.

Source: design/adrs/0209-the-orchestrator-runs-for-a-reason-and-ends-started-by-flai-serve-on-the-events.md § Decision

### 2026-10-29T17:20:43Z agent-S-0418
Resolved: Option A applied: S-0418 re-scoped to ADR-0209 and ADR-0211, T-1741 and T-1746 cancelled, worked after S-0536's acceptance.
