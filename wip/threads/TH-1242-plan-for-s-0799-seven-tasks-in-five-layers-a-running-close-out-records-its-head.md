---
id: TH-1242
title: "Plan for S-0799: seven tasks in five layers; a running close-out records its head and step, and flai queue, the preview, verify.status, and the review page name it"
anchor:
  path: wip/kanban/stories/S-0799-after-a-sync-moves-a-story-s-head-the-review-page-flai-queue-and-the-acceptance.md
  item: S-0799
status: resolved
participants: [planner-S-0799, orchestrator]
created: 2026-11-16T01:48:17Z
updated: 2026-11-16T01:50:32Z
---

# TH-1242 Plan for S-0799: seven tasks in five layers; a running close-out records its head and step, and flai queue, the preview, verify.status, and the review page name it

On wip/kanban/stories/S-0799-after-a-sync-moves-a-story-s-head-the-review-page-flai-queue-and-the-acceptance.md.

## Entries

### 2026-11-16T01:48:17Z planner-S-0799
S-0799 is planned: 28 file touches, seven tasks in five layers, a forecast of 1h15m, and a cost of delay of 50 USD a week.

## Tasks and layers

| Layer | Task | What it does |
|-------|------|--------------|
| 1 | T-3409 | The close-out job's file records the head it runs at, and `closeout` reads a story's running job with its head and step |
| 2 | T-3414 | `verify.LastReportAt` gives a report `moved` (its commit is not the head) and `running` (the job at the head) |
| 3 | T-3415 | `flai queue` gives the new reason `verify_running`, naming the job and step, and an old `verify_failed` says it is from before the head moved |
| 3 | T-3416 | The acceptance preview carries `verify_running` beside its blockers, which are unchanged |
| 3 | T-3417 | `flai verify --last` and `verify.status` carry `moved` and `running` |
| 4 | T-3418 | The review page's Verification row and the merge queue panel show them |
| 5 | T-3419 | flai-cli.md, flai.md, and flaiover.md describe them, and I-0384 is closed |

Each task waits only for the output it reads. The three tasks of layer 3 share no file.

## Assumptions

- The job file records no step today. T-3409 reads the step from the job's log, where `flai verify` prints `passed <step>` as each step ends. The story's agent may write the step to the job file instead if that is more reliable.
- A running job blocks the same as the report it replaces: no blocker or eligibility rule changes. Only what flai says changes.
- `verify_running` is a new queue reason kind, not new text on `verify_failed`, so the dashboard can show a running close-out as such and not as a failure.
- No ADR: the story changes what flai reports, not a decision ADR-0229 or ADR-0291 records.
- No `after`. S-0796, in progress, changes `flai/internal/preview/accept.go`, `design/system/flai-cli.md`, and `docs/users/flai.md`. S-0798 changes the close-out. A trial merge judges those overlaps. T-3416 and T-3419 say to sync first.

## Figures

- Forecast 1h15m, raised from flai's 58m: the job must newly record its head and step, and the close-out runs the flaiover tiers beside integration and smoke.
- Cost of delay 50 USD a week, from the inputs flai carried from I-0384. It stands as computed.

Nothing here needs the operator's answer before the story is finalized.

### 2026-11-16T01:50:15Z orchestrator
Approved. The plan fits S-0799.

- The seven tasks cover both criteria. T-3415, T-3416, T-3417, and T-3418 each reproduce I-0384 at their own surface. T-3419 closes the issue.
- The five layers hold. T-3415, T-3416, and T-3417 share no file. Each `after` names the output it reads.
- Every touch names a file, as ADR-0188 asks. The overlaps with S-0796 and S-0798 are left to the trial merge (ADR-0230), and T-3416 and T-3419 say to sync first.
- The figures stand: 1h15m, with its reason, and 50 USD a week from the inputs flai carried from I-0384. I set nothing.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-16T01:50:32Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks, touches, and figures fit S-0799.
