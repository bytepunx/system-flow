---
id: S-0420
type: story
nature: feature
title: The dashboard's Planner, Orchestrator, and Analyzer pages become one role page per defined role, with its settings panel, its activity document, and its runs, and agent.status carries every role by name
status: backlog
parent: E-0024
owner: alex
created: 2026-10-16T05:40:31Z
updated: 2026-10-25T21:20:31Z
transitions: []
tags: [dashboard, cli]
topics: [dashboard, planning, orchestration]
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md, flai/internal/hostapi/contract_test.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/hostapi_test.go, flai/internal/hostapi/settings.go, flai/internal/hostapi/settings_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/serve/agents.go, flai/internal/serve/stream.go, flai/internal/serve/stream_test.go, flaiover/src/lib/components/RolePage.svelte, flaiover/src/lib/components/RolePage.svelte.test.ts, flaiover/src/lib/components/StrategicAgentPanel.svelte, flaiover/src/lib/components/StrategicSettings.svelte, flaiover/src/lib/components/StrategicSettings.svelte.test.ts, flaiover/src/lib/roles.test.ts, flaiover/src/lib/roles.ts, flaiover/src/lib/settings.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/sitemenu.ts, "flaiover/src/routes/api/roles/[name]/roles.test.ts", "flaiover/src/routes/api/roles/[name]/+server.ts", flaiover/src/routes/workflow/analyzer/analyzer.svelte.test.ts, flaiover/src/routes/workflow/analyzer/+page.svelte, flaiover/src/routes/workflow/orchestrator/orchestrator.svelte.test.ts, flaiover/src/routes/workflow/orchestrator/+page.svelte, flaiover/src/routes/workflow/planner/+page.svelte, flaiover/src/routes/workflow/planner/planner.svelte.test.ts, "flaiover/src/routes/workflow/[role]/+page.svelte", flai/internal/serve/agents_test.go, flaiover/src/lib/server/agent.ts, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/serve/role_schedule.go, flai/internal/serve/role_schedule_test.go, docs/operators/settings.md, flaiover/src/lib/strategic.ts, flaiover/src/lib/strategic.test.ts, "flaiover/src/routes/api/agent-stream/[story]/+server.ts", "flaiover/src/routes/api/agent-stream/[story]/agent-stream.test.ts", flaiover/src/lib/components/NavRail.svelte, flaiover/src/lib/components/NavRail.svelte.test.ts, flaiover/src/routes/+layout.svelte, flaiover/src/routes/layout.svelte.test.ts, design/system/project-manifest.md, design/system/flai-cli.md]
after: [S-0419]
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
      seconds: 302
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 109
          output: 1938
          cache_read: 26318370
          cache_write: 25326
          cost: 5.8168
draft: true
cost_of_delay:
  value: 73.64
  by: planner-S-0420
  at: 2026-10-21T12:28:25Z
forecast:
  duration: 1h35m
  delivery: 2026-10-26T09:31:00Z
  basis: "Its own forecast of 1h35m; 15th in the pull order with an in-progress limit of 5, behind S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418 and S-0419."
  by: flai
  at: 2026-10-25T21:20:31Z
---
# S-0420 The dashboard's Planner, Orchestrator, and Analyzer pages become one role page per defined role, with its settings panel, its activity document, and its runs, and agent.status carries every role by name

## Goal

A role a project defines shows on the dashboard as the built-in ones do. The Workflow menu lists every role `agent.status` reports, and one page, `/workflow/<name>`, shows any of them: its host action's state, the run under way and its stream, its activity document, its runs with their cost and outcome, the controls its trigger allows, and its settings panel. The Planner, Orchestrator, and Analyzer pages become that page for their roles, at the same addresses and with what they show today.

## Acceptance criteria

- [ ] `agent.status` carries every role by name under `roles`, each with its action's state, its run under way, its past runs (up to 20), and `held`; `agent.stream` reads any role's stream by name; `plans`, `orchestrator`, `analyzer`, `past_orchestrators`, and `past_analyzers` stay for a dashboard older than this flai.
- [ ] `/workflow/<name>` shows a role's action state, the run under way with its stream, its activity log, its runs with cost and outcome, and **Run** (with an item when its file takes one), **Stop**, and **Start** as its trigger allows, each through `role.run`, `role.stop`, and `role.start`.
- [ ] Its settings panel edits `roles.<name>.agent`, the role's schedule when its trigger has one, and the permissions its file declares through `flai manifest set` while the `settings` host action is on, and shows them read-only otherwise; `flai manifest set` takes those keys.
- [ ] `/workflow/planner`, `/workflow/orchestrator`, and `/workflow/analyzer` are that page for the three roles and show what they show today, their component tests passing; the Workflow menu lists each defined role.
- [ ] Component and API tests cover a defined role's page from a fixture `agent.status`, and `design/system/flaiover-dashboard.md` § Strategic agents' pages and settings, and `docs/users/flaiover.md`, describe the role page.

