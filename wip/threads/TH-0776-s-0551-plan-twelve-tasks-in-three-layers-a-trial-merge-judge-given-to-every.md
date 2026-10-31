---
id: TH-0776
title: "S-0551 plan: twelve tasks in three layers, a trial-merge judge given to every caller of the holds"
anchor:
  path: wip/kanban/stories/S-0551-a-story-starts-on-a-trial-merge-instead-of-an-overlap-hold-and-after-holds-a.md
  item: S-0551
status: open
participants: [planner-S-0551]
created: 2026-10-31T07:46:00Z
updated: 2026-10-31T07:46:00Z
---

# TH-0776 S-0551 plan: twelve tasks in three layers, a trial-merge judge given to every caller of the holds

On wip/kanban/stories/S-0551-a-story-starts-on-a-trial-merge-instead-of-an-overlap-hold-and-after-holds-a.md.

## Entries

### 2026-10-31T07:46:00Z planner-S-0551
S-0551 now has twelve tasks in three layers, predicted touches, and a 2h forecast. The cost of delay value waits on TH-0772.

## Layers

| Layer | Task | What | After |
|-------|------|------|-------|
| 1 | T-2388 | Two ADRs (trial merge and `after` refining ADR-0046; the in-progress share metric), `workflow.md`, `strategic-agents.md`, `metrics.md` | none |
| 1 | T-2389 | `storygit/trial.go`: trial merge of the in-progress branch at main's head, conflicting paths | none |
| 1 | T-2391 | `workitem/hold.go`: the overlap verdict takes a trial-merge judge | none |
| 1 | T-2393 | `flai check` advisory on an `after` its Notes give no reason for | none |
| 2 | T-2395 | `flai serve` wiring and the journal's conflicting paths | T-2389, T-2391 |
| 2 | T-2404 | `wait_for_work`, story start, inbox, board, promotion candidates | T-2389, T-2391 |
| 2 | T-2406 | Planner role file, baseline, start prompt, `.claude/agents` copies | T-2388 |
| 2 | T-2408 | work-management and strategic-agents conventions, template, changelog | T-2388 |
| 2 | T-2410 | `flai stats`: share with n in progress, ready wait p50 and p85 | T-2388 |
| 3 | T-2412 | Review the open chain S-0420 to S-0422 | T-2393, T-2408 |
| 3 | T-2414 | Dashboard flow chart | T-2410 |
| 3 | T-2416 | `docs/users/flai.md` and `flai-cli.md` | T-2393, T-2395, T-2404, T-2410 |

## Assumptions

- A ready story has no branch, so "the ready story's base" is main's head: the trial merge asks whether the in-progress branch merges cleanly there.
- When the trial merge cannot run (git older than 2.38, no branch yet), the story is held as today.
- `workitem` cannot import `storygit`, so the judge is a function the callers pass; with none, the verdict is today's.
- The ready wait is the existing queue time (ADR-0190); only the in-progress share is a new metric, and it gets its own ADR because `metrics.md` changes only with one.
- S-0417 to S-0419 are archived, so T-2412 reviews S-0420 to S-0422. It asks you on each story's own thread to add a reason or drop the `after`, since their words are yours.
- `.claude/agents/planner.md` is written through the permission prompt, which asks you.

## Proposals

1. **Drop the four declared folder touches** `flai/internal/serve`, `flai/internal/sync`, `flai/internal/claims`, and `design/workflows/roles`. The tasks name the files below them. Two of them name no package in the tree. The serve folder holds S-0546, S-0549, and S-0550 while S-0551 runs. I kept them because they are yours.
2. **Split off the stats and the dashboard** (T-2410, T-2414) into a story of their own, if you want the hold change to land sooner. They share no file with the hold work and need only T-2388's metric ADR.
