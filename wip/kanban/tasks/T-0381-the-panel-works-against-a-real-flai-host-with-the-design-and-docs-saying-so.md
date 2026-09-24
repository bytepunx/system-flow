---
id: T-0381
type: task
nature: improvement
title: The panel works against a real flai host, with the design and docs saying so
status: in-progress
parent: S-0107
owner: alex
created: 2026-09-24T01:35:33Z
updated: 2026-09-24T01:40:31Z
transitions:
  - to: ready
    at: 2026-09-24T01:37:11Z
    by: system-flow
  - to: in-progress
    at: 2026-09-24T01:40:31Z
    by: system-flow
stream: S-0107
tags: []
touches: [flai/internal/hostapi, flai/cmd, docs, design/system]
---
# T-0381 The panel works against a real flai host, with the design and docs saying so

## Work

- Once S-0106 is on main, rebase `story/S-0107` onto it. Add any `host.start` and `host.stop` hostapi method that S-0106 did not add, with its contract test (thread TH-0004).
- Run a scratch host with this tree's flai and a scratch flaiover in a temp config, away from the operator's host and dashboard. Check each control and the upgrade path, or a dry run of it where a real release cannot be installed.
- `design/system/flaiover-dashboard.md` and the flaiover user docs describe the processes area.
- All three test tiers, flaiover checks, and `flai check --strict`.

## Done when

- Every acceptance criterion of S-0107 is verified live. `make test`, `make integration`, and `make smoke` pass, and check is clean.

## Notes
