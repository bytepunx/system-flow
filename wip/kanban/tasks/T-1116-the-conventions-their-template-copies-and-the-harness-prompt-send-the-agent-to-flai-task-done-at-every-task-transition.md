---
id: T-1116
type: task
nature: improvement
title: The conventions, their template copies, and the harness prompt send the agent to flai task done at every task transition
status: done
parent: S-0269
owner: alex
created: 2026-10-06T22:53:50Z
updated: 2026-10-07T03:02:56Z
transitions:
  - to: ready
    at: 2026-10-07T02:58:35Z
    by: agent-S-0269
  - to: in-progress
    at: 2026-10-07T02:58:35Z
    by: agent-S-0269
  - to: done
    at: 2026-10-07T03:02:56Z
    by: agent-S-0269
stream: S-0269
tags: [conventions, template]
touches: [design/conventions/work-management.md, design/conventions/git.md, template/root/design/conventions/work-management.md, template/root/design/conventions/git.md, template/CHANGELOG.md, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/system/workflow.md]
after: [T-1102, T-1107]
usage:
  source: log
  seconds: 261
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 65
      output: 508
      cache_read: 2230863
      cache_write: 95210
      cost: 1.165
---
# T-1116 The conventions, their template copies, and the harness prompt send the agent to flai task done at every task transition

## Work

Send the agent to the one call.

- In `design/conventions/git.md`, the per-task cycle ADR-0069 set (commit, `flai stream sync`, resolve, test, commit any fix) becomes:
  - `flai task done T-nnnn -m "<message>"`, or the MCP tool `task_done`;
  - on a refused sync, resolve what it lists, then call it again;
  - run the task's tests, and close any fix with a commit.
- In `design/conventions/work-management.md`, make the same change wherever it lists the steps of a task transition. That covers "Finish each task with the cycle in `git.md`", narrative logging at a transition, and calling `inbox` at every task transition, all of which the call does now. Do not touch the baseline above each file's project marker beyond what the cycle's wording needs.
- Copy both files' changes into `template/root/design/conventions/`, and add an entry to `template/CHANGELOG.md`.
- In `flai/internal/harness/harness.go`, the story agent's prompt (`Prompt`, near "When a task is done, commit its changes") says: close each task with `flai task done <task> -m` (or the MCP tool `task_done`), resolve what a refused sync lists, then call it again. `harness_test.go` checks the new wording, and that the prompt no longer lists the steps one by one.

This task waits for T-1102 and T-1107, so that the names it writes exist. It is layer 4.

## Done when

- [ ] `git.md`, `work-management.md`, and their template copies send the agent to `flai task done` or `task_done` at every task transition.
- [ ] `template/CHANGELOG.md` records the change.
- [ ] The harness prompt names the command, `harness_test.go` covers it, and `scripts/flai-test.sh` and the markdown lint pass.

## Notes

Drafted by the planner. S-0271 and S-0274 edit the same convention files and the same prompt; whichever story is accepted second syncs over the first.
