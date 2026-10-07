---
id: TH-0294
title: S-0293 and S-0298 conflict when merged
anchor:
  path: wip/kanban/stories/S-0298-the-host-updates-tab-allows-an-operator-to-rollback-to-a-previous-version.md
  item: S-0298
status: open
participants: [flai]
created: 2026-10-07T15:00:24Z
updated: 2026-10-07T15:00:24Z
---

# TH-0294 S-0293 and S-0298 conflict when merged

On wip/kanban/stories/S-0298-the-host-updates-tab-allows-an-operator-to-rollback-to-a-previous-version.md.

## Entries

### 2026-10-07T15:00:24Z flai
A trial merge of story/S-0293 with story/S-0298 at flai stream sync conflicts in:

- `design/adrs/README.md`
- `design/system/flai-cli.md`

Whichever of S-0293 and S-0298 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.
