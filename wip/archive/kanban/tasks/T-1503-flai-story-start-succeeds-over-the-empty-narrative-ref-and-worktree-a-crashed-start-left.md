---
id: T-1503
type: task
nature: remediation
title: flai story start succeeds over the empty narrative, ref, and worktree a crashed start left
status: done
parent: S-0375
owner: alex
created: 2026-10-13T16:36:55Z
updated: 2026-10-13T17:01:35Z
transitions:
  - to: ready
    at: 2026-10-13T16:58:33Z
    by: agent-S-0375
  - to: in-progress
    at: 2026-10-13T16:58:33Z
    by: agent-S-0375
  - to: done
    at: 2026-10-13T17:01:35Z
    by: agent-S-0375
stream: S-0375
tags: [flai, test]
touches: [flai/internal/storystart/start_test.go]
after: [T-1501, T-1502]
usage:
  source: log
  seconds: 182
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 37
      output: 13059
      cache_read: 1718417
      cache_write: 70786
      cost: 1.0497
---
# T-1503 flai story start succeeds over the empty narrative, ref, and worktree a crashed start left

## Work

The end-to-end reproduction of I-0137, for the story's first criterion. Beside `TestStartOpensTheBranchAndWorktree` in `flai/internal/storystart/start_test.go`, with real git and the file's `newFixture`:

- Leave what S-0337's crash left: the story in progress, a zero-byte `wip/agents/S-nnnn.md`, a zero-byte `refs/heads/story/S-nnnn` and reflog, and a worktree folder with a zero-byte `.git` and zero-byte files.
- Run `storystart.Start` on the story. It succeeds, the narrative has the template's sections, the branch resolves to main's head, and the worktree is a working checkout of it.

Waits for T-1501 and T-1502: the test passes only once `OpenStream` writes over an empty narrative and `storygit.Open` repairs the ref and worktree. If `storystart/start.go` needs a change to pass, add it to this task's touches with `flai touches` and say why in the narrative's `## Decisions`.

## Done when

- The test reproduces I-0137's state and `Start` recovers from it.
- `flai test flai/internal/storystart` passes.

## Notes

Drafted by the planner from I-0137's instance.
