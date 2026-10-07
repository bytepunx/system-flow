---
id: TH-0242
title: S-0269 and S-0277 conflict when merged
anchor:
  path: wip/kanban/stories/S-0269-one-command-closes-a-task-flai-task-done-commits-syncs-moves-logs-widens-touches-checks-and-answers-the-inbox.md
  item: S-0269
status: open
participants: [flai]
created: 2026-10-07T03:13:11Z
updated: 2026-10-07T03:13:11Z
---

# TH-0242 S-0269 and S-0277 conflict when merged

On wip/kanban/stories/S-0269-one-command-closes-a-task-flai-task-done-commits-syncs-moves-logs-widens-touches-checks-and-answers-the-inbox.md.

## Entries

### 2026-10-07T03:13:11Z flai
A trial merge of story/S-0269 with story/S-0277 at flai stream sync conflicts in:

- `design/issues/I-0076-flai-check-finds-wip-overlap-outside-the-story-at-close-out.md`
- `design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md`

Whichever of S-0269 and S-0277 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.
