---
id: T-1616
type: task
nature: feature
title: An ADR records reading a pull request's checks and reviews into kept state, and flai sending a story back to in-progress with a thread on a failing check or a requested change
status: done
parent: S-0394
owner: alex
created: 2026-10-16T04:36:42Z
updated: 2026-10-17T16:14:50Z
transitions:
  - to: ready
    at: 2026-10-17T16:12:02Z
    by: agent-S-0394
  - to: in-progress
    at: 2026-10-17T16:12:02Z
    by: agent-S-0394
  - to: done
    at: 2026-10-17T16:14:50Z
    by: agent-S-0394
stream: S-0394
tags: [cli]
touches: [design/adrs]
usage:
  source: log
  seconds: 168
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 14
      output: 5822
      cache_read: 833559
      cache_write: 25181
      cost: 0.4257
---
# T-1616 An ADR records reading a pull request's checks and reviews into kept state, and flai sending a story back to in-progress with a thread on a failing check or a requested change

## Work

- Write the ADR with `flai adr new`. It refines S-0393's ADR, ADR-0055, and ADR-0108. It records:
  - that flai serve reads the checks and reviews of each story in review, and keeps them under `.flai-cache`
  - that `flai show`, `item_get`, and the dashboard read only the kept state
  - that a failing required check, or a requested change, opens one thread per failure and commit and sends the story back to in-progress, by `flai`, as a new maker of that transition
  - that flai serve then starts the story's agent
  - the alternative, leaving the story in review, and why it lost
- First in its layer: it waits for no task.

## Done when

The ADR is accepted, and `flai check --strict` passes on it.

## Notes
