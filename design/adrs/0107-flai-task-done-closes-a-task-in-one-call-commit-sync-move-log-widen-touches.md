---
id: ADR-0107
title: "flai task done closes a task in one call: commit, sync, move, log, widen touches, check, and inbox, stopping at the first step that fails"
status: proposed
date: 2026-10-07
supersedes: []
superseded_by: []
refines: [ADR-0069]
---

# ADR-0107 flai task done closes a task in one call: commit, sync, move, log, widen touches, check, and inbox, stopping at the first step that fails

## Context

[ADR-0069](0069-agents-commit-each-task-on-the-story-branch-after-flai-stream-sync-and-never.md) made each task's close a cycle the agent runs by hand: commit on `story/S-nnnn`, `flai stream sync`, resolve, test, commit any fix. Around it the conventions add `flai move T-nnnn done`, a narrative log entry, `flai touches` for paths the diff added, `flai check`, and `inbox`. E-0017's classification of 108 story runs counted about 2,900 model turns spent on these steps, each re-reading the agent's whole context, and none of them a decision: the repository determines every outcome. S-0269 moves them into flai as one call.

## Decision

A task is closed with one call, `flai task done`, run from the story's worktree, which commits, syncs, moves the task to done, logs, widens touches, checks, and answers the inbox, stopping at the first step that fails.

```bash
flai task done T-nnnn -m "<commit message>" [--log "<entry>"] [--json]
```

Over MCP it is the tool `task_done`; on the host channel it is the write `task.done`, which runs `flai task done --json`. All three answer the same result, built by one package, `flai/internal/taskdone`.

It runs these steps, in this order:

1. **Commit.** `git add -A` and `git commit -m` in the story's worktree. Nothing to commit is not a failure: the step records no commit and the run goes on. While a rebase is unfinished the step refuses, as `flai stream sync` does.
2. **Sync.** The story's `flai stream sync`: rebase onto the main branch, the trial merge, and what lies outside the claim. A sync that refuses, or stops on conflicts, stops the run with the conflicting paths and how to continue.
3. **Move.** The task to `done`, under `flai move`'s rules, with any story or epic that moves with it. A task already done is not moved again, so that the call can be repeated after a stop.
4. **Log.** The narrative log entry: the message's subject line, or `--log` when given.
5. **Touches.** The paths the step 1 commit changed that the task's `touches`, or the story's, do not cover are added to each, as `flai touches` records them. With no commit, nothing is widened.
6. **Check.** `flai check --strict` scoped to the story, as the close-out scopes it ([ADR-0085](0085-a-close-out-s-flai-check-reports-findings-outside-the-story-as-notes-and.md)). Findings in the story stop the run with them; findings outside it are notes.
7. **Inbox.** The agent's inbox, as the MCP tool `inbox` answers it, advancing the same cursor.

It stops at the first step that fails, answers that step's findings, and leaves the steps after it undone. The answer gives:

- `commit`: the commit's hash and subject, or none;
- `sync`: synced, or refused with why, or the conflicting paths and how to continue;
- `move`: the task's new state, and each story or epic moved with it;
- `log`: the entry written;
- `touches`: the paths added to the task and to the story;
- `check`: the findings and notes;
- `inbox`: the inbox;
- `stopped`: the step it stopped at, empty when it ran to the end.

`flai task done` exits 0 when every step ran, 3 when the sync stopped the run, 4 when the check did, and 1 when another step did. The host channel maps 3 to its conflict error and 4 to its refused error, each with the result as data.

The step after the sync that needs judgment stays the agent's: running the task's tests with `flai test`, fixing what they find, and ticking criteria with `flai criteria tick`. A fix is closed by calling `flai task done` again, which commits it, syncs, and checks. On a stopped sync the agent resolves what it lists in the worktree, `git add`s it, runs `git rebase --continue`, and calls it again.

This refines ADR-0069: commit per task on the story branch, sync before testing, and never rebasing or merging by hand all stand.

## Consequences

- The sync moves out of `flai/cmd` into `flai/internal/storygit` and the agent's inbox out of `flai/internal/mcpserver` into `flai/internal/inbox`, so that the CLI, MCP, and the host channel reach them.
- `design/conventions/git.md`, `work-management.md`, their template copies, the MCP server's instructions, and the `claude-code` harness prompt send the agent to `flai task done` or `task_done` at every task transition.
- A failed check comes after the move and the log: the task is done, and the agent fixes what the check found and calls again.
- One call is one turn: E-0017 measures the ceremony turns it removes with `flai stats` (S-0293).

## Alternatives considered

- A script under `scripts/` chaining the commands: each step still prints for the agent to read, and MCP and the host channel cannot reach it.
- Keep the steps as separate calls and shorten their output: the turns, and the context each re-reads, stay.
- Move the task before the commit, so the commit records the move: `wip/` is written in the main checkout, never on the story branch (ADR-0019), so there is nothing of the move to commit there.
