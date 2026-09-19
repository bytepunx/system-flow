---
id: T-0224
type: task
nature: feature
title: Operators documentation with the safer choice first, the dashboard and CLI design, and the security notes
status: done
parent: S-0062
owner: alex
created: 2026-09-19T08:21:01Z
updated: 2026-09-19T08:33:45Z
transitions:
  - to: ready
    at: 2026-09-19T08:29:07Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T08:29:08Z
    by: system-flow
  - to: done
    at: 2026-09-19T08:33:45Z
    by: system-flow
stream: S-0062
tags: []
---

# T-0224 Operators documentation with the safer choice first, the dashboard and CLI design, and the security notes

## Work
`docs/operators/index.md`: the opt-in, with a key dedicated to this repository first (how to make one, where to add it as a deploy key with write access), then the operator's own key with what it exposes in plain words, the repository rules to set, what the dashboard token now means, and how to revoke. `flaiover-dashboard.md`, `flai-cli.md`, `docs/users/flai.md`, and every sentence that says the container holds no credential or that the dashboard cannot push are brought into line. Tick the criteria for what was verified.

## Done when
- The documents agree with the behaviour and with ADR-0026
- The criteria are ticked and `flai check --strict` is clean

## Notes
