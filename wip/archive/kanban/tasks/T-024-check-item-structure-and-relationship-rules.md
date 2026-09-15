---
id: T-024
type: task
nature: feature
title: "Check: item, structure, and relationship rules"
status: done
parent: S-008
owner: alex
created: 2026-09-15T18:00:58Z
updated: 2026-09-15T18:09:33Z
transitions:
  - to: ready
    at: 2026-09-15T18:09:33Z
    by: agent
  - to: in-progress
    at: 2026-09-15T18:09:33Z
    by: agent
  - to: done
    at: 2026-09-15T18:09:33Z
    by: agent
stream: S-008
tags: [cli, check]
---

# T-024 Check: item, structure, and relationship rules

## Work
internal/check: load every item across kanban and archive and report: Validate() failures, filename not matching ID, duplicate IDs, missing or wrong-typed parents, archived items that are not closed, closed items lingering in kanban, transitions out of chronological order or outside created/updated, transition sequences the state machine does not allow, stories at ready or beyond without tasks or acceptance criteria, done items with open children or unchecked criteria, missing body headings, parents whose Stories/Tasks list omits a child, tasks whose stream differs from their parent.

## Done when
Each rule has a name, a level, and a test case in a bad fixture.

## Notes
