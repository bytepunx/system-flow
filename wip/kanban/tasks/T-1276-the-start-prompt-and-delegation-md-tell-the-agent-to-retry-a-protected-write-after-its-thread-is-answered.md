---
id: T-1276
type: task
nature: improvement
title: The start prompt and delegation.md tell the agent to retry a protected write after its thread is answered
status: in-progress
parent: S-0309
owner: alex
created: 2026-10-07T23:27:31Z
updated: 2026-10-08T05:53:19Z
transitions:
  - to: ready
    at: 2026-10-08T05:53:19Z
    by: agent-S-0309
  - to: in-progress
    at: 2026-10-08T05:53:19Z
    by: agent-S-0309
stream: S-0309
tags: [harness, conventions, permission-prompt]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/delegation.md]
after: [T-1274]
---
# T-1276 The start prompt and delegation.md tell the agent to retry a protected write after its thread is answered

## Work

The story agent's start prompt in `flai/internal/harness/harness.go` and the project additions of `design/conventions/delegation.md` tell the agent, when the operator may be away, to stage each protected file in the worktree's `.flai-cache/` folder and open a thread with `cp` commands. That workaround came from the thirty-minute wait I-0103 records. It waits for T-1274 because the ADR decides what replaces it.

Rewrite both to say what the ADR decides: `permission_prompt` holds the write a few minutes, then refuses it naming an open thread; the agent goes on with work that does not need the write, and makes the same write again once the thread is answered. When nothing else is left, it writes the narrative's Current state and Next steps naming the thread and ends; flai serve starts it again on the answer. Keep the rule that sub-agents never make these writes.

Update the prompt's test in `flai/internal/harness/harness_test.go`, the one that checks the `.flai-cache/` wording, to check the new wording. The template's `delegation.md` carries none of this text, so it is not changed.

## Done when

- The prompt and `delegation.md` no longer tell the agent to stage files with `cp` commands, and say to retry after the answer.
- `flai test flai/internal/harness design/conventions/delegation.md` passes.

## Notes
