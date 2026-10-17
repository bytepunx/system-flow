---
id: T-1892
type: task
nature: feature
title: An integration test asks the dev server on this repository's flaiover for each listed component's module and fails on one that does not compile
status: done
parent: S-0444
owner: alex
created: 2026-10-16T21:40:07Z
updated: 2026-10-17T05:31:08Z
transitions:
  - to: ready
    at: 2026-10-17T05:28:26Z
    by: agent-S-0444
  - to: in-progress
    at: 2026-10-17T05:28:26Z
    by: agent-S-0444
  - to: done
    at: 2026-10-17T05:31:08Z
    by: agent-S-0444
stream: S-0444
tags: [cli]
touches: [flai/internal/components/flaiover_test.go]
after: [T-1890]
usage:
  source: log
  seconds: 162
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 40
      output: 12976
      cache_read: 2318808
      cache_write: 67012
      cost: 1.1349
---
# T-1892 An integration test asks the dev server on this repository's flaiover for each listed component's module and fails on one that does not compile

## Work

`flaiover_test.go`, skipped under `-short`, without `node`, and without `flaiover/node_modules`: find the repository's root from the test's folder and load its manifest. Start the dev server on flaiover as `flai components serve` does, and request each listed component's module through Vite's transform. Fail with the component's path and Vite's error for any that does not compile.

Waits for T-1890, whose manifest entry the test reads.

## Done when

`flai test --all flai/internal/components` runs it and it passes on this host.

## Notes

None.
