---
id: S-0467
type: story
nature: improvement
title: A story in review waits days for the operator when the orchestrator's acceptance stops on a dirty main, a rebase on a shared path, or a gate it raised itself
status: ready
owner: alex
created: 2026-10-21T11:57:59Z
updated: 2026-10-22T21:44:12Z
transitions:
  - to: ready
    at: 2026-10-22T20:49:53Z
    by: alex
tags: [flai, template]
topics: [orchestration]
touches: [flai/cmd/accept.go, flai/cmd/accept_batch_dirty.go, flai/cmd/accept_batch_dirty_test.go, flai/cmd/accept_uncommitted_test.go, flai/cmd/accept_orchestrator_test.go, flai/cmd/accept_own_thread_test.go, flai/cmd/accept_conflict_test.go, flai/cmd/branch.go, flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, flai/internal/preview/orchestrator.go, flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go, flai/internal/threads/threads.go, flai/internal/threads/threads_test.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, ".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md, design/adrs/0186-the-orchestrator-s-acceptance-lands-on-a-dirty-main-as-the-queue-does.md, design/adrs/README.md, design/system/strategic-agents.md, design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0150-a-story-in-review-waits-days-for-the-operator-when-the-orchestrator-s-acceptance-stops-on-a-dirty-main-a-rebase-on-a-shared-path-or-a-gate-it-raised-itself.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 290
  priced_by: harness
  turns:
    - day: 2026-10-22
      work: 10
  models:
    - model: claude-haiku-5-5
      input: 80
      output: 40242
      cache_read: 4975842
      cache_write: 202818
      cost: 0.3874
    - model: claude-opus-5-5
      input: 22
      output: 4645
      cache_read: 568705
      cache_write: 61338
      cost: 0.6974
  strategic:
    - kind: orchestrator
      seconds: 77
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 44
          output: 679
          cache_read: 7979314
          cache_write: 22767
          cost: 1.7669
    - kind: analyzer
      seconds: 114
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 98
          output: 6305
          cache_read: 370755
          cache_write: 41688
          cost: 1.2427
        - model: claude-haiku-5-5
          input: 8
          output: 6718
          cache_read: 204810
          cache_write: 23259
          cost: 0.0083
cost_of_delay:
  inputs:
    time_lost_per_cycle: 96h
    by: flai
    at: 2026-10-21T11:57:59Z
  value: 14400
  by: planner-S-0467
  at: 2026-10-22T21:02:20Z
forecast:
  duration: 2h
  delivery: 2026-10-23T00:03:00Z
  basis: "Its own forecast of 2h; 2nd in the pull order with an in-progress limit of 5, behind S-0466, S-0476 and S-0468."
  by: flai
  at: 2026-10-22T21:44:12Z
finalized:
  by: alex
  at: 2026-10-21T20:51:58Z
---
# S-0467 A story in review waits days for the operator when the orchestrator's acceptance stops on a dirty main, a rebase on a shared path, or a gate it raised itself

## Goal

This story remediates [I-0150](../../../design/issues/I-0150-a-story-in-review-waits-days-for-the-operator-when-the-orchestrator-s-acceptance-stops-on-a-dirty-main-a-rebase-on-a-shared-path-or-a-gate-it-raised-itself.md), "A story in review waits days for the operator when the orchestrator's acceptance stops on a dirty main, a rebase on a shared path, or a gate it raised itself", which the analysis in [2026-10-15-bottlenecks.md](../../../design/analysis/2026-10-15-bottlenecks.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0150 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0150 is closed with `flai issue close I-0150 --reason` saying what fixed it

## Tasks
- T-2037 An ADR settles the three remedies: the orchestrator's acceptance lands on a dirty main as the queue does, aborts a rebase it cannot settle and sends the story back, and is not gated by a thread of its own
- T-2039 flai accept --by orchestrator lands on a dirty main as the batch's landing does, refuses only a change that would overwrite an uncommitted file, and finishes an acceptance of its own that stopped before its commit
- T-2042 A thread the orchestrator opened with no entry by anyone else does not gate the story: the acceptance preview notes it, the merge queue's eligibility leaves it out, and the acceptance resolves it
- T-2047 An orchestrator's acceptance whose rebase stops on a conflict it cannot settle aborts the rebase, sends the story back to in-progress with the paths, and flai serve starts its agent as for a story the queue sent back
- T-2049 The orchestrator's prompt, agent definition, and convention say it accepts on flai's blockers and the criteria alone, finishes its own stopped acceptance, and expects the send-back, with the template copies and changelog
- T-2052 The design and the user docs describe the orchestrator's acceptance on a dirty main, the send-back on a stop it cannot settle, and the thread that does not gate it
- T-2053 I-0150 is closed with the reason naming what fixed each of its three causes

## Notes

The analyzer found I-0150: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0150. time_lost_per_cycle 96h carried over from I-0150's Impact section. Its Impact time_lost_per_cycle was taken rather than the 96h derived from its cost and count.

Strategic usage carried over from I-0150: analyzer 1m54s, 1.2510 USD.

### Planning

Planned by planner-S-0467 on 2026-10-22. The proposed solution, from I-0150's instances (TH-0387, TH-0389, TH-0391, TH-0394, TH-0397), is in T-2037's ADR and the plan's thread: the orchestrator's acceptance lands on a dirty main as the merge queue's landing does (ADR-0178) and finishes a stopped acceptance of its own; a rebase it cannot settle is aborted and the story sent back to its agent, as the queue's batch does (ADR-0175, ADR-0176); and a thread the orchestrator alone wrote does not gate the story.

Touches, file by file; the story declared none, so `flai touches suggest` was seeded with the acceptance files the design names:

- Layout, the acceptance and its preview: `flai/cmd/accept.go`, `flai/cmd/accept_batch_dirty.go`, `flai/cmd/branch.go`, `flai/internal/preview/accept.go`, `flai/internal/preview/orchestrator.go`, and their tests `flai/cmd/accept_batch_dirty_test.go`, `flai/cmd/accept_uncommitted_test.go`, `flai/cmd/accept_orchestrator_test.go`, `flai/cmd/accept_conflict_test.go`, `flai/internal/preview/accept_test.go`; `flai/cmd/accept_own_thread_test.go` is new.
- Layout, the queue's eligibility and the threads: `flai/internal/mergequeue/eligible.go`, `flai/internal/mergequeue/eligible_test.go`, `flai/internal/threads/threads.go`, `flai/internal/threads/threads_test.go`.
- Layout, the rebase's abort and the send-back's start: `flai/internal/storygit/sync.go`, `flai/internal/storygit/sync_test.go`, `flai/internal/serve/agents.go`, `flai/internal/serve/agents_test.go` (co-change 11% and 8%).
- Co-change, the orchestrator's prompt: `flai/internal/harness/harness.go` (seed) and `flai/internal/harness/harness_test.go` (60%).
- Design and co-change, the words the orchestrator works by: `.claude/agents/orchestrator.md`, `template/root/.claude/agents/orchestrator.md`, `design/conventions/strategic-agents.md` (9%), `template/root/design/conventions/strategic-agents.md` (9%), `template/CHANGELOG.md` (21%).
- Design, the documents that describe the acceptance: `design/system/strategic-agents.md` § Accepting a story, `design/system/workflow.md` (6%), `design/system/flai-cli.md` (20%), `docs/users/flai.md` (16%); `design/adrs/README.md` (4%) with the new ADR, predicted as `design/adrs/0186-the-orchestrator-s-acceptance-lands-on-a-dirty-main-as-the-queue-does.md` from the next number and the slug rule, which `flai task done` widens to the name `flai adr new` gives.
- Criteria: `design/issues/I-0150-...md` and `design/issues/summary.md`, which `flai issue close` writes.

No folder touch was kept: every file is named. Left out as co-change noise: `docs/users/flai-reference.md` (generated, no new flag), `design/conventions/work-management.md` and `delegation.md`, `flai/internal/guard/guard.go` (flai aborts the rebase itself, so the orchestrator needs no git write), and `flai/internal/hostapi/writes.go`.

Forecast: `flai forecast` gave 54m at size 36 (2 criteria, 34 touches, large band, median 89 s per unit over 87 improvement stories on claude-opus-5-5). Raised to 2h: the seven tasks form five serial layers across six Go packages, the send-back needs a start path in `flai serve`, and the `.claude/` write waits on the operator's permission prompt. Delivery 2026-10-22T23:20:00Z is the start now plus 2h times the cycle factor 1.15, a lane being free under the in-progress limit of 5.

Cost of delay: `flai cod` gave 14400 USD a week from the operator's input, 96h lost per 168h cycle at 150 USD an hour, one cycle a week. It stands: the input is flai's carry-over from I-0150's impact, and the 259 h of review wait in W42 bears it out.

Tags `flai` and `template` were added because the story touches both components and `flai check` requires a component tag; topic `orchestration` was added so the story's pack briefs `design/system/strategic-agents.md` and ADR-0093.
