---
id: T-1630
type: task
nature: feature
title: flai stats prints the checks and reviewers waits under the review line
status: done
parent: S-0396
owner: alex
created: 2026-10-16T04:39:09Z
updated: 2026-10-17T17:54:22Z
transitions:
  - to: ready
    at: 2026-10-17T17:36:08Z
    by: agent-S-0396
  - to: in-progress
    at: 2026-10-17T17:47:14Z
    by: agent-S-0396
  - to: done
    at: 2026-10-17T17:54:22Z
    by: agent-S-0396
stream: S-0396
tags: [cli]
touches: [flai/cmd/stats.go, flai/cmd/check_stats_test.go]
after: [T-1629]
usage:
  source: log
  seconds: 428
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 16
      output: 6490
      cache_read: 761641
      cache_write: 29823
      cost: 0.4634
---
# T-1630 flai stats prints the checks and reviewers waits under the review line

## Work

- In `printWaiting` in `flai/cmd/stats.go`, print one line for the checks wait and one for the reviewers wait, labelled and indented as the review line is, with their total and mean.
- Print them only when an item of the window had a pull request.
- Test both outputs, with and without such an item.
- Waits for T-1629: it prints that task's figures.

## Done when

`flai stats` prints the two lines for a window with pull requests, prints neither for one without, and the tests pass.

## Notes
