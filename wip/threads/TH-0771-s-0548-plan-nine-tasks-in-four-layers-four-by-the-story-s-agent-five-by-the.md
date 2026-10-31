---
id: TH-0771
title: "S-0548 plan: nine tasks in four layers (four by the story's agent, five by the planner); cost of delay inputs needed"
anchor:
  path: wip/kanban/stories/S-0548-a-queue-verification-takes-under-five-minutes-the-slow-and-network-bound-tests.md
  item: S-0548
status: resolved
participants: [planner-S-0548, alex]
created: 2026-10-31T07:28:57Z
updated: 2026-10-31T07:29:34Z
---

# TH-0771 S-0548 plan: nine tasks in four layers (four by the story's agent, five by the planner); cost of delay inputs needed

On wip/kanban/stories/S-0548-a-queue-verification-takes-under-five-minutes-the-slow-and-network-bound-tests.md.

## Entries

### 2026-10-31T07:28:57Z planner-S-0548
**Question: cost of delay inputs.** Neither S-0548 nor E-0031 has cost of delay inputs, so `flai cod S-0548` refuses and I wrote no value.

Recommended, on S-0548:

| Input | Value | Why |
|---|---|---|
| `penalty_per_week` | 14 | about 44 batches a week (15.2 h at 10.4 min a fortnight) × 0.32 USD of gateway spend |
| `time_lost_per_cycle` | 12m | about 6 min of go-test per close-out (394 s against about 30 s) plus about 5.5 min per batch above five minutes |
| `revenue_per_week` | none | |

Set them with `flai edit S-0548 --penalty-per-week 14 --time-lost-per-cycle 12m`, or answer with other figures and I write the value.

**Plan.** The story's agent started at 07:21:39 while I planned, and wrote T-2361 to T-2364 from the story's first four criteria. I wrote T-2365 to T-2369 for what those tasks leave out. I first drafted T-2366 and T-2367 as duplicates of T-2361 and T-2362, then rewrote them as new work. Nothing was cancelled.

| Layer | Task | Criteria | Waits for |
|---|---|---|---|
| 1 | T-2361 mcpserver tests run with `t.Parallel` (agent's) | 3, 4 | none |
| 1 | T-2362 serve tests I-0226 and I-0233 wait for what they assert (agent's) | 7 | none |
| 1 | T-2363 go-test tier `-timeout`, `flai check` warns (agent's) | 8 | none |
| 1 | T-2365 dashboard releases test reaches a closed listener | 1 | none |
| 2 | T-2366 four whole-repository tests skip under `-short`, go-test tier `-short`, check and lint once a batch | 5, 6 | T-2361 (`prime_test.go`), T-2363 (`system-flow.yaml`) |
| 2 | T-2367 cmd tests run with `t.Parallel` | 3, 4 | T-2365 (`dashboard_test.go`) |
| 2 | T-2368 `flai test <tier>` runs a named tier | 2 | T-2363 (`docs/users/flai.md`, `flai-cli.md`) |
| 3 | T-2369 gateway tier out of smoke | 2 | T-2366, T-2368 |
| 4 | T-2364 close the four issues, record timings (agent's) | 7, 9 | T-2361, T-2362, T-2365, T-2366, T-2367, T-2369 |

I re-enriched T-2364's touches (the four issue files) and its `after`, so that it measures last. Forecast: 6h, raised from `flai forecast`'s 39m; S-0548's `### Planning` says why.

**Assumptions.**

- "A schedule the host sets" is the host's own scheduler (cron or a systemd timer) running `flai test gateway`, with no schedule code in flai (T-2369).
- `flai test gateway` needs new code: `flai test` takes no tier name today (T-2368).
- "The testing convention" means the tiers line of `design/conventions/code-quality.md` and the `flai test` line of `tooling.md`. The declared touch `design/conventions/testing.md` does not exist; I kept it.
- Criterion 4's `./internal/cmd` is the package `./cmd`.
- Running the close-out's go-test with `-short` reverses I-0201's choice. T-2366 records that in `## Decisions`.
- `.github/workflows/flai.yml` and `scripts/integration.sh` already pass `-timeout` (30m, 20m). T-2363 says so and changes neither.

**Proposals for the story's agent and the operator.**

- T-2361's Done when says three minutes; criterion 4 now says two. I propose the agent align it, since it wrote the task.
- T-2364 should also record criterion 9's batch time, and `cmd`'s time and five slowest tests from T-2367.
- S-0543 (backlog) remediates I-0233, and S-0540 (backlog) remediates I-0227. Both are absorbed here: I propose cancelling them once S-0548 is accepted. That is the operator's call.
- S-0552, in progress, overlaps on `flai/cmd`, `flai/internal/serve`, `system-flow.yaml`, and `docs/operators`.

**Defect.** This planner's start prompt said "the orchestrator asked for it again for what the operator wrote on asked". The journal shows the operator's `plan.run` at 07:21:28 with the trigger `asked`, and no thread `asked` exists. The installed flai reads the trigger `asked` as a thread name.

### 2026-10-31T07:29:34Z alex
Resolved.
