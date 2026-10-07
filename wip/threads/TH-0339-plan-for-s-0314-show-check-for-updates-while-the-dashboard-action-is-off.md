---
id: TH-0339
title: "Plan for S-0314: show Check for updates while the dashboard action is off"
anchor:
  path: wip/kanban/stories/S-0314-the-flaiover-guide-says-check-for-updates-always-works-but-the-dashboard-area-hides-it-while-the-dashboard-host-action-is-off.md
  item: S-0314
status: resolved
participants: [planner-S-0314, orchestrator]
created: 2026-10-07T23:42:15Z
updated: 2026-10-07T23:42:23Z
---

# TH-0339 Plan for S-0314: show Check for updates while the dashboard action is off

On wip/kanban/stories/S-0314-the-flaiover-guide-says-check-for-updates-always-works-but-the-dashboard-area-hides-it-while-the-dashboard-host-action-is-off.md.

## Entries

### 2026-10-07T23:42:15Z planner-S-0314
Plan: make the code match the guide. `HostPanel.svelte` shows Check for updates whatever the `dashboard` host action says, as Versions already does, and the design's sentence is corrected.

## Tasks and layers

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1290 HostPanel shows Check for updates with the action off, with a test that reproduces I-0115 | `HostPanel.svelte`, `HostPanel.svelte.test.ts` | none |
| 1 | T-1291 `design/system/flaiover-dashboard.md` says which controls show with the action off | `design/system/flaiover-dashboard.md` | none |
| 2 | T-1292 close I-0115 with `flai issue close --reason` | the I-0115 file, `design/issues/summary.md` | T-1290, T-1291 |

T-1290 and T-1291 share no path and can run together.

## Figures

- Forecast 12m, delivery 2026-10-08T06:25:00Z, from `flai forecast` after the touches were set; kept.
- Cost of delay value 12.50 USD a week, from `flai cod` on the inputs flai set from I-0115; kept.

## Assumptions

- The fix is to show the button, not to change the guide: `dashboard.check` is an ungated read on the host (`design/system/dashboard-host-channel.md` item 10) and the route runs it with `repo().run`, so no host or route change is needed.
- With the action off, the check's "Not running; Restart starts it." should not name the hidden Restart; T-1290 rewords that case.
- `docs/users/flaiover.md` already reads true after the fix and stays unchanged.
- `HostProcesses.svelte` (co-change) is left out: its Check for upgrade is already outside its gate.

No task to split, merge, or drop.

### 2026-10-07T23:42:22Z orchestrator
Approved.

- T-1290 removes the cause with a test that reproduces I-0115 (criterion 1). T-1292 closes I-0115 (criterion 2).
- The fix follows the recorded design: `dashboard-host-channel.md` makes `dashboard.check` an ungated read. Showing the button with the action off brings the code in line with that design and with the guide, so no decision changes.
- T-1291 corrects the one design sentence that says otherwise. T-1290 and T-1291 share no path.
- The 12m forecast and the 12.50 USD a week value stand.

Source: design/system/dashboard-host-channel.md

### 2026-10-07T23:42:23Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories
