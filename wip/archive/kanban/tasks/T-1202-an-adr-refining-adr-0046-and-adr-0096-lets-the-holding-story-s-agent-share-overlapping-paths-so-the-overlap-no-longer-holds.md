---
id: T-1202
type: task
nature: feature
title: An ADR refining ADR-0046 and ADR-0096 lets the holding story's agent share overlapping paths, so the overlap no longer holds
status: done
parent: S-0334
owner: alex
created: 2026-10-07T20:15:50Z
updated: 2026-10-08T09:45:01Z
transitions:
  - to: ready
    at: 2026-10-08T09:43:49Z
    by: agent-S-0334
  - to: in-progress
    at: 2026-10-08T09:43:49Z
    by: agent-S-0334
  - to: done
    at: 2026-10-08T09:45:01Z
    by: agent-S-0334
stream: S-0334
tags: [flai]
touches: [design/adrs]
usage:
  source: log
  seconds: 72
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 17
      output: 7433
      cache_read: 1217529
      cache_write: 35978
      cost: 0.6094
---
# T-1202 An ADR refining ADR-0046 and ADR-0096 lets the holding story's agent share overlapping paths, so the overlap no longer holds

## Work

Record the decision before building it: it changes when a story is held. It waits for nothing.

- flai serve messages each story in progress that holds a ready story on overlap alone, once per pair while the hold lasts, with the paths and the held story's goal.
- That conversation has a ready story on one side, which S-0330's ADR refuses for a send between agents: allow it for flai's request about a hold alone, and say who answers for the held story until it starts (flai, or the operator).
- The holding agent answers by narrowing its touches, which clears the hold as today, or by sharing paths with a stated split of who changes what. Only that agent, or the operator, may share.
- A share is recorded in the conversation, ends when either story leaves ready or the open columns, and makes an overlap on its paths not hold, as the shared paths do (ADR-0096); the trial merge at sync and the notice at acceptance still report it.
- Why the holding agent decides rather than the held story's: the held story has no agent until it starts.

## Done when

- The ADR is accepted in `design/adrs` and names ADR-0046, ADR-0096, and the ADRs of S-0330 and S-0332 as what it refines.
- `flai check --strict` passes.

## Notes
