---
id: T-1879
type: task
nature: feature
title: flai/internal/components starts a dev server in a story's worktree, installing the sub-project's dependencies by its lock file when they are missing
status: done
parent: S-0442
owner: alex
created: 2026-10-16T21:37:19Z
updated: 2026-10-17T02:53:25Z
transitions:
  - to: ready
    at: 2026-10-17T02:38:22Z
    by: agent-S-0442
  - to: in-progress
    at: 2026-10-17T02:38:23Z
    by: agent-S-0442
  - to: done
    at: 2026-10-17T02:53:25Z
    by: agent-S-0442
stream: S-0442
tags: [cli]
touches: [flai/internal/components/worktree.go, flai/internal/components/worktree_test.go, flai/internal/components/devserver.go, flai/internal/components/devserver_test.go]
usage:
  source: log
  seconds: 902
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 23
      output: 9410
      cache_read: 1671638
      cache_write: 43972
      cost: 0.7803
---
# T-1879 flai/internal/components starts a dev server in a story's worktree, installing the sub-project's dependencies by its lock file when they are missing

## Work

- `worktree.go`: `Root(project, story)` resolves the story's worktree as ADR-0145's saves do, refused for a story without one. `EnsureDeps(sub-project)` finds the lock file and runs `pnpm install --frozen-lockfile`, `npm ci`, or `yarn install --frozen-lockfile` when `node_modules` is missing or older than the lock file, with a timeout, returning the last lines on failure.
- `devserver.go`: `Start` takes the root to run in, so a story's dev server gets its own config folder (`.flai-cache/components/<project>/<sub-project>/<story>/`) and port.

Waits for nothing in this story; it extends S-0438's runner.

## Done when

Tests with a fake package manager cover each lock file, deps already current, an install failure, and a story with no worktree. `flai test flai/internal/components` passes.

## Notes

None.
