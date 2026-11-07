---
id: T-2512
type: task
nature: remediation
title: flai check reports a schedule entry under each-ready-story on a role whose item does not list story
status: done
parent: S-0571
owner: alex
created: 2026-11-07T03:46:41Z
updated: 2026-11-07T03:55:24Z
transitions:
  - to: ready
    at: 2026-11-07T03:50:52Z
    by: agent-S-0571
  - to: in-progress
    at: 2026-11-07T03:50:54Z
    by: agent-S-0571
  - to: done
    at: 2026-11-07T03:55:24Z
    by: agent-S-0571
stream: S-0571
tags: [flai, check, roles]
touches: [flai/internal/roles/roles.go, flai/internal/roles/roles_test.go]
usage:
  source: log
  seconds: 270
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 11
      output: 2286
      cache_read: 1315547
      cache_write: 24452
      cost: 0.4983
---
# T-2512 flai check reports a schedule entry under each-ready-story on a role whose item does not list story

## Work

`each-ready-story` starts a role on ready stories only, so a role whose `item` lists epics or tasks but not `story` would still start nothing on that entry. In `flai/internal/roles/roles.go`, where the trigger check reports a schedule entry on a role that takes an item without `when: each-ready-story`, also report `roles.file` for an entry with `when: each-ready-story` on a role whose `item` is not empty and does not list `story`. Name the role's items and say the condition starts the role on stories alone (ADR-0250 §5). Leave the built-in roles' launcher checks as they are. Add a case to `roles_test.go` beside the S-0554 cases for the existing message.

This task waits for none. It touches no path of the scheduler task's, so the two run together.

## Done when

- `flai test flai/internal/roles` passes.
- A role file with `item: [epic]` and `{schedule: roles.<name>.schedule, when: each-ready-story}` gets one `roles.file` finding from `flai check`. One with `item: [story]` gets none.

## Notes
