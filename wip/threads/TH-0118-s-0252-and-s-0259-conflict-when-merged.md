---
id: TH-0118
title: S-0252 and S-0259 conflict when merged
anchor:
  path: wip/kanban/stories/S-0259-the-workflow-menu-has-a-planner-page-showing-its-status-activity-log-and-runs.md
  item: S-0259
status: resolved
participants: [flai]
created: 2026-10-05T00:28:42Z
updated: 2026-10-05T00:28:53Z
---

# TH-0118 S-0252 and S-0259 conflict when merged

On wip/kanban/stories/S-0259-the-workflow-menu-has-a-planner-page-showing-its-status-activity-log-and-runs.md.

## Entries

### 2026-10-05T00:28:42Z flai
A trial merge of story/S-0252 with story/S-0259 at flai stream sync conflicts in:

- `design/issues/I-0057-flai-check-strict-stops-a-story-s-close-out-on-wip-findings-outside-the-story.md`
- `design/issues/summary.md`

Whichever of S-0252 and S-0259 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-05T00:28:53Z flai
Resolved: story/S-0252 and story/S-0259 merge cleanly at the sync of S-0259
