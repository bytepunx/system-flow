---
id: T-2632
type: task
nature: improvement
title: Check the work-item rules and the markdown lint over a given set of items
status: done
parent: S-0599
owner: alex
created: 2026-11-07T23:48:14Z
updated: 2026-11-08T00:05:09Z
transitions:
  - to: ready
    at: 2026-11-07T23:52:13Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T23:52:15Z
    by: system-flow
  - to: done
    at: 2026-11-08T00:05:09Z
    by: agent-S-0599
stream: S-0599
tags: [cli]
touches: [flai/internal/check/items.go, flai/internal/check/items_test.go]
usage:
  source: log
  seconds: 774
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 19
      output: 3261
      cache_read: 3417748
      cache_write: 55461
      cost: 1.0718
---
# T-2632 Check the work-item rules and the markdown lint over a given set of items

## Work

Export an entry in `flai/internal/check/check.go`, such as `check.Items(repo, items, now)`, that runs only what a forecast write can break, over the items given:

- `oneItem` for each item: its front matter, its forecast, and its body sections.
- The markdown lint (`mdlint`) of each item's file, with the project's configuration.

It reads no other item, walks no folder, starts no `node`, and scans no digest. It returns a `*Result` in the shape `Run` gives, so a caller can compare its findings as `itemedit.Apply` does.

First layer: it waits for no task, and shares no path with the other first-layer task.

## Done when

- A test in `flai/internal/check/items_test.go` shows `Items` reports a bad forecast and a markdown finding on the items given, and nothing about an item it was not given.
- On the same items, its findings match what `Run` reports for those items' files under those rules.
- `flai test flai/internal/check` passes.

## Notes
