---
id: T-1513
type: task
nature: remediation
title: Close I-0140, saying the OS file lock in flai/internal/filelock fixed it
status: done
parent: S-0377
owner: alex
created: 2026-10-13T16:48:28Z
updated: 2026-10-13T17:07:28Z
transitions:
  - to: ready
    at: 2026-10-13T17:07:20Z
    by: agent-S-0377
  - to: in-progress
    at: 2026-10-13T17:07:22Z
    by: agent-S-0377
  - to: done
    at: 2026-10-13T17:07:28Z
    by: agent-S-0377
stream: S-0377
tags: [flai, issues]
touches: [design/issues/I-0140-testactivitywriterstaketurns-fails-about-one-run-in-three-two-writers-that-both-find-the-activity-lock-stale-can-both-take-it-and-an-entry-is-lost.md, design/issues/summary.md]
after: [T-1511, T-1512]
usage:
  source: log
  seconds: 4
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 4
      output: 48
      cache_read: 368044
      cache_write: 2349
      cost: 0.1702
---
# T-1513 Close I-0140, saying the OS file lock in flai/internal/filelock fixed it

## Work

In the story's worktree, run `flai issue close I-0140 --reason "<reason>"`. The reason says that `lockActivity`, the guard's `lockRecord`, and flai serve's `lockFile` now take an operating-system file lock through `flai/internal/filelock`. That lock is freed when its holder exits and is never broken by its age, so two writers cannot both hold it. It names the test that now passes fifty runs in a row. flai writes the issue and regenerates `design/issues/summary.md`.

Then tick criterion 1 and criterion 2 with `flai criteria tick S-0377`, each once verified.

It waits for T-1511 and T-1512: the issue closes only once every lock that had the cause is fixed.

## Done when

- I-0140 has status `closed`, with the reason above.
- `design/issues/summary.md` no longer lists I-0140.
- Both criteria of S-0377 are ticked.

## Notes
