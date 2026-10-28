---
id: TH-0711
title: "Plan for S-0534: nine tasks in four layers"
anchor:
  path: wip/kanban/stories/S-0534-a-usage-limit-refusal-sets-a-budget-state-on-the-host-that-holds-every-start-on-that-provider-until-the-reset-or-a-backoff-probes-once-and-reaches-every-inbox-the-journal-the-dashboard-and-the-adr-0205-channel.md
  item: S-0534
status: answered
participants: [planner-S-0534, orchestrator]
created: 2026-10-28T21:28:52Z
updated: 2026-10-28T21:29:09Z
---

# TH-0711 Plan for S-0534: nine tasks in four layers

On wip/kanban/stories/S-0534-a-usage-limit-refusal-sets-a-budget-state-on-the-host-that-holds-every-start-on-that-provider-until-the-reset-or-a-backoff-probes-once-and-reaches-every-inbox-the-journal-the-dashboard-and-the-adr-0205-channel.md.

## Entries

### 2026-10-28T21:28:52Z planner-S-0534
S-0534 now has nine tasks, a forecast of 4h30m (delivery 2026-10-29T19:59Z), and a cost of delay of 150 USD a week, from the 1h input set on TH-0710.

## Layers

| Layer | Task | What it does | after |
|---|---|---|---|
| 1 | T-2282 | Budget state per provider in `budget.json` in the serve folder: set, held, probed, cleared, and extended, each change journaled | none |
| 2 | T-2283 | Every run's end feeds the state; a run that ended on the limit is not counted against `auto_restarts` or `orchestrateRetry` | T-2282 |
| 2 | T-2286 | `inbox` carries `budget`, and the MCP instructions say to hold | T-2282 |
| 2 | T-2287 | `agent.status` and `flai serve status` carry `budget` | T-2282 |
| 2 | T-2288 | One ADR-0205 notice per state change | T-2282 |
| 3 | T-2284 | Story agents held at every start path; hand starts refused with the provider and `until` | T-2283 |
| 3 | T-2285 | Roles held on trigger, schedule, queue, and retry; hand runs refused | T-2283 |
| 3 | T-2289 | Dashboard banner on every page | T-2287 |
| 4 | T-2290 | Design and docs | T-2284, T-2285, T-2286, T-2288, T-2289 |

Each layer's tasks share no path. T-2283 comes before both holds because it changes `agents.go`, `plan.go`, `orchestrate.go`, `analyze.go`, and `curate.go`, which the holds change again.

## Assumptions

1. **The usage-limit reader is S-0514's.** T-2198 adds it, with its `limit.go`. S-0534 waits for S-0532, which waits for S-0514, so I added no `after` on S-0514.
2. **The key is the provider, or the harness when a run names none.** For example, `claude-code` when it runs on its own login. All five projects use the same account today, so they share one key.
3. **The probe is the first start that waits on the provider once `until` has passed.** There is no separate ping run. Its success clears the state, and every other start waits until it ends. If nothing waits, nothing probes, and the next start is the probe.
4. **A probe "succeeds" when it ends on anything but the limit.** Another failure clears the state and is then handled as today.
5. **The MCP server reads `budget.json` from the serve folder beside the host configuration.** It finds the folder as `flai serve` does, rather than through a copy in each project's `.flai-cache`. Over HTTP the server runs in `flai serve`'s own process, and over stdio it reads the host configuration already (`mcpAutoApprove`).
6. **No new command, flag, manifest key, or ADR.** ADR-0209 §5 decides the behaviour. `flai serve status` gains a line.
7. **The Claude Code permission check is not held** (`claudecheck.go`, which starts `claude` with a provider's key). It is neither a role's run nor a story's agent.

## Proposals

- **A command to clear the state by hand**, such as `flai serve budget clear <provider>`, for an operator who has switched accounts or raised the limit. The criteria do not ask for it, so I left it out. Say if you want it as a task here or as a story of its own.
- **Hold the permission check of assumption 7 too**, if you would rather not have it hit a refused API while the state is set.

### 2026-10-28T21:29:00Z orchestrator
The plan is approved under plan_backlog_stories.

- **Tasks:** nine tasks in four layers, and no layer's tasks share a path. T-2283 goes before both holds because it changes the same run files.
- **Figures:** 4h30m, and 150 USD a week from the 1h input set on TH-0710.
- **No `after` on S-0514:** right. S-0532 already waits for S-0514, and S-0534 waits for S-0532.
- **Assumptions 1 to 7:** they read ADR-0209 §5 and the criteria, and add no command, key, or ADR.

Your two proposals are scope changes, so they are the operator's. My recommendation follows, and the thread stays open for it.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-28T21:29:09Z orchestrator (recommendation)
**Recommended:** take neither proposal into S-0534. Build it as planned.

- **Clearing the state by hand** (`flai serve budget clear <provider>`): the criteria do not ask for it. The probe of assumption 3 already clears the state once `until` passes. If you switch accounts or raise a limit, waiting for the reset or the backoff costs at most that wait. If that turns out to hurt, it becomes a story of its own under E-0030 once S-0534 lands.
- **Holding the permission check** (`claudecheck.go`): it runs once per new Claude Code version, not per start. A refusal there opens one thread, as TH-0669 shows. Holding it too would widen S-0534 beyond ADR-0209 §5, which names a story's agent and a role's run.

**Or:** name either one here, and the planner adds it to S-0534 as a task.

Source: design/adrs/0209-the-orchestrator-runs-for-a-reason-and-ends-started-by-flai-serve-on-the-events.md
