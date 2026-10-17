---
id: S-0420
type: story
nature: feature
title: The dashboard's Planner, Orchestrator, and Analyzer pages become one role page per defined role, with its settings panel, its activity document, and its runs, and agent.status carries every role by name
status: backlog
parent: E-0024
owner: alex
created: 2026-10-16T05:40:31Z
updated: 2026-10-17T01:56:07Z
transitions: []
tags: [dashboard, cli]
topics: [dashboard, planning, orchestration]
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md, flai/internal/hostapi/contract_test.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/hostapi_test.go, flai/internal/hostapi/settings.go, flai/internal/hostapi/settings_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/serve/agents.go, flai/internal/serve/stream.go, flai/internal/serve/stream_test.go, flaiover/src/lib/agent.ts, flaiover/src/lib/components/RolePage.svelte, flaiover/src/lib/components/RolePage.svelte.test.ts, flaiover/src/lib/components/SiteMenu.svelte, flaiover/src/lib/components/StrategicAgentPanel.svelte, flaiover/src/lib/components/StrategicSettings.svelte, flaiover/src/lib/components/StrategicSettings.svelte.test.ts, flaiover/src/lib/roles.test.ts, flaiover/src/lib/roles.ts, flaiover/src/lib/settings.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/sitemenu.ts, "flaiover/src/routes/api/roles/[name]/roles.test.ts", "flaiover/src/routes/api/roles/[name]/+server.ts", flaiover/src/routes/workflow/analyzer/analyzer.svelte.test.ts, flaiover/src/routes/workflow/analyzer/+page.svelte, flaiover/src/routes/workflow/orchestrator/orchestrator.svelte.test.ts, flaiover/src/routes/workflow/orchestrator/+page.svelte, flaiover/src/routes/workflow/planner/+page.svelte, flaiover/src/routes/workflow/planner/planner.svelte.test.ts, "flaiover/src/routes/workflow/[role]/+page.svelte", flai/internal/serve/agents_test.go, flaiover/src/lib/server/agent.ts]
after: [S-0419]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 1h4m
  delivery: 2026-10-17T18:03:00Z
  basis: "Its own forecast of 1h4m; 34th in the pull order with an in-progress limit of 5, behind S-0440, S-0441, S-0442, S-0444, S-0443, S-0241, S-0392, S-0393, S-0394, S-0395, S-0396, S-0397, S-0398, S-0399, S-0400, S-0401, S-0402, S-0403, S-0404, S-0405, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418 and S-0419."
  by: flai
  at: 2026-10-17T01:56:07Z
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

- **Touches.** The story declared none. Every touch comes from its six tasks.
  - **From layout:** the host API's `hostapi.go`, `settings.go`, and `contract_test.go`; serve's `stream.go` and `agents.go`; the manifest's settings table; and in `flaiover/src`, a new `RolePage.svelte`, `roles.ts`, the `/workflow/[role]` route, and the `/api/roles/[name]` route, the three workflow pages and their tests, `sitemenu.ts`, `SiteMenu.svelte`, `StrategicSettings.svelte`, `StrategicAgentPanel.svelte`, `settings.ts`, and `agent.ts`.
  - **From design:** `flaiover-dashboard.md` and `docs/users/flaiover.md`.
  - **Co-change, added:** `flaiover/src/lib/server/agent.ts` (9%), which carries the agent stream a role page reads. Also `flai/internal/serve/agents_test.go` (9%), for `roles` in the agent state.
  - **Co-change, not added:** the CLI docs and `serve_actions.go`. No command changes here.
  - No folder touch.
- **Forecast.** Kept at flai's 1h4m. The forecast's basis gives the figures, and its delivery is flai's play-out moved by any raise times the cycle factor.
- **Cost of delay.** Not set. Neither the story nor E-0024 has inputs, so `flai cod` refuses. The epic's input is asked on TH-0471, question 1, and the value follows its answer.
