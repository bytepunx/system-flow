---
id: T-1422
type: task
nature: improvement
title: Record in an ADR that a close-out records no board.wip-limit and a check scoped to a story leaves it out
status: done
parent: S-0348
owner: alex
created: 2026-10-08T08:58:51Z
updated: 2026-10-08T09:19:35Z
transitions:
  - to: ready
    at: 2026-10-08T09:19:02Z
    by: agent-S-0348
  - to: in-progress
    at: 2026-10-08T09:19:05Z
    by: agent-S-0348
  - to: done
    at: 2026-10-08T09:19:35Z
    by: agent-S-0348
stream: S-0348
tags: [flai, check, adr]
touches: [design/adrs]
usage:
  source: log
  seconds: 30
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 7
      output: 1723
      cache_read: 377003
      cache_write: 15662
      cost: 0.2188
---
# T-1422 Record in an ADR that a close-out records no board.wip-limit and a check scoped to a story leaves it out

## Work

Write and accept an ADR with `flai adr new`, refining ADR-0085 and following ADR-0122, ADR-0123, and ADR-0125: a `board.wip-limit` finding is left out of a check scoped to a story (`flai check --story`), so a close-out neither notes it nor records it in an issue. The whole-repository `flai check --strict` in the main checkout still reports it, and still fails on ready or in-progress over its limit (ADR-0073 is unchanged).

Its context is I-0123: S-0339's close-out found `wip/kanban/board.md: 6 stories in in-progress, limit 5`. The board's counts are the main checkout's state, which no story branch changes or clears, and the breach names no story. Its consequences name what reports a breach instead: `flai check` in the main checkout, the pull hold, and `flai serve`, which starts no agent over the limit.

If, on reading the instances, a note without an issue (as `wip.overlap` gets under ADR-0115) fits better than leaving the finding out, decide that instead and say why in the ADR.

It waits for nothing. The code and documentation tasks wait for it so they can cite its number.

## Done when

- The ADR is accepted under `design/adrs/`, indexed in `design/adrs/README.md`, with topics set.
- It states which findings the scoped check leaves out (`board.wip-limit` for every column, or only ready and in-progress), and why.

## Notes
