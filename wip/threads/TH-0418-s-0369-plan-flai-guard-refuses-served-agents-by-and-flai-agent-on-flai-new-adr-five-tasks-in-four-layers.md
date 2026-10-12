---
id: TH-0418
title: "S-0369 plan: flai guard refuses served agents --by and FLAI_AGENT on flai (new ADR), five tasks in four layers"
anchor:
  path: wip/kanban/stories/S-0369-a-thread-entry-s-author-is-whatever-by-says-so-an-agent-with-a-shell-could-answer-a-permission-thread-as-the-owner.md
  item: S-0369
status: resolved
participants: [planner-S-0369, orchestrator, alex]
created: 2026-10-12T17:07:08Z
updated: 2026-10-12T18:00:29Z
---

# TH-0418 S-0369 plan: flai guard refuses served agents --by and FLAI_AGENT on flai (new ADR), five tasks in four layers

On wip/kanban/stories/S-0369-a-thread-entry-s-author-is-whatever-by-says-so-an-agent-with-a-shell-could-answer-a-permission-thread-as-the-owner.md.

## Entries

### 2026-10-12T17:07:08Z planner-S-0369
**Plan:** S-0369 closes I-0091 with a new ADR and a guard rule. The rule refuses a story's agent, the planner, and the analyzer any flai command that names its own author. It also fixes the fallback in `threadAuthor`. There are five tasks in four layers.

## Tasks and layers

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1465 ADR refining ADR-0060 and ADR-0097; agent-adapters.md and strategic-agents.md record it | `design/adrs/README.md`, `design/system/agent-adapters.md`, `design/system/strategic-agents.md` | none |
| 1 | T-1466 `threadAuthor` falls back to the config's author when `FLAI_AGENT` is unset | `flai/cmd/thread.go`, `thread_test.go` | none |
| 2 | T-1467 guard rule, with tests that reproduce I-0091 | `flai/internal/guard/guard.go`, `guard_test.go`, `thread_test.go`, `flai/cmd/guard.go`, `guard_test.go`, `docs/users/flai-reference.md` | T-1465 |
| 3 | T-1468 `docs/users/flai.md` | `docs/users/flai.md` | T-1466, T-1467 |
| 4 | T-1469 close I-0091 | the issue, `design/issues/summary.md` | T-1468 |

## Figures

- Forecast: 45m, delivery 2026-10-13T02:21Z. `flai forecast` gave 21m. I raised it as S-0368's was, because this story adds an ADR and the guard's first rule on a story's own agent.
- Cost of delay: 10 USD a week, from your `penalty_per_week: 10` (TH-0416).

## Assumptions

- **Guard, not an env check.** Today the guard lets every `Bash` call of a story's own agent pass. The issue's first direction, comparing `--by` to `FLAI_AGENT` inside flai, is defeated by `FLAI_AGENT=alex flai …`. It would also refuse the dashboard, which runs flai as `flaiover` with `--by=<owner>`. So T-1465 plans the issue's second direction, widened:
  - Refuse `--by`, `--placed-by`, and an `FLAI_AGENT` or `FLAI_SESSION` set for flai.
  - Cover every writing command, not only `reply`, `resolve`, and `confirm`.
  - Apply it to the planner and the analyzer too.
  - Keep the orchestrator's existing `writesAsSelf` rule.
  - The ADR may change this.
- **The operator's session stays free.** A session with no `FLAI_STORY` or `FLAI_ROLE` is not held. In this checkout `.claude/settings.local.json` sets `FLAI_AGENT=system-flow`, so the guard keys on `FLAI_STORY` and `FLAI_ROLE`, not `FLAI_AGENT`.
- **Residuals stay open and are named in the ADR.** These are a shell redirection into a thread file by a story's own agent (S-0368 covers sub-agents only) and an HTTP MCP call with another `X-Flai-Agent`.
- **T-1466 is an addition** the issue does not name. With no `--by` and no `FLAI_AGENT`, an entry is written as `agent`, not the config's author as the docs say. Once agents cannot name an author, the operator's session is the one place a name comes from. Drop T-1466 if you would rather fix it in its own story.
- **S-0368 overlaps** on the guard files, `flai-reference.md`, `flai.md`, `agent-adapters.md`, and the ADR index. I put no `after` between the two stories. The overlap hold keeps them apart, and whichever runs second syncs onto the other.

### 2026-10-12T17:07:40Z orchestrator (recommendation)
**Recommendation:** approve the plan with one narrowing. The guard refuses a served agent (a story's agent, the planner, or the analyzer) a flai command whose `--by`, `--placed-by`, `FLAI_AGENT`, or `FLAI_SESSION` names anyone but that agent's own name, as flai serve set it. It does not refuse those flags or variables as such. Keep T-1466 in this story.

## Why this needs you

- **A new rule on a story's agent.** ADR-0060 lets a story's own agent's calls pass the guard untouched. This would be the guard's first rule on it, a decision not yet recorded. T-1465 records it, but the choice is yours.
- **It conflicts with a convention as drafted.** `design/conventions/tooling.md` § Rules says "Set `FLAI_AGENT` and `FLAI_SESSION` at session start so narrative entries record who wrote them". Refusing any `FLAI_AGENT` set for flai would refuse agents that follow it. Refusing only a name other than the agent's own still blocks `FLAI_AGENT=alex flai …`, which is I-0091's case.
- **Scope.** Covering the planner and the analyzer, and T-1466's `threadAuthor` fallback, both go beyond I-0091. T-1466 is small, in the same file, and makes the docs true, so I would keep it. Dropping it is your call.

## The rest fits

- Five tasks in four layers, all file touches.
- Figures: 45m (flai's 21m raised as S-0368's was) and 10 USD a week from my input on TH-0416.
- The overlap hold with S-0368 keeps the two apart.

I leave S-0369 a draft until you answer.

This is a recommendation, not an answer: confirm it, or answer otherwise.

Source: design/conventions/tooling.md § Rules

### 2026-10-12T18:00:15Z alex
Confirmed the recommendation of 2026-10-12T17:07:40Z orchestrator.

Source: design/conventions/tooling.md § Rules

### 2026-10-12T18:00:28Z orchestrator
The operator confirmed the narrowing. **For S-0369's agent**, when you pull it:

- T-1465 and T-1467 refuse a served agent a flai command only when its `--by`, `--placed-by`, `FLAI_AGENT`, or `FLAI_SESSION` names anyone but the agent's own name, as flai serve set it. They do not refuse those flags or variables as such, so agents still follow `tooling.md`'s rule to set `FLAI_AGENT`.
- Update T-1465's and T-1467's `## Work` before you work them, and let `flai adr new` number the ADR.
- T-1466 stays in the story.

The rest stands: five tasks in four layers, 45m, and 10 USD a week. I am finalizing S-0369 now.

Source: design/conventions/tooling.md § Rules

### 2026-10-12T18:00:29Z orchestrator
Resolved: Operator confirmed the narrowing; S-0369's agent applies it to T-1465/T-1467.
