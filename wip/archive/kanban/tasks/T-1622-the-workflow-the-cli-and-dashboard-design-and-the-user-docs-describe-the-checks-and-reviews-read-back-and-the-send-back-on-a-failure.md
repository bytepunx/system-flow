---
id: T-1622
type: task
nature: feature
title: The workflow, the CLI and dashboard design, and the user docs describe the checks and reviews read back and the send back on a failure
status: done
parent: S-0394
owner: alex
created: 2026-10-16T04:37:14Z
updated: 2026-10-17T17:05:03Z
transitions:
  - to: ready
    at: 2026-10-17T16:49:38Z
    by: agent-S-0394
  - to: in-progress
    at: 2026-10-17T16:49:39Z
    by: agent-S-0394
  - to: done
    at: 2026-10-17T17:05:03Z
    by: agent-S-0394
stream: S-0394
tags: [cli, dashboard]
touches: [design/system/workflow.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flaiover.md, docs/users/flai-reference.md, design/system/agent-adapters.md, docs/operators/index.md]
after: [T-1616, T-1618, T-1620, T-1621]
usage:
  source: log
  seconds: 924
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 161
      output: 67983
      cache_read: 9733374
      cache_write: 294032
      cost: 4.971
---
# T-1622 The workflow, the CLI and dashboard design, and the user docs describe the checks and reviews read back and the send back on a failure

## Work

- `design/system/workflow.md`, § Transitions and who makes them: the review to in-progress row gains flai, on a failing required check or a requested change under pull-request acceptance, linking T-1616's ADR.
- `design/system/flai-cli.md` and `docs/users/flai.md`: the pull request section of `flai show`, `flai forge status`, and the kept state.
- `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md`: the pull request on the story and review pages, and `forge.status`.
- Regenerate `docs/users/flai-reference.md`.
- Waits for T-1616 (the ADR), T-1618 (the CLI's output), T-1620 (the agent's start), and T-1621 (the pages): the docs describe what they built.

## Done when

Every document above describes the read back and the send back, the reference test passes, and the markdown lint is clean.

## Notes
