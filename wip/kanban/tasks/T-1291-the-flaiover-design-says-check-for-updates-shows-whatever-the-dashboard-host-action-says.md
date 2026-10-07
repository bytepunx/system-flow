---
id: T-1291
type: task
nature: remediation
title: The flaiover design says Check for updates shows whatever the dashboard host action says
status: done
parent: S-0314
owner: alex
created: 2026-10-07T23:41:57Z
updated: 2026-10-07T23:44:57Z
transitions:
  - to: ready
    at: 2026-10-07T23:43:47Z
    by: agent-S-0314
  - to: in-progress
    at: 2026-10-07T23:43:47Z
    by: agent-S-0314
  - to: done
    at: 2026-10-07T23:44:57Z
    by: agent-S-0314
stream: S-0314
tags: [dashboard, docs]
touches: [design/system/flaiover-dashboard.md]
usage:
  source: log
  seconds: 61
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 6
      output: 68
      cache_read: 312730
      cache_write: 3051
      cost: 0.1393
---
# T-1291 The flaiover design says Check for updates shows whatever the dashboard host action says

## Work

In `design/system/flaiover-dashboard.md` § Managing the dashboard from its own page (S-0081), the sentence "`HostPanel.svelte` … shows the running image and, once enabled, the four actions" is wrong once the fix lands. Say instead that it shows the running image, Check for updates, and Versions always, and Restart, Upgrade, and Stop once the `dashboard` action is enabled, naming S-0314. Bump the front matter's `updated`. Read § Versions in the same file and keep it consistent. `docs/users/flaiover.md` § Host already says Check for updates always works: confirm it still reads true and leave it unchanged.

Waits for nothing: first layer, beside the HostPanel task, with which it shares no path.

## Done when

- `design/system/flaiover-dashboard.md` says which Dashboard-area controls show with the action off and which once it is on, as `HostPanel.svelte` does after the fix
- Its `updated` is today's date
- `flai test design/system/flaiover-dashboard.md` passes the markdown lint

## Notes
