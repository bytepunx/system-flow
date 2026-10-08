---
id: S-0320
type: story
nature: remediation
title: flai serve's test that the orchestrator is started again when it ends fails under the close-out's full integration run
status: review
owner: alex
created: 2026-10-07T18:59:51Z
updated: 2026-10-08T05:51:50Z
transitions:
  - to: ready
    at: 2026-10-08T00:18:40Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-08T04:42:15Z
    by: agent-S-0320
  - to: review
    at: 2026-10-08T05:51:50Z
    by: agent-S-0320
tags: [flai, serve, tests]
topics: [continuous-improvement]
touches: [flai/internal/serve/orchestrate_test.go, design/issues/I-0106-flai-serve-s-test-that-the-orchestrator-is-started-again-when-it-ends-fails-under-the-close-out-s-full-integration-run.md, design/issues/summary.md, design/issues/I-0118-flai-check-finds-markdown-md034-outside-the-story-at-close-out.md, design/issues/I-0119-the-close-out-s-last-check-that-the-branch-contains-main-fails-when-flai-commits-wip-on-main-during-its-run.md, design/issues/I-0113-a-failed-integration-tier-in-the-close-out-shows-only-the-last-lines-of-go-test-so-the-failing-test-is-not-named.md, design/issues/I-0117-a-story-s-integration-tier-lints-main-s-committed-wip-so-markdown-any-agent-commits-there-fails-every-story-s-close-out-until-main-commits-a-fix.md, design/issues/I-0086-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2435
  turns:
    - day: 2026-10-08
      ceremony: 4
      hand_edits: 4
      work: 71
  models:
    - model: claude-opus-5-5
      input: 168
      output: 32421
      cache_read: 13499524
      cache_write: 205976
      cost: 4.9968
  strategic:
    - kind: orchestrator
      seconds: 2295
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 105
          output: 1610
          cache_read: 16706962
          cache_write: 64899
          cost: 4.1327
        - model: claude-sonnet-5-5
          input: 6
          output: 44
          cache_read: 39191
          cache_write: 36122
          cost: 0.0649
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-07T18:59:51Z
  value: 25
  by: planner-S-0320
  at: 2026-10-08T00:18:16Z
forecast:
  duration: 25m
  delivery: 2026-10-08T05:05:00Z
  basis: "Its own forecast of 25m; 1st in the pull order with an in-progress limit of 3, behind S-0232, S-0316 and S-0318."
  by: flai
  at: 2026-10-08T04:38:47Z
finalized:
  by: orchestrator
  at: 2026-10-08T00:18:36Z
---
# S-0320 flai serve's test that the orchestrator is started again when it ends fails under the close-out's full integration run

## Goal

This story remediates [I-0106](../../../design/issues/I-0106-flai-serve-s-test-that-the-orchestrator-is-started-again-when-it-ends-fails-under-the-close-out-s-full-integration-run.md), "flai serve's test that the orchestrator is started again when it ends fails under the close-out's full integration run". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0106 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0106 is closed with `flai issue close I-0106 --reason` saying what fixed it

## Tasks
- T-1307 TestTheOrchestratorIsStartedAgainWhenItEnds sets the lab's clock from the failed run's recorded end, not from the host's clock plus a margin
- T-1308 Close I-0106 saying what fixed it

## Notes

Cost of delay inputs set by flai from I-0106. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-07T07:31:02Z, 0.5 days before this story; under one cycle counts as one).

### Planning

The proposed solution, from I-0106's one instance and the code. The test failed at `started again within a minute of a failure` (`flai/internal/serve/orchestrate_test.go`). Then its clock was the host's plus `orchestrateRetry - 1s`. `AgentRun.Ended` is written to the second (`time.RFC3339`, `flai/internal/serve/agents.go`), and `orchestrator.due` adds `orchestrateRetry` to it. So the truncated fraction of a second, plus the host's time before the look, made the failed run due again early. That explains the failure run alone too. Commit 7772bd34 (2026-10-07T07:56Z, after the instance at 07:31Z) widened the margin to twenty seconds for I-0102. That makes the failure rarer, not impossible. T-1307 removes the dependence on the host: the two looks around the retry take their time from the failed run's recorded `Ended`. T-1308 then closes the issue. No production code changes: `due` is right, and a retry up to a second early does no harm.

Touches, file by file, none a folder:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/serve/orchestrate_test.go` | layout | The failing test and its lab, `orchestrateLab`, are here. T-1307 changes both. |
| `design/issues/I-0106-…md` | design | Criterion 2 closes it. |
| `design/issues/summary.md` | design | `flai issue close` regenerates it. |

The story declared no touches. `flai touches suggest` from these three lists issue files and design docs that change with `summary.md`; this story changes none of them. `design/issues` is in the manifest's `claims.shared`, so the two issue paths hold no story. No open story touches `flai/internal/serve/orchestrate_test.go`.

Figures:

- **Forecast: 25m, up from flai's 14m.** flai sized it at 160 s per unit times size 5 (2 criteria, 3 touches). S-0310, a like remediation, took 28m in progress: one timing-dependent test, its issue closed, and a close-out that runs the full tiers. This story is the same shape. The delivery is flai's 2026-10-08T07:07Z plus the 11m added.
- **Cost of delay: 25.00 USD a week, as `flai cod` gives it.** That is 10m lost per 168h cycle at 150 USD an hour, from flai's input. It stands. A failed close-out costs about that much each time, and the twenty-second margin already makes one rarer.
