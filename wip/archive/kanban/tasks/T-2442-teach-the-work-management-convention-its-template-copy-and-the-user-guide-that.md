---
id: T-2442
type: task
nature: improvement
title: Teach the work-management convention, its template copy, and the user guide that a notice waits for a paused agent's answer
status: done
parent: S-0545
owner: alex
created: 2026-11-02T20:06:16Z
updated: 2026-11-06T22:10:08Z
transitions:
  - to: ready
    at: 2026-11-06T22:00:48Z
    by: agent-S-0545
  - to: in-progress
    at: 2026-11-06T22:00:49Z
    by: agent-S-0545
  - to: done
    at: 2026-11-06T22:10:08Z
    by: agent-S-0545
stream: S-0545
tags: [template]
touches: [design/conventions/work-management.md, template/root/design/conventions/work-management.md, docs/users/flai.md, template/CHANGELOG.md]
after: [T-2440]
usage:
  source: log
  seconds: 559
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 24
      output: 5950
      cache_read: 3105139
      cache_write: 46081
      cost: 1.0849
  strategic:
    - kind: orchestrator
      seconds: 3
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 3
          cache_read: 30803
          cache_write: 2234
          cost: 0.0078
---
# T-2442 Teach the work-management convention, its template copy, and the user guide that a notice waits for a paused agent's answer

## Work

Bring the rules agents read in line with T-2440's ADR. In `design/conventions/work-management.md` and its baseline copy `template/root/design/conventions/work-management.md`, change the bullets on messages between stories and on a task-done notice: while your story waits on your own question to the operator, a task-done notice does not start you again; the answer does, and your first `inbox` lists the notices, which you answer then, settled when they do not break your work. Keep both copies' baseline text the same. Record the convention change under `## Unreleased` in `template/CHANGELOG.md`. In `docs/users/flai.md` § Messages between stories, say what starts a paused agent again. Read the instructions in `flai/internal/mcpserver/server.go` that say a new message starts an agent again, and correct them in this task if they no longer hold.

Waits for T-2440, whose decision it states. It changes no path T-2441 or T-2443 changes, so it runs beside them.

## Done when

- Both copies of `work-management.md` state the rule with the same baseline text and link the ADR.
- `template/CHANGELOG.md` and `docs/users/flai.md` state the change.
- The markdown lint and `flai check --strict` pass.

## Notes
