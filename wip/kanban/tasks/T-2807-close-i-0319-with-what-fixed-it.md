---
id: T-2807
type: task
nature: remediation
title: Close I-0319 with what fixed it
status: backlog
parent: S-0640
owner: alex
created: 2026-11-09T22:08:00Z
updated: 2026-11-09T22:08:00Z
transitions: []
stream: S-0640
tags: [flaiover, testing]
touches: [design/issues/I-0319-a-flaiover-props-panel-test-times-out-at-5-000-ms-under-load-and-fails-the.md, design/issues/summary.md]
after: [T-2805]
---
# T-2807 Close I-0319 with what fixed it

## Work

Run `flai issue close I-0319 --reason "<text>"`. The reason names the change T-2805 made to the props-panel test, and the test's duration before and after it from T-2805's `## Notes`. It writes the issue file and regenerates `design/issues/summary.md`.

It waits for T-2805, because the reason states what T-2805 changed and measured.

## Done when

- I-0319's status is closed, with a reason that names the fix and both durations.
- `design/issues/summary.md` no longer lists I-0319.
- `flai check --strict` passes.

## Notes
