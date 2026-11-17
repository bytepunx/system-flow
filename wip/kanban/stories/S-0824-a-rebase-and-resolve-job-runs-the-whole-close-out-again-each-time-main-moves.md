---
id: S-0824
type: story
nature: improvement
title: A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted
status: backlog
owner: alex
created: 2026-11-17T07:13:20Z
updated: 2026-11-17T12:57:33Z
transitions: []
tags: [flai, flaiover]
topics: [verify, acceptance, serve, dashboard, template]
touches: [design/adrs, design/adrs/README.md, flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go, flai/internal/closeout/closeout.go, flai/internal/closeout/closeout_test.go, flai/internal/closeout/script_test.go, flai/internal/serve/resolvejob.go, flai/internal/serve/resolvejob_test.go, flai/internal/serve/verifyjob.go, flai/internal/serve/verifyjob_container.go, flai/internal/serve/verifyjob_test.go, flai/internal/hostapi/resolve_status.go, flai/internal/hostapi/resolve_status_test.go, scripts/close-out.sh, template/root/scripts/close-out.sh, template/CHANGELOG.md, flaiover/src/lib/review.ts, flaiover/src/lib/review.test.ts, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flai-reference.md, docs/users/flaiover.md, design/issues/I-0398-a-rebase-and-resolve-job-runs-the-whole-close-out-again-each-time-main-moves.md, design/issues/summary.md]
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
      seconds: 102
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 71
          output: 9474
          cache_read: 5771157
          cache_write: 46405
          cost: 1.8896
    - kind: orchestrator
      seconds: 43
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 10
          output: 152
          cache_read: 467966
          cache_write: 101252
          cost: 0.1768
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1h
    by: flai
    at: 2026-11-17T07:13:20Z
  value: 150
  by: planner-S-0824
  at: 2026-11-17T08:32:55Z
forecast:
  duration: 2h30m
  delivery: 2026-11-17T22:34:00Z
  basis: "Its own forecast of 2h30m; 14th in the pull order with an in-progress limit of 3, behind S-0806, S-0829, S-0834, S-0811, S-0840, S-0809, S-0842, S-0812, S-0808, S-0807, S-0810, S-0241, S-0813, S-0814, S-0815 and S-0822."
  by: flai
  at: 2026-11-17T12:57:33Z
---
# S-0824 A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted

## Goal

This story remediates [I-0398](../../../design/issues/I-0398-a-rebase-and-resolve-job-runs-the-whole-close-out-again-each-time-main-moves.md), "A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted". The issue recommends this solution:

Decide, in an ADR refining ADR-0313 and ADR-0156, what a rebase and resolve job owes when main moves under it, and build it. Options to weigh: the job passes at the base its last round ran on and leaves the rest to the acceptance's own `verify_behind_main` gate (ADR-0291, ADR-0303), which sends back only the tiers the new paths overlap; the later rounds apply ADR-0303's overlap rule in place of every tier main's paths select; or flai serve holds acceptances and resolve jobs in one order, so that main does not move under a job it is about to accept. ADR-0312's acceptance suites (S-0808) bear on the same cost. The job's record and the review page say which round the job is on and why it started another.

## Acceptance criteria
- [ ] The cause I-0398 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0398 is closed with `flai issue close I-0398 --reason` saying what fixed it

## Tasks
- T-3551 An ADR decides what a rebase and resolve job's close-out owes when main moves under it, and how its record names the round
- T-3552 flai verify --sync-only, when asked, passes a branch main outran while verify_behind_main would raise no blocker at the head it verified
- T-3553 closeout.Status and the job's record say which round a close-out is on and why it started another, and resolve.status answers them
- T-3554 A rebase and resolve job's close-out passes over a main that moved outside the story's change set, in both scripts and the container close-out
- T-3555 The review page shows which round a close-out or Rebase and resolve job is on and why it started each round after the first
- T-3556 The design and the user docs describe what a rebase and resolve job owes a moved main and how it names its rounds, and I-0398 is closed

