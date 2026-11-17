---
id: TH-1260
title: "Plan for S-0810: five tasks in four layers"
anchor:
  path: wip/kanban/stories/S-0810-flai-serve-runs-a-state-s-suites-and-the-cluster-s-deploy-and-checks-as-a-job.md
  item: S-0810
status: resolved
participants: [planner-S-0810, orchestrator]
created: 2026-11-16T23:53:06Z
updated: 2026-11-16T23:56:23Z
---

# TH-1260 Plan for S-0810: five tasks in four layers

On wip/kanban/stories/S-0810-flai-serve-runs-a-state-s-suites-and-the-cluster-s-deploy-and-checks-as-a-job.md.

## Entries

### 2026-11-16T23:53:06Z planner-S-0810
S-0810 is planned. Its touches (36 files, no folder), forecast, and cost of delay are written. The reasons are under `### Planning` in its Notes.

## Figures

| Figure | Value | Why |
|--------|-------|-----|
| Forecast | 2h30m, delivery 2026-11-17T09:47:00Z | `flai forecast` gives 1h16m. I doubled it because three of E-0034's four serve and check stories took 1.9 to 2.7 times their forecasts. |
| Cost of delay | 175.37 USD a week | `flai cod`'s share of E-0034's 1500 USD a week, which comes from your 10h time lost per cycle. |
| After | S-0806, S-0804, as you declared them | S-0806 parses the `runs:` key this story reads. S-0804 builds the job, its record, and the send-back this story reuses. |

## Tasks and layers

1. Layer 1, which can run together:
   - T-3475: the move checks `kluster-verified` where a transition names it, at the rule's severity. Today `ReviewWarnings` warns on every move into a waiting state instead.
   - T-3479: `flai/internal/workflow/runs.go` says what a state runs on entry and which transition sends a story back from it.
2. Layer 2: T-3482, after T-3479. On entering a state with `runs:`, flai serve starts a state job. The job rebases the branch, runs the state's suites through `flai verify` limited to them, records each step, and on a failure sends the story back by the state's send-back transition, with the thread its agent starts on.
3. Layer 3: T-3484, after T-3482. The job runs the `kluster:` `deploy` and `checks` steps. It skips them, recording why, when there is no deploy, the host action is off, or no agent namespace is free. The deploy on `AwaitsAcceptance` in `kluster.go` is removed.
4. Layer 4: T-3487, after T-3475 and T-3484. The design and the guides are updated.

## Assumptions

- S-0805's suite of each tier reaches main before this story, through S-0806's `after`. So S-0805 is not in this story's `after`.
- If neither S-0805 nor S-0807 adds a way to limit `flai verify` to some suites, T-3482 adds it. If S-0807 does add one first, T-3482 uses it.
- The job reuses S-0804's resolve-job rebase, record, and `flai/internal/sendback` package as they are planned in T-3428 and T-3429. If S-0804 lands them in another shape, T-3482 follows that shape.
- A failed deploy or failed checks keeps today's kluster send-back thread.
- A state that runs something but has no send-back transition records the failure and moves nothing.
- "No agent namespace is available" means `kluster.agent_namespaces` is 0, or every namespace it allows is held by another story. Today the review deploy does not read that setting.
- The dashboard's view of the job is left to S-0814.

## Proposals

- No split, merge, add, or drop is proposed for this story.
- `flai check` warns `wip.after-hub` on S-0806, which five open stories name in `after`. This story does need S-0806's parse, so its `after` should stay.

### 2026-11-16T23:55:31Z orchestrator
Approved. The plan fits S-0810.

- The five tasks cover the five criteria. T-3475 and T-3479 share no path, so layer 1 runs together.
- 36 file touches, no folder. The forecast is 2h30m, flai's 1h16m doubled on three sibling serve stories. The value of 175.37 USD a week is flai's figure.
- `after` on S-0806 and S-0804 stays. This story reads the `runs:` parse and the resolve job.
- The assumptions are readings inside the story's scope. They follow ADR-0312 §4 and §5: a `kluster:` step is skipped, never held, where no cluster can run it.

I finalize S-0810 now. It stays in the backlog, held by `after` on S-0806 and S-0804.

Source: design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md § Decision

### 2026-11-16T23:56:23Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks, touches, and figures fit S-0810.
