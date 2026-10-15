---
id: T-1597
type: task
nature: improvement
title: Document how task_done closes a task with the installed flai, for users, operators, and agents
status: done
parent: S-0390
owner: alex
created: 2026-10-15T02:59:48Z
updated: 2026-10-15T03:15:11Z
transitions:
  - to: ready
    at: 2026-10-15T03:03:08Z
    by: agent-S-0390
  - to: in-progress
    at: 2026-10-15T03:03:09Z
    by: agent-S-0390
  - to: done
    at: 2026-10-15T03:15:11Z
    by: agent-S-0390
stream: S-0390
tags: [docs]
touches: [docs/users/flai.md, docs/operators/index.md, design/system/agent-narrative.md]
after: [T-1595]
usage:
  source: log
  seconds: 722
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 19
      output: 5804
      cache_read: 1182197
      cache_write: 35717
      cost: 0.5878
---
# T-1597 Document how task_done closes a task with the installed flai, for users, operators, and agents

## Work

Say what T-1596 built where readers look for it. It waits for T-1596 so the docs describe what was built, including the refusal's wording.

- `docs/users/flai.md` § Closing a task: the bullet on ADR-0142 gains what `task_done` does when the session's flai is behind and the installed flai has caught up.
- `docs/operators/index.md` § Keeping the host's flai current: the bullet "What it stops" says that a session started before the upgrade no longer stops at its next task close once the host is upgraded.
- `design/system/agent-narrative.md` § What an agent is told through MCP: the `task_done` line says the same, if it describes the refusal.

## Done when

- The three sections describe the behaviour T-1596 built, and link the new ADR.
- The markdown lint passes on the three files.

## Notes

`docs/users/flai-reference.md` is left out because no flag or help text changes.
