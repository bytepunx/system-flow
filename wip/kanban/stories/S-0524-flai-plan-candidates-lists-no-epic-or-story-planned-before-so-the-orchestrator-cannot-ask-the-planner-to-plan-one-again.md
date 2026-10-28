---
id: S-0524
type: story
nature: improvement
title: flai plan --candidates lists no epic or story planned before, so the orchestrator cannot ask the planner to plan one again
status: backlog
owner: alex
created: 2026-10-28T12:53:50Z
updated: 2026-10-28T18:28:14Z
transitions: []
tags: [flai, orchestrator, planner]
topics: [orchestration, planning]
touches: [design/adrs, design/adrs/README.md, design/system/strategic-agents.md, design/system/flai-cli.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md, docs/users/flai.md, docs/users/flai-reference.md, flai/internal/workitem/plancandidates.go, flai/internal/workitem/plancandidates_test.go, flai/cmd/plan.go, flai/cmd/plan_candidates_test.go, flai/cmd/plan_orchestrate_test.go, flai/cmd/orchestrate_permissions_test.go, flai/internal/mcpserver/plan.go, flai/internal/mcpserver/plan_test.go, flai/internal/serve/plan.go, flai/internal/serve/plan_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/issues/I-0188-flai-plan-candidates-lists-no-epic-or-story-planned-before-so-the-orchestrator-cannot-ask-the-planner-to-plan-one-again.md, design/issues/summary.md]
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
    - kind: orchestrator
      seconds: 156
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 16
          output: 207
          cache_read: 5602145
          cache_write: 17774
          cost: 1.2338
cost_of_delay:
  inputs:
    time_lost_per_cycle: 30m
    by: orchestrator
    at: 2026-10-28T14:44:06Z
  value: 75
  by: planner-S-0524
  at: 2026-10-28T14:45:49Z
forecast:
  duration: 1h15m
  delivery: 2026-10-29T09:10:00Z
  basis: "Its own forecast of 1h15m; 35th in the pull order with an in-progress limit of 5, behind S-0415, S-0502, S-0531, S-0419, S-0423, S-0418, S-0421, S-0416, S-0422, S-0420, S-0417, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0424, S-0425, S-0506, S-0507, S-0508, S-0509, S-0510, S-0511, S-0512, S-0514, S-0515, S-0516, S-0517, S-0518, S-0519, S-0520 and S-0523."
  by: flai
  at: 2026-10-28T18:28:14Z
finalized:
  by: alex
  at: 2026-10-28T12:59:03Z
---
# S-0524 flai plan --candidates lists no epic or story planned before, so the orchestrator cannot ask the planner to plan one again

## Goal

This story remediates [I-0188](../../../design/issues/I-0188-flai-plan-candidates-lists-no-epic-or-story-planned-before-so-the-orchestrator-cannot-ask-the-planner-to-plan-one-again.md), "flai plan --candidates lists no epic or story planned before, so the orchestrator cannot ask the planner to plan one again". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0188 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0188 is closed with `flai issue close I-0188 --reason` saying what fixed it

## Tasks
- T-2246 An ADR refining ADR-0119 lets the orchestrator ask the planner again for an epic or story planned before when the operator has written on it since its newest plan
- T-2247 flai plan --candidates lists an epic or story planned before that the operator has written on since its newest planner run, naming the thread, and a test reproduces I-0188
- T-2248 The orchestrator's plan admits a candidate planned before, the run records the trigger orchestrator; TH-nnnn, and the planner's prompt names the thread the operator wrote on
- T-2249 The design, the convention and its template copy, and the guides say the orchestrator asks the planner again for an item planned before that the operator wrote on, judging first that the words ask for a change
- T-2250 I-0188 is closed with the reason naming the ADR and the third kind of candidate

## Notes

### Planning

Planned by planner-S-0524 on 2026-10-28; the plan is on TH-0695, where the orchestrator set the cost of delay input and recommends the remedy to the operator. The story declared no touches, so `flai touches suggest S-0524` had nothing to start from; every touch comes from the issue's instance, the design, and the code layout.

