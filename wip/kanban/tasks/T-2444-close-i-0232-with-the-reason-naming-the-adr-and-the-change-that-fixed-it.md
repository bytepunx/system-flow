---
id: T-2444
type: task
nature: remediation
title: Close I-0232 with the reason naming the ADR and the change that fixed it
status: backlog
parent: S-0545
owner: alex
created: 2026-11-02T20:06:49Z
updated: 2026-11-02T20:06:49Z
transitions: []
stream: S-0545
tags: [flai]
touches: [design/issues/I-0232-flai-serve-wakes-a-story-s-agent-that-is-paused-on-a-thread-for-every-task-done-message-from-another-story-each-needing-only-a-settled-reply.md, design/issues/summary.md]
after: [T-2442, T-2443]
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 3
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 0
          output: 2
          cache_read: 30802
          cache_write: 2233
          cost: 0.0078
---
# T-2444 Close I-0232 with the reason naming the ADR and the change that fixed it

## Work

In the story's worktree, run `flai issue close I-0232 --reason "<reason>"`, the reason naming T-2440's ADR and saying that flai serve no longer starts an agent paused on its own question for a task-done notice, which waits for the answer. Tick criterion 2 with `flai criteria tick`, and criterion 1 once T-2443's test is in.

Waits for T-2443, which removes the cause, and T-2442, so the reason can say the convention changed too.

## Done when

- I-0232 is closed with that reason, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` passes.

## Notes
