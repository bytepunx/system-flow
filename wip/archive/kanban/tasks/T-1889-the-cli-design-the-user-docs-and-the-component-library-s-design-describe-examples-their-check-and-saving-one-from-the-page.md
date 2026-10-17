---
id: T-1889
type: task
nature: feature
title: The CLI design, the user docs, and the component library's design describe examples, their check, and saving one from the page
status: done
parent: S-0443
owner: alex
created: 2026-10-16T21:38:41Z
updated: 2026-10-17T04:46:47Z
transitions:
  - to: ready
    at: 2026-10-17T04:41:00Z
    by: agent-S-0443
  - to: in-progress
    at: 2026-10-17T04:41:00Z
    by: agent-S-0443
  - to: done
    at: 2026-10-17T04:46:47Z
    by: agent-S-0443
stream: S-0443
tags: [cli, dashboard]
touches: [design/system/component-library.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flaiover.md, design/system/flaiover-dashboard.md, flai/cmd/check.go, flai/cmd/check_stats_test.go, design/system/workflow.md, docs/users/flai-reference.md]
after: [T-1888, T-1894]
usage:
  source: log
  seconds: 347
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 116
      output: 45787
      cache_read: 5756670
      cache_write: 195735
      cost: 3.1804
---
# T-1889 The CLI design, the user docs, and the component library's design describe examples, their check, and saving one from the page

## Work

Describe the examples files under `design/components/`, `components.examples`, examples in `flai components list`, and saving one from the Components page, in the four files named first. Add `PUT /api/components/examples` to the API table of `design/system/flaiover-dashboard.md`.

Waits for T-1888, the last of what they describe, and T-1894, so that they link the folder's README.

## Done when

The five files describe examples as built, and `flai test` on them passes.

## Notes

None.
