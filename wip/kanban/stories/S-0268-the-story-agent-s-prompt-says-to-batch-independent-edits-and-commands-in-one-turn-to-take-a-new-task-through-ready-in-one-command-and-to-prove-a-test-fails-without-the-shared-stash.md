---
id: S-0268
type: story
nature: improvement
title: The story agent's prompt says to batch independent edits and commands in one turn, to take a new task through ready in one command, and to prove a test fails without the shared stash
status: backlog
owner: alex
created: 2026-10-05T00:06:36Z
updated: 2026-10-05T00:06:36Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0268 The story agent's prompt says to batch independent edits and commands in one turn, to take a new task through ready in one command, and to prove a test fails without the shared stash

## Goal

In S-0248's run the main agent made 64 model calls, median 2.8 seconds, 3.5 minutes in all, for three edits and a test: 32 of its 48 API messages carried one tool call, three consecutive edits to one file went out as three turns, and a `flai move T-nnnn in-progress` straight from backlog was refused and retried through ready. To prove its new test fails without the fix, it pushed and dropped a `git stash` in the stash stack every worktree on the host shares, which another session's stash could have been caught by. The prompt flai serve builds for a story agent (`flai/internal/harness/harness.go`) and the conventions it points at (`design/conventions/delegation.md`, `git.md`, and the template's copies) say none of this. Each is a sentence; together they save a minute or two per story and remove one hazard.

## Acceptance criteria
- [ ] The story agent's prompt says to make independent edits and commands in one turn and to move a task it has just written to ready and in-progress in one command
- [ ] `design/conventions/git.md` and the template's copy say how to show a test fails without the change under test (check the files out from the main branch into a scratch copy, or build the old binary) and that `git stash` is not used in a worktree, since the stash stack is shared across the host's worktrees
- [ ] The prompt's tests in `flai/internal/harness` cover the new sentences, and the design (`design/system/flai-cli.md` or where the prompt is described) records them

## Tasks

## Notes

Found by the operator's review of the S-0248 agent log on 2026-10-04 (`~/.flai/serve/agents/sf-S-0248-20261004T231346Z.log`, 23:14:54Z and 23:15:45Z to 23:16:14Z).
