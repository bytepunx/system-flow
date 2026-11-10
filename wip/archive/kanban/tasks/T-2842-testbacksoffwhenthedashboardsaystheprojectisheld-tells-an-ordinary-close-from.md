---
id: T-2842
type: task
nature: remediation
title: TestBacksOffWhenTheDashboardSaysTheProjectIsHeld tells an ordinary close from the held back-off by a wide gap
status: done
parent: S-0651
owner: alex
created: 2026-11-10T07:01:00Z
updated: 2026-11-10T07:23:48Z
transitions:
  - to: ready
    at: 2026-11-10T07:22:58Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T07:22:58Z
    by: system-flow
  - to: done
    at: 2026-11-10T07:23:48Z
    by: agent-S-0651
stream: S-0651
tags: [cli, testing]
touches: [flai/internal/channel/channel_test.go]
usage:
  source: log
  seconds: 50
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 701
      cache_read: 261448
      cache_write: 5893
      cost: 0.1135
---
# T-2842 TestBacksOffWhenTheDashboardSaysTheProjectIsHeld tells an ordinary close from the held back-off by a wide gap

## Work

The test sets `MinBackoff` 10 ms, `MaxBackoff` 40 ms, and `HeldBackoff` 600 ms, which `jitter` spreads over 300 to 600 ms. It asserts that an ordinary close is dialled again within 250 ms, and under a load average of about 160 one waited 252 ms (I-0309, instance of 2026-11-09T22:38:42Z).

- Widen the gap between the two waits, so that load cannot carry an ordinary close's retry past the bound: a longer `HeldBackoff`, such as 3 s, whose jittered lower half, 1.5 s, stays well above the bound on the ordinary close, raised to about 1 s. Keep the held-back-off assertion at the jittered lower half of the new `HeldBackoff`, and keep the `trying again in` text of the state in step with it.
- Keep the test's two claims: a held close waits the held back-off, and any other close is retried at once.
- Reproduce the failure before the change under a CPU load, as T-2839 does, and record in the task's Notes what failed and how.

It waits for no task: it changes only this file, which no other task of the story touches.

## Done when

- `go test -race -run '^TestBacksOffWhenTheDashboardSaysTheProjectIsHeld$' -count=20 ./flai/internal/channel/` passes under the load that made the old test fail, recorded in the task's Notes.
- `flai test flai/internal/channel/channel_test.go` passes.

## Notes
