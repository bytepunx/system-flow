---
id: S-0694
type: story
nature: improvement
title: "ADR-0254 §5a's send-back of a story whose kluster deploy or checks failed, and its kluster acceptance blocker, are not built"
status: done
owner: alex
created: 2026-11-10T08:54:25Z
updated: 2026-11-11T14:52:10Z
transitions:
  - to: ready
    at: 2026-11-11T01:45:35Z
    by: alex
  - to: in-progress
    at: 2026-11-11T08:56:31Z
    by: agent-S-0694
  - to: review
    at: 2026-11-11T14:50:29Z
    by: agent-S-0694
  - to: done
    at: 2026-11-11T14:52:10Z
    by: orchestrator
tags: [flai, serve, kluster, accept]
topics: [kluster, deployment, review]
touches: [flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go, flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go, flai/internal/serve/agents_restart.go, flai/internal/serve/agents_test.go, flai/internal/harness/harness.go, flai/internal/preview/accept.go, flai/internal/preview/orchestrator.go, flai/internal/preview/accept_test.go, flai/cmd/accept.go, flai/cmd/move.go, flai/cmd/accept_kluster_test.go, flai/internal/hostapi/reads.go, flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go, design/system/kluster.md, design/system/workflow.md, design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, design/issues/I-0268-adr-0254-5a-s-send-back-of-a-story-whose-kluster-deploy-or-checks-failed-and.md, design/issues/summary.md, flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go, flai/internal/serve/mergequeue.go, flai/cmd/queue.go, flai/internal/serve/checks_test.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4451
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-11
      ceremony: 3
      test_runs: 2
      work: 123
  models:
    - model: claude-opus-5-5
      input: 662
      output: 57451
      cache_read: 41565248
      cache_write: 1214483
      cost: 14.7058
  strategic:
    - kind: orchestrator
      seconds: 14
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 4
          output: 21
          cache_read: 201197
          cache_write: 23365
          cost: 0.0619
    - kind: analyzer
      seconds: 42
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 33
          output: 3103
          cache_read: 238317
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 15
          output: 2841
          cache_read: 759050
          cache_write: 13855
          cost: 0.278
cost_of_delay:
  inputs:
    time_lost_per_cycle: 20m
    by: flai
    at: 2026-11-10T08:54:25Z
  value: 50
  by: planner-S-0694
  at: 2026-11-10T10:32:23Z
forecast:
  duration: 2h
  delivery: 2026-11-11T11:42:00Z
  basis: "Its own forecast of 2h; 1st in the pull order with an in-progress limit of 3, behind S-0698 and S-0702."
  by: flai
  at: 2026-11-11T08:56:25Z
finalized:
  by: orchestrator
  at: 2026-11-10T11:10:10Z
---
# S-0694 ADR-0254 §5a's send-back of a story whose kluster deploy or checks failed, and its kluster acceptance blocker, are not built

## Goal

This story remediates [I-0268](../../../design/issues/I-0268-adr-0254-5a-s-send-back-of-a-story-whose-kluster-deploy-or-checks-failed-and.md), "ADR-0254 §5a's send-back of a story whose kluster deploy or checks failed, and its kluster acceptance blocker, are not built", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0268 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0268 is closed with `flai issue close I-0268 --reason` saying what fixed it

## Tasks
- T-3024 The kluster record of a story's namespace keeps its checks outcome at the deployed commit, with a gate acceptance can read
- T-3025 flai serve sends a story whose kluster deploy or checks failed back to in-progress, starts its agent, and keeps its namespace
- T-3028 flai accept refuses a story with the blocker kluster while the action is on, and --without-kluster lands it on the operator's word
- T-3029 The merge queue lists a story the kluster gate refuses as ineligible, with the reason kluster, instead of batching it
- T-3035 The kluster design, the workflow, the orchestrator's blockers, and the user guide say a failed deploy sends the story back and blocks its acceptance
- T-3038 Close I-0268, naming the send-back, the kluster blocker, and the tests that hold them

## Notes

The analyzer found I-0268: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0268. time_lost_per_cycle 20m: 10m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-11-07T03:22:56Z, 3.2 days before this story; under one cycle counts as one).

