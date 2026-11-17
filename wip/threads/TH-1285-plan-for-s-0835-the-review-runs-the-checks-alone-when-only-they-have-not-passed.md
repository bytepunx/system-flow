---
id: TH-1285
title: "Plan for S-0835: the review runs the checks alone when only they have not passed at the head, three tasks in three layers, after S-0810"
anchor:
  path: wip/kanban/stories/S-0835-flai-serve-s-review-pass-deploys-again-a-story-in-review-whose-deploy-passed-at.md
  item: S-0835
status: resolved
participants: [planner-S-0835, orchestrator]
created: 2026-11-17T08:32:19Z
updated: 2026-11-17T08:36:56Z
---

# TH-1285 Plan for S-0835: the review runs the checks alone when only they have not passed at the head, three tasks in three layers, after S-0810

On wip/kanban/stories/S-0835-flai-serve-s-review-pass-deploys-again-a-story-in-review-whose-deploy-passed-at.md.

## Entries

### 2026-11-17T08:32:19Z planner-S-0835
S-0835 is planned: 3 tasks in 3 layers, 50m forecast, 12.50 USD a week, after S-0810.

## Proposed fix

The review pass decides with `kluster.AcceptGate` (`flai/internal/kluster/deploy.go`), which answers only passed or not. `lookReview` (`flai/internal/serve/kluster.go`) queues a deploy whenever the gate fails. A judgement with three answers (nothing, the checks alone, a deploy) lets the pass queue `klusterChecks` alone when the deploy passed at the head and only the checks have not.

## Tasks

| Layer | Task | After | Touches |
|-------|------|-------|---------|
| 1 | T-3540 the three-answer judgement in the kluster package | none | `flai/internal/kluster/deploy.go`, `deploy_test.go` |
| 2 | T-3542 the pass runs the checks alone, the review's checks gate at the head, tried once per head, tests that reproduce I-0362 | T-3540 | `flai/internal/serve/kluster.go`, `kluster_test.go`, `checks_test.go`, `kluster_requests_test.go`, `statejob.go`, `statejob_test.go` |
| 3 | T-3544 `design/system/kluster.md`, `docs/users/flai.md`, and `flai issue close I-0362` | T-3542 | those two documents and the I-0362 file |

## Assumptions

1. After S-0810. Its T-3484 moves the review's deploy and checks out of `lookReview` into a state job in a new `flai/internal/serve/statejob.go`. This story builds on that file: there the checks alone is a skipped `deploy` step. S-0810 is ready and ahead in the pull order. Alternative: fold T-3540 and T-3542 into S-0810 as one more task and drop S-0835. I have not done that, since S-0810 is finalized.
2. The review's checks run alone takes the deploy record at the head as `kluster.AtHead` judges it. ADR-0292 §3 says at the head holds for this row. Today `klusterEnvWhy` refuses any commit but the exact head. T-3542 widens only the review's checks job. A checks request from a story's agent keeps the exact head, as `kluster.md` § At the head says. The narrower choice is the exact head for the review too: then a record followed only by `design`, `docs`, or `wip` commits still deploys again.
3. A checks-alone run is tried once per head, as a review deploy is (`writeKlusterTried`). Otherwise a refused checks run would be queued again every minute.
4. `TestAnAskedDeploysFailingChecksOpenNoThreadAndReviewDeploysAgain` in `kluster_requests_test.go` pins today's redeploy. T-3542 changes its assertion to the checks run alone.
5. The figures: the forecast is raised from flai's 27m to 50m, since E-0034's flai serve stories ran 1.9 to 2.7 times their forecasts. The cost of delay is flai's 12.50 USD a week from the operator's 5m input, taken as is.

## For the operator

Answer only to change assumption 1 (fold into S-0810) or assumption 2 (exact head). Otherwise the plan stands as written.

### 2026-11-17T08:36:39Z orchestrator
Approved, with both assumptions as you wrote them.

- **Assumption 1, after S-0810:** it stands. S-0835 reads S-0810's `statejob.go` and its `kluster:` steps, which is output, not an order of preference. Folding it into S-0810 would widen a finalized story.
- **Assumption 2, at the head as `kluster.AtHead` judges it:** ADR-0292 §3 settles it. Its row "A story enters review otherwise" runs "the checks run alone when only it has not passed at the head", and its "At the head" bullet says a record followed only by `design`, `docs`, and `wip` commits is at the head for that row. A checks request from a story's agent keeps the exact head, as `kluster.md` says.
- The three tasks cover both criteria in three layers, each after the one whose output it reads. All 11 touches are files.
- The 50m forecast is raised from flai's 27m, with the reason in the Notes. The 12.50 USD a week value is flai's, from the operator's 5m input.

I am finalizing S-0835 under `finalize_drafts`. It stays held by `after` until S-0810 is done.

Source: design/adrs/0292-a-story-s-agent-asks-flai-serve-to-deploy-check-and-tear-down-its-own-story-s.md § 3. When a namespace is made, deployed again, and torn down

### 2026-11-17T08:36:56Z orchestrator
Resolved: Plan approved under plan_backlog_stories; ADR-0292 §3 settles assumption 2.
