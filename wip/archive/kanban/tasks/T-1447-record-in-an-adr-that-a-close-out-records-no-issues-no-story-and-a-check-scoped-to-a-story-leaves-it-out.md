---
id: T-1447
type: task
nature: improvement
title: Record in an ADR that a close-out records no issues.no-story and a check scoped to a story leaves it out
status: done
parent: S-0365
owner: alex
created: 2026-10-12T14:56:25Z
updated: 2026-10-12T14:59:25Z
transitions:
  - to: ready
    at: 2026-10-12T14:58:36Z
    by: agent-S-0365
  - to: in-progress
    at: 2026-10-12T14:58:39Z
    by: agent-S-0365
  - to: done
    at: 2026-10-12T14:59:25Z
    by: agent-S-0365
stream: S-0365
tags: [flai, check, adr]
touches: [design/adrs, design/adrs/README.md]
usage:
  source: log
  seconds: 46
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 15
      output: 3846
      cache_read: 783770
      cache_write: 34971
      cost: 0.4697
---
# T-1447 Record in an ADR that a close-out records no issues.no-story and a check scoped to a story leaves it out

## Work

Write and accept an ADR with `flai adr new --commit`. It refines ADR-0085 and follows ADR-0122, ADR-0123, ADR-0125, and ADR-0133. Its decision: an `issues.no-story` finding outside the story is left out of a check scoped to a story (`flai check --story`). So a close-out neither notes it nor records it in an issue. One on an issue file the story changes, such as an issue it bumped, stays inside, as `item.archive` and `board.wip-limit` do.

Its context is I-0134. Four close-outs (S-0343, S-0289, S-0342, S-0361) found `issues.no-story` on I-0052, an issue none of them touched. The finding names an issue and no story. It clears only when someone makes a story for that issue, which is the operator's choice. Two places already put that choice to the operator, and the ADR's consequences name them:

- the review page lists every open issue no open story links, unchecked
- `flai check` in the main checkout still warns `issues.no-story`

The rule in `design/conventions/continuous-improvement.md` to make the story with `flai issue story` when `flai check` warns stays true for that check.

The alternative is a note without an issue, as ADR-0115 gives a `wip.overlap` naming the story. It would keep the nudge in front of every story's agent. Decide it instead if, on reading the instances, it fits better, and say why in the ADR.

It waits for nothing. The code and documentation tasks wait for it so they can cite its number.

## Done when

- The ADR is accepted under `design/adrs/`, indexed in `design/adrs/README.md`, with topics set.
- It says whether an `issues.no-story` outside the story is left out or noted, and why, and what still reports it.

## Notes
