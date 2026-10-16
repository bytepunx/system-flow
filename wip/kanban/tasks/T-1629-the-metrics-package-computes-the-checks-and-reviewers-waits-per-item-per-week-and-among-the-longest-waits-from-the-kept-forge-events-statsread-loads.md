---
id: T-1629
type: task
nature: feature
title: The metrics package computes the checks and reviewers waits per item, per week, and among the longest waits, from the kept forge events statsread loads
status: backlog
parent: S-0396
owner: alex
created: 2026-10-16T04:39:03Z
updated: 2026-10-16T04:39:03Z
transitions: []
stream: S-0396
tags: [cli]
touches: [flai/internal/metrics/waiting.go, flai/internal/metrics/waiting_test.go, flai/internal/statsread/statsread.go, flai/internal/statsread/statsread_test.go]
after: [T-1628]
---
# T-1629 The metrics package computes the checks and reviewers waits per item, per week, and among the longest waits, from the kept forge events statsread loads

## Work

- In `flai/internal/statsread/statsread.go`, load S-0394's kept forge events with the other inputs, as the conflict log is loaded. A missing file reads as no events.
- In `flai/internal/metrics/waiting.go`, compute what T-1628 defines:
  - `wait_checks_seconds` and `wait_reviewers_seconds` per item
  - `checks` and `reviewers` per week
  - the waits of the kinds `checks` and `reviewers` among `longest()`
- Test each definition's cases in `waiting_test.go`:
  - a pass
  - a pending check across a push of a new head
  - changes requested
  - a story sent back
  - no events
  - a project in `local` mode
- Make sure the other waits' figures are unchanged.
- Waits for T-1628: it implements that definition.

## Done when

`flai stats --json` carries the new fields as `metrics.md` defines them, the existing waiting tests pass unchanged, and the new ones pass.

## Notes
