---
id: T-2518
type: task
nature: improvement
title: Record in an ADR that a close-out records no narrative.missing and a check scoped to a story leaves it out
status: done
parent: S-0555
owner: alex
created: 2026-11-07T03:47:09Z
updated: 2026-11-07T03:59:20Z
transitions:
  - to: ready
    at: 2026-11-07T03:57:45Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T03:57:47Z
    by: system-flow
  - to: done
    at: 2026-11-07T03:59:20Z
    by: agent-S-0555
stream: S-0555
tags: [flai, check, adr]
touches: [design/adrs/README.md, design/adrs/0256-a-close-out-records-no-narrative-missing-and-a-check-scoped-to-a-story-leaves.md]
usage:
  source: log
  seconds: 93
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 16
      output: 2681
      cache_read: 2049976
      cache_write: 42224
      cost: 0.736
---
# T-2518 Record in an ADR that a close-out records no narrative.missing and a check scoped to a story leaves it out

## Work

Write and accept an ADR with `flai adr new --commit`. It refines ADR-0085 and follows ADR-0125, ADR-0137, and ADR-0255. Its decision is that a `narrative.missing` finding outside the story is left out of a check scoped to a story (`flai check --story`). A close-out then neither notes it nor records it in an issue. A `narrative.missing` on the story itself stays inside, as `ScopeToStories` already counts the story's own item file.

Its context is I-0235. Five close-outs, those of S-0548, S-0547, S-0550, S-0546, and S-0551 between 2026-10-31 and 2026-11-01, found `narrative.missing` on S-0549, a story none of them touched. The operator moved S-0549 to in progress at 2026-10-31T09:47:53Z. That move goes through `flai move` (`flai/cmd/move.go`), which opens no narrative. Only `flai story start` and `flai stream open` call `OpenStream` (`flai/internal/workitem/narrative.go`). Narratives live in `wip/` on the main branch, so no other story's branch can open one. Only that story's agent or the operator can. The ADR's consequences name two places that still report it:

- `flai check` in the main checkout still reports `narrative.missing`
- the story's own close-out stops on its own missing narrative

The alternative is to make a move to in progress, from `flai move`, the dashboard, or `item_move`, open the narrative, as `flai story start` does. That removes the cause. But it changes what a move does for every caller. It also opens a worktree and a branch for a story nobody has started. Decide it instead if it fits better on reading the instances, and say why in the ADR. If it is decided, say so on the story's thread before T-2521 is worked, since T-2521's work then changes.

It waits for nothing in this story.

## Done when

- The ADR is accepted under `design/adrs/`, indexed in `design/adrs/README.md`, with its topics set.
- It says whether a `narrative.missing` outside the story is left out, why, and what still reports it.

## Notes
