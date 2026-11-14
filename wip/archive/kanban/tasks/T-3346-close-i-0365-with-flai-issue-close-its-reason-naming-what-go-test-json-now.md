---
id: T-3346
type: task
nature: remediation
title: Close I-0365 with flai issue close, its reason naming what go-test-json now names for a crashed test binary
status: done
parent: S-0782
owner: alex
created: 2026-11-13T22:53:19Z
updated: 2026-11-13T23:08:34Z
transitions:
  - to: ready
    at: 2026-11-13T23:08:23Z
    by: agent-S-0782
  - to: in-progress
    at: 2026-11-13T23:08:23Z
    by: agent-S-0782
  - to: done
    at: 2026-11-13T23:08:34Z
    by: agent-S-0782
stream: S-0782
tags: [flai, issues]
touches: [design/issues/I-0365-flai-test-reports-a-go-test-failure-in-flai-cmd-that-names-no-test-and-the-same.md]
after: [T-3344]
usage:
  source: log
  seconds: 11
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 6
      output: 1767
      cache_read: 434516
      cache_write: 21285
      cost: 0.2676
---
# T-3346 Close I-0365 with flai issue close, its reason naming what go-test-json now names for a crashed test binary

## Work

Criterion 2. Run `flai test flai/cmd/guard.go` once in the worktree, the path the I-0365 instance ran again, and record its time and result in the narrative.

Then close the issue with `flai issue close I-0365 --reason "<reason>"`. The reason names S-0782, that a package which fails with no test failing is now reported by the test its crashed binary was running, with the crash line and the exit, else by the package's last output lines, and that the flake behind the instance kept no output to name: a recurrence now names it and is recorded as a new instance.

It waits for T-3344: the reason names what T-3344 changed.

## Done when

- `flai test flai/cmd/guard.go` ran in the worktree, and its time and result are in the narrative.
- I-0365 is closed, and its reason names S-0782 and the go-test-json change.

## Notes
