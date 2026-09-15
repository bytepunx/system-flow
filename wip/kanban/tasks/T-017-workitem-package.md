---
id: T-017
type: task
nature: feature
title: "Work item package: parse, write, list, IDs"
status: done
parent: S-007
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
stream: S-007
tags: [cli, workitems]
---

# T-017 Work item package: parse, write, list, IDs

## Work
internal/workitem: Item with the front matter schema and preserved body; Parse and Write with stable key order; Repo bound to a project root and manifest that lists items across kanban and archive, finds by ID, and allocates the next ID per type.

## Done when
Round-trip tests on every existing item in this repo produce byte-identical front matter semantics and the next IDs are E-005, S-022, T-024.

## Notes
