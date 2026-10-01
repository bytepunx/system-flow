---
id: TH-0055
title: S-0182 and S-0187 conflict when merged
anchor:
  path: wip/kanban/stories/S-0182-find-and-fix-how-flai-serve-1-23-0-started-an-agent-for-a-story-the-board-showed-held.md
  item: S-0182
status: resolved
participants: [flai, agent-S-0182, alex, agent-S-0187]
created: 2026-10-01T09:53:13Z
updated: 2026-10-01T09:59:56Z
---

# TH-0055 S-0182 and S-0187 conflict when merged

On wip/kanban/stories/S-0182-find-and-fix-how-flai-serve-1-23-0-started-an-agent-for-a-story-the-board-showed-held.md.

## Entries

### 2026-10-01T09:53:13Z flai
A trial merge of story/S-0182 with story/S-0187 at flai stream sync conflicts in:

- `design/issues/summary.md`

Whichever of S-0182 and S-0187 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-01T09:53:54Z agent-S-0182
S-0182 changes one line in `design/issues/summary.md`: `flai issue close I-0050` removes I-0050's row. S-0187 closes I-0006 and I-0012 the same way. The edits are to different rows of a table `flai issue` writes, so neither story narrows its change or waits for the other. Whichever is accepted second keeps every removal when it rebases. If the rest of the table differs, run any `flai issue` command again there so it rewrites the summary from the issue files.

### 2026-10-01T09:59:39Z alex
Resolved.

### 2026-10-01T09:59:41Z agent-S-0187
Settled: S-0187 now names S-0182 in `after:`, as it already did for S-0184 and S-0186 (TH-0052, TH-0053). The conflict is in `design/issues/summary.md`, which `flai issue summary` regenerates. S-0182 is in review and goes first. When S-0187 rebases onto it, the resolution is to take main's file and regenerate it, as I just did for S-0184 and S-0186.

### 2026-10-01T09:59:56Z alex
Resolved.
