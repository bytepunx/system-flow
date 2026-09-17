---
id: T-0019
type: task
nature: feature
title: epic new, story new, task new
status: done
parent: S-0007
owner: agent
created: 2026-09-15T17:42:41Z
updated: 2026-09-15T17:52:35Z
transitions:
  - to: ready
    at: 2026-09-15T17:42:41Z
    by: agent
  - to: in-progress
    at: 2026-09-15T17:52:35Z
    by: agent
  - to: done
    at: 2026-09-15T17:52:35Z
    by: agent
stream: S-0007
tags: [cli, workitems]
---

# T-0019 epic new, story new, task new

## Work
Render from the template's items/ (when cached) or embedded defaults; allocate ID and slug; write to the right folder; append the child to the parent's Stories or Tasks section; --nature, --owner, --tags.

## Done when
Creating an epic, story, and task in a temp project yields valid files and linked parents.

## Notes
