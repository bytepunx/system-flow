---
id: TH-0177
title: S-0221 and S-0292 conflict when merged
anchor:
  path: wip/kanban/stories/S-0221-the-orchestrator-accepts-stories-in-review-when-permitted.md
  item: S-0221
status: resolved
participants: [flai, agent-S-0221]
created: 2026-10-06T11:39:23Z
updated: 2026-10-06T11:39:29Z
---

# TH-0177 S-0221 and S-0292 conflict when merged

On wip/kanban/stories/S-0221-the-orchestrator-accepts-stories-in-review-when-permitted.md.

## Entries

### 2026-10-06T11:39:23Z flai
A trial merge of story/S-0221 with story/S-0292 at flai stream sync conflicts in:

- `design/issues/summary.md`

Whichever of S-0221 and S-0292 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-06T11:39:29Z agent-S-0221
Resolved: design/issues/summary.md is written by flai issue commands: S-0221 changed it only by bumping I-0058, and S-0292 by its own issue record. Neither narrows; whichever is accepted second takes both rows (each row is the issue's own count line), as any summary.md conflict is settled. Nothing for the designer.
