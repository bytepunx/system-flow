---
id: S-0526
type: story
nature: remediation
title: component-library.md still says the main checkout is never installed into, though flai accept now installs there by ADR-0185
status: done
owner: alex
created: 2026-10-28T12:54:04Z
updated: 2026-10-28T13:08:13Z
transitions:
  - to: ready
    at: 2026-10-28T13:03:22Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-28T13:05:49Z
    by: agent-S-0526
  - to: review
    at: 2026-10-28T13:07:36Z
    by: agent-S-0526
  - to: done
    at: 2026-10-28T13:08:13Z
    by: alex
tags: [design, components, issues]
topics: [components]
touches: [design/system/component-library.md, design/issues/I-0192-component-library-md-still-says-the-main-checkout-is-never-installed-into-though-flai-accept-now-installs-there-by-adr-0185.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 116
  priced_by: harness
  turns:
    - day: 2026-10-28
      ceremony: 1
      hand_edits: 1
      work: 17
  models:
    - model: claude-opus-5-5
      input: 44
      output: 6344
      cache_read: 2408164
      cache_write: 146538
      cost: 1.781
  strategic:
    - kind: orchestrator
      seconds: 135
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 48
          output: 752
          cache_read: 10647859
          cache_write: 25188
          cost: 2.3432
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: orchestrator
    at: 2026-10-28T13:01:58Z
  value: 12.5
  by: planner-S-0526
  at: 2026-10-28T13:03:05Z
forecast:
  duration: 20m
  delivery: 2026-10-28T13:31:00Z
  basis: "Its own forecast of 20m; 1st in the pull order with an in-progress limit of 5, behind S-0502."
  by: flai
  at: 2026-10-28T13:05:47Z
finalized:
  by: alex
  at: 2026-10-28T12:58:19Z
---
# S-0526 component-library.md still says the main checkout is never installed into, though flai accept now installs there by ADR-0185

## Goal

This story remediates [I-0192](../../../design/issues/I-0192-component-library-md-still-says-the-main-checkout-is-never-installed-into-though-flai-accept-now-installs-there-by-adr-0185.md), "component-library.md still says the main checkout is never installed into, though flai accept now installs there by ADR-0185". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0192 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0192 is closed with `flai issue close I-0192 --reason` saying what fixed it

## Tasks
- T-2204 The install step in component-library.md names flai accept's install as the one exception to the main checkout never being installed into, with I-0179 as the reason and a link to ADR-0185
- T-2205 I-0192 is closed with the reason that the install step now names flai accept's install as the exception

## Notes

### Planning

Planned by planner-S-0526 on 2026-10-28. The plan and its assumptions are on TH-0682, which the orchestrator approved.

The solution proposed from I-0192's instance: the last sentence of **The install step** in `design/system/component-library.md` § A story's worktree (S-0442) › Sources is rewritten to name `flai accept`'s install, by the manifest's `install` command when the accepted change touched its `install.paths`, as the one exception to the main checkout never being installed into, with I-0179 as the reason and a link to [ADR-0185](../../../design/adrs/0185-flai-accept-runs-a-sub-project-s-install-command-in-the-main-checkout-when-the.md), and to say that install has no time limit (`runInstalls` in `flai/cmd/accept_install.go` runs under `context.Background()`). No code changes. No test fits: the defect is one sentence of the living design, checked for form, not truth. `docs/users/flai.md` § Accept and release and `docs/operators/settings.md` already describe the install correctly.

Touches, and where each came from:

- `design/system/component-library.md`: design, the document I-0192 names; the starting path given to `flai touches suggest`, since the story declared no touches.
- `design/issues/I-0192-...md`, `design/issues/summary.md`: declared by the second criterion; `flai issue close` writes both. `summary.md` was also 10% co-change of the design document.
- Left out of the suggestion: `docs/users/flai.md` (43%), `docs/users/flaiover.md`, `design/adrs/README.md`, `design/system/flaiover-dashboard.md`, `design/system/flai-cli.md`, `design/system/dashboard-host-channel.md`, `docs/users/flai-reference.md`, `docs/operators/index.md`, and the `flai/internal/components` files: co-change of past component-library stories, but nothing in them changes here; the docs already say what the design should.

No folder touch is kept.

Forecast: `flai forecast` gave 11m (median 131 s per unit of size over 12 done small remediation stories on claude-opus-5-5, size 5: 2 criteria, 3 touches). Raised to 20m: a story this small is mostly fixed overhead, the start, the prime, two task closes with their syncs, and the close-out's tiers, which the per-unit median spreads over larger stories. Delivery is flai's 2026-10-29T09:04Z moved by the added 9m, 40th in the pull order with an in-progress limit of 5.

Cost of delay: `flai cod` gave 12.50 USD a week from the input the orchestrator set on TH-0682 at the planner's recommendation, `time_lost_per_cycle: 5m` (5m per 168h cycle at 150 USD an hour, 1.00 cycles a week). The value stands as computed; the input is not changed.
