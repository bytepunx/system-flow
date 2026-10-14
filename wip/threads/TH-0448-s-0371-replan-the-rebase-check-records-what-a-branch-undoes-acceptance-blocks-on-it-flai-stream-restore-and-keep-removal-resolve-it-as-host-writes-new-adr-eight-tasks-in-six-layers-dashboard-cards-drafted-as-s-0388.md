---
id: TH-0448
title: "S-0371 replan: the rebase check records what a branch undoes, acceptance blocks on it, flai stream restore and --keep-removal resolve it, as host writes (new ADR), eight tasks in six layers; dashboard cards drafted as S-0388"
anchor:
  path: wip/kanban/stories/S-0371-flai-stream-sync-replays-a-story-s-revert-of-its-own-change-over-an-identical-change-on-main-undoing-main-s.md
  item: S-0371
status: open
participants: [planner-S-0371, orchestrator]
created: 2026-10-14T21:02:51Z
updated: 2026-10-14T21:03:04Z
---

# TH-0448 S-0371 replan: the rebase check records what a branch undoes, acceptance blocks on it, flai stream restore and --keep-removal resolve it, as host writes (new ADR), eight tasks in six layers; dashboard cards drafted as S-0388

On wip/kanban/stories/S-0371-flai-stream-sync-replays-a-story-s-revert-of-its-own-change-over-an-identical-change-on-main-undoing-main-s.md.

## Entries

### 2026-10-14T21:02:51Z planner-S-0371
**Recommendation: approve the replan as drafted, and let me add three acceptance criteria to S-0371 (below).** It follows the shape you confirmed on TH-0420: you resolve each undone path from the dashboard, and flai does the git work.

## What changed from the first plan

- **The `git checkout` advice is gone.** Two flai commands replace it: `flai stream restore S-nnnn --path <p>`, and `flai accept --keep-removal <p> --reason`.
- **The findings are recorded per story** under `.flai-cache/`. The close-out syncs before review, so the acceptance's own rebase usually finds nothing. Without a record, the preview could not show the path.
- **The acceptance preview gives each path as a blocker and as data** (`undoes_main`, with main's lines the branch removes). `flai accept --dry-run` and the host read `accept.preview` share it. A refused acceptance exits 4 with the paths, as the markdown refusal does.
- **Two host writes:** `stream.restore`, and `keep_removals` on `accept.run`. flai and `accept.run` both refuse a keep with `by: orchestrator`.

## Tasks and layers

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-1476 ADR, refining ADR-0069, ADR-0093, ADR-0098 | `design/adrs/README.md` |
| 2 | T-1477 the check and the record, I-0124 reproduced | `storygit/undone.go`, `undone_test.go`, `sync.go`, `sync_test.go` |
| 3 | T-1574 `flai stream restore` (new) | `cmd/stream_restore.go`, `stream_restore_test.go`, `stream.go` |
| 3 | T-1479 preview blocker, refusal, `--keep-removal` | `preview/accept.go`, `cmd/accept.go`, `branch.go`, `accept_undone_test.go` |
| 4 | T-1478 sync and task done warning, naming restore | `cmd/stream_sync.go`, `stream_sync_test.go`, `stream.go` |
| 4 | T-1575 host writes (new) | `hostapi/writes.go`, `writes_test.go`, `flaiover/src/lib/server/agent.ts` |
| 5 | T-1480 docs, reference and flag index regenerated | `flai-cli.md`, `workflow.md`, `strategic-agents.md`, `docs/users/flai.md`, `flai-reference.md`, `docs/operators/settings.md` |
| 6 | T-1481 close I-0124 | the issue file, `summary.md` |

- T-1478 waits for T-1574. Both change `stream.go`, and the warning names T-1574's command.
- The tasks in each layer share no path. Every touch is a file.
- I rewrote all six earlier tasks to this shape. I wrote them, so no one else's words changed.

## Figures

- **Forecast 1h15m**, delivery 2026-10-14T22:20Z. flai gave 40m at size 27. I raised it for eight tasks in six layers, a new command, a record, an override, and two host writes.
- **Cost of delay 25 USD a week**, unadjusted, from flai's 10m input.

## Assumptions

- **The record is host-local.** It lives under `.flai-cache/`, like `flai verify`'s last result, so an acceptance on another clone does not see it. Its own rebase still runs the check.
- **`T-1575` adds `stream.restore` to `flaiover/src/lib/server/agent.ts`.** `TestTheDashboardAsksForExactlyWhatFlaiOffers` fails while flai offers a method the dashboard does not list. That file then holds a flaiover story that touches it while S-0371 is in progress.
- **`--keep-removal` also reaches `flai move S-nnnn done`**, which shares accept's flags.
- **No convention change.** The warning and the blocker say what to do.

## For you to decide

1. **Three criteria for S-0371.** You finalized it while I planned, so I have not changed its words. I propose:
   - `flai stream sync` and `flai task done` warn on each path a rebased branch undoes, and record it.
   - The acceptance preview lists each recorded path as a blocker with its lines as data, and `flai accept` refuses it unless the path is restored with `flai stream restore` or kept with `--keep-removal` and a reason, which the orchestrator cannot give.
   - The host writes `stream.restore` and `accept.run`'s `keep_removals` expose both to the dashboard.
2. **S-0388, drafted.** It covers the review page's cards (Restore, Keep with a reason, Send back) and the board's badge. It comes after S-0371, as TH-0420 said. S-0378 is done, so it waits for S-0371 alone. It has no tasks, forecast, or cost of delay yet: its own planner run adds them.

I propose no split, merge, or drop of S-0371's tasks.

### 2026-10-14T21:03:04Z orchestrator (recommendation)
**Recommendation:** approve the replan as drafted, and add the planner's three criteria to S-0371. Today its criteria still describe only "the cause no longer occurs" and closing I-0124. Without the three, its agent and I would accept against criteria that do not name the restore, the keep, or the host writes.

## Why this needs you

- **S-0371 is finalized, so its words are yours.** The planner may not rewrite them, and neither may I (strategic-agents.md § Rules).
- **Timing.** S-0371 is first in ready. It is held only by its overlap with S-0304 (MS-0112). If its agent starts before the criteria change, it works to the old ones. Add them now, or move S-0371 back to backlog until you do.

## What fits

- **It follows TH-0420 as you confirmed it.** Each undone path is a preview blocker that you resolve with Restore (`flai stream restore`) or Keep with a reason (`--keep-removal`), or by sending the story back. Both are host writes, and a keep is refused to me.
- **Eight tasks in six layers.** Each layer's tasks share no path, and every touch is a file. T-1478 waits for T-1574 on `stream.go`, with the reason given.
- **The figures are flai's.** Forecast 1h15m, raised from flai's 40m with a stated reason. Cost of delay 25 USD a week, from flai's 10m input.
- **S-0388** carries the review-page cards and the board badge, after S-0371. Under `plan_backlog_stories`, I start its planner once flai lists it.

Source: design/conventions/strategic-agents.md § Rules
