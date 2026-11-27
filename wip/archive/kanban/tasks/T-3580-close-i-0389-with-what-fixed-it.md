---
id: T-3580
type: task
nature: improvement
title: Close I-0389 with what fixed it
status: done
parent: S-0847
owner: alex
created: 2026-11-17T12:29:54Z
updated: 2026-11-27T23:24:18Z
transitions:
  - to: ready
    at: 2026-11-27T23:24:06Z
    by: agent-S-0847
  - to: in-progress
    at: 2026-11-27T23:24:06Z
    by: agent-S-0847
  - to: done
    at: 2026-11-27T23:24:18Z
    by: agent-S-0847
stream: S-0847
tags: [issues]
touches: [design/issues/I-0389-flai-serve-opens-a-new-board-starvation-thread-at-each-flip-of-the-board-48-in.md, design/issues/summary.md]
after: [T-3573]
usage:
  source: log
  seconds: 12
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 2108
      cache_read: 569171
      cache_write: 14496
      cost: 0.262
---
# T-3580 Close I-0389 with what fixed it

## Work

This task waits for T-3573: the reason it closes the issue with names the fix that task builds and the test that reproduces I-0389.

Run `flai issue close I-0389 --reason "<reason>"`. The reason names the ADR T-3569 wrote, the settle window, and the test in `flai/internal/serve/starved_test.go` that reproduces the flips. flai writes the issue file and `design/issues/summary.md`; edit neither by hand.

## Done when

- I-0389's status is closed, with a reason that names the ADR and the test.
- `design/issues/summary.md` no longer lists I-0389 as open.

## Notes
