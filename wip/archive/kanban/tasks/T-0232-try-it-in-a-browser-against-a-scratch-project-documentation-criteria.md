---
id: T-0232
type: task
nature: feature
title: Try it in a browser against a scratch project, documentation, criteria
status: done
parent: S-0059
owner: alex
created: 2026-09-19T09:21:37Z
updated: 2026-09-19T09:34:30Z
transitions:
  - to: ready
    at: 2026-09-19T09:30:18Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T09:30:18Z
    by: system-flow
  - to: done
    at: 2026-09-19T09:34:30Z
    by: system-flow
stream: S-0059
tags: []
---

# T-0232 Try it in a browser against a scratch project, documentation, criteria

## Work
On the branch's dev server against a scratch project: create an epic and a story from the board, see the files flai would have made (ID, front matter, parent link, commit with the designer's identity and the trailer), the board showing the story in backlog without a reload, moving it to ready from its page; a body that the check refuses; a read-only dashboard showing no action. `flaiover-dashboard.md`, `flai-cli.md`, `docs/users`. Tick the criteria for what was seen.

## Done when
- The browser check is recorded in the narrative
- The documents describe creation from the board and the new flags
- The criteria are ticked and `flai check --strict` is clean

## Notes