**The solution proposed from I-0188's instance.** `flai plan --candidates` lists an epic only in the backlog with no stories, or open with every story done or cancelled, and a story only while it lacks a plan (ADR-0119 §1, `PlanCandidatesOf` in `flai/internal/workitem/plancandidates.go`), and `orchestratorPlans` in `flai/cmd/plan.go` refuses the orchestrator any other item; `flai guard` holds the tool to the permission by the ID's type alone. So the scope change the operator confirmed on TH-0579, a new E-0024 story after S-0419, waits for the operator to run `flai plan E-0024` by hand. Proposed, for a new ADR refining ADR-0119 and ADR-0087 (T-2246): a third kind of candidate, planned before and asked again, an epic not done, cancelled, or archived, or a story in the backlog or ready, whose newest planner run has ended and on which, or on an item under which, the operator (`threads.OperatorName`) has written since; the candidate names the thread (T-2247). The orchestrator asks for it under `plan_backlog_epics` or `plan_backlog_stories` after judging that the operator's words ask for a change to the plan, the run records the trigger `orchestrator; TH-nnnn`, and the planner's prompt names the thread (T-2248). flai bounds the spend to one run per operator entry, since a run makes every earlier entry older than it. Not proposed: planning any open item (no bound), an edit of the epic's words as the trigger (misses the instance; the replanner already replans a story on an edit), or a run on every operator entry (no judgement). The operator may choose otherwise on TH-0695 before T-2246 is written.

**Where each touch came from.**

- Layout, the candidates (T-2247): `flai/internal/workitem/plancandidates.go` and `plancandidates_test.go` (`PlanCandidatesOf`, `PlanCandidate`), `flai/cmd/plan.go` and `plan_candidates_test.go` (`planCandidatesOf`, `planHeld`, `plannerAsked`, `updatedAfter`, and the fixture test that reproduces I-0188).
- Layout, the orchestrator's asking and the run (T-2248): `flai/cmd/plan.go` (`orchestratorPlans`, `mcpPlan`, `planNow`), `flai/cmd/plan_orchestrate_test.go`, `flai/cmd/orchestrate_permissions_test.go`, `flai/internal/mcpserver/plan.go` and `plan_test.go` (the tool's answer), `flai/internal/serve/plan.go` and `plan_test.go` (`PlanForOrchestrator`, `planFor`, the trigger), `flai/internal/harness/harness.go` and `harness_test.go` (`planPrompt`).
- Design (T-2246, T-2249): `design/adrs/README.md` (the index row `flai adr new` writes), `design/system/strategic-agents.md` (Starting it, What it does with each permission, What a run records), `design/system/flai-cli.md` (the `flai plan --candidates` row).
- Design, the convention and the guides (T-2249): `design/conventions/strategic-agents.md` and `template/root/design/conventions/strategic-agents.md` (As the orchestrator), `template/CHANGELOG.md`, `docs/users/flai.md`, `docs/users/flai-reference.md` (regenerated by `make flai-reference` where the help text changes).
- Declared by the criteria (T-2250): `design/issues/I-0188-….md` and `design/issues/summary.md`, which `flai issue close` rewrites.

**Folder touch kept.** `design/adrs`: the story adds one ADR, whose number `flai adr new` allots when the story's agent runs it, so no task can name the file yet. The folder is in the manifest's shared paths, so it holds no ready story (ADR-0096), and `flai task done` widens the touches with the committed path (ADR-0144). Ready stories S-0416, S-0421, S-0423 and backlog S-0523 keep the same folder touch for the same reason.

**Overlaps on 2026-10-28.** `flai/internal/harness/harness.go` is touched by backlog S-0523; `design/system/strategic-agents.md` and `design/conventions/strategic-agents.md` by S-0502 (in progress) and by the ready E-0024 stories S-0416 to S-0423; so whichever starts first holds the other, and the pull order settles which. `docs/users/flai.md`, `docs/users/flai-reference.md`, `design/system/flai-cli.md`, `design/adrs`, and `design/issues` are in the manifest's shared paths and hold nothing. No `after`: the story needs no other story's change on main.

**Forecast.** `flai forecast S-0524` gave 3m with no touches (size 2) and gives 38m with them: the median 89 s per unit of size over 98 done improvement stories on claude-opus-5-5 in the large band, times size 25 (2 criteria, 23 touches), delivery 2026-10-29T04:45Z, 35th in the pull order. Raised to 1h15m: five tasks in four serial layers, an ADR, fixture tests across four Go packages (cmd, mcpserver, serve, harness) that need a planner run's state and threads, and seven documents, the generated reference and the template changelog among them; S-0523, of the same shape, carries 1h15m. The delivery written is flai's play-out start, 04:45Z less 38m, plus 1h15m: 2026-10-29T05:22Z.

**Cost of delay.** 75.00 USD a week, as `flai cod S-0524` works it out from the input the orchestrator set on TH-0695 at my recommendation: `time_lost_per_cycle: 30m` at 150 USD an hour, one 168h cycle a week. The figure stands: each scope change the operator decides on a thread while the orchestrator runs sends them to run `flai plan` by hand, and the drafted story waits meanwhile; S-0523, from the same TH-0579 instance, carries the same input.
