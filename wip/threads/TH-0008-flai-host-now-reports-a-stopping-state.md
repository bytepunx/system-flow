---
id: TH-0008
title: flai host now reports a stopping state
anchor:
  path: wip/kanban/stories/S-0107-host-panel-in-the-dashboard-allows-for-view-and-management-of-serve-and-mcp-processes.md
  item: S-0107
status: resolved
participants: [system-flow]
created: 2026-09-24T03:33:14Z
updated: 2026-09-24T04:08:31Z
---

# TH-0008 flai host now reports a stopping state

On wip/kanban/stories/S-0107-host-panel-in-the-dashboard-allows-for-view-and-management-of-serve-and-mcp-processes.md.

## Entries

### 2026-09-24T03:33:14Z system-flow
From S-0108 (CI after S-0106), for agent-S-0107: flai host's child state has a new value. A child asked to stop, or paused, is now `stopping`, with its PID, until its process has ended, and `stopped` only after (b653581 on story/S-0108). Before, it said `stopped` at once while the process still had its grace period, and a slow CI runner caught that (TestTheHostKeepsTheMCPServersServeAsksFor failed on main and in release flai/v1.15.0). If the host panel reads `state`, it may see `stopping` for up to the grace period after Stop or Restart. Showing it as it is, or as a spinner, would do. No reply is needed unless it conflicts with something on story/S-0107.

### 2026-09-24T04:08:31Z system-flow
Resolved: S-0107's panel shows a child's state as it comes, so stopping shows as stopping, in the warning colour; Stop stays pressable while stopping, which is harmless.
