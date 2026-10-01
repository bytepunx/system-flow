---
id: T-0663
type: task
nature: feature
title: flai stream open reopens a stream whose worktree is gone, from the local or remote story branch, and the narrative records its host
status: done
parent: S-0177
owner: alex
created: 2026-10-01T10:06:53Z
updated: 2026-10-01T10:11:31Z
transitions:
  - to: ready
    at: 2026-10-01T10:07:15Z
    by: agent-S-0177
  - to: in-progress
    at: 2026-10-01T10:08:14Z
    by: agent-S-0177
  - to: done
    at: 2026-10-01T10:11:31Z
    by: agent-S-0177
stream: S-0177
tags: []
touches: [flai/cmd/stream.go, flai/cmd/branch.go, flai/internal/storygit, flai/internal/workitem]
usage:
  source: log
  seconds: 197
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 36
      output: 150
      cache_read: 2568074
      cache_write: 38507
      cost: 1.0407
---
# T-0663 flai stream open reopens a stream whose worktree is gone, from the local or remote story branch, and the narrative records its host

## Work

- `flai stream open` on a story whose narrative exists and whose worktree does not: leave the narrative as it is, and check out `story/S-nnnn` from the local branch, else fetch it from the remote when it is there, else create it from the main branch. Refuse, as now, when the worktree exists too.
- The narrative's front matter records `host`, the name of the host whose `flai stream open` wrote it; reopening on another host updates it.
- Tests: a reopen with the branch only on the remote, with no branch anywhere, and with the worktree present.
- `design/system/flai-cli.md`, `design/system/agent-narrative.md`, `docs/users/flai.md`, and the generated reference say so.

## Done when

The tests pass, the docs say what the command does, and the change is committed on `story/S-0177`.
