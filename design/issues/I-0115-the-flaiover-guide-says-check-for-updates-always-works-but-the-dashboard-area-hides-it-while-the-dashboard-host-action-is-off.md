---
id: I-0115
title: The flaiover guide says Check for updates always works, but the Dashboard area hides it while the dashboard host action is off
class: defect
status: closed
count: 1
cost: 5m
first_reported: 2026-10-07T15:11:31Z
last_reported: 2026-10-07T15:11:31Z
updated: 2026-10-07T23:45:16Z
---

# I-0115 The flaiover guide says Check for updates always works, but the Dashboard area hides it while the dashboard host action is off

## Description
The flaiover guide says Check for updates always works, but the Dashboard area hides it while the dashboard host action is off

## Instances

### 2026-10-07T15:11:31Z
Story: S-0298.
Found by T-1046's sub-agent in S-0298: docs/users/flaiover.md (Host) says the Dashboard area's Check for updates always works, but HostPanel.svelte shows it only inside {#if status.dashboard_enabled}. Already so on main before S-0298. dashboard.check is a read that the host API does not gate, so either the button should show with the action off, as S-0298's Versions does, or the guide should say it is hidden.

## Remediation

Story S-0314 remediates this issue, created from it at 2026-10-07T18:59:44Z.
Closed 2026-10-07T23:45:16Z: fixed by S-0314: HostPanel.svelte shows Check for updates outside the dashboard_enabled gate, beside Versions, as docs/users/flaiover.md says, since dashboard.check is a read the host API does not gate; HostPanel.svelte.test.ts pins the buttons with the action off and a check made with it off, and design/system/flaiover-dashboard.md no longer says the page shows the four actions only once enabled
