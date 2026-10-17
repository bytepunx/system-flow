---
id: T-1887
type: task
nature: feature
title: components.list and flai components list carry each component's examples and the path of its examples file
status: done
parent: S-0443
owner: alex
created: 2026-10-16T21:38:33Z
updated: 2026-10-17T04:29:26Z
transitions:
  - to: ready
    at: 2026-10-17T04:20:48Z
    by: agent-S-0443
  - to: in-progress
    at: 2026-10-17T04:20:49Z
    by: agent-S-0443
  - to: done
    at: 2026-10-17T04:29:26Z
    by: agent-S-0443
stream: S-0443
tags: [cli]
touches: [flai/internal/hostapi/components.go, flai/internal/hostapi/components_test.go, flai/cmd/components.go, flai/cmd/components_test.go]
after: [T-1885]
usage:
  source: log
  seconds: 517
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 41
      output: 16232
      cache_read: 2040754
      cache_write: 69389
      cost: 1.1274
---
# T-1887 components.list and flai components list carry each component's examples and the path of its examples file

## Work

Add `examples` and `examples_path` to each component in `components.list`'s answer and in `flai components list --json`, and the count in its text output, from T-1885's `Examples` and `ExamplesPath`. The path is given whether or not the file exists, so the page knows where a first example goes.

Waits for T-1885, whose reader they call.

## Done when

Tests cover a component with examples and one without, in both. `flai test flai/internal/hostapi flai/cmd` passes.

## Notes

None.
