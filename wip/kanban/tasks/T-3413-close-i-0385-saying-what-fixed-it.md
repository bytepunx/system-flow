---
id: T-3413
type: task
nature: remediation
title: Close I-0385 saying what fixed it
status: done
parent: S-0800
owner: alex
created: 2026-11-16T01:45:22Z
updated: 2026-11-16T02:19:52Z
transitions:
  - to: ready
    at: 2026-11-16T01:47:32Z
    by: agent-S-0800
  - to: in-progress
    at: 2026-11-16T02:19:38Z
    by: agent-S-0800
  - to: done
    at: 2026-11-16T02:19:52Z
    by: agent-S-0800
stream: S-0800
tags: [issues]
touches: [design/issues/I-0385-flai-stream-sync-reports-a-rebased-branch-undoes-main-s-line-when-the-branch.md, design/issues/summary.md]
after: [T-3410, T-3412]
usage:
  source: log
  seconds: 14
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 8
      output: 2099
      cache_read: 586209
      cache_write: 21680
      cost: 0.2995
---
# T-3413 Close I-0385 saying what fixed it

## Work

Run `flai issue close I-0385 --reason "<what fixed it>"`, naming the new ADR, the change to the undone check in `flai/internal/storygit/undone.go`, and the real-git test that reproduces the issue. It regenerates `design/issues/summary.md`.

It waits for T-3410, whose test shows the cause no longer occurs, and for T-3412, so that the reason can point at the design as it now reads.

## Done when

- I-0385's status is closed, with a reason that names the ADR, the code change, and the test.
- `design/issues/summary.md` no longer lists I-0385, and `flai check --strict` is clean.

## Notes
