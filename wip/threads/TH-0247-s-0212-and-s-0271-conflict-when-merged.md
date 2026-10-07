---
id: TH-0247
title: S-0212 and S-0271 conflict when merged
anchor:
  path: wip/kanban/stories/S-0271-criteria-and-narrative-state-are-commands-flai-story-tick-checks-a-criterion-and-flai-stream-state-writes-current-state-and-next-steps.md
  item: S-0271
status: resolved
participants: [flai]
created: 2026-10-07T07:27:31Z
updated: 2026-10-07T07:31:32Z
---

# TH-0247 S-0212 and S-0271 conflict when merged

On wip/kanban/stories/S-0271-criteria-and-narrative-state-are-commands-flai-story-tick-checks-a-criterion-and-flai-stream-state-writes-current-state-and-next-steps.md.

## Entries

### 2026-10-07T07:27:31Z flai
A trial merge of story/S-0212 with story/S-0271 at flai stream sync conflicts in:

- `design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md`

Whichever of S-0212 and S-0271 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-07T07:31:32Z flai
Resolved: story/S-0212 and story/S-0271 merge cleanly at the sync of S-0212
