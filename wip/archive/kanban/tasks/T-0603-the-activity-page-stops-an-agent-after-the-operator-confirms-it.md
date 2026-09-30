---
id: T-0603
type: task
nature: feature
title: The activity page stops an agent after the operator confirms it
status: done
parent: S-0170
owner: alex
created: 2026-09-30T00:48:33Z
updated: 2026-09-30T00:59:20Z
transitions:
  - to: ready
    at: 2026-09-30T00:48:49Z
    by: agent-S-0170
  - to: in-progress
    at: 2026-09-30T00:55:03Z
    by: agent-S-0170
  - to: done
    at: 2026-09-30T00:59:20Z
    by: agent-S-0170
stream: S-0170
tags: []
touches: [flaiover/src]
usage:
  source: log
  seconds: 257
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 38
      output: 13822
      cache_read: 3405135
      cache_write: 41607
      cost: 1.2905
---

# T-0603 The activity page stops an agent after the operator confirms it

## Work

- `POST /api/items/[id]/agent` accepts `stop`, asking flai's `agent.stop`.
- The activity page shows Stop beside each agent that runs or waits for an answer, when the dashboard may write and the agent action is on.
- Stop opens a confirmation that says what stopping does (the process and what it started are ended, uncommitted work stays in the worktree, the story stays where it is and gets no new agent until Retry or a move back to ready) and acts only on confirm.
- A stopped agent reads as stopped by the operator.
- Component and route tests.

## Done when

- `npm run check`, `npm run lint`, and `npm test` in `flaiover/` pass, with tests for the button, the confirmation, and the route.

## Notes
