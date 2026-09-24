---
id: T-0426
type: task
nature: feature
title: Runbooks for installing, updating, and deleting flai and flaiover
status: done
parent: S-0028
owner: alex
created: 2026-09-24T09:28:02Z
updated: 2026-09-24T09:38:16Z
transitions:
  - to: ready
    at: 2026-09-24T09:28:10Z
    by: agent-S-0028
  - to: in-progress
    at: 2026-09-24T09:35:27Z
    by: agent-S-0028
  - to: done
    at: 2026-09-24T09:38:16Z
    by: agent-S-0028
stream: S-0028
tags: []
touches: [docs/operators]
---
# T-0426 Runbooks for installing, updating, and deleting flai and flaiover

## Work
- `docs/operators/runbooks/README.md` indexes the runbooks.
- `install.md`, `update.md`, and `delete.md`, each with a section for flai and one for flaiover: prerequisites, the commands in order, how to check it worked, and what to do when a step fails. Each step is taken from the code and the existing guide, not invented.

## Done when
- The three pages exist, every command in them exists with the flags they give it, and every path they name is where flai puts it.

## Notes
