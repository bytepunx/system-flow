---
id: S-0264
type: story
nature: remediation
title: flaiover's HostProcesses test that waits past the old host after an upgrade fails under machine load
status: done
owner: alex
created: 2026-10-05T00:03:14Z
updated: 2026-10-06T22:56:13Z
transitions:
  - to: ready
    at: 2026-10-06T22:47:41Z
    by: alex
  - to: in-progress
    at: 2026-10-06T22:47:49Z
    by: agent-S-0264
  - to: review
    at: 2026-10-06T22:55:31Z
    by: agent-S-0264
  - to: done
    at: 2026-10-06T22:56:13Z
    by: alex
tags: [flaiover]
topics: [testing]
touches: [flaiover/src/lib/components/HostProcesses.svelte.test.ts, design/issues/I-0053-flaiover-s-hostprocesses-test-that-waits-past-the-old-host-after-an-upgrade-fails-under-machine-load.md, design/issues/summary.md, design/issues/I-0073-flai-check-finds-threads-archived-outside-the-story-at-close-out.md, design/issues/I-0076-flai-check-finds-wip-overlap-outside-the-story-at-close-out.md, design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 473
  models:
    - model: claude-opus-5-5
      input: 62
      output: 17242
      cache_read: 2247907
      cache_write: 95678
      cost: 1.5601
    - model: claude-sonnet-5-5
      input: 12
      output: 3120
      cache_read: 155514
      cache_write: 44175
      cost: 0.1728
cost_of_delay:
  inputs:
    time_lost_per_cycle: 4m
    by: flai
    at: 2026-10-05T00:03:14Z
  value: 10
  by: planner-S-0264
  at: 2026-10-05T05:50:19Z
forecast:
  duration: 30m
  delivery: 2026-10-07T06:20:00Z
  basis: "Its own forecast of 30m; 23rd in the pull order with an in-progress limit of 3, behind S-0299, S-0301, S-0300, S-0212, S-0213, S-0214, S-0215, S-0216, S-0228, S-0232, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0245, S-0246, S-0251, S-0254 and S-0261."
  by: flai
  at: 2026-10-06T22:46:26Z
finalized:
  by: alex
  at: 2026-10-06T22:47:39Z
---
# S-0264 flaiover's HostProcesses test that waits past the old host after an upgrade fails under machine load

## Goal

This story remediates [I-0053](../../../design/issues/I-0053-flaiover-s-hostprocesses-test-that-waits-past-the-old-host-after-an-upgrade-fails-under-machine-load.md), "flaiover's HostProcesses test that waits past the old host after an upgrade fails under machine load". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0053 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0053 is closed with `flai issue close I-0053 --reason` saying what fixed it

## Tasks
- T-0977 HostProcesses' reconnect tests wait for the outcome rather than a fixed 40 ms, with a test that reproduces a slow poll
- T-0979 The HostProcesses tests pass while the Go tiers load the machine, and I-0053 is closed

## Notes

Cost of delay inputs set by flai from I-0053. time_lost_per_cycle 4m: 4m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-01T11:35:40Z, 3.5 days before this story; under one cycle counts as one).

### Planning

Proposed solution, from the instance: the reconnect tests in `HostProcesses.svelte.test.ts` run on real timers cut short by `fast` (`disconnectTimeoutMs: 5, reconnectPollMs: 5, reconnectGiveUpMs: 50`) and assert after a fixed wall-clock `settleThrough(40)`. Under load a 5 ms timer stretches past 40 ms, so the assertion runs before the second poll answers, and `waitForReconnect`'s `Date.now()` deadline of 50 ms can pass first and report a give-up. The tests should wait for their outcome (`vi.waitFor`) with a give-up margin of seconds, keeping real timers (the HostPanel test explains why fake timers were avoided). The component already takes its timings as props and needs no change.

Touches:

- `flaiover/src/lib/components/HostProcesses.svelte.test.ts`: design, the file I-0053 names, where the fixed waits are; T-0977.
- `design/issues/I-0053-...md`: design, the second criterion closes it; T-0979.
- `design/issues/summary.md`: layout, `flai issue close` rewrites it; T-0979.
- Left out: `flai touches suggest` from the test and component lists `flaiover/src/routes/api/host/+server.ts` and `host.test.ts` (co-change, 2 of 4 commits); the fix is in the test alone, so neither changes. `HostProcesses.svelte` is left out as the timings are already props; T-0977 adds it if the reproduction shows otherwise.

Tasks: T-0977 (layer 1) adds the reproduction and the waits; T-0979 (layer 2, after T-0977) runs the tier under Go-tier load and closes the issue.

Figures:

- Forecast: flai gave 12m (134 s per unit over 4 done remediation stories, size 5). Raised to 30m: a load-dependent flake has to be shown passing under load, which means several unit-tier runs alongside the Go tiers, plus the close-out's full run. Delivery is flai's queue estimate (32nd in the pull order) moved by the added time.
- Cost of delay: 10 USD a week, as `flai cod` gives from the 4m-per-cycle input at 150 USD an hour. It stands: one instance, and a rerun of the close-out is what each one costs.
- Topic `testing` added: the story is about a test's timing.
