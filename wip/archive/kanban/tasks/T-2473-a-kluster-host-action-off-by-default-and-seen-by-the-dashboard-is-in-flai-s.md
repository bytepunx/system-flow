---
id: T-2473
type: task
nature: feature
title: A kluster host action, off by default and seen by the dashboard, is in flai's list of host actions with a plain description
status: done
parent: S-0577
owner: alex
created: 2026-11-06T21:42:33Z
updated: 2026-11-07T01:08:44Z
transitions:
  - to: ready
    at: 2026-11-07T00:55:42Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T00:55:43Z
    by: system-flow
  - to: done
    at: 2026-11-07T01:08:44Z
    by: agent-S-0577
stream: S-0577
tags: [kluster, hostapi, host-actions]
touches: [flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/cmd/serve_actions_test.go]
usage:
  source: log
  seconds: 781
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 11
      output: 3364
      cache_read: 1502498
      cache_write: 16797
      cost: 0.4863
---
# T-2473 A kluster host action, off by default and seen by the dashboard, is in flai's list of host actions with a plain description

## Work

In `flai/internal/hostapi/writes.go`, add `ActionKluster` (`kluster`) beside `ActionComponents`, and its entry in the `Actions` map. Its description says, in the map's plain words, that flai serve deploys each story in review into its namespace on the project's kluster, tears it down when it leaves review, promotes main on acceptance, and lets the dashboard manage the cluster. Leave it out of `shellOnly`, so `DashboardSees` lists it, and off by default like every action. `flai serve enable kluster` then turns it on with no further change, and the Features page lists it.

Test in `writes_test.go` that the action is listed, is seen by the dashboard, is off for a project with no host configuration, and that `settings.action` turns it on and off.

This is the first layer: it waits for no task, and shares no path with the manifest task.

## Done when

- `kluster` is a host action, off by default, with its description, and the dashboard sees it.
- `settings.action` turns it on and off in a test.
- `flai test` on the two paths passes.

## Notes
