---
id: T-3495
type: task
nature: improvement
title: "ADR: the section merge merges by its words a line both sides change, refining ADR-0287"
status: done
parent: S-0817
owner: alex
created: 2026-11-17T02:31:30Z
updated: 2026-11-17T02:35:47Z
transitions:
  - to: ready
    at: 2026-11-17T02:35:05Z
    by: agent-S-0817
  - to: in-progress
    at: 2026-11-17T02:35:11Z
    by: agent-S-0817
  - to: done
    at: 2026-11-17T02:35:47Z
    by: agent-S-0817
stream: S-0817
tags: [adr, merge]
touches: [design/adrs/README.md]
usage:
  source: log
  seconds: 35
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 6
      output: 37
      cache_read: 691864
      cache_write: 3516
      cost: 0.2645
---
# T-3495 ADR: the section merge merges by its words a line both sides change, refining ADR-0287

## Work

Record the decision with `flai adr new`, refining ADR-0287 (and through it ADR-0175 and ADR-0276). The decision sentence: where the section merge finds both sides changing the same line, and the sentence merge cannot settle it or the line is a table row, it merges that line again by its words, keeping both changes when they touch different words.

Say in the ADR:

- The context, from I-0387: every acceptance stop since ADR-0287 is on `design/system/flai-cli.md` (S-0689, S-0761, S-0797). In S-0797's stop, main and the branch each inserted words at different places in one table row and in one sentence of `docs/users/flai.md`, removing nothing.
- The tokens: words, runs of spaces and tabs, and the line break. They join back into the line. A table row's pipes are words, so a row keeps its cells.
- What still stops: both sides change or remove the same word; both insert at the same point between two words, since no order of the two is right; a heading; a fence line or a line inside a fence.
- That a merged line holds every word of main's line, so ADR-0297 keeps it from reading as undone.
- That the trial merge still counts the file as a conflict, as ADR-0276 leaves it.

The ADR file is named by `flai adr new`; the story's touches predict `design/adrs/0314-…`, and `flai task done` widens this task's touches with the file it writes.

## Done when

- The ADR exists with status accepted, refines ADR-0287, and opens its `## Decision` with the decision sentence.
- `design/adrs/README.md` lists it, and `flai check --strict` is clean.

## Notes
