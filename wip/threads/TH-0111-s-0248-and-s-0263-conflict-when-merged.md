---
id: TH-0111
title: S-0248 and S-0263 conflict when merged
anchor:
  path: wip/kanban/stories/S-0263-add-plan-to-the-context-menu-for-epic-and-story-cards.md
  item: S-0263
status: answered
participants: [flai, agent-S-0263]
created: 2026-10-04T23:47:22Z
updated: 2026-10-04T23:47:49Z
---

# TH-0111 S-0248 and S-0263 conflict when merged

On wip/kanban/stories/S-0263-add-plan-to-the-context-menu-for-epic-and-story-cards.md.

## Entries

### 2026-10-04T23:47:22Z flai
A trial merge of story/S-0248 with story/S-0263 at flai stream sync conflicts in:

- `design/issues/I-0057-flai-check-strict-stops-a-story-s-close-out-on-wip-findings-outside-the-story.md`
- `design/issues/summary.md`

Whichever of S-0248 and S-0263 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-04T23:47:46Z flai
Resolved: story/S-0248 and story/S-0263 merge cleanly at the sync of S-0263

### 2026-10-04T23:47:49Z agent-S-0263
S-0263 narrowed its change: its bump of I-0057 is reverted on story/S-0263 (101b19b), so S-0248 keeps the issue and the summary, and the next sync found the two merging cleanly. S-0263's occurrence (close-out stopped on main's own S-0173, E-0015, and TH-0094 findings) is kept in its narrative, to add to I-0057 once S-0248 is accepted.
