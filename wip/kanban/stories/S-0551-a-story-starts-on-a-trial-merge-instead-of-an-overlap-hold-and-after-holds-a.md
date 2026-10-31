---
id: S-0551
type: story
nature: improvement
title: A story starts on a trial merge instead of an overlap hold, and after holds a story for a dependency only, so the worktrees run stories in parallel as people on one remote do
status: review
parent: E-0031
owner: alex
created: 2026-10-30T01:52:25Z
updated: 2026-10-31T10:33:22Z
transitions:
  - to: ready
    at: 2026-10-31T07:24:00Z
    by: alex
  - to: in-progress
    at: 2026-10-31T08:40:02Z
    by: agent-S-0551
  - to: review
    at: 2026-10-31T10:33:22Z
    by: agent-S-0551
tags: [flai]
topics: [git, orchestration, planning, conventions, dashboard]
touches: [flai/internal/serve, design/workflows/roles, design/system/workflow.md, design/system/strategic-agents.md, design/conventions/work-management.md, docs/users/flai.md, design/adrs, design/adrs/README.md, design/system/metrics.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, flai/internal/storygit/trial.go, flai/internal/storygit/trial_test.go, flai/internal/workitem/hold.go, flai/internal/workitem/hold_test.go, flai/internal/check/check.go, flai/internal/check/check_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/serve/hold_test.go, flai/internal/serve/restart.go, flai/internal/serve/start.go, flai/internal/mcpserver/work.go, flai/internal/mcpserver/work_test.go, flai/internal/mcpserver/orchestrate.go, flai/internal/workitem/boardview.go, flai/internal/storystart/start_test.go, design/workflows/roles/planner.md, template/root/design/workflows/roles/planner.md, flai/internal/roles/baseline/planner.md, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, template/root/.claude/agents/planner.md, ".claude/agents/planner.md", design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, flai/internal/metrics/concurrency.go, flai/internal/metrics/concurrency_test.go, flai/internal/metrics/metrics.go, flai/cmd/stats.go, flai/cmd/stats_test.go, flaiover/src/lib/viz/catalog.ts, flaiover/src/lib/viz/catalog.test.ts, flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts, flaiover/src/lib/viz/explain.ts, flaiover/src/lib/sitemenu.ts, "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts", docs/users/flaiover.md, flai/internal/workitem/repo.go, flai/internal/workitem/share.go, flai/internal/check/after_reason_test.go, flai/internal/check/scope.go, flai/internal/check/scope_test.go, design/issues/I-0227-flai-test-s-go-test-tier-runs-the-mcpserver-package-with-race-into-go-test-s-10-minute-default-timeout-on-a-loaded-host-and-names-only-the-package-so-a-story-that-changes-mcpserver-cannot-pass-its-tests.md, design/issues/summary.md, flai/internal/itemedit/itemedit.go, flai/internal/itemnew/itemnew.go, flai/internal/itemnew/itemnew_test.go, flai/cmd/project.go, flai/internal/mcpserver/server.go, docs/users/flai-reference.md, flai/cmd/promote.go, flai/cmd/story_start.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/story_start.go, design/issues/I-0235-flai-check-finds-narrative-missing-outside-the-story-at-close-out.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 6817
  priced_by: harness
  turns:
    - day: 2026-10-31
      ceremony: 3
      hand_edits: 4
      work: 84
  models:
    - model: claude-opus-5-5
      input: 1026
      output: 325361
      cache_read: 63980067
      cache_write: 2346530
      cost: 32.3828
  strategic:
    - kind: orchestrator
      seconds: 173
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 47
          output: 429
          cache_read: 1651579
          cache_write: 253880
          cost: 0.4338
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2h
    by: alex
    at: 2026-10-31T09:10:28Z
  value: 300
  by: planner-S-0551
  at: 2026-10-31T09:12:03Z
forecast:
  duration: 2h
  delivery: 2026-10-31T13:33:00Z
  basis: "Its own forecast of 2h; 4th in the pull order with an in-progress limit of 5, behind S-0548, S-0550, S-0552, S-0546, S-0547 and S-0549."
  by: flai
  at: 2026-10-31T08:20:45Z
---
# S-0551 A story starts on a trial merge instead of an overlap hold, and after holds a story for a dependency only, so the worktrees run stories in parallel as people on one remote do

## Goal

The worktrees isolate every story, and the scheduler then serializes them: with an in-progress limit of 5, nothing was in progress for 42 percent of the fortnight, one story for 29 percent, two for 26 percent, three for 2 percent, and four or five never, mean 1.6 when anything ran. The journal records 124 overlap holds and 113 `after` holds against 62 real conflicts, the overlap holds on `check.go`, the strategic-agents design and convention, the manifest and dashboard designs, and the `after` chains planner-made, S-0417 to S-0422 each waiting on the one before, as I-0152 already noted. People sharing one remote merge optimistically and resolve the rare conflict; this story makes flai do the same. A ready story whose touches overlap an in-progress story's is started when a trial merge of the in-progress story's branch into a worktree at the ready story's base is clean, and held only when it is not, with the journal naming the conflicting file; `after` is honoured as today but the planner sets it only when a story reads what another writes, named in the story's notes, and `flai check` warns of an `after` with no reason; and the ready wait at p85, 5.4 hours in the fortnight, is a figure `flai stats` reports.

## Acceptance criteria
- [x] `wait_for_work` and the orchestrator's pull start a story overlapping an in-progress one when a trial merge of the in-progress branch at the ready story's base is clean, and hold it with the conflicting path in the journal otherwise, with a test of each
- [x] The planner's role file and the work-management convention say `after` names a story this one reads the output of, with the reason in the story's notes; `flai check` warns of an `after` without one, and the open chains S-0417 to S-0422 are reviewed against it
- [x] `flai stats` reports stories in progress at once as a share of the window and the ready wait at p50 and p85, and the dashboard's flow page shows them
- [x] workflow.md and strategic-agents.md describe the trial merge in place of the overlap hold, and ADR-0046's `after` is refined by a new ADR where the rule changes

## Tasks

The planner drafts them.
- T-2388 Write the ADRs and the design: a trial merge replaces the overlap hold, after names a read dependency with its reason, and the in-progress share is a metric
- T-2389 storygit runs a trial merge of an in-progress story's branch at a ready story's base and returns the conflicting paths
- T-2391 The overlap verdict takes a trial-merge judge, holds a ready story only on a conflicting trial merge, and names the conflicting paths in its reason
- T-2393 flai check advises of a story's after that names a story its Notes give no reason for
- T-2395 flai serve starts a ready story that overlaps an in-progress one on a clean trial merge, and journals the conflicting paths when it holds
- T-2404 wait_for_work, story start, the inbox, the board, and the orchestrator's candidates judge an overlap by the trial merge
- T-2406 The planner's role file, its baseline, and its start prompt say after names a story whose output this one reads, with the reason in its Notes
- T-2408 The work-management and strategic-agents conventions state the trial merge and the after rule, in the tree and the template, with a changelog entry
- T-2410 flai stats reports the share of the window with n stories in progress at once, and the ready wait at p50 and p85
- T-2412 Review the open after chain S-0420 to S-0422 against the rule, keeping each after that reads another story's output with its reason and proposing the rest
- T-2414 The dashboard's flow charts show the share of the window with n stories in progress and the ready wait at p50 and p85
- T-2416 The flai user guide and the CLI design describe the trial merge, the after advisory, and the new stats

## Notes

Evidence: the board's transitions and the serve journal's `held` entries 2026-10-15 to 2026-10-29, counted in the operator's conversation; the 2026-10-15 bottlenecks analysis and I-0152 for the `after` chains.

### Planning

Touches, by source:

- Declared by the operator: `flai/internal/serve`, `flai/internal/sync`, `flai/internal/claims`, `design/workflows/roles`, `design/system/workflow.md`, `design/system/strategic-agents.md`, `design/conventions/work-management.md`, `docs/users/flai.md`. The folders `flai/internal/sync` and `flai/internal/claims` named no package in the tree and have since left the story's touches: the claim and hold code is in `flai/internal/workitem/hold.go`, and the trial merge in `flai/internal/storygit`. The folders `flai/internal/serve` and `design/workflows/roles` stay because the operator declared them; the tasks name the files below them, and TH-0776 proposes dropping them.
- Folder kept by the planner: `design/adrs`. The story adds two ADRs whose file names came from `flai adr new` and no task could name when it was planned; the task names `design/adrs/README.md`.
- Design: `design/system/metrics.md` (a metric changes only with an ADR), `design/system/flai-cli.md`, `design/system/flaiover-dashboard.md`, `design/adrs/README.md`, `design/conventions/strategic-agents.md` (it states the `after` rule), and the template copies of the conventions and the planner's role file, with `template/CHANGELOG.md`.
- Layout, from the explorer's map of the code: `flai/internal/storygit` (the `TrialMerge` helper already in `sync.go`, used between two story branches; a new `trial.go` takes a ready story's base), `flai/internal/workitem/hold.go` (the overlap and `after` verdicts), its callers in `flai/internal/serve/agents.go`, `restart.go`, `start.go`, `flai/internal/mcpserver/work.go` and `orchestrate.go`, `flai/internal/workitem/promote.go` and `boardview.go`, `flai/internal/storystart/start.go`, and `flai/internal/inbox/inbox.go`; `flai/internal/check/check.go` (the `wip.after-hub` advisory is the precedent); the planner's start prompt in `flai/internal/harness/harness.go`, its baseline in `flai/internal/roles/baseline/planner.md`, and its generated definitions under `.claude/agents` and `template/root/.claude/agents`; `flai/internal/metrics` and `flai/cmd/stats.go`; the flow charts in `flaiover/src/lib/viz` and `flaiover/src/lib/sitemenu.ts`.
- Co-change (`flai touches suggest`): `design/system/flai-cli.md` (54 percent), `template/CHANGELOG.md`, `design/adrs/README.md`, `template/root/design/conventions/work-management.md`, `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`. Left out: `docs/users/flai-reference.md` (no flag changes), `docs/operators` (no setting is planned), `design/issues/summary.md`.
- Story files: the files of S-0420, S-0421, and S-0422, which the chain-review task edits through `flai edit`, are in that task's touches only. S-0417 to S-0419 are archived, so the open chain is S-0420 to S-0422.

The ready wait at p50 and p85 is already computed as queue time ([ADR-0190](../../../design/adrs/0190-an-item-s-queue-time-is-the-seconds-it-spent-in-ready-before-it-started-so-a.md)). The stats task reports it under that name and adds the share of the window with n stories in progress, which is new.

The topic `dashboard` is added for T-2414's flow chart.

Forecast: `flai forecast` gave 1h30m, 89 s per unit of size over 108 done large improvement stories on claude-opus-5-5, times size 60 (4 criteria, 56 touches). Raised to 2h: twelve tasks in three layers, hold tests that need temporary git repositories, a `.claude/` write that waits on the operator, and S-0548's and S-0552's claims on serve and stats files to agree in conversation. flai sets the delivery from the pull order.

Cost of delay: 300 USD a week, as `flai cod` works it out from the operator's `time_lost_per_cycle` of 2h (TH-0772): 2h lost per 168h cycle at the manifest's 150 USD an hour, one cycle a week. It stands unadjusted: the input is the operator's, and nothing in the story argues for a different cycle.
