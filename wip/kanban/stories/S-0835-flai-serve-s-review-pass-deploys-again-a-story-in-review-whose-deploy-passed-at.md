---
id: S-0835
type: story
nature: remediation
title: "flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone"
status: ready
owner: alex
created: 2026-11-17T08:17:30Z
updated: 2026-11-28T00:03:04Z
transitions:
  - to: ready
    at: 2026-11-27T23:54:12Z
    by: claude-concierge
tags: [flai]
topics: [kluster, deployment]
touches: [flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go, flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go, flai/internal/serve/checks_test.go, flai/internal/serve/kluster_requests_test.go, flai/internal/serve/statejob.go, flai/internal/serve/statejob_test.go, design/system/kluster.md, docs/users/flai.md, design/issues/I-0362-flai-serve-s-review-pass-deploys-again-a-story-in-review-whose-deploy-passed-at.md]
after: [S-0810]
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
      seconds: 58
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 2
          output: 3
          cache_read: 196833
          cache_write: 2940
          cost: 0.0136
        - model: claude-opus-5-5
          input: 54
          output: 3457
          cache_read: 4620195
          cache_write: 60350
          cost: 1.4187
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-17T08:17:30Z
  value: 12.5
  by: planner-S-0835
  at: 2026-11-17T08:32:07Z
forecast:
  duration: 50m
  delivery: 2026-11-28T01:08:00Z
  basis: "Its own forecast of 50m; 1st in the pull order with an in-progress limit of 3, behind S-0837, S-0844 and S-0846."
  by: flai
  at: 2026-11-28T00:03:04Z
finalized:
  by: orchestrator
  at: 2026-11-17T08:36:57Z
---
# S-0835 flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone

## Goal

This story remediates [I-0362](../../../design/issues/I-0362-flai-serve-s-review-pass-deploys-again-a-story-in-review-whose-deploy-passed-at.md), "flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0362 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0362 is closed with `flai issue close I-0362 --reason` saying what fixed it

## Tasks
- T-3540 The kluster package says what a story entering review needs from its namespace's record: nothing, the checks alone, or a deploy
- T-3542 flai serve runs the checks alone against a story in review whose deploy passed at its head and whose checks have not, and tries each head's checks once
- T-3544 kluster.md and the user guide say the review runs the checks alone when only they have not passed at the head, and I-0362 is closed

## Notes

Cost of delay inputs set by flai from I-0362. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-13T13:28:22Z, 3.8 days before this story; under one cycle counts as one).

### Planning

Written by planner-S-0835 on 2026-11-17.

The proposed fix, from I-0362's one instance (S-0661) and an explorer's read of the code on main: the review pass decides with `kluster.AcceptGate`, which answers only passed or not, and queues a `klusterDeploy` whenever it fails. A judgement with three answers (nothing, the checks alone, a deploy) lets the pass queue `klusterChecks` alone when the deploy passed at the head and only the checks have not. The checks job the review queues then takes the record at the head as `kluster.AtHead` judges it, as ADR-0292 §3 says for this row, while an agent's checks request keeps the exact head. Each head's checks-alone run is tried once, as a review deploy is.

After:

- S-0810: its T-3484 removes the review's deploy from `lookReview` and runs the cluster's `deploy` and `checks` as steps of a state job in a new `flai/internal/serve/statejob.go`. This story reads that file and its `kluster:` steps: the checks alone is a skipped `deploy` step there. S-0810 is ready and ahead in the pull order; were this story built first, S-0810 would have to carry the fix into code it rewrites.

Touches, all files, no folder. The story declared none, so `flai touches suggest` was run from `flai/internal/serve/kluster.go` and `flai/internal/kluster/deploy.go`, the two files the explorer found the decision in.

- Layout, the decision: `flai/internal/kluster/deploy.go` (`AcceptGate`) and `deploy_test.go` (`TestAcceptGate`), co-change 50%.
- Layout, the pass: `flai/internal/serve/kluster.go` (`lookReview`, `klusterEnvWhy`, `writeKlusterTried`, `runChecks`); `kluster_test.go`, co-change 50%; `checks_test.go`, 25%; `kluster_requests_test.go`, 10%, whose `TestAnAskedDeploysFailingChecksOpenNoThreadAndReviewDeploysAgain` pins today's deploy.
- After S-0810: `flai/internal/serve/statejob.go` and `statejob_test.go`, which S-0810 writes.
- Design: `design/system/kluster.md`, whose § What happens when marks the checks run alone not built for the review; `docs/users/flai.md`, which says a story entering review is deployed again in place.
- Criteria: the I-0362 file, which `flai issue close` writes.
- Co-change left out: `flai/internal/serve/checks.go` (20%), `kluster_requests.go` (15%), `flai/cmd/checks.go`, `flai/internal/inbox/inbox.go`, and the kluster request files (10%). The fix changes no request, no command, and no event.

Forecast: 50m, delivery 2026-11-17T17:18:00Z. `flai forecast` gives 27m, the median 147 s a unit over 68 remediation stories times size 11 (2 criteria, 9 touches before the two of S-0810's). Raised, since E-0034's flai serve stories ran 1.9 to 2.7 times their forecasts (S-0743 2h02m on 45m, S-0755 1h30m on 47m, S-0741 1h19m on 29m, from S-0810's planning), and this one changes a gate and a test that pins the gate's behaviour. Delivery is flai's 16:55Z plus the 23m added.

Cost of delay: 12.50 USD a week, as `flai cod` gives it from the operator's input, 5m of time lost per 168h cycle at 150 USD an hour. Taken as is: the one instance cost 5m, and nothing in the issue says it costs more.

Task layers:

1. T-3540 (the three-answer judgement in the kluster package): no `after`.
2. T-3542 (the pass runs the checks alone, the review's checks gate at the head, tried once, and the tests that reproduce I-0362), after T-3540, whose judgement it reads.
3. T-3544 (`kluster.md`, the user guide, and the issue's close), after T-3542, whose behaviour it describes.
