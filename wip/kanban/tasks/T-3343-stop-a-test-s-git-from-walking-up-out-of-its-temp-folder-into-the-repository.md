---
id: T-3343
type: task
nature: remediation
title: Stop a test's git from walking up out of its temp folder into the repository that holds it
status: done
parent: S-0781
owner: alex
created: 2026-11-13T22:53:10Z
updated: 2026-11-13T23:10:53Z
transitions:
  - to: ready
    at: 2026-11-13T23:00:10Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T23:00:10Z
    by: system-flow
  - to: done
    at: 2026-11-13T23:10:53Z
    by: agent-S-0781
stream: S-0781
tags: [flai, tests, git]
touches: [flai/cmd/cmd_test.go, flai/internal/mcpserver/server_test.go, flai/cmd/story_start_test.go]
usage:
  source: log
  seconds: 643
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 41
      output: 12326
      cache_read: 2054454
      cache_write: 91687
      cost: 1.2385
---
# T-3343 Stop a test's git from walking up out of its temp folder into the repository that holds it

## Work

The four branches I-0371 names, `story/S-0001`, `story/S-0002`, `story/S-0003`, and `story/S-0005`, are not those stories' branches. Each points at 5c3b47d0c, a commit of S-0777's own, and `git worktree list` shows their worktrees under `.flai-cache/worktrees/S-0777/bin/S-0777/real/<Test…>/002/.flai-cache/worktrees/`. S-0777 ran the Go tests by hand with `TMPDIR` inside its worktree, as `wip/archive/agents/S-0777.md` records. `tempProject` (`flai/cmd/workitems_test.go`) makes a project with `t.TempDir()` and a manifest and runs no `git init`. So `storygit.InWorkTree` (`flai/internal/storygit/storygit.go`) ran `git rev-parse --is-inside-work-tree`, which walked up into S-0777's worktree and answered true. Then `storygit.Open` (`flai/internal/storygit/open.go`) ran `git worktree add -b story/S-nnnn` against the main repository. The tests that did it are the `TestInTheSessionOfAgentA` subtests of `flai/cmd/stream_open_test.go` and `TestTheOrchestratorSetsCostOfDelayInputsByItsPermissions`.

Make git stop at the temp root for the whole package, in the `TestMain` of `flai/cmd/cmd_test.go` and of `flai/internal/mcpserver/server_test.go`: set `GIT_CEILING_DIRECTORIES` to the temp root (`os.TempDir()`, resolved through its symbolic links, since `real` above is a link's target). Set it in `TestMain`, not with `t.Setenv`, which the parallel `together` subtests refuse. A temp project that runs `git init` itself still works, since git looks no higher than the project. If another package's tests start or open a story in a temp project with no `git init` of its own, give its `TestMain` the same setting and add its file to this task's touches.

Reproduce the leak in `flai/cmd/story_start_test.go` where a test fits. For example, `git init` an outer repository in a temp folder, make a project in a folder below it with no git of its own, start a story there with the package's environment, and check that it answers "not a git work tree" and the outer repository gains no `story/` branch. The test must fail without the ceiling.

It waits for nothing in this story.

## Done when

- `TestMain` in `flai/cmd` and `flai/internal/mcpserver` sets `GIT_CEILING_DIRECTORIES` to the resolved temp root.
- A test in `flai/cmd` shows that a story started in a temp project inside another repository makes no branch there, and it fails without the ceiling.
- `flai test` passes on the changed files.

## Notes
