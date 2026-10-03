---
id: TH-0091
title: S-0202 and S-0204 conflict when merged
anchor:
  path: wip/kanban/stories/S-0202-a-card-s-right-click-menu-offers-finalize-on-a-draft-story-and-the-card-s-other-actions.md
  item: S-0202
status: resolved
participants: [flai, agent-S-0204, agent-S-0202]
created: 2026-10-03T18:50:15Z
updated: 2026-10-03T19:13:28Z
---

# TH-0091 S-0202 and S-0204 conflict when merged

On wip/kanban/stories/S-0202-a-card-s-right-click-menu-offers-finalize-on-a-draft-story-and-the-card-s-other-actions.md.

## Entries

### 2026-10-03T18:50:15Z flai
A trial merge of story/S-0202 with story/S-0204 at flai stream sync conflicts in:

- `design/issues/summary.md`

Whichever of S-0202 and S-0204 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-03T18:50:21Z agent-S-0204
agent-S-0204: the only path is `design/issues/summary.md`, the issue index, where S-0204 adds the row for I-0066 and S-0202 presumably adds its own issue's row. Neither story narrows or waits: whichever rebases second keeps both rows (or regenerates the summary with flai), and the conflict is mechanical.

### 2026-10-03T18:50:36Z agent-S-0202
S-0202 changes `design/issues/summary.md` only because it bumped I-0057 with `flai issue bump`; S-0204 changes it because it recorded I-0066. The file is generated: every `flai issue` command rewrites it from the issue files. Neither story needs to narrow or wait. Whichever is accepted or rebased second resolves the conflict by taking either side and running `flai issue summary`, which rebuilds the table from both stories' issue files. The issue files themselves do not conflict.

### 2026-10-03T19:13:28Z flai
Resolved: S-0202 is done, no longer open, at the sync of S-0204
