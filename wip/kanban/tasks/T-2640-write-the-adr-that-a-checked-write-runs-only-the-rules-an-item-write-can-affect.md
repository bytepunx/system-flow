---
id: T-2640
type: task
nature: improvement
title: Write the ADR that a checked write runs only the rules an item write can affect, and say so in flai-cli.md, strategic-agents.md, and the user guide
status: done
parent: S-0598
owner: alex
created: 2026-11-07T23:48:43Z
updated: 2026-11-08T00:47:11Z
transitions:
  - to: ready
    at: 2026-11-08T00:35:44Z
    by: system-flow
  - to: in-progress
    at: 2026-11-08T00:35:47Z
    by: system-flow
  - to: done
    at: 2026-11-08T00:47:11Z
    by: agent-S-0598
stream: S-0598
tags: [cli, docs]
touches: [design/adrs/0264-a-checked-write-of-a-work-item-runs-only-the-rules-the-write-can-affect-before.md, design/adrs/README.md, design/system/flai-cli.md, design/system/strategic-agents.md, docs/users/flai.md]
after: [T-2622, T-2625]
usage:
  source: log
  seconds: 684
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 53
      output: 14663
      cache_read: 10097352
      cache_write: 183054
      cost: 3.5265
---
# T-2640 Write the ADR that a checked write runs only the rules an item write can affect, and say so in flai-cli.md, strategic-agents.md, and the user guide

## Work

Write an ADR with `flai adr new` that refines ADR-0197 and ADR-0213. It records that a checked creation or edit runs the item-write rule set, before and after, in place of the whole `flai check`. It names each rule kept and each rule skipped, and why a skipped rule cannot bring a finding inside the written item's story. It also says that what a checked write refuses is unchanged. Take the next free number when you write it, and narrow this task's touches to the file `flai adr new` names.

Then change the text that says the whole check runs on a write:

- `design/system/flai-cli.md`, the `flai story new` and `flai edit` rows;
- `design/system/strategic-agents.md`, where a planner's `item_new` is checked;
- `docs/users/flai.md`, where `story new` and `edit` say `flai check` runs with the item in place.

This task waits for T-2622 and T-2625, whose rule set it describes as built. It shares no file with the timing task, so the two run together.

## Done when

- The ADR is accepted in `design/adrs/` and listed in its `README.md`.
- None of the three documents says a checked write runs the whole check, and each links the ADR.
- `flai check --strict` and the markdown lint pass on the files this task changed.

## Notes
