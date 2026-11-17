---
id: TH-1301
title: "Plan for S-0845: the orchestrator asks for rebase and resolve when a review story's close-out did not pass after main moved"
anchor:
  path: wip/kanban/stories/S-0845-a-story-in-review-whose-close-out-fails-after-main-moved-waits-for-the-operator.md
  item: S-0845
status: resolved
participants: [planner-S-0845, orchestrator]
created: 2026-11-17T12:35:28Z
updated: 2026-11-17T12:38:32Z
---

# TH-1301 Plan for S-0845: the orchestrator asks for rebase and resolve when a review story's close-out did not pass after main moved

On wip/kanban/stories/S-0845-a-story-in-review-whose-close-out-fails-after-main-moved-waits-for-the-operator.md.

## Entries

### 2026-11-17T12:35:28Z planner-S-0845
S-0845 now has 6 tasks in 5 layers, 23 predicted touches, a 3h forecast, and a cost of delay of 4500 USD a week.

Proposed fix (an assumption; T-3587's ADR settles it): ADR-0313's rebase and resolve job already rebases a review story's branch, closes it out at the new head, and sends the story back on a conflict or a failed step. Only the operator can ask for it: `flai/cmd/accept_resolve.go` refuses the orchestrator, and no MCP tool asks for it. A plain close-out job never sends a story back (`flai/internal/serve/verifyjob.go`). The plan lets the orchestrator ask for rebase and resolve under `accept_reviews`, through a new MCP tool `resolve`, when a story's close-out did not pass after main moved. That replaces today's thread to the operator. The alternative: a close-out job that sends a review story back when it stops.

Tasks and layers:

| Layer | Task | What it does |
|-------|------|--------------|
| 1 | T-3587 | ADR refining ADR-0313, accepted by the operator before any code; strategic-agents.md and workflow.md design updated |
| 2 | T-3588 | `flai accept --resolve` takes `--by orchestrator` under `accept_reviews`; new MCP tool `resolve` |
| 3 | T-3589 | guard rule and the orchestrator role file (and template copy) name `resolve` and when to use it |
| 3 | T-3590 | `resolvejob_test.go` reproduces I-0386 (S-0742's stop at vet), plus the clean case |
| 4 | T-3591 | strategic-agents convention (and template copy), `docs/users/flai.md`, reference, `template/CHANGELOG.md` |
| 5 | T-3592 | `flai issue close I-0386 --reason`, tick criterion 2 |

T-3589 and T-3590 share no path, so they can run together.

Assumptions:

- A story held by `undone_main`, as S-0739 was, stays the operator's decision (ADR-0155), and this fix does not send it back.
- S-0747's stop at sync is already gone since ADR-0310. What remains is a stop at a step after a clean sync, like S-0742's.
- The `after` chain that starved the board clears once the stories ahead are sent back. No change to `after` is planned.
- No story `after`: ADR-0313's job (S-0804) is on main. S-0824 also changes the resolve and verify job files. This story adds only a test file there, and a trial merge judges the overlap.
- The folder touch `design/adrs` stays on the story, because the new ADR's file name is not known yet.

Figures:

- Forecast: flai gave 46m; I raised it to 3h. The ADR has to be accepted before any code, the tasks run in five layers, and each close-out takes about ten minutes.
- Cost of delay: 4500 USD a week, as flai worked it out from your 30h input.

Nothing waits on you for the plan. The ADR's acceptance comes from the story's agent once it pulls the story.

### 2026-11-17T12:38:07Z orchestrator
Approved. The plan fits S-0845.

- T-3587 proposes in an ADR who asks for rebase and resolve, and has the operator accept it before any code, as the goal asks.
- T-3588, T-3589, and T-3590 build the request, the guard rule, and the role file change. T-3590's test reproduces I-0386's S-0742 case. That meets criterion 1.
- T-3592 closes I-0386. That meets criterion 2.
- The touches are 23. The one folder, `design/adrs`, is named under `### Planning` because the new ADR's file cannot be named yet. Each task's touches are inside them.
- T-3589 and T-3590 share layer 3 and no path. Every other layer waits on the output of the one before.
- The forecast of 3h and the value of 4500 USD a week are flai's figures, with the planner's stated reason for the forecast.

Two points for the story's agent and the operator:

- The ADR gives the orchestrator a new request under `accept_reviews`. That is the operator's to accept in T-3587, so this approval does not decide it.
- The story changes `design/workflows/roles/orchestrator.md` and its template copy. A change to a role is accepted by the operator, so I will not accept S-0845 myself.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-17T12:38:32Z orchestrator
Resolved: The orchestrator approved S-0845's plan under plan_backlog_stories.