## Tasks
- T-1762 agent.status carries every role by name under roles with its action, run, past runs, and hold, and agent.stream reads any role's stream by name
- T-1763 The host's settings methods list and write a role's agent, schedule, and declared permissions through flai manifest set
- T-1764 A RolePage component and the /workflow/[role] route show any role's action, run and stream, activity log, runs, and its Run, Stop, and Start
- T-1765 The role page's settings panel edits a role's agent, schedule, and declared permissions while settings is on and shows them read-only otherwise
- T-1766 The Planner, Orchestrator, and Analyzer pages become the role page for their roles at the same addresses, and the Workflow menu lists every defined role
- T-1767 The dashboard's design and the user guide describe the role page, its controls, and its settings panel

## Notes

- Drafted by the planner for E-0024. Waits for the story that moves the three roles onto the defined-role path, so that every role reaches `agent.status` the same way.
- The strategic cost and use charts by role are the metrics story's.

### Planning

- **Plan thread.** TH-0580, replanned 2026-10-21, approved and resolved by the orchestrator. It settled three things:
  - `roles.<name>.schedule` is a manifest key that overrides the role file's schedule, as `analysis.schedule` sets the analyzer's. T-1763 adds it.
  - T-1764's and T-1766's Work bullets were rewritten to `strategic.ts` with the agent-stream route, and to the rail's Agents group.
  - `flaiover/src/lib/agent.ts` and `flaiover/src/lib/components/SiteMenu.svelte` were dropped from the touches.
- **Reading of the criteria.** The "Workflow menu" is the rail's Agents group, where S-0458 put the three pages. The Flow group's Workflow entry stays a placeholder.
- **Layers.** Five layers. T-1762 runs first. T-1763 and T-1764 run together; they share no path. T-1765, T-1766, and T-1767 follow in turn. T-1763 waits for T-1762, because both change `flai/cmd/serve_actions.go`.
- **What earlier stories give.**
  - S-0414 gives the role loader.
  - S-0417 gives `role.run`, the `roles.<name>.agent` and `roles.<name>.permissions` keys, and runs under `roles` in `agent.status`.
  - S-0418 gives `role.stop`, `role.start`, and the schedule trigger.
  - T-1762 extends S-0417's `roles` record with the action state, `held`, 20 past runs, and the built-in roles.
  - The three pages keep their API routes and their `PlannerPanel`, `OrchestratorPanel`, and `AnalyzerPanel`, plugged into `RolePage`, so their tests keep their assertions.
- **Touches.** 48, all files; no folder touch.
  - **Declared, kept:** 33 of the 35 earlier touches. The two dropped, on TH-0580, were `lib/agent.ts`, whose `Role` is a sub-agent's harness and model, and `SiteMenu.svelte`, which does not exist.
  - **From layout, added (code read 2026-10-21):**
    - `flai/cmd/serve_actions.go` and its test, on T-1762 and T-1763. It builds `agent.status`'s state and `settings.get`'s strategic block.
    - `flai/internal/serve/role_schedule.go` and its test, on T-1763, for `roles.<name>.schedule`.
    - `flaiover/src/lib/strategic.ts` and its test, on T-1764. They hold `STREAM_ROLES` and the run types.
    - `flaiover/src/routes/api/agent-stream/[story]/+server.ts` and its test, on T-1764. The route refuses any role but `orchestrate` and `analyze`.
    - `flaiover/src/lib/components/NavRail.svelte` and its test, `flaiover/src/routes/+layout.svelte`, and `layout.svelte.test.ts`, on T-1766. Since S-0458 the rail renders `SITE_MENU`.
  - **From design, added:** `docs/operators/settings.md` on T-1763, which `TestSettingsIndexIsComplete` needs for the new key. `design/system/project-manifest.md` on T-1767, for the keys.
  - **Co-change, added:** `design/system/flai-cli.md` (30%) on T-1767. It describes `agent.status`'s fields, and it is a shared path, so it holds nothing.
  - **Co-change, not added:**
    - `flai/internal/hostapi/writes.go` and its test (16%), since `role.run`, `role.stop`, and `role.start` come from S-0417 and S-0418.
    - `docs/users/flai.md`, since no command changes.
    - The three pages' API routes under `flaiover/src/routes/api/`, since the pages keep them.
- **Forecast.** flai's 1h35m is kept: 107 s a unit over 85 done feature stories on this model in the large band, times size 53 (5 criteria and 48 touches). The delivery is flai's play-out, S-0419's delivery at 22:45Z plus 1h35m times the cycle factor 1.12: 2026-10-22T00:31Z.
- **Cost of delay.** 73.64 USD a week. I worked it out by hand from TH-0471's answer 1, as the planners of S-0417 to S-0419 did, because `flai cod` refuses while E-0024 has no inputs.
  - 4h lost per cycle, at `planning.hour_rate` 150 over the default 168h cycle, gives E-0024 600 USD a week.
  - S-0420's share is its 1h35m of the 12h54m forecast for E-0024's twelve open stories: 95 / 774 × 600.
  - Setting 4h on E-0024 is asked on TH-0567, question 3. Once it is set, `flai cod S-0420` gives the same figure while the forecasts stand.