## Notes

Cost of delay inputs set by flai from I-0398. time_lost_per_cycle 1h carried over from I-0398's Impact section. Its Impact time_lost_per_cycle was taken rather than the 22m derived from its cost and count.

### Planning

The plan assumes the first option, and the plan's thread asks the operator to confirm it. In a rebase and resolve job, the close-out's last check passes when main moved but `verify_behind_main` would raise no blocker for the head it verified. That holds when no path main brought in since the round's merge base overlaps the story's change set by ADR-0303's rule, or when the overlapping paths send back no tier. The acceptance's own gate then judges what main brings in later, as it does today. Another round runs only when the gate would block, and its round line says why. ADR-0156 already rejected the third option, holding acceptances, because it delays the operator's acceptances. The second option is left to T-3551's ADR to weigh, since a later round that overlaps `flai/**` selects the integration tier either way.

`flai touches suggest` listed nothing, since the story declared no touches. Where each touch came from:

- `design/adrs/` (layout) stays a folder: the new ADR's file name comes from its title, which the story's agent settles in T-3551. `design/adrs/README.md` (layout) indexes it.
- `flai/internal/preview/accept.go` and its test (design: ADR-0291, ADR-0303). `verifyBehindMain`, `mainPathsSinceVerify`, and `overlappingMainPaths` are the gate the last check reuses.
- `flai/internal/verify/story.go`, `flai/cmd/verify.go`, and their tests (design: ADR-0156; layout). These hold `verify.SyncOnly` and `flai verify --sync-only`, the close-out's last check.
- `flai/internal/closeout/closeout.go`, `closeout_test.go`, and `script_test.go` (design: ADR-0313 §5, §7; layout). These hold the job's record, `closeout.Status`, which reads the round lines, and the scripts' tests.
- `flai/internal/serve/resolvejob.go`, `resolvejob_test.go`, `verifyjob.go`, `verifyjob_container.go`, and `verifyjob_test.go` (design: ADR-0313, ADR-0274; layout). These hold the resolve job's environment, the container close-out's own three-round loop, and the test that reproduces I-0398.
- `flai/internal/hostapi/resolve_status.go` and its test (layout): `resolve.status`, which answers the round to the dashboard.
- `scripts/close-out.sh`, `template/root/scripts/close-out.sh`, and `template/CHANGELOG.md` (design: ADR-0156 §1 says both scripts loop; layout).
- `flaiover/src/lib/review.ts`, `Review.svelte`, and their tests (goal: the review page says which round the job is on and why; layout).
- Docs (design): `design/system/flai-cli.md` § The close-out job, `design/system/flaiover-dashboard.md`, `docs/users/flai.md` § The close-out job and § Resuming across a clean rebase, `docs/users/flai-reference.md` § flai verify, and `docs/users/flaiover.md` § Reviewing a story.
- The I-0398 file and `design/issues/summary.md` (criterion 2).

No `after`:

- S-0808 (ADR-0312's acceptance suites) changes what `verify_behind_main` judges. This story reads the gate as it is on main when the story starts, not an output S-0808 writes.
- S-0822 and S-0823 also change `resolvejob.go`, `closeout.go`, and the review page. A trial merge judges that overlap.

Figures:

- Forecast: `flai forecast` gives 1h2m, from 108 s per unit of size over 186 done large-band improvement stories, times a size of 34 (2 criteria, 32 touches). I raised the duration to 2h30m. Reason: S-0804, which built the resolve job across the same packages, the scripts, the review page, and an ADR, spent 2h17m in progress, and this story has six tasks in four layers. flai's replan keeps that duration and sets the delivery from the pull order.
- Cost of delay: `flai cod` gives 150 USD a week, from 1h lost per 168h cycle at 150 USD an hour. The value stands: the 1h per cycle is the operator's input from the Impact section, and the value follows from it unchanged.
