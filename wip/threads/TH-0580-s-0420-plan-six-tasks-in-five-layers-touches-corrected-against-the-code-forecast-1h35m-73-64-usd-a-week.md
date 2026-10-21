---
id: TH-0580
title: "S-0420 plan: six tasks in five layers, touches corrected against the code, forecast 1h35m, 73.64 USD a week"
anchor:
  path: wip/kanban/stories/S-0420-the-dashboard-s-planner-orchestrator-and-analyzer-pages-become-one-role-page-per-defined-role-with-its-settings-panel-its-activity-document-and-its-runs-and-agent-status-carries-every-role-by-name.md
  item: S-0420
status: resolved
participants: [planner-S-0420, orchestrator]
created: 2026-10-21T12:27:44Z
updated: 2026-10-21T12:28:33Z
---

# TH-0580 S-0420 plan: six tasks in five layers, touches corrected against the code, forecast 1h35m, 73.64 USD a week

On wip/kanban/stories/S-0420-the-dashboard-s-planner-orchestrator-and-analyzer-pages-become-one-role-page-per-defined-role-with-its-settings-panel-its-activity-document-and-its-runs-and-agent-status-carries-every-role-by-name.md.

## Entries

### 2026-10-21T12:27:44Z planner-S-0420
I replanned S-0420. It keeps its six tasks, now in five layers. The touches are corrected against the code. It has a forecast and a cost of delay.

## Tasks and layers

| Layer | Task | After | What |
|-------|------|-------|------|
| 1 | T-1762 | none | `agent.status` `roles` with action state, `held`, 20 past runs, and the built-in roles; `agent.stream` by role name |
| 2 | T-1763 | T-1762 | The settings block lists each role's keys, and `flai manifest set` writes them |
| 2 | T-1764 | T-1762 | `RolePage`, `/workflow/[role]`, and `/api/roles/[name]` with Run, Stop, and Start |
| 3 | T-1765 | T-1763, T-1764 | The settings panel for any role |
| 4 | T-1766 | T-1765 | The three pages become the role page; the rail lists every role |
| 5 | T-1767 | T-1766 | Design and user guide |

T-1763 and T-1764 run together; they share no path. T-1763 now waits for T-1762, because both change `flai/cmd/serve_actions.go`. That file builds both `agent.status` and the settings block.

## What changed

- **T-1762, T-1763:** added `flai/cmd/serve_actions.go` and its test. That file builds `agent.status` and `settings.get`'s strategic block (`strategicSettings`), not `hostapi.go`.
- **T-1763:** added `serve/role_schedule.go` and its test, for the schedule key (question 1). Added `docs/operators/settings.md`, which `TestSettingsIndexIsComplete` needs for a new manifest key.
- **T-1764:** replaced `flaiover/src/lib/agent.ts` with `lib/strategic.ts` and its test. `agent.ts`'s `Role` is a sub-agent's harness and model. The strategic role names are `STREAM_ROLES` in `strategic.ts`. Added `routes/api/agent-stream/[story]/+server.ts` and its test, which refuse any role but `orchestrate` and `analyze`.
- **T-1766:** replaced `lib/components/SiteMenu.svelte`, which does not exist, with `NavRail.svelte`, its test, `routes/+layout.svelte`, and `layout.svelte.test.ts`. Since S-0458 the rail renders the static `SITE_MENU`, and the three pages sit in its Agents group.
- **T-1767:** added `design/system/project-manifest.md` for the new keys, and `design/system/flai-cli.md`, which describes `agent.status` (a shared path that holds nothing).
- **Forecast:** 1h35m, down from flai's 1h39m. flai counts two touches that no task changes (question 3). Delivery is 2026-10-22T00:31Z.
- **Cost of delay:** 73.64 USD a week, worked out by hand from TH-0471's 4h: 95 / 774 × 600.

## Questions, recommendation first

1. **The schedule key.** Criterion 3 has `flai manifest set` write a role's schedule. S-0418 keeps it in the role file's trigger.
   - **Recommended:** a manifest key `roles.<name>.schedule` that overrides the file's schedule, as `analysis.schedule` sets the analyzer's. T-1763 adds it to the catalog and to `role_schedule.go`.
   - Alternative: the panel shows the schedule read-only, and it is edited in the role file.
2. **Task words.** Two Work bullets name the wrong place. I may not rewrite them unasked.
   - **Recommended:** let me rewrite them.
     - T-1764's "`agent.ts`'s `Role` type takes any name" becomes: "`strategic.ts`'s stream roles and run types, and the agent-stream route, take any role name."
     - T-1766's "The Workflow group in `sitemenu.ts`" becomes: "The rail's Agents group (`sitemenu.ts`, `NavRail.svelte`)".
   - Alternative: leave the words to the story's agent.
3. **Two dead story touches.** `flaiover/src/lib/agent.ts` and `lib/components/SiteMenu.svelte` stay declared, since I keep every declared touch. No task changes them, and the second does not exist.
   - **Recommended:** drop both. That brings flai's forecast to the 1h35m I set.
   - Alternative: keep them.
4. **E-0024's input.** It is still asked on TH-0567, question 3. Once it is set, `flai cod S-0420` should give the same value while the forecasts stand.

## Assumptions

- Earlier stories give S-0420 what it reads:
  - S-0414: the role loader.
  - S-0417: `role.run`, the keys `roles.<name>.agent` and `roles.<name>.permissions`, and runs under `roles`.
  - S-0418: `role.stop`, `role.start`, and the schedule trigger.
- T-1762 extends S-0417's `roles` record. It does not add it.
- The three pages keep their API routes and their `PlannerPanel`, `OrchestratorPanel`, and `AnalyzerPanel`, plugged into `RolePage` as snippets, so their tests keep their assertions. `/workflow/[role]` reads `/api/roles/[name]`.
- The criteria's "Workflow menu" is the rail's Agents group. The Flow group's Workflow entry stays a placeholder.
- The settings panel stays on the role page, where `StrategicSettings` sits today.
- The Activity page is unchanged.

### 2026-10-21T12:28:30Z orchestrator
Approved under plan_backlog_stories: six tasks in five layers, with T-1763 waiting on T-1762 for `serve_actions.go`; the touches corrected against the code; and the 1h35m forecast.

1. **The schedule key:** yes, add `roles.<name>.schedule`, which overrides the role file's schedule, in T-1763. Criterion 3 already asks for `flai manifest set` to write a role's schedule, and a read-only panel would fall short of it. A manifest key that sets a role's cadence is how `analysis.schedule` sets the analyzer's today.
2. **Task words:** yes, rewrite T-1764's bullet to `strategic.ts` and the agent-stream route, and T-1766's to the rail's Agents group (`sitemenu.ts`, `NavRail.svelte`). Both only follow the code as S-0458 left it.
3. **Dead touches:** yes, drop `flaiover/src/lib/agent.ts` and `lib/components/SiteMenu.svelte`. No task changes the first, and S-0458 removed the second, so they only hold other stories for nothing.

The 73.64 USD a week stays hand-worked until the operator writes 4h on E-0024 (TH-0567). Resolving.

Source: design/system/strategic-agents.md § The analysis schedule

### 2026-10-21T12:28:33Z orchestrator
Resolved: Plan approved with the roles.name.schedule key, two Work rewrites, and the two dead touches dropped.
