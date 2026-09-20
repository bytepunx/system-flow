---
id: T-0295
type: task
nature: feature
title: flai edit changes an item's title, nature, tags, touches, parent, and body in one checked, committed step
status: done
parent: S-0085
owner: alex
created: 2026-09-20T15:21:41Z
updated: 2026-09-20T15:30:05Z
transitions:
  - to: ready
    at: 2026-09-20T15:21:43Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T15:21:43Z
    by: system-flow
  - to: done
    at: 2026-09-20T15:30:05Z
    by: system-flow
stream: S-0085
tags: []
---
# T-0295 flai edit changes an item's title, nature, tags, touches, parent, and body in one checked, committed step

## Work
internal/itemedit and the flai edit command. Any of the fields and the body in one change. A retitle keeps the front matter, the heading, the file name, the parent's list, the narrative, and links to the old file name in step. A new parent must exist, be of the right type, and be open; the item leaves the old parent's list and joins the new one. An archived item is refused. With a hash, a change made meanwhile is a conflict. flai check runs with the change in place; what it introduces refuses the edit and everything is put back. One commit for every file touched, with trailers, unless the project turned autocommit off. flai edit --show prints the fields, the body, the hash, and whether it may be edited.

## Done when
- Tests with a real repository: each field, a retitle with links and a narrative, a move between parents, a refusal that restores every file, a conflict, an archived item, arguments that are not data
- The Go tests and lint pass

## Notes
