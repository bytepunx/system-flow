---
id: T-0161
type: task
nature: improvement
title: Docs, all checks, and a look at the real board at several widths
status: done
parent: S-0048
owner: alex
created: 2026-09-18T20:41:14Z
updated: 2026-09-18T20:48:41Z
transitions:
  - to: ready
    at: 2026-09-18T20:44:10Z
    by: alex
  - to: in-progress
    at: 2026-09-18T20:44:10Z
    by: alex
  - to: done
    at: 2026-09-18T20:48:41Z
    by: alex
stream: S-0048
tags: []
touches: [docs/users]
---

# T-0161 Docs, all checks, and a look at the real board at several widths

## Work
Update `docs/users/flaiover.md` (what a card shows) and the `/board` row in `design/system/flaiover-dashboard.md`. Run `make flaiover-test` and `flai check --strict`. Run the dev server against this repository and look at `/board` in a browser at phone, tablet, and desktop widths with "epics and tasks too" on and off: the indicator sits bottom right, long titles wrap above the bottom row, nothing overflows the card. Keep a screenshot in the scratchpad, not in the repository. Tick the story's criteria for what was seen.

## Done when
- flaiover lint, type check, and tests pass; `flai check --strict` is clean
- The board was looked at in a browser at three widths and the result is in the narrative
- Every criterion on S-0048 is checked, or unchecked with the reason in the story notes

## Notes
