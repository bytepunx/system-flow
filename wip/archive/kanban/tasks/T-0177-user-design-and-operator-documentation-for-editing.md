---
id: T-0177
type: task
nature: feature
title: User, design, and operator documentation for editing
status: done
parent: S-0040
owner: alex
created: 2026-09-19T02:11:30Z
updated: 2026-09-19T02:27:41Z
transitions:
  - to: ready
    at: 2026-09-19T02:27:00Z
    by: alex
  - to: in-progress
    at: 2026-09-19T02:27:01Z
    by: alex
  - to: done
    at: 2026-09-19T02:27:41Z
    by: alex
stream: S-0040
tags: []
touches: [docs, design/system]
---

# T-0177 User, design, and operator documentation for editing

## Work
Update `docs/users/flaiover.md` (editing, what cannot be edited and why, conflicts, the commit), `docs/users/flai.md` (`flai doc show` and `save`), `docs/operators/index.md` (`dashboard.autocommit`, and that commits made from the dashboard are not pushed), and check `design/system/flaiover-dashboard.md` and `flai-cli.md` from T-0172 against what was built.

## Done when
- Each document states the behaviour as built
- Markdown lint and `flai check --strict` pass

## Notes
