---
id: S-0808
type: story
nature: feature
title: The acceptance's gates, the merge queue's eligibility, queue_tiers, and verify_behind_main, judge only the suites the accept transition's verify lists
status: ready
parent: E-0034
owner: alex
created: 2026-11-16T23:42:04Z
updated: 2026-11-17T12:57:33Z
transitions:
  - to: ready
    at: 2026-11-17T07:47:04Z
    by: alex
tags: [flai]
topics: [verify, workflow]
touches: [flai/internal/workitem/transition.go, flai/internal/workitem/transition_test.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, flai/cmd/accept.go, flai/cmd/accept_suites_test.go, flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go, docs/users/flai-reference.md, docs/users/flai.md, docs/operators/index.md, design/system/workflow.md, design/system/flai-cli.md]
after: [S-0807]
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
      seconds: 248
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 438
          output: 1494
          cache_read: 20253284
          cache_write: 1139225
          cost: 1.4559
        - model: claude-opus-5-5
          input: 428
          output: 24997
          cache_read: 32455725
          cache_write: 1723153
          cost: 9.9993
    - kind: orchestrator
      seconds: 183
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 56
          output: 343
          cache_read: 3054624
          cache_write: 268127
          cost: 1.018
cost_of_delay:
  value: 185.64
  by: planner-S-0808
  at: 2026-11-16T23:47:49Z
forecast:
  duration: 1h15m
  delivery: 2026-11-17T16:26:00Z
  basis: "Its own forecast of 1h15m; 6th in the pull order with an in-progress limit of 3, behind S-0806, S-0829, S-0834, S-0811, S-0840, S-0809, S-0842 and S-0812."
  by: flai
  at: 2026-11-17T12:57:33Z
finalized:
  by: orchestrator
  at: 2026-11-16T23:48:11Z
---
# S-0808 The acceptance's gates, the merge queue's eligibility, queue_tiers, and verify_behind_main, judge only the suites the accept transition's verify lists

## Goal

[ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md) §3. The operator's note of 2026-11-16: every story in review has to be rebased and verified again even when the change to main cannot break it. With the gates reading the file, a project whose acceptance needs fewer suites is held for fewer.

## Acceptance criteria

- [ ] `verify_behind_main` sends back only tiers of the suites the transition with the `accept` effect lists in `verify:`, and lists no blocker when it lists none of the tiers main's paths select.
- [ ] `queue_tiers` and the merge queue's eligibility read the same suites.
- [ ] A story whose last verify did not run a suite the acceptance needs is refused with words naming the suite and the close-out that runs it.
- [ ] On the baseline every gate judges what it judged before this story, shown by the tests of ADR-0291, ADR-0303, and ADR-0221 passing unchanged.

## Tasks
- T-3433 Read the suites the acceptance needs from the workflow, and judge a verify report against them
- T-3434 queue_tiers and verify_behind_main judge only the acceptance's suites, and an acceptance waits for a suite its verify did not run
- T-3435 The merge queue's eligibility needs only the acceptance's suites to have passed at the story's head
- T-3436 The workflow design and the guides say the acceptance's gates judge the suites its transition's verify lists

## Notes

- After S-0807: the gates read the suites the verify report records, which S-0807 adds.

### Planning

- `after`: S-0807. This story reads the suites S-0807's verify report records for each suite's tiers. Through S-0807 it also reads the transition's `verify:` key S-0806 parses and the tier's `suite:` S-0805 adds. No other `after` is named: the pull order sets the rest.
- Touches, file by file, no folder:
  - layout and design: `flai/internal/workitem/transition.go`, where the accept transition's helpers live, and its new test `transition_test.go`, for `AcceptSuites`.
  - layout and design: `flai/internal/verify/story.go`, which holds `StaleTiers` and `QueueTiersNotRun`, and its test.
  - layout and design: `flai/internal/preview/accept.go`, which holds the `queue_tiers` and `verify_behind_main` gates, and its test.
  - layout and design: `flai/internal/mergequeue/eligible.go`, which holds the queue's eligibility, `gitState.verified`, and its test.
  - co-change: `flai/cmd/accept.go` (31% of the commits that changed the seed files) for its help text, and the new `flai/cmd/accept_suites_test.go`, so that the ADR-0291, ADR-0303, and ADR-0221 tests stay unchanged.
  - co-change: `docs/users/flai-reference.md` (19%), regenerated from the help.
  - design: `design/system/workflow.md`, `design/system/flai-cli.md`, `docs/users/flai.md`, and `docs/operators/index.md`, the documents that describe `queue_tiers` and `verify_behind_main` today.
  - `flai touches suggest` also listed `flai/internal/preview/orchestrator.go`, `flai/cmd/verify.go`, and `flai/internal/mergequeue/batch.go`. They are left out: the gates read `repo.Flow()` themselves, so no option or batch code changes.
- Forecast 1h15m, raised from flai's 35m. flai's figure is 110 s per unit of size over 142 large-band feature stories, times size 19. The one-gate stories S-0683 and S-0792 took 55m and 48m of agent time. This story changes three gates, adds a refusal, and updates four documents. Delivery 2026-11-17T06:30Z, after the serial chain S-0805, S-0806, and S-0807 it waits on is accepted. flai's 03:05Z does not count that chain.
- Cost of delay 185.64 USD a week, as `flai cod` gives it after the forecast. It is S-0808's share of E-0034's 1500 USD a week, the operator's input, 1h15m of the 10h6m forecast over the epic's 20 open stories without inputs. The story has no inputs of its own, and none were written. The figure moves as the other planners revise their stories' forecasts.
