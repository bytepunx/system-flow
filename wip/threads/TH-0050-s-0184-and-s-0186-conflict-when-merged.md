---
id: TH-0050
title: S-0184 and S-0186 conflict when merged
anchor:
  path: wip/kanban/stories/S-0186-the-root-page-asks-for-the-projects-once-per-load.md
  item: S-0186
status: resolved
participants: [flai, agent-S-0186, alex]
created: 2026-10-01T09:24:20Z
updated: 2026-10-01T09:53:23Z
---

# TH-0050 S-0184 and S-0186 conflict when merged

On wip/kanban/stories/S-0186-the-root-page-asks-for-the-projects-once-per-load.md.

## Entries

### 2026-10-01T09:24:20Z flai
A trial merge of story/S-0184 with story/S-0186 at flai stream sync conflicts in:

- `design/issues/summary.md`

Whichever of S-0184 and S-0186 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-01T09:24:45Z agent-S-0186
This is the same kind of conflict as TH-0048. The only path that conflicts is `design/issues/summary.md`, which flai regenerates from the issue files. S-0184 closes I-0025 and I-0029, and S-0186 closes I-0033, so the issue files themselves don't conflict.

Whichever story rebases second resolves it without changing either story's scope: `flai issue summary`, `git add design/issues/summary.md`, `git rebase --continue`. S-0184 is in review and S-0186 is about to go there. If S-0184 is accepted first, S-0186's next `flai stream sync`, or its acceptance, does this.

### 2026-10-01T09:53:23Z alex
Resolved.
