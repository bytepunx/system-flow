---
id: TH-0287
title: S-0293 and S-0298 conflict when merged
anchor:
  path: wip/kanban/stories/S-0298-the-host-updates-tab-allows-an-operator-to-rollback-to-a-previous-version.md
  item: S-0298
status: answered
participants: [flai, agent-S-0298]
created: 2026-10-07T14:32:32Z
updated: 2026-10-07T14:33:00Z
---

# TH-0287 S-0293 and S-0298 conflict when merged

On wip/kanban/stories/S-0298-the-host-updates-tab-allows-an-operator-to-rollback-to-a-previous-version.md.

## Entries

### 2026-10-07T14:32:32Z flai
A trial merge of story/S-0293 with story/S-0298 at flai stream sync conflicts in:

- `design/adrs/README.md`
- `design/system/flai-cli.md`

Whichever of S-0293 and S-0298 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-07T14:33:00Z agent-S-0298
Additive, nothing to narrow. S-0298 adds ADR-0117's row to `design/adrs/README.md`, and S-0293 adds ADR-0116's; whichever is accepted second keeps both in number order. In `design/system/flai-cli.md`, S-0298 adds to the commands table's `flai dashboard`, `flai host`, `flai serve`, `flai hostapi`, `flai mcp` (the `versions` tool), and `flai self-upgrade` rows, and a `### Listing and installing a published release` subsection under Versions; a merge keeps both sides of any row both stories extend.
