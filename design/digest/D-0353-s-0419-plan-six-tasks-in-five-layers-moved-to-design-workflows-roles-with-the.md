---
id: D-0353
title: "S-0419 plan: six tasks in five layers, moved to design/workflows/roles, with the embedded baseline first: are you able to record the two items that just came up as issues? you should be allowed to…"
kind: answer
sources: [TH-0579, S-0524, I-0188]
topics: [cli, planning, orchestration, template, conventions, flai, go, code, all]
paths: [design/conventions/strategic-agents.md, design/workflows/roles/analyzer.md, design/workflows/roles/orchestrator.md, design/workflows/roles/planner.md, design/system/strategic-agents.md, flai/internal/check/check.go, flai/internal/check/roles_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/roles/baseline.go, flai/internal/roles/baseline_test.go, flai/internal/roles/baseline/planner.md, flai/internal/roles/baseline/orchestrator.md, flai/internal/roles/baseline/analyzer.md, flai/internal/serve/agents.go, flai/internal/serve/analysis_schedule.go, flai/internal/serve/analyze.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_hold.go, flai/internal/serve/plan.go, flai/internal/serve/replan.go, flai/internal/serve/role_continuous.go, flai/internal/serve/role.go, flai/internal/serve/role_test.go, flai/internal/serve/serve.go, template/CHANGELOG.md, template/root/design/conventions/strategic-agents.md, template/root/design/workflows/roles/analyzer.md, template/root/design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/planner.md, flai/cmd/plan.go, flai/internal/workitem/plancandidates.go]
date: 2026-10-21
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0353 S-0419 plan: six tasks in five layers, moved to design/workflows/roles, with the embedded baseline first: are you able to record the two items that just came up as issues? you should be allowed to…

The operator, alex, answered TH-0579 at 2026-10-21T12:40:26Z.

What was asked:

```text
S-0419 is replanned on your TH-0471 answers. It keeps its six tasks. The embedded baseline now comes before the prompt move. The story has a forecast and a cost of delay value.
```

The answer:

```text
are you able to record the two items that just came up as issues? you should be allowed to update the cost of delay and you should be allowed to invoke the planner on epics/stories even if they’ve already been planned before
```

The last request, starting the planner on items planned before, was recorded as I-0188. S-0524 built it in a bounded form under ADR-0222. `flai plan --candidates` now also lists an epic or story planned before when the operator has written on a thread on it, or on an item under it, since its newest planner run ended. The orchestrator asks the planner again only when that entry asks for a change to the plan. S-0524's narrative says: "The operator asked for this capability on TH-0579 (I-0188's instance), and the proposal is its bounded form."

Source: TH-0579; the Decisions of wip/archive/agents/S-0524.md.
