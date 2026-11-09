---
id: T-2766
type: task
nature: remediation
title: The guard passes git tag, worktree, and branch in the forms that only read
status: done
parent: S-0623
owner: alex
created: 2026-11-09T21:55:57Z
updated: 2026-11-09T22:10:22Z
transitions:
  - to: ready
    at: 2026-11-09T21:58:24Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T21:58:25Z
    by: system-flow
  - to: done
    at: 2026-11-09T22:10:22Z
    by: agent-S-0623
stream: S-0623
tags: [flai, guard]
touches: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go]
usage:
  source: log
  seconds: 717
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 46
      output: 7645
      cache_read: 5593256
      cache_write: 167801
      cost: 2.3706
---
# T-2766 The guard passes git tag, worktree, and branch in the forms that only read

## Work

In `flai/internal/guard/guard.go`, judge `tag`, `worktree`, and `branch` by their form rather than by name alone, next to `gitReads`, as `flagReads` does for flai's commands. `git tag` reads when it lists: no positional word, or `-l`/`--list` (whose positionals are patterns), with list options such as `--sort`, `--contains`, `--points-at`, `--merged`, `-n`, and none of `-a`, `-s`, `-u`, `-m`, `-F`, `-d`, `-f`, `-e`, `--annotate`, `--sign`, `--delete`, `--force`. `git worktree` reads with `list`. `git branch` reads when it lists: no positional word, or `--list`/`-l`, `--show-current`, `-a`, `-r`, `-v`, `--contains`, `--merged`, and none of `-d`, `-D`, `-m`, `-M`, `-c`, `-C`, `-f`, `-u`, `--set-upstream-to`, `--unset-upstream`, `--edit-description`. The global options `-C`, `-c`, `--git-dir` stay skipped as the `case "git"` branch skips them. Every other form keeps its refusal, and the orchestrator's `git tag <name>` keeps `publishesThrough`.

In `flai/internal/guard/guard_test.go`, add a table test that reproduces I-0289: `git tag --sort=-creatordate`, `git tag --list`, `git tag -l 'flai/v*'`, `git worktree list`, `git -C /w branch --show-current` pass for a sub-agent, the planner, and the orchestrator; `git tag v1.2.3`, `git tag -d v1`, `git worktree add x`, `git branch -D x` are refused as before. Waits for no task.

## Done when

- `git tag --sort=-creatordate` and `git worktree list` pass the orchestrator's guard, and the new test fails without the change.
- The writing forms are refused with the reasons they had, and `TestTheOrchestratorPublishesOnlyThroughReleasePublish` still passes.
- `flai test flai/internal/guard` passes.

## Notes
