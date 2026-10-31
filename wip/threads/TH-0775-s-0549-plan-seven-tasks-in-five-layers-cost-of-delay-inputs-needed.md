---
id: TH-0775
title: "S-0549 plan: seven tasks in five layers; cost of delay inputs needed"
anchor:
  path: wip/kanban/stories/S-0549-the-merge-queue-checks-the-main-checkout-s-uncommitted-files-and-the-lint-of.md
  item: S-0549
status: open
participants: [planner-S-0549, orchestrator]
created: 2026-10-31T07:43:36Z
updated: 2026-10-31T07:46:35Z
---

# TH-0775 S-0549 plan: seven tasks in five layers; cost of delay inputs needed

On wip/kanban/stories/S-0549-the-merge-queue-checks-the-main-checkout-s-uncommitted-files-and-the-lint-of.md.

## Entries

### 2026-10-31T07:43:36Z planner-S-0549
**Question: cost of delay inputs.** Neither S-0549 nor E-0031 has inputs, so `flai cod S-0549` refuses and I wrote no value.

Recommended, on S-0549:

| Input | Value | Why |
|---|---|---|
| `time_lost_per_cycle` | 20m | about 40 minutes a fortnight of finding and fixing a dirty main by hand, as T-2332's fix on 2026-10-29 and S-0541's 28 failed landings took |
| `penalty_per_week` | none | the verification time a dirty main wastes is the host's, not money |
| `revenue_per_week` | none | |

At 150 USD an hour and a 168h cycle that is 50 USD a week. Set it with `flai edit S-0549 --time-lost-per-cycle 20m`, or answer with other figures (12m as recommended for S-0548, or 1h to count the stories a failing batch holds back) and I write the value.

**Plan.** Seven tasks, all new; the story had none.

| Layer | Task | Criteria | Waits for |
|---|---|---|---|
| 1 | T-2381 An ADR for the preflight, the held state, the thread, and the landing again | 4 | none |
| 1 | T-2382 The lint findings move from `flai/cmd/accept_markdown.go` to `flai/internal/verify/findings.go` | 1 | none |
| 2 | T-2383 `RunBatch` runs the preflight and holds the queue on a finding, in a new `mergequeue/preflight.go` | 1 | T-2381, T-2382 |
| 3 | T-2384 A verified batch whose landing failed on main's state is kept and lands again from its verified commits | 3 | T-2383 (`state.go`, `log.go`) |
| 3 | T-2385 `flai board --json` and the dashboard's merge queue panel show the hold | 4 | T-2383 |
| 4 | T-2386 `flai serve` journals the hold, opens and resolves one thread, wakes the queue on main's changes, and journals a landing again | 1, 2, 3 | T-2383, T-2384 |
| 5 | T-2387 `flai-cli.md`, `docs/users/flai.md`, and `docs/operators/index.md` | 4 | T-2383, T-2384, T-2386 |

Forecast: 5h, raised from `flai forecast`'s 16m. The story's `### Planning` says why and where each touch came from.

**Assumptions.**

- The preflight holds the queue on error-level `flai check` findings only. A batch's own check is scoped to its stories (ADR-0183), and a warning elsewhere should not stop every batch. T-2381 decides it.
- The preflight hold is a field of its own, apart from the operator's `flai queue hold`, so `flai queue resume` cannot clear a finding that still stands.
- The one thread is anchored on `system-flow.yaml`, as `flai serve`'s upgrade notice is. flai resolves it when the hold clears.
- "A `wip/` or main change" means any change in the main checkout, or a move of main's head. Today the watcher wakes the queue only for `/kanban/` and `/threads/` paths (`serve/serve.go`), so T-2386 widens that while the queue is held.
- A kept batch lands again only while main's head is still the base it was verified against. When main moved, it is dropped as today.
- I-0203's case, a stale generated file committed on main, is caught only if `flai check` reports it. Accepting outside the queue is S-0552's concern, not this story's.

**Overlaps.** S-0552 (in progress) shares `flai/internal/mergequeue` under MS-0263's split. T-2383, T-2384, and T-2385 keep their edits to `batch.go`, `land.go`, `log.go`, and `eligible.go` small. `docs/operators/index.md` is claimed by S-0548 and S-0552 too.

**Proposal.** `flai queue`'s text output (`flai/cmd/queue.go`) could print the preflight hold too. I left it out because S-0552 changes that file. `flai queue --json` and `flai board --json` carry the hold anyway. Tell me if you want it as a task.

**Defect.** My start prompt said to read the thread `asked`. No such thread exists: `asked` is the operator's plan trigger, read as a thread name. planner-S-0548 reported the same on TH-0771.

### 2026-10-31T07:46:35Z orchestrator (recommendation)
**Recommendation:** set `time_lost_per_cycle: 20m` on S-0549, as the planner recommends:

```bash
flai edit S-0549 --time-lost-per-cycle 20m
```

## The figure, and not 12m or 1h

- **The goal is time lost.** About 40 minutes a fortnight of finding and fixing a dirty main by hand. The verification time wasted is the host's, not money.
- **The story gives no reason for an alternative.** It does not count the stories a failing batch holds back, which would argue for 1h.
- **On the story, not E-0031**, as on TH-0773 and TH-0774.

## Why this needs you

S-0549 is ready. You moved it there at 2026-10-31T07:24:09Z. flai refuses the orchestrator inputs on a story past the backlog: it refused S-0547's just now with "S-0547 is ready, not in the backlog: the orchestrator gives cost of delay inputs to a story in the backlog alone". A confirmation of this entry does not write the input. Running the command does.

## The plan: approved

- **The tasks fit the criteria.** T-2382 and T-2383 cover criterion 1, T-2386 criteria 1 to 3, T-2384 criterion 3, T-2381, T-2385, and T-2387 criterion 4.
- **Assumptions: right.** Holding on error-level findings only, a hold apart from `flai queue hold`, landing again only while main is the verified base: each keeps the queue safe. T-2381's ADR records them.
- **Proposal, `flai queue` text output: leave it out.** No criterion asks for it, `flai queue --json` and `flai board --json` carry the hold, and S-0552 changes that file.
- **Touches and forecast fit.** The 5h forecast over flai's 16m states its reason.

The thread stays open for the cost of delay input only.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision
