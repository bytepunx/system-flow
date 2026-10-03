---
id: T-0744
type: task
nature: feature
title: An ADR records the planning fields, and the design, the template's conventions, and the user guide describe them
status: done
parent: S-0199
owner: alex
created: 2026-10-03T05:41:26Z
updated: 2026-10-03T06:20:04Z
transitions:
  - to: ready
    at: 2026-10-03T06:12:24Z
    by: agent-S-0199
  - to: in-progress
    at: 2026-10-03T06:12:24Z
    by: agent-S-0199
  - to: done
    at: 2026-10-03T06:20:04Z
    by: agent-S-0199
stream: S-0199
tags: []
touches: [design/adrs, design/system/work-hierarchy.md, design/system/project-manifest.md, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, template/template.yaml]
after: [T-0742]
usage:
  source: log
  seconds: 460
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 91
      output: 33560
      cache_read: 5426361
      cache_write: 136893
      cost: 2.5643
---
# T-0744 An ADR records the planning fields, and the design, the template's conventions, and the user guide describe them

## Work

Record the schema with `flai adr new`: draft, cost of delay inputs and value, forecast, `by` and `at`, `estimate` as the human's figure, the manifest's `planning`, the refused move and who may finalize. Describe the fields in `design/system/work-hierarchy.md` (schema, rules, the older-flai note), `design/system/project-manifest.md`, `design/system/flai-cli.md` and `design/system/workflow.md` where the move rules are, `docs/users/flai.md`, and the template's conventions where they describe items, copied to `design/conventions` above the marker. Waits for T-0742, so that it describes what the commands do; shares no path with T-0743, so the two can run together.

## Done when

The ADR is accepted and linked from the living design, and every document that describes an item's fields, `flai edit`, or the move to ready says what the story added.

## Notes
