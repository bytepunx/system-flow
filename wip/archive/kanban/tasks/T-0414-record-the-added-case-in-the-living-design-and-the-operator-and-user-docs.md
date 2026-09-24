---
id: T-0414
type: task
nature: remediation
title: Record the added case in the living design and the operator and user docs
status: done
parent: S-0116
owner: alex
created: 2026-09-24T08:46:28Z
updated: 2026-09-24T08:52:13Z
transitions:
  - to: ready
    at: 2026-09-24T08:51:31Z
    by: agent-S-0116
  - to: in-progress
    at: 2026-09-24T08:51:32Z
    by: agent-S-0116
  - to: done
    at: 2026-09-24T08:52:13Z
    by: agent-S-0116
stream: S-0116
tags: []
touches: [design/system, docs/operators, docs/users]
---
# T-0414 Record the added case in the living design and the operator and user docs

## Work

A changed agent starting another run is a case added to ADR-0041, not a reversal: a living-design edit. Update `design/system/flai-cli.md` where it says a story is started once each time it enters ready, `docs/operators/index.md` ("When one starts"), and `docs/users` where they say the same.

## Done when

- Every place that says a failed story waits until it is moved to ready again also names changing its agent.
- `flai check --strict` and the markdown lint are clean.

## Notes
