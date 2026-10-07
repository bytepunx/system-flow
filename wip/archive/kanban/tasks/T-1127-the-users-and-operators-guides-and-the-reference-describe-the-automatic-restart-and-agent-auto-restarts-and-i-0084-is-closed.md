---
id: T-1127
type: task
nature: remediation
title: The users' and operators' guides and the reference describe the automatic restart and agent.auto_restarts, and I-0084 is closed
status: done
parent: S-0294
owner: alex
created: 2026-10-06T23:08:07Z
updated: 2026-10-07T02:01:21Z
transitions:
  - to: ready
    at: 2026-10-07T01:59:54Z
    by: agent-S-0294
  - to: in-progress
    at: 2026-10-07T01:59:54Z
    by: agent-S-0294
  - to: done
    at: 2026-10-07T02:01:21Z
    by: agent-S-0294
stream: S-0294
tags: [docs, serve]
touches: [docs/users/flai.md, docs/users/flai-reference.md, docs/operators/index.md, docs/operators/settings.md, design/issues/I-0084-claude-code-ends-a-headless-agent-ten-minutes-after-its-turn-ends-even-while-its-background-sub-agent-is-still-working-and-flai-serve-leaves-the-story-in-progress-with-no-agent.md, design/issues/summary.md]
after: [T-1126]
usage:
  source: log
  seconds: 87
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 29
      output: 10434
      cache_read: 1917836
      cache_write: 57601
      cost: 0.9584
---
# T-1127 The users' and operators' guides and the reference describe the automatic restart and agent.auto_restarts, and I-0084 is closed

## Work

This is criterion 2, and the user-facing half of criterion 1.

- **`docs/operators/index.md`.** The **Restart** paragraph (ADR-0043) says that flai serve restarts a story's agent on its own, under T-1124's conditions, up to `agent.auto_restarts` times. It then opens a thread, and the operator restarts as before. It links the new ADR.
- **`docs/operators/settings.md`.** The host config table gains a row for `agent.auto_restarts`: its default, `flai serve agent set --auto-restarts`, and what 0 does.
- **`docs/users/flai.md`.** Where `flai serve agent` is described, add `--auto-restarts` and the thread at the limit.
- **`docs/users/flai-reference.md`.** Regenerate it with `make flai-reference`, so that it carries the new flag.
- **The issue.** From the story's worktree, run `flai issue close I-0084 --reason` saying what fixed it:
  - S-0285 (ADR-0092) gave the agent a safe way to wait;
  - this story has flai serve restart an agent that ended with its story in progress, up to the limit, and then tell the operator.

  The command updates `design/issues/summary.md`.

## Done when

- The guides and the reference describe the automatic restart and `agent.auto_restarts` as the code does.
- I-0084 is closed with a reason that names S-0285 and this story's ADR.
- `flai check --strict` and the markdown lint pass.

## Notes

Drafted by the planner. The operator confirmed TH-0208's recommendation on 2026-10-06.
