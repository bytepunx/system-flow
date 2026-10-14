---
id: T-1584
type: task
nature: improvement
title: ADR records that a rebase merges design/adrs/README.md's index by rows, and the design and user docs say so
status: done
parent: S-0306
owner: alex
created: 2026-10-14T23:04:04Z
updated: 2026-10-14T23:07:38Z
transitions:
  - to: ready
    at: 2026-10-14T23:04:26Z
    by: agent-S-0306
  - to: in-progress
    at: 2026-10-14T23:04:27Z
    by: agent-S-0306
  - to: done
    at: 2026-10-14T23:07:38Z
    by: agent-S-0306
stream: S-0306
tags: []
touches: [design/adrs/README.md, design/system/flai-cli.md, docs/users/flai.md]
usage:
  source: log
  seconds: 178
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 39
      output: 227
      cache_read: 1389629
      cache_write: 61117
      cost: 0.6695
---
# T-1584 ADR records that a rebase merges design/adrs/README.md's index by rows, and the design and user docs say so

## Work

Record the decision with `flai adr new` (no `--commit`), refining ADR-0098 and ADR-0126: when a rebase that `flai stream sync`, `flai task done`, or `flai accept` runs stops on `design/adrs/README.md`, flai merges the index row by row from the three versions and continues, and the trial merge leaves the file out of a pair's conflicts. Alternatives: regenerating the index whole from the ADR files (I-0099's recommendation; loses the hand-written status notes such as "superseded in part by 0016" and the text above the table), a git merge driver (as ADR-0098 rejects it). Then say the same in `design/system/flai-cli.md` (the `flai stream sync` row, beside ADR-0098 and ADR-0126) and `docs/users/flai.md` (a bullet beside **Issue files** in the sync section), linking the ADR.

Waits for nothing: it records the design this story settles before the code, and shares no path with the code tasks.

## Done when

- The new ADR file and its index row exist, the two documents describe the behaviour and link it, and `flai test` on the changed markdown passes.

## Notes
