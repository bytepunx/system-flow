---
id: T-2707
type: task
nature: remediation
title: The design and the user guide say a redirection neither allows nor refuses a strategic agent's command
status: done
parent: S-0613
owner: alex
created: 2026-11-09T08:12:04Z
updated: 2026-11-09T09:07:42Z
transitions:
  - to: ready
    at: 2026-11-09T09:05:58Z
    by: agent-S-0613
  - to: in-progress
    at: 2026-11-09T09:05:59Z
    by: agent-S-0613
  - to: done
    at: 2026-11-09T09:07:42Z
    by: agent-S-0613
stream: S-0613
tags: [flai, guard, docs]
touches: [design/system/strategic-agents.md, docs/users/flai.md]
after: [T-2706]
usage:
  source: log
  seconds: 103
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 31
      output: 4130
      cache_read: 4265252
      cache_write: 114135
      cost: 1.5939
---
# T-2707 The design and the user guide say a redirection neither allows nor refuses a strategic agent's command

## Work

Say what T-2706 built, in the documents that describe the guard's checks of a strategic agent's shell commands:

- `design/system/strategic-agents.md` § The orchestrator › Its permissions and the guard: the guard judges each simple command of a shell call by its arguments, with its redirections and their targets left out (S-0613), so that `2>&1`, `2>/dev/null`, or `>/dev/null` after a command neither allows nor refuses it, and the same holds for every role a role file holds (§ Role files). Bump the front matter's `updated`.
- `docs/users/flai.md` § Role files, where it says `flai guard` reads each role's rules from its file: one sentence saying the same for a reader, with `flai order --by throughput --apply 2>&1` as the example.

It waits for T-2706 because it documents the rule T-2706 settles: which callers judge arguments rather than words.

## Done when

- [ ] Both sections say it, and `flai test design/system/strategic-agents.md docs/users/flai.md` passes the markdown lint

## Notes
