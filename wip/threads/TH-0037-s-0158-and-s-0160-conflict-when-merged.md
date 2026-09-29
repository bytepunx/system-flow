---
id: TH-0037
title: S-0158 and S-0160 conflict when merged
anchor:
  path: wip/kanban/stories/S-0160-a-flai-process-starts-in-milliseconds-whatever-the-host-s-path.md
  item: S-0160
status: resolved
participants: [flai]
created: 2026-09-29T19:40:16Z
updated: 2026-09-29T19:40:44Z
---

# TH-0037 S-0158 and S-0160 conflict when merged

On wip/kanban/stories/S-0160-a-flai-process-starts-in-milliseconds-whatever-the-host-s-path.md.

## Entries

### 2026-09-29T19:40:16Z flai
A trial merge of story/S-0158 with story/S-0160 at flai stream sync conflicts in:

- `design/system/server-performance.md`

Whichever of S-0158 and S-0160 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-09-29T19:40:44Z flai
Resolved: story/S-0158 and story/S-0160 merge cleanly at the sync of S-0160
