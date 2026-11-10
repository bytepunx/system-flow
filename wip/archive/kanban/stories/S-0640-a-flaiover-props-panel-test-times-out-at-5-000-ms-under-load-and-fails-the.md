---
id: S-0640
type: story
nature: remediation
title: A flaiover props-panel test times out at 5,000 ms under load and fails the flaiover tier, serial or beside integration
status: done
owner: alex
created: 2026-11-09T16:55:57Z
updated: 2026-11-10T05:54:32Z
transitions:
  - to: ready
    at: 2026-11-09T22:20:01Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-10T05:12:17Z
    by: agent-S-0640
  - to: review
    at: 2026-11-10T05:51:59Z
    by: agent-S-0640
  - to: done
    at: 2026-11-10T05:54:32Z
    by: alex
tags: [flaiover, testing, flaky-test]
topics: [testing]
touches: [flaiover/src/routes/design/components/components.svelte.test.ts, design/issues/I-0319-a-flaiover-props-panel-test-times-out-at-5-000-ms-under-load-and-fails-the.md, design/issues/summary.md, design/issues/I-0330-a-flaiover-design-page-test-takes-up-the-draft-another-tab-writes-times-out-at.md, design/issues/I-0323-flai-check-finds-item-duplicate-id-outside-the-story-at-close-out.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2402
  priced_by: harness
  turns:
    - day: 2026-11-10
      ceremony: 4
      hand_edits: 3
      work: 68
  models:
    - model: claude-opus-5-5
      input: 158
      output: 16804
      cache_read: 16857188
      cache_write: 246593
      cost: 5.6809
  strategic:
    - kind: orchestrator
      seconds: 66
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 8
          output: 78
          cache_read: 664480
          cache_write: 7688
          cost: 0.179
        - model: claude-sonnet-5-5
          input: 2
          output: 11
          cache_read: 14233
          cache_write: 8148
          cost: 0.0166
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-11-09T16:55:57Z
  value: 37.5
  by: planner-S-0640
  at: 2026-11-09T22:07:30Z
forecast:
  duration: 45m
  delivery: 2026-11-10T00:19:00Z
  basis: "Its own forecast of 45m; 1st in the pull order with an in-progress limit of 5, behind S-0633, S-0636 and S-0638."
  by: flai
  at: 2026-11-09T23:22:43Z
finalized:
  by: orchestrator
  at: 2026-11-09T22:19:29Z
---
# S-0640 A flaiover props-panel test times out at 5,000 ms under load and fails the flaiover tier, serial or beside integration

## Goal

This story remediates [I-0319](../../../design/issues/I-0319-a-flaiover-props-panel-test-times-out-at-5-000-ms-under-load-and-fails-the.md), "A flaiover props-panel test times out at 5,000 ms under load and fails the flaiover tier, serial or beside integration". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0319 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0319 is closed with `flai issue close I-0319 --reason` saying what fixed it

## Tasks
- T-2805 Post the props-panel test's 55-call burst without settling between calls, and measure the test before and after
- T-2807 Close I-0319 with what fixed it

## Notes

Cost of delay inputs set by flai from I-0319. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-09T10:34:28Z, 0.3 days before this story; under one cycle counts as one).

### Planning

The proposed cause, from the test itself: `flaiover/src/routes/design/components/components.svelte.test.ts` › "lists the calls of the component's function props, and clears them on another pick" mounts the page, makes 2 clicks, and posts 59 harness messages. Each of those 62 steps calls the file's `settle()`, which waits out 20 `setTimeout(0)` turns and then flushes. That is about 1,240 timer turns in one test, each at least a millisecond and many more under a load average of 50 to 125, so the test crosses vitest's 5,000 ms on a busy host whatever runs beside it. The page's own handler, `called` in `+page.svelte`, only prepends to a log capped at 50 entries, so the page is not the slow part. The proposed fix posts the 55-message burst in one go and settles once after it. It keeps every assertion and also checks that a burst of calls lands whole. Raising the test's timeout would only hide the cost, so it is not proposed.

Touches:

| Path | Source | Why |
|------|--------|-----|
| `flaiover/src/routes/design/components/components.svelte.test.ts` | design (I-0319's instances) | the test that times out, changed by T-2805 |
| `design/issues/I-0319-a-flaiover-props-panel-test-times-out-at-5-000-ms-under-load-and-fails-the.md` | design (criterion 2) | `flai issue close` writes it, in T-2807 |
| `design/issues/summary.md` | co-change (2 of 8 commits) | `flai issue close` regenerates it, in T-2807 |

`flai touches suggest` also listed `+page.svelte` (5 of 8), `ComponentFrame.svelte` (4), `ComponentList.svelte`, `PropsPanel.svelte.test.ts`, `protocol.ts`, `protocol.test.ts`, and `server/repo.ts` (2 each). They are left out: they co-changed with the test when the Components page was built (S-0440 to S-0476), and the fix changes no page code. The story's agent adds `+page.svelte` only if the burst shows the page dropping calls.

Forecast: `flai forecast` gave 5m from the median of 17 small remediation stories at size 2. It is raised to 45m because the close-out alone runs the flaiover tier, about 6m, and the change needs vitest runs before and after to show the time it saves. The delivery moves 40m past `flai forecast`'s 23:50 for the same reason.

Cost of delay: `flai cod` gives 37.50 USD a week from the 15m time_lost_per_cycle flai set from I-0319, at 150 USD an hour, and it is written as it stands. I-0319 has had 2 instances and 11m of cost since those inputs were set, so the inputs may now be low; they are the operator's to change, not the planner's.

No `after`: the story reads no other story's output.

### Outcome

Criterion 1: the timeout needs a host load of 50 or more, so no deterministic test reproduces it, and none is added. The evidence is the test's own duration. It took 2,174 to 2,447 ms before the fix and 672 to 682 ms after, at a load of 12 to 19 on 24 cores. The test now also checks that a 55-call burst lands whole, all 50 entries newest first.
