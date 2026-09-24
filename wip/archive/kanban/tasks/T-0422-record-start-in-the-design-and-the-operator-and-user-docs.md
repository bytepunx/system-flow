---
id: T-0422
type: task
nature: feature
title: Record start in the design and the operator and user docs
status: done
parent: S-0115
owner: alex
created: 2026-09-24T09:07:10Z
updated: 2026-09-24T09:20:30Z
transitions:
  - to: ready
    at: 2026-09-24T09:19:33Z
    by: agent-S-0115
  - to: in-progress
    at: 2026-09-24T09:19:34Z
    by: agent-S-0115
  - to: done
    at: 2026-09-24T09:20:30Z
    by: agent-S-0115
stream: S-0115
tags: []
touches: [design/system, docs/operators, docs/users]
---
# T-0422 Record start in the design and the operator and user docs

## Work

Add `flai serve agent start` and the Start agent button beside restart in `design/system/flai-cli.md`, `design/system/flaiover-dashboard.md`, `design/system/dashboard-host-channel.md`, `docs/operators/index.md`, `docs/users/flai.md`, and `docs/users/flaiover.md`, and regenerate `docs/users/flai-reference.md`.

## Done when

- Start is documented for operators and users, with its refusals.
- `flai check --strict` and the markdown lint are clean.

## Notes
