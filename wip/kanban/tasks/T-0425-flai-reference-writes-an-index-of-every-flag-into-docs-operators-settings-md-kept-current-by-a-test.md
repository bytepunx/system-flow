---
id: T-0425
type: task
nature: feature
title: flai reference writes an index of every flag into docs/operators/settings.md, kept current by a test
status: done
parent: S-0028
owner: alex
created: 2026-09-24T09:28:02Z
updated: 2026-09-24T09:35:27Z
transitions:
  - to: ready
    at: 2026-09-24T09:28:09Z
    by: agent-S-0028
  - to: in-progress
    at: 2026-09-24T09:32:11Z
    by: agent-S-0028
  - to: done
    at: 2026-09-24T09:35:27Z
    by: agent-S-0028
stream: S-0028
tags: []
touches: [docs/operators, flai/cmd, scripts]
---
# T-0425 flai reference writes an index of every flag into docs/operators/settings.md, kept current by a test

## Work
- Extend the hidden `flai reference` command to render an index of every flag of every available command (flag, the commands that take it, linked to their section of `docs/users/flai-reference.md`) between generated markers in `docs/operators/settings.md`.
- `scripts/flai-reference.sh` (`make flai-reference`) writes both pages; a test fails when the committed index is stale, as the reference's does.

## Done when
- `docs/operators/settings.md` lists every flag `docs/users/flai-reference.md` lists, and the hand-written part of the page is left as it was.
- The staleness test and `scripts/flai-test.sh` pass.

## Notes
