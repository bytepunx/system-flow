---
id: T-1307
type: task
nature: remediation
title: TestTheOrchestratorIsStartedAgainWhenItEnds sets the lab's clock from the failed run's recorded end, not from the host's clock plus a margin
status: done
parent: S-0320
owner: alex
created: 2026-10-08T00:17:55Z
updated: 2026-10-08T04:44:13Z
transitions:
  - to: ready
    at: 2026-10-08T04:43:56Z
    by: agent-S-0320
  - to: in-progress
    at: 2026-10-08T04:43:57Z
    by: agent-S-0320
  - to: done
    at: 2026-10-08T04:44:13Z
    by: agent-S-0320
stream: S-0320
tags: [flai, serve, tests]
touches: [flai/internal/serve/orchestrate_test.go]
usage:
  source: log
  seconds: 16
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 3
      output: 582
      cache_read: 242154
      cache_write: 3695
      cost: 0.0896
---
# T-1307 TestTheOrchestratorIsStartedAgainWhenItEnds sets the lab's clock from the failed run's recorded end, not from the host's clock plus a margin

## Work

I-0106's one instance (S-0212, 2026-10-07T07:31Z) failed at `started again within a minute of a failure`. That was before commit 7772bd34 (07:56Z) widened the test's margin from one second to twenty (I-0102). The cause is the margin itself. `AgentRun.Ended` is written to the second (`time.RFC3339` in `flai/internal/serve/agents.go`), and `orchestrator.due` adds `orchestrateRetry` to it. The lab's clock is the host's plus `lab.shift`. So with a shift of `orchestrateRetry - 1s`, the run was due again whenever the truncated fraction of a second plus the host's own time before the look reached one second. That happens often enough to fail a run alone. Twenty seconds makes it rarer; it does not make the test independent of the host.

- Give `orchestrateLab` a way to pin its clock to an instant, for example `lab.at(t time.Time)` beside `shift`, keeping `shift` for the tests that use it.
- In `TestTheOrchestratorIsStartedAgainWhenItEnds`, parse the failed run's `Ended` and set the clock to `Ended + orchestrateRetry - time.Second`: the run is not started again. Then set it to `Ended + orchestrateRetry`: it is. Drop the host-relative shift and the I-0102 comment from those two looks.
- Check the other `lab.shift` uses in `flai/internal/serve/orchestrate_test.go` (`2*`, `4*`, and `6*orchestrateRetry`). Each leaves at least a minute for the host, so they stay unless the reading finds otherwise; say so in the commit.
- Waits for nothing: it is the first layer.

## Done when

- The test's two looks around the retry take their time from the failed run's recorded `Ended`, so neither depends on how long the host took.
- The boundary is exact: the look one second before `Ended + orchestrateRetry` starts nothing, and the look at it starts the third run.
- `go test -race -count=20 -run TestTheOrchestratorIsStartedAgainWhenItEnds ./internal/serve/` passes in `flai/`, and `flai test flai/internal/serve/orchestrate_test.go` passes.

## Notes
