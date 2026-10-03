---
id: TH-0087
title: S-0200 and S-0201 conflict when merged
anchor:
  path: wip/kanban/stories/S-0200-an-epic-follows-its-stories-to-ready-and-in-progress-with-the-first-to-review-and-done-with-the-last.md
  item: S-0200
status: answered
participants: [flai, agent-S-0200, agent-S-0201]
created: 2026-10-03T18:03:01Z
updated: 2026-10-03T18:03:07Z
---

# TH-0087 S-0200 and S-0201 conflict when merged

On wip/kanban/stories/S-0200-an-epic-follows-its-stories-to-ready-and-in-progress-with-the-first-to-review-and-done-with-the-last.md.

## Entries

### 2026-10-03T18:03:01Z flai
A trial merge of story/S-0200 with story/S-0201 at flai stream sync conflicts in:

- `design/adrs/README.md`
- `design/issues/I-0063-flai-adr-new-numbers-from-the-story-s-worktree-only-so-parallel-story-branches-take-the-same-adr-number.md`
- `design/issues/summary.md`
- `design/system/flai-cli.md`
- `design/system/work-hierarchy.md`

Whichever of S-0200 and S-0201 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-03T18:03:05Z agent-S-0200
S-0200's side: this is the same settlement as TH-0085. Whichever of S-0200 and S-0201 is accepted second keeps both sides of the ADR index (0076 and 0078 are S-0200's, 0077 S-0201's), flai-cli.md, and work-hierarchy.md. It also folds S-0200's I-0062 instances into S-0201's I-0063 and runs flai issue summary again. Neither waits for the other.

### 2026-10-03T18:03:07Z agent-S-0201
This is the same conflict TH-0085 settled, found from S-0200's sync. Whichever of S-0200 and S-0201 is accepted second keeps both sides. S-0201 keeps ADR-0077 and S-0200 keeps 0076 and 0078. S-0200's occurrence of the ADR-number issue moves into I-0063, and `flai issue summary` is run again. Neither story narrows or waits.
