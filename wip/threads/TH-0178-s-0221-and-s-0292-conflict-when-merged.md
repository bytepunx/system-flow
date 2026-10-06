---
id: TH-0178
title: S-0221 and S-0292 conflict when merged
anchor:
  path: wip/kanban/stories/S-0292-flai-serve-s-orchestrator-test-that-waits-for-the-stopped-run-s-activity-fails-under-the-close-out-s-full-parallel-test-run.md
  item: S-0292
status: resolved
participants: [flai, alex]
created: 2026-10-06T11:40:43Z
updated: 2026-10-06T11:43:00Z
---

# TH-0178 S-0221 and S-0292 conflict when merged

On wip/kanban/stories/S-0292-flai-serve-s-orchestrator-test-that-waits-for-the-stopped-run-s-activity-fails-under-the-close-out-s-full-parallel-test-run.md.

## Entries

### 2026-10-06T11:40:43Z flai
A trial merge of story/S-0221 with story/S-0292 at flai stream sync conflicts in:

- `design/issues/summary.md`

Whichever of S-0221 and S-0292 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-06T11:43:00Z alex
Resolved.
