---
id: T-2840
type: task
nature: remediation
title: Close I-0330 with what fixed it
status: backlog
parent: S-0655
owner: alex
created: 2026-11-10T07:00:55Z
updated: 2026-11-10T07:00:55Z
transitions: []
stream: S-0655
tags: [flaiover, testing, flaky-test]
touches: [design/issues/I-0330-a-flaiover-design-page-test-takes-up-the-draft-another-tab-writes-times-out-at.md, design/issues/summary.md]
after: [T-2838]
---
# T-2840 Close I-0330 with what fixed it

## Work

Close I-0330 with `flai issue close I-0330 --reason`. The reason names T-2838, its commit, the cause it found, and the test's duration before and after at the load it was measured at. This task waits for T-2838 because the reason reports T-2838's fix and figures. It is the second layer. `flai issue close` writes the issue file and regenerates `design/issues/summary.md`.

## Done when

- I-0330 is closed, and its reason says what fixed it, with T-2838's figures.
- `design/issues/summary.md` no longer lists I-0330 as open.
- `flai check --strict` is clean for the story.

## Notes
