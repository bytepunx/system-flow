---
id: T-0157
type: task
nature: remediation
title: User, operator, and tech documentation for the mount and the setting
status: done
parent: S-0050
owner: alex
created: 2026-09-18T19:47:24Z
updated: 2026-09-18T19:57:53Z
transitions:
  - to: ready
    at: 2026-09-18T19:56:47Z
    by: alex
  - to: in-progress
    at: 2026-09-18T19:56:48Z
    by: alex
  - to: done
    at: 2026-09-18T19:57:53Z
    by: alex
stream: S-0050
tags: []
touches: [docs, design/tech]
---

# T-0157 User, operator, and tech documentation for the mount and the setting

## Work
Update `docs/users/flai.md` (`flai dashboard`, `flai stream open`, `flai config`), `docs/users/flaiover.md` (acceptance from the board, what to do when it cannot run), the operator documentation's settings index with the new key, what turning it on does to the clone, and how to turn it back off, and `design/tech` for the git versions involved (2.48 for the opt-in).

## Done when
- Each document states the behaviour as built
- Markdown lint passes

## Notes
