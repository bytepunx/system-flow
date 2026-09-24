---
id: T-0377
type: task
nature: feature
title: The dashboard asks the host through flai serve for status, version, upgrade, and restart
status: done
parent: S-0106
owner: alex
created: 2026-09-24T01:30:20Z
updated: 2026-09-24T01:44:48Z
transitions:
  - to: ready
    at: 2026-09-24T01:42:35Z
    by: system-flow
  - to: in-progress
    at: 2026-09-24T01:42:35Z
    by: system-flow
  - to: done
    at: 2026-09-24T01:44:48Z
    by: system-flow
stream: S-0106
tags: []
touches: [flai/internal/hostapi]
---
# T-0377 The dashboard asks the host through flai serve for status, version, upgrade, and restart

## Work

- `hostapi` methods `host.status` and `host.check` (reads) and `host.restart` and `host.upgrade` (writes), each running `flai host ... --json`.
- The writes gate on a new host action `host`, detached from the request, because restarting serve ends the connection they came on.

## Done when

- The contract test covers the four methods; hostapi tests pass.

## Notes
