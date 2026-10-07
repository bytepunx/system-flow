---
id: TH-0280
title: S-0279 and S-0293 conflict when merged
anchor:
  path: wip/kanban/stories/S-0293-flai-stats-classifies-a-story-run-s-tool-calls-so-the-ceremony-turns-e-0017-removes-are-measured-per-story-and-over-time.md
  item: S-0293
status: answered
participants: [flai, agent-S-0293]
created: 2026-10-07T09:48:01Z
updated: 2026-10-07T09:48:10Z
---

# TH-0280 S-0279 and S-0293 conflict when merged

On wip/kanban/stories/S-0293-flai-stats-classifies-a-story-run-s-tool-calls-so-the-ceremony-turns-e-0017-removes-are-measured-per-story-and-over-time.md.

## Entries

### 2026-10-07T09:48:01Z flai
A trial merge of story/S-0279 with story/S-0293 at flai stream sync conflicts in:

- `design/adrs/README.md`

Whichever of S-0279 and S-0293 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-07T09:48:10Z agent-S-0293
Additive, nothing to narrow. S-0293 adds ADR-0116's row to `design/adrs/README.md`, and S-0279 adds ADR-0115's. Whichever story is accepted second keeps both rows, in number order. In `design/system/work-hierarchy.md`, S-0293 only adds a `turns` bullet after `empty_wakes` and a `turns:` block to the usage example, so a merge keeps both sides there too.
