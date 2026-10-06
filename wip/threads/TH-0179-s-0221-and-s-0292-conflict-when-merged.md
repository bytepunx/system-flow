---
id: TH-0179
title: S-0221 and S-0292 conflict when merged
anchor:
  path: wip/kanban/stories/S-0292-flai-serve-s-orchestrator-test-that-waits-for-the-stopped-run-s-activity-fails-under-the-close-out-s-full-parallel-test-run.md
  item: S-0292
status: answered
participants: [flai, agent-S-0292]
created: 2026-10-06T11:43:22Z
updated: 2026-10-06T11:43:32Z
---

# TH-0179 S-0221 and S-0292 conflict when merged

On wip/kanban/stories/S-0292-flai-serve-s-orchestrator-test-that-waits-for-the-stopped-run-s-activity-fails-under-the-close-out-s-full-parallel-test-run.md.

## Entries

### 2026-10-06T11:43:22Z flai
A trial merge of story/S-0221 with story/S-0292 at flai stream sync conflicts in:

- `design/issues/summary.md`

Whichever of S-0221 and S-0292 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-06T11:43:32Z agent-S-0292
Same conflict as TH-0175 and TH-0178: only the updated line in design/issues/summary.md's front matter (I-0074, which S-0278 remediates). Neither story can narrow it. Whichever of S-0221 and S-0292 is accepted second keeps both stories' rows and the later updated line.
