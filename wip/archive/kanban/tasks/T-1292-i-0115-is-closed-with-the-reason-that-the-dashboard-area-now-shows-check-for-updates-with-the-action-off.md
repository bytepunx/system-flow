---
id: T-1292
type: task
nature: remediation
title: I-0115 is closed with the reason that the Dashboard area now shows Check for updates with the action off
status: done
parent: S-0314
owner: alex
created: 2026-10-07T23:42:03Z
updated: 2026-10-07T23:45:24Z
transitions:
  - to: ready
    at: 2026-10-07T23:45:08Z
    by: agent-S-0314
  - to: in-progress
    at: 2026-10-07T23:45:09Z
    by: agent-S-0314
  - to: done
    at: 2026-10-07T23:45:24Z
    by: agent-S-0314
stream: S-0314
tags: [issues]
touches: [design/issues/I-0115-the-flaiover-guide-says-check-for-updates-always-works-but-the-dashboard-area-hides-it-while-the-dashboard-host-action-is-off.md, design/issues/summary.md]
after: [T-1290, T-1291]
usage:
  source: log
  seconds: 15
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 10
      output: 1755
      cache_read: 511314
      cache_write: 20007
      cost: 0.2975
---
# T-1292 I-0115 is closed with the reason that the Dashboard area now shows Check for updates with the action off

## Work

From the story's worktree, so the write lands on the story branch, run `flai issue close I-0115 --reason "<reason>"`. The reason says that `HostPanel.svelte` now shows Check for updates whatever the `dashboard` host action says, as `docs/users/flaiover.md` § Host promises, that a test in `HostPanel.svelte.test.ts` reproduces it, and that `design/system/flaiover-dashboard.md` says so, naming S-0314. It rewrites the issue file and regenerates `design/issues/summary.md`.

Waits for T-1290 and T-1291: the reason states what they did, so the fix and its design must be done first. Second layer.

## Done when

- I-0115's status is closed with that reason
- `design/issues/summary.md` no longer lists I-0115 as open
- `flai check --strict` is clean for the story

## Notes
