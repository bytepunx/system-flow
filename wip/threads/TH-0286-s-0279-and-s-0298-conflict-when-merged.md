---
id: TH-0286
title: S-0279 and S-0298 conflict when merged
anchor:
  path: wip/kanban/stories/S-0298-the-host-updates-tab-allows-an-operator-to-rollback-to-a-previous-version.md
  item: S-0298
status: resolved
participants: [flai, agent-S-0298]
created: 2026-10-07T14:32:32Z
updated: 2026-10-07T14:42:46Z
---

# TH-0286 S-0279 and S-0298 conflict when merged

On wip/kanban/stories/S-0298-the-host-updates-tab-allows-an-operator-to-rollback-to-a-previous-version.md.

## Entries

### 2026-10-07T14:32:32Z flai
A trial merge of story/S-0279 with story/S-0298 at flai stream sync conflicts in:

- `design/adrs/README.md`

Whichever of S-0279 and S-0298 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-07T14:32:58Z agent-S-0298
Additive, nothing to narrow: S-0298 adds ADR-0117's row to `design/adrs/README.md`, and S-0279 adds ADR-0115's. Whichever is accepted second keeps both rows in number order.

### 2026-10-07T14:42:46Z flai
Resolved: S-0279 is done, no longer open, at the sync of S-0298
