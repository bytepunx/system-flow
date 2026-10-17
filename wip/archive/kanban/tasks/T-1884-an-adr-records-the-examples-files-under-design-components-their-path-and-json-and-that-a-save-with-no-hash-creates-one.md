---
id: T-1884
type: task
nature: feature
title: An ADR records the examples files under design/components, their path and JSON, and that a save with no hash creates one
status: done
parent: S-0443
owner: alex
created: 2026-10-16T21:38:18Z
updated: 2026-10-17T04:04:03Z
transitions:
  - to: ready
    at: 2026-10-17T03:58:25Z
    by: agent-S-0443
  - to: in-progress
    at: 2026-10-17T03:58:26Z
    by: agent-S-0443
  - to: done
    at: 2026-10-17T04:04:03Z
    by: agent-S-0443
stream: S-0443
tags: [cli, dashboard]
touches: [design/adrs, design/adrs/README.md, design/system/component-library.md]
usage:
  source: log
  seconds: 337
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 50
      output: 19636
      cache_read: 2468792
      cache_write: 83943
      cost: 1.3639
---
# T-1884 An ADR records the examples files under design/components, their path and JSON, and that a save with no hash creates one

## Work

Write and accept the ADR with `flai adr new`, on TH-0483 answer 4 ("lets put supporting files under design/components"). It gives the layout, `design/components/<sub-project>/<component path, .svelte replaced by .examples.json>` unless a better one is found and argued, and the JSON shape. It records that the files are in `file.get`'s and `file.save`'s reach as part of the design folder, so ADR-0147's reach is unchanged, and that a save with no hash creates a file. It names the convention's second exception to "Everything is markdown" and `components.examples` with what it reports. Add the format to `design/system/component-library.md`.

Waits for nothing.

## Done when

The ADR is accepted and indexed, the design document has the format, and `flai check --strict` passes.

## Notes

`design/adrs` is a folder touch for the new ADR's file, in `claims.shared`.
