---
id: TH-0524
title: S-0404 is verified and its criteria are met, but metrics.md now says the conflict log drops a cut-off last line, which its code does not do
anchor:
  path: wip/kanban/stories/S-0404-metrics-md-defines-the-queue-wait-per-story-the-batch-size-and-the-bisections-per-week-by-an-adr-and-flai-stats-and-the-dashboard-s-charts-report-them.md
  item: S-0404
status: open
participants: [orchestrator]
created: 2026-10-20T01:08:41Z
updated: 2026-10-20T01:08:41Z
---

# TH-0524 S-0404 is verified and its criteria are met, but metrics.md now says the conflict log drops a cut-off last line, which its code does not do

On wip/kanban/stories/S-0404-metrics-md-defines-the-queue-wait-per-story-the-batch-size-and-the-bisections-per-week-by-an-adr-and-flai-stats-and-the-dashboard-s-charts-report-them.md.

## Entries

### 2026-10-20T01:08:41Z orchestrator
**Recommended:** S-0404's agent reverts the change to the Conflict log row in `design/system/metrics.md`. It should also reword ADR-0177's "kept as the conflict log is" to say the queue's log follows the conflict log only in being append-only and never trimmed. Then let the close-out's verify pass at the new head, and I will accept.

## The one blocker

The diff changes the Coordination table's **Conflict log** row in `design/system/metrics.md` from

> Append-only and never trimmed.

to

> Append-only and never trimmed; a last line a write was cut off in is skipped when read and dropped by the next append.

The conflict log's code does neither:

- `AppendConflict` in `flai/internal/storygit/conflicts.go` only appends.
- `ReadConflicts` there returns an error on a line it cannot parse.

The branch changes nothing under `flai/internal/storygit`. `metrics.md` is the contract between `flai stats` and the dashboard, and changes only as an ADR records. ADR-0177 decides the queue's log, not a change to the conflict log.

If you mean the conflict log to get the same cut-off handling, that is a change outside this story's touches. It would be a story of its own.

## Everything else is ready

- **Verify:** I ran `flai verify S-0404` again at the head, f24cb90e4437, against main 7fd31f6. It passed every step: flaiover 1m22s, integration 5m25s and smoke 4m36s. The branch contains main, and the dry run shows no other blocker.
- **Criteria 1 to 5:** the verifier matched each to changed files, and every box is ticked.
  - 1: ADR-0177, accepted, and the new Merge queue section of metrics.md.
  - 2: `mergequeue.Update`'s event log under `.flai-cache/merge-queue/events.jsonl`.
  - 3: `queue.weeks[]` and `items[].queue_wait_seconds` in `flai stats`.
  - 4: the `queue` chart beside Agent Waiting.
  - 5: the dashboard and CLI designs and the user guides.
