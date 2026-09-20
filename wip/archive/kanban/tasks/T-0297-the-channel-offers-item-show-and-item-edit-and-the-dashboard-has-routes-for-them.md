---
id: T-0297
type: task
nature: feature
title: The channel offers item.show and item.edit, and the dashboard has routes for them
status: done
parent: S-0085
owner: alex
created: 2026-09-20T15:21:42Z
updated: 2026-09-20T15:36:25Z
transitions:
  - to: ready
    at: 2026-09-20T15:30:05Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T15:30:05Z
    by: system-flow
  - to: done
    at: 2026-09-20T15:36:25Z
    by: system-flow
stream: S-0085
tags: []
---
# T-0297 The channel offers item.show and item.edit, and the dashboard has routes for them

## Work
Two named methods: item.show (a read) and item.edit (a write with a request ID, conflict and refusal as for doc.save). Values only ever reach the command line as flag values or on standard input. flaiover: GET and PUT on the item's edit route, statuses 409 and 422 with their payloads. REQUIRED_METHODS and the contract test.

## Done when
- hostapi tests for the command lines and refusals; route tests through flai hostapi
- The Go and flaiover tests pass

## Notes
