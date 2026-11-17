---
id: T-3539
type: task
nature: remediation
title: Close I-0359 with the reason that the publish commit waits out a held index lock
status: done
parent: S-0833
owner: alex
created: 2026-11-17T08:30:10Z
updated: 2026-11-17T09:33:17Z
transitions:
  - to: ready
    at: 2026-11-17T09:32:49Z
    by: agent-S-0833
  - to: in-progress
    at: 2026-11-17T09:32:49Z
    by: agent-S-0833
  - to: done
    at: 2026-11-17T09:33:17Z
    by: agent-S-0833
stream: S-0833
tags: [issues, release]
touches: [design/issues/I-0359-flai-release-pending-fails-when-another-process-holds-the-main-checkout-s-git.md, design/issues/summary.md]
after: [T-3535]
usage:
  source: log
  seconds: 28
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 14
      output: 2476
      cache_read: 1001463
      cache_write: 32967
      cost: 0.4962
---
# T-3539 Close I-0359 with the reason that the publish commit waits out a held index lock

## Work

Run `flai issue close I-0359 --reason "<reason>"` in the story's worktree. The reason names what fixed it: `flai release --pending` runs its publish's `git add -A` and `git commit` through `storygit.RunPastIndexLock`, as I-0100's fix does for `flai accept`, and `flai/cmd/release_lock_test.go` reproduces the held lock. `flai issue close` writes the issue file and regenerates `design/issues/summary.md`.

This task waits for T-3535, because the issue is closed only once the fix and its test are committed and pass.

## Done when

- I-0359's status is `closed`, with the reason naming the fix and the test.
- `design/issues/summary.md` no longer lists I-0359 as open.
- `flai check --strict` passes on the story.

## Notes
