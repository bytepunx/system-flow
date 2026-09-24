---
id: T-0407
type: task
nature: feature
title: "Operators guide: warn against exposing the dashboard"
status: done
parent: S-0019
owner: alex
created: 2026-09-24T08:28:56Z
updated: 2026-09-24T08:31:24Z
transitions:
  - to: ready
    at: 2026-09-24T08:29:00Z
    by: agent-S-0019
  - to: in-progress
    at: 2026-09-24T08:30:40Z
    by: agent-S-0019
  - to: done
    at: 2026-09-24T08:31:24Z
    by: agent-S-0019
stream: S-0019
tags: []
touches: [docs/operators, docs/README.md]
---
# T-0407 Operators guide: warn against exposing the dashboard

## Work
Put a plain warning at the head of "Running the dashboard" in `docs/operators/index.md`: the dashboard is published on every interface by default, speaks plain HTTP, and its token can change the project (and, with host actions on, push or run commands as the operator); never publish it to the internet, bind it to loopback or put a TLS tunnel in front beyond a trusted network. Correct "Security posture" (it still says "the project token"; the token is per user since ADR-0033) and link the users guide. Update `docs/README.md`'s row if it no longer describes the folder.

## Done when
- The operators guide warns against exposure before the first command, and "Security posture" matches the current design.
- `scripts/lint-md.sh` and `flai check --strict` pass.

## Notes
