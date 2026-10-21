---
id: S-0475
type: story
nature: remediation
title: The Components page cannot list, start, or stop anything through the dashboard, because the channel's project key and the components methods' sub-project share the one params field project
status: done
owner: alex
created: 2026-10-21T11:59:09Z
updated: 2026-10-21T12:13:56Z
transitions:
  - to: ready
    at: 2026-10-21T12:08:55Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-21T12:09:13Z
    by: agent-S-0475
  - to: review
    at: 2026-10-21T12:13:02Z
    by: agent-S-0475
  - to: done
    at: 2026-10-21T12:13:56Z
    by: orchestrator
tags: []
topics: [cli, dashboard]
touches: [design/issues/I-0163-the-components-page-cannot-list-start-or-stop-anything-through-the-dashboard-because-the-channel-s-project-key-and-the-components-methods-sub-project-share-the-one-params-field-project.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 245
  priced_by: harness
  turns:
    - day: 2026-10-21
      ceremony: 2
      hand_edits: 1
      work: 18
  models:
    - model: claude-opus-5-5
      input: 44
      output: 6307
      cache_read: 2327980
      cache_write: 124638
      cost: 1.589
  strategic:
    - kind: orchestrator
      seconds: 18
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 8
          output: 87
          cache_read: 1538910
          cache_write: 6214
          cost: 0.3412
cost_of_delay:
  inputs:
    penalty_per_week: 1000
    time_lost_per_cycle: 20m
    by: flai
    at: 2026-10-21T11:59:09Z
  value: 50
  by: planner-S-0475
  at: 2026-10-21T12:06:14Z
forecast:
  duration: 15m
  delivery: 2026-10-22T04:49:00Z
  basis: "flai forecast's 5m (median 131 s per size unit over 8 done remediation stories on claude-opus-5-5, size 2) raised to 15m because the story runs the S-0448 regression tests in flai/internal/hostapi, flai/internal/channel, and flaiover's agent and components tests, which its size does not count; delivery kept from flai forecast, 38th in the pull order."
  by: planner-S-0475
  at: 2026-10-21T12:06:14Z
finalized:
  by: alex
  at: 2026-10-21T12:04:35Z
---
# S-0475 The Components page cannot list, start, or stop anything through the dashboard, because the channel's project key and the components methods' sub-project share the one params field project

## Goal

This story remediates [I-0163](../../../design/issues/I-0163-the-components-page-cannot-list-start-or-stop-anything-through-the-dashboard-because-the-channel-s-project-key-and-the-components-methods-sub-project-share-the-one-params-field-project.md), "The Components page cannot list, start, or stop anything through the dashboard, because the channel's project key and the components methods' sub-project share the one params field project". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0163 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0163 is closed with `flai issue close I-0163 --reason` saying what fixed it

## Tasks
- T-1987 Confirm that S-0448's fix holds: the issue's reproduction answers and the regression tests through the channel pass
- T-1990 Close I-0163 with a reason naming S-0448's fix, its commits, and the releases that carry it

## Notes

Cost of delay inputs set by flai from I-0163. time_lost_per_cycle 20m: 20m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-18T02:25:28Z, 3.4 days before this story; under one cycle counts as one). penalty_per_week 1000 carried over from I-0163's Impact section.

### Planning

The cause is fixed already. S-0448, accepted 2026-10-18, made the components methods take their sub-project as `sub_project` (382fab56f) and made the dashboard's agent put the channel's project key after a call's params (d94c114d4). Both reached releases: flai v1.63.1 and flaiover v0.52.5. The host runs flai 1.67.9. S-0448 added the regression tests: `TestComponentsAnswerThroughTheChannelWithTheProjectKeyInTheirParams` in `flai/internal/hostapi/components_test.go`, `TestNoMethodTakesAParameterNamedProject` in `flai/internal/hostapi/params_test.go`, and the S-0448 cases in `flaiover/src/lib/server/agent.test.ts` and `flaiover/src/routes/api/components/components.test.ts`. I-0163 stayed open because S-0448 never closed it, and flai's issue step made this story from it on 2026-10-21. What is left is to confirm the fix and close the issue.

Tasks, in two layers:

1. T-1987 confirms the fix: the issue's reproduction and the regression tests. It touches nothing.
2. T-1990 closes I-0163, after T-1987.

Topics `cli` and `dashboard` were added because T-1987 runs the tests of `flai/internal/hostapi` and of flaiover.

Touches:

| Path | Source | Why |
|------|--------|-----|
| `design/issues/I-0163-…-project.md` | design | `flai issue close I-0163` writes the issue's status and reason |
| `design/issues/summary.md` | layout | closing an issue regenerates the summary of open issues |

The story declared no touches. `flai touches suggest` on the issue file found no co-change. No folder touch is kept. No code path is claimed: the regression tests exist, so the plan expects no code change. If T-1987 finds the cause still occurs, its agent widens the touches to the files it fixes.

Figures:

- Forecast 15m, raised from flai's 5m. flai's figure is the median of 8 remediation stories at size 2. This story also runs the S-0448 regression tests in Go and in flaiover, which the size does not count. Delivery 2026-10-22T04:49:00Z is flai's, from the pull order.
- Cost of delay value 50 USD a week, lowered from flai's 1050. The 1000 USD penalty a week was for a Components page that fails on every call, and that stopped once flaiover v0.52.5 and flai v1.63.1 shipped S-0448's fix. What delay still costs is the time lost the inputs name, 20m a cycle at 150 USD an hour, while the open issue keeps being read and reported. The inputs are left as flai set them.

### Accepted by the orchestrator

- Verified: 3258d4bde522a848ccf2c0f002fc327f1f718bfe
- Batch: 20261021-121345
- Enqueued verified: 3258d4bde522a848ccf2c0f002fc327f1f718bfe
- At: 2026-10-21T12:13:56Z

Verdict: PASS. Both criteria are met, and no convention is broken. The fix is S-0448's, already on main with its regression tests.
- 1: design/issues/I-0163-the-components-page-cannot-list-start-or-stop-anything-through-the-dashboard-because-the-channel-s-project-key-and-the-components-methods-sub-project-share-the-one-params-field-project.md
- 2: design/issues/I-0163-the-components-page-cannot-list-start-or-stop-anything-through-the-dashboard-because-the-channel-s-project-key-and-the-components-methods-sub-project-share-the-one-params-field-project.md, design/issues/summary.md
