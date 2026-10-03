---
id: T-0769
type: task
nature: feature
title: Cost of delay inputs carry their own by and at, apart from the value's
status: ready
parent: S-0204
owner: alex
created: 2026-10-03T18:35:52Z
updated: 2026-10-03T18:36:18Z
transitions:
  - to: ready
    at: 2026-10-03T18:36:18Z
    by: agent-S-0204
stream: S-0204
tags: []
touches: [flai/internal/workitem, flai/internal/itemedit, design/adrs, design/system/work-hierarchy.md]
after: [T-0767]
---
# T-0769 Cost of delay inputs carry their own by and at, apart from the value's

## Work

As TH-0090 settles: `cost_of_delay.inputs` carry their own `by` and `at`, stamped when an input changes, and the block's `by` and `at` become the value's, stamped only when `value` changes, so the item page can show the planner's stamp and tell a value older than its inputs. A new ADR records it beside ADR-0074. `workitem` validates it (inputs need `by`/`at`, a value needs the block's), `itemedit` and `create.go` stamp it, `front-matter-fields.txt` lists the keys, and `design/system/work-hierarchy.md` shows it (and loses the conflict markers main carries there, I-0066). A block written by a flai with one stamp reads as the value's when it has a value, else the inputs'. Waits for T-0767, which changes how creation stamps the inputs, and for the thread's answer.

## Done when

- Behaviour tests in `flai/internal/workitem` and `flai/internal/itemedit` cover: an input edit stamps the inputs and leaves the value's stamp; a value edit stamps the block and leaves the inputs'; clearing all three inputs keeps the value and its stamp; the one-stamp form still reads.
- The ADR, `work-hierarchy.md`, and `front-matter-fields.txt` say the same.

## Notes
