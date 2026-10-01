---
id: TH-0054
title: S-0182 and S-0186 conflict when merged
anchor:
  path: wip/kanban/stories/S-0182-find-and-fix-how-flai-serve-1-23-0-started-an-agent-for-a-story-the-board-showed-held.md
  item: S-0182
status: resolved
participants: [flai, agent-S-0182]
created: 2026-10-01T09:53:13Z
updated: 2026-10-01T09:58:02Z
---

# TH-0054 S-0182 and S-0186 conflict when merged

On wip/kanban/stories/S-0182-find-and-fix-how-flai-serve-1-23-0-started-an-agent-for-a-story-the-board-showed-held.md.

## Entries

### 2026-10-01T09:53:13Z flai
A trial merge of story/S-0182 with story/S-0186 at flai stream sync conflicts in:

- `design/issues/summary.md`

Whichever of S-0182 and S-0186 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-01T09:53:53Z agent-S-0182
S-0182 changes one line in `design/issues/summary.md`: `flai issue close I-0050` removes I-0050's row. S-0186 closes I-0033 the same way. The two edits are to different rows of a table `flai issue` writes, so neither story narrows its change or waits for the other. Whichever is accepted second keeps both removals when it rebases. If the rest of the table differs, run any `flai issue` command again there so it rewrites the summary from the issue files.

### 2026-10-01T09:58:02Z flai
Resolved: S-0186 is done, no longer open, at the sync of S-0182
