---
id: S-0655
type: story
nature: remediation
title: "A flaiover design-page test, \"takes up the draft another tab writes\", times out at 5,000 ms under load and stops the flaiover tier"
status: ready
owner: alex
created: 2026-11-10T05:54:35Z
updated: 2026-11-10T07:08:49Z
transitions:
  - to: ready
    at: 2026-11-10T07:08:49Z
    by: orchestrator
tags: [flaiover, testing, flaky-test]
topics: [testing]
touches: [flaiover/src/routes/design/design.svelte.test.ts, design/issues/I-0330-a-flaiover-design-page-test-takes-up-the-draft-another-tab-writes-times-out-at.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: planner
      seconds: 254
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 284
          output: 29347
          cache_read: 16087933
          cache_write: 867284
          cost: 6.9457
    - kind: orchestrator
      seconds: 31
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 40
          cache_read: 279468
          cache_write: 5646
          cost: 0.0767
cost_of_delay:
  inputs:
    time_lost_per_cycle: 16m
    by: flai
    at: 2026-11-10T05:54:35Z
  value: 40
  by: planner-S-0655
  at: 2026-11-10T06:59:13Z
forecast:
  duration: 45m
  delivery: 2026-11-10T09:54:00Z
  basis: "Its own forecast of 45m; 10th in the pull order with an in-progress limit of 5, behind S-0628, S-0636, S-0638, S-0647, S-0656, S-0241, S-0637, S-0648, S-0649, S-0650, S-0651, S-0652, S-0653 and S-0654."
  by: flai
  at: 2026-11-10T07:02:19Z
finalized:
  by: orchestrator
  at: 2026-11-10T07:08:13Z
---
# S-0655 A flaiover design-page test, "takes up the draft another tab writes", times out at 5,000 ms under load and stops the flaiover tier

## Goal

This story remediates [I-0330](../../../design/issues/I-0330-a-flaiover-design-page-test-takes-up-the-draft-another-tab-writes-times-out-at.md), "A flaiover design-page test, "takes up the draft another tab writes", times out at 5,000 ms under load and stops the flaiover tier". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0330 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0330 is closed with `flai issue close I-0330 --reason` saying what fixed it

## Tasks
- T-2838 Measure where the design page's draft-from-another-tab test spends its time, and cut it below vitest's 5,000 ms under load
- T-2840 Close I-0330 with what fixed it

## Notes

Cost of delay inputs set by flai from I-0330. time_lost_per_cycle 16m: 16m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-10T05:37:30Z, 0 days before this story; under one cycle counts as one).

### Planning

The cause is not known yet, and T-2838 finds it. The failing test is `flaiover/src/routes/design/design.svelte.test.ts` › "editing the design system on the design page (S-0383)" › "takes up the draft another tab writes, and drops its edits when another tab drops it". It calls `open()` once, which mounts the whole Design page, and then `settle()` three times. Each `settle()` waits out 20 `setTimeout(0)` turns and flushes, so the test runs about 60 timer turns, against the 1,240 that caused I-0319 (fixed by S-0640). The issue puts the time in mounting and rendering the page at a load of about 110 on 24 cores. Every other test in the same `describe` mounts the page too, so the fix T-2838 finds may apply to the whole block. Raising the timeout alone would only hide the cost, so the plan prefers cutting it, as S-0640 did. If the time is the page's own work, which the test cannot cut, T-2838 may give the slow tests a stated timeout and record why.

Touches:

| Path | Source | Why |
|------|--------|-----|
| `flaiover/src/routes/design/design.svelte.test.ts` | design (I-0330's instance) | the test that times out, changed by T-2838 |
| `design/issues/I-0330-a-flaiover-design-page-test-takes-up-the-draft-another-tab-writes-times-out-at.md` | design (criterion 2) | `flai issue close` writes it, in T-2840 |
| `design/issues/summary.md` | co-change (S-0640's close of I-0319) | `flai issue close` regenerates it, in T-2840 |

`flai touches suggest` from the test and the issue also listed `flaiover/src/routes/design/+page.svelte`, in 5 of 8 commits. It is left out because it changed alongside the test while the Design page was being built (S-0382 to S-0384), and the plan changes no page code. The story's agent adds it only if the measurement shows the page's mount doing work it should not. `flai task done` would widen the touches then.

Tasks: T-2838 (layer 1) measures and fixes the test. T-2840 (layer 2, after T-2838) closes I-0330 with T-2838's figures.

Forecast: `flai forecast` gave 5m, from the median of 17 small remediation stories at size 2. It is raised to 45m. S-0640 has the same shape: one test timing out under load, measured, fixed in the test, then the issue closed. It took 40m from start to review, because the close-out alone runs the flaiover tier and the change needs vitest runs before and after to show the time saved. flai places the delivery from the 45m and the pull order.

Cost of delay: `flai cod` gives 40.00 USD a week, from the 16m `time_lost_per_cycle` that flai set from I-0330, at 150 USD an hour. It is written as it stands. The inputs are the operator's.

No `after`: the story reads no other story's output. S-0640's fix of I-0319 is already on main, and T-2838 reads it as a pattern only.
