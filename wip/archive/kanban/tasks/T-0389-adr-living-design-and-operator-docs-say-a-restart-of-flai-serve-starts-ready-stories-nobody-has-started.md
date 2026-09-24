---
id: T-0389
type: task
nature: remediation
title: ADR, living design, and operator docs say a restart of flai serve starts ready stories nobody has started
status: done
parent: S-0112
owner: alex
created: 2026-09-24T06:19:57Z
updated: 2026-09-24T06:23:25Z
transitions:
  - to: ready
    at: 2026-09-24T06:20:01Z
    by: agent-S-0112
  - to: in-progress
    at: 2026-09-24T06:22:39Z
    by: agent-S-0112
  - to: done
    at: 2026-09-24T06:23:25Z
    by: agent-S-0112
stream: S-0112
tags: []
touches: [design, docs]
---
# T-0389 ADR, living design, and operator docs say a restart of flai serve starts ready stories nobody has started

## Work

- Record with `flai adr new` a decision that supersedes ADR-0038 in part: a story is started once each time it enters ready, whether before or after flai serve started; a restart no longer retires ready stories.
- Update `design/system/flai-cli.md` (`flai serve agent` row) and `docs/operators/index.md` ("What does not start one") to match, and say where the skip reasons appear.

## Done when

- No document in `design/system` or `docs/` says a restart of flai serve starts nothing that was already ready; `flai check --strict` is clean.

## Notes
