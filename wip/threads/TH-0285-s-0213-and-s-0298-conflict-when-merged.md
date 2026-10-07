---
id: TH-0285
title: S-0213 and S-0298 conflict when merged
anchor:
  path: wip/kanban/stories/S-0298-the-host-updates-tab-allows-an-operator-to-rollback-to-a-previous-version.md
  item: S-0298
status: resolved
participants: [flai, agent-S-0298]
created: 2026-10-07T14:32:32Z
updated: 2026-10-07T14:51:27Z
---

# TH-0285 S-0213 and S-0298 conflict when merged

On wip/kanban/stories/S-0298-the-host-updates-tab-allows-an-operator-to-rollback-to-a-previous-version.md.

## Entries

### 2026-10-07T14:32:32Z flai
A trial merge of story/S-0213 with story/S-0298 at flai stream sync conflicts in:

- `design/adrs/README.md`
- `design/system/flaiover-dashboard.md`

Whichever of S-0213 and S-0298 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-07T14:32:57Z agent-S-0298
Additive, nothing to narrow. S-0298 adds ADR-0117's row to `design/adrs/README.md`, and S-0213 adds its own ADR's row; the numbers differ, so whichever is accepted second keeps both rows in number order. In `design/system/flaiover-dashboard.md`, S-0298 adds "Versions" to the `/host` row, the `versions` action and optional `tag`/`version` to the `/api/dashboard` and `/api/host` bullets, and a new `### Deploying a chosen release (S-0298)` section under Writes; a merge keeps both sides.

### 2026-10-07T14:51:27Z flai
Resolved: story/S-0213 and story/S-0298 merge cleanly at the sync of S-0298
