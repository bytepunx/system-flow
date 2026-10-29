---
id: T-2342
type: task
nature: feature
title: The Dashboard area shows the running dashboard version, and Check for updates names the running and the available version
status: done
parent: S-0544
owner: alex
created: 2026-10-29T17:40:43Z
updated: 2026-10-29T17:47:42Z
transitions:
  - to: ready
    at: 2026-10-29T17:43:48Z
    by: agent-S-0544
  - to: in-progress
    at: 2026-10-29T17:43:52Z
    by: agent-S-0544
  - to: done
    at: 2026-10-29T17:47:42Z
    by: agent-S-0544
stream: S-0544
tags: []
touches: [flaiover/src/routes/api/dashboard/+server.ts, flaiover/src/routes/api/dashboard/dashboard.test.ts, flaiover/src/lib/components/HostPanel.svelte, flaiover/src/lib/components/HostPanel.svelte.test.ts, docs/users/flaiover.md, design/system/flaiover-dashboard.md]
after: [T-2341]
usage:
  source: log
  seconds: 230
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 26
      output: 5435
      cache_read: 1895499
      cache_write: 43679
      cost: 0.8373
---
# T-2342 The Dashboard area shows the running dashboard version, and Check for updates names the running and the available version

## Work

`GET /api/dashboard` and its `check` answer `version`, the version of the dashboard answering, which is the container running. `dashboard.check` already answers `release`, the signed release the configured tag resolves to, and `image`, its reference. `HostPanel.svelte` shows the running version in its status line, and the check result says `Dashboard <running> is running; <available> is available.` or `Dashboard <running> is the latest.`, naming the image where no release is known. Docs: the Dashboard paragraph in `docs/users/flaiover.md` and S-0081's bullet in `design/system/flaiover-dashboard.md`.

Waits for T-2341: both change the same two documents.

## Done when

- The check result names the running and the available version, and a test pins each case.
- The status line shows the running version, and a test pins it.
- The docs say so.

## Notes