Strategic usage carried over from I-0268: analyzer 42s, 0.8156 USD.

### Planning

The solution is ADR-0254 §5a as written, which the operator settled on TH-0832 ([D-0481](../../../design/digest/D-0481-accept-adr-0254-how-flai-runs-a-project-s-kluster-for-a-story-in-review-1-yes.md)), and `design/system/kluster.md` § A failed deploy holds the story. Nothing in it waits on a new decision.

Touches. The story declared none, and `flai touches suggest` needs a path to start from, so every touch comes from the issue's instances, ADR-0254 §5a, and the code layout an explorer read:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/kluster/deploy.go`, `deploy_test.go` | layout | `kluster.Record` per namespace; gains the checks outcome at the deployed commit and a reader acceptance can call |
| `flai/internal/serve/kluster.go`, `kluster_test.go` | design (issue) | `deploy()` and `checks()` report a failure and move nothing; `lookLeft()` and `teardown()` must keep a sent-back story's namespace |
| `flai/internal/serve/agents_restart.go`, `agents_test.go` | layout | ADR-0167's and the acceptance send-back's agent start, which finds no kluster send-back today |
| `flai/internal/harness/harness.go` | layout | `harness.SentBack` names why an agent was started again |
| `flai/internal/preview/accept.go`, `orchestrator.go`, `accept_test.go` | design (issue) | the blockers beside `queue_tiers` and `AcceptOptions` |
| `flai/cmd/accept.go`, `move.go`, `accept_kluster_test.go` | design (issue) | `--without-kluster` beside `--without-queue-tiers`, shared with `flai move`; the test file is new |
| `flai/internal/hostapi/reads.go`, `flai/internal/mcpserver/queue.go` | layout | the dashboard's and the orchestrator's previews pass the host action |
| `flai/internal/mergequeue/eligible.go`, `eligible_test.go` | design (issue) | the queue's acceptance is refused too; a reason keeps it from batching a story it cannot land |
| `design/system/kluster.md`, `workflow.md`, `strategic-agents.md`, `flai-cli.md` | design | each names the blockers or says the send-back is not built |
| `docs/users/flai.md`, `docs/users/flai-reference.md`, `docs/operators/settings.md` | layout | the user guide, and the generated reference and flag table |
| `design/issues/I-0268-…`, `design/issues/summary.md` | declared by the criteria | the close |

No folder touch is kept.

Forecast. `flai forecast` gave 2m from size 2 (2 criteria, 0 touches), which counts none of the work. Adjusted to 2h: S-0579, which built this kluster pass, took 3033 s against a 1h forecast, and this story changes three paths (serve, accept, the merge queue) and their docs. Delivery 16:10Z is flai's queue estimate of 14:10Z plus the 2h.

Cost of delay. `flai cod` gives 50 USD a week from the operator's 20m lost per 168h cycle at 150 USD an hour. It stands: the gap matters only while the `kluster` host action is on, which is off by default, and the inputs already count two occurrences.

No `after`: the story reads nothing another open story writes.

### Accepted by the orchestrator

- Verified: 9cb9d38759f1b97ea57c3ce5cb9ef7636719b52e
- At: 2026-11-11T14:52:10Z

Verdict: meets both criteria; the close-out passed every step at 9cb9d38759f1, the head of story/S-0694.
- 1: flai/internal/kluster/deploy.go, flai/internal/serve/kluster.go, flai/internal/serve/agents_restart.go, flai/internal/harness/harness.go, flai/internal/preview/accept.go, flai/internal/preview/orchestrator.go, flai/cmd/accept.go, flai/cmd/move.go, flai/cmd/queue.go, flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/batch.go, flai/internal/serve/mergequeue.go, flai/internal/hostapi/reads.go, flai/internal/kluster/deploy_test.go, flai/internal/serve/kluster_test.go, flai/internal/serve/agents_test.go, flai/internal/preview/accept_test.go, flai/cmd/accept_kluster_test.go, flai/internal/mergequeue/eligible_test.go, flai/internal/mergequeue/batch_test.go, flai/internal/serve/checks_test.go
- 2: design/issues/I-0268-adr-0254-5a-s-send-back-of-a-story-whose-kluster-deploy-or-checks-failed-and.md, design/issues/summary.md
