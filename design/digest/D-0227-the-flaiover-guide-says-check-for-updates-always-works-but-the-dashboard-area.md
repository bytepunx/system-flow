---
id: D-0227
title: The flaiover guide says Check for updates always works, but the Dashboard area hides it while the dashboard host action is off
kind: cause
sources: [I-0115]
topics: [flaiover, dashboard, sveltekit, code, all]
paths: [flaiover/src/lib/components/HostPanel.svelte, flaiover/src/lib/components/HostPanel.svelte.test.ts, design/system/flaiover-dashboard.md, design/issues/I-0115-the-flaiover-guide-says-check-for-updates-always-works-but-the-dashboard-area-hides-it-while-the-dashboard-host-action-is-off.md, design/issues/summary.md, design/issues/I-0104-flai-task-done-commits-everything-in-the-worktree-so-two-tasks-of-one-layer-cannot-be-closed-apart.md]
date: 2026-10-07
status: absorbed
absorbed_by: "design/system/flaiover-dashboard.md#writes-s-0075-adr-0029"
asked_again: []
cited: []
---

# D-0227 The flaiover guide says Check for updates always works, but the Dashboard area hides it while the dashboard host action is off

I-0115 was closed on 2026-10-07.

The cause:

```text
The flaiover guide says Check for updates always works, but the Dashboard area hides it while the dashboard host action is off
```

The fix:

```text
fixed by S-0314: HostPanel.svelte shows Check for updates outside the dashboard_enabled gate, beside Versions, as docs/users/flaiover.md says, since dashboard.check is a read the host API does not gate; HostPanel.svelte.test.ts pins the buttons with the action off and a check made with it off, and design/system/flaiover-dashboard.md no longer says the page shows the four actions only once enabled
```

Source: I-0115.
