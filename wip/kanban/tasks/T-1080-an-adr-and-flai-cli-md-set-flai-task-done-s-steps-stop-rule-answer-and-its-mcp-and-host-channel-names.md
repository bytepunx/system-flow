---
id: T-1080
type: task
nature: improvement
title: An ADR and flai-cli.md set flai task done's steps, stop rule, answer, and its MCP and host-channel names
status: done
parent: S-0269
owner: alex
created: 2026-10-06T22:52:39Z
updated: 2026-10-07T02:00:18Z
transitions:
  - to: ready
    at: 2026-10-07T01:43:26Z
    by: agent-S-0269
  - to: in-progress
    at: 2026-10-07T01:43:26Z
    by: agent-S-0269
  - to: done
    at: 2026-10-07T02:00:18Z
    by: agent-S-0269
stream: S-0269
tags: [cli, design]
touches: [design/adrs, design/system/flai-cli.md]
usage:
  source: log
  seconds: 1012
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 47
      output: 20233
      cache_read: 2723955
      cache_write: 79672
      cost: 1.414
---
# T-1080 An ADR and flai-cli.md set flai task done's steps, stop rule, answer, and its MCP and host-channel names

## Work

Write down the contract the code tasks build to, before any code.

- Add an ADR, refining ADR-0069. A task is closed with one call, `flai task done T-nnnn -m "<message>" [--log "<entry>"] [--json]`, run from the story's worktree. It does these steps in order:
  1. commit (`git add -A`, `git commit`) in the worktree;
  2. `flai stream sync` of the story;
  3. move the task to done, under the move's rules;
  4. append the narrative log entry, which is the message unless `--log` gives another;
  5. widen the task's touches, and the story's, with the paths the task's commit changed that they do not cover;
  6. run `flai check`, scoped to the story;
  7. answer the agent's inbox.
- It stops at the first step that fails, with that step's findings, and leaves the steps after it undone. It still refuses while a rebase is unfinished. Nothing to commit is not a failure.
- The answer has the same fields on the CLI, over MCP, and on the host channel. Over MCP the call is `task_done`; on the host channel it is `task.done`, which runs `flai task done --json`. The answer gives:
  - the commit (its hash, or none);
  - the sync's result (synced, or the conflicts and how to continue);
  - the move, and any epic or story it moved with it;
  - the touches widened;
  - the check's findings;
  - the inbox;
  - the step it stopped at, if any.
- In `design/system/flai-cli.md`, describe the command in the commands table. Add `task.done` to the `flai hostapi` method list.
- Say what the step after the sync leaves to the agent. Running the task's tests stays the agent's job (S-0273 adds `flai test`), and so does ticking criteria (`flai criteria tick`).

This task waits for nothing. It is in layer 1, with the two extractions.

## Done when

- [ ] A proposed ADR in `design/adrs` refines ADR-0069 with the steps, their order, the stop rule, and the answer's fields.
- [ ] `design/system/flai-cli.md` describes `flai task done`, the MCP tool `task_done`, and the host method `task.done`.
- [ ] `flai check --strict` and the markdown lint pass on both files.

## Notes

Drafted by the planner. `design/adrs` is a folder touch because the ADR's number is taken when it is written.
