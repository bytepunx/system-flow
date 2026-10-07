---
id: T-1194
type: task
nature: improvement
title: An ADR refining ADR-0046 sends conflicts and grown overlaps to the two stories as a message, with escalation to the operator
status: done
parent: S-0332
owner: alex
created: 2026-10-07T20:14:56Z
updated: 2026-10-07T22:50:05Z
transitions:
  - to: ready
    at: 2026-10-07T22:49:40Z
    by: system-flow
  - to: in-progress
    at: 2026-10-07T22:49:41Z
    by: system-flow
  - to: done
    at: 2026-10-07T22:50:05Z
    by: agent-S-0332
stream: S-0332
tags: [flai]
touches: [design/adrs]
usage:
  source: log
  seconds: 24
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 7
      output: 2772
      cache_read: 424179
      cache_write: 13448
      cost: 0.2209
---
# T-1194 An ADR refining ADR-0046 sends conflicts and grown overlaps to the two stories as a message, with escalation to the operator

## Work

Record the change before building it. It waits for nothing.

- A trial-merge conflict at `flai stream sync` and a claim grown into another story's (S-0244) open, or add to, one conversation per pair of stories, `about` the paths, instead of a thread on one story.
- The conversation closes when a later sync merges the two cleanly or the other story is no longer open, as the conflict thread does today.
- Escalation: either agent opens a thread on the operator naming both stories and linking the conversation, with what they could not agree.
- Conflict threads already open stay as they are, and close as they do today.
- The `overlapped` change both stories get stays, so an older flai and the dashboard still see it.

## Done when

- The ADR is accepted in `design/adrs` and names ADR-0046 and S-0330's ADR as what it refines.
- `flai check --strict` passes.

## Notes
