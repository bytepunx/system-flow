---
id: S-0229
type: story
nature: feature
title: "The strategic agents' settings are edited in the dashboard: permissions, policy, release policy, schedules, and agents"
status: ready
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:18Z
updated: 2026-10-06T19:35:23Z
transitions:
  - to: ready
    at: 2026-10-05T06:14:02Z
    by: alex
tags: [dashboard, flai]
touches: [flaiover/src/routes, flai/internal/hostapi, flai/internal/manifest, design/system/flaiover-dashboard.md, docs/users/flaiover.md, flaiover/src/routes/settings, flaiover/src/lib/components, design/system/project-manifest.md, docs/operators/settings.md, flai/cmd/manifest.go, flai/cmd/manifest_test.go, flai/cmd/root.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/cmd/guard_permissions_change_test.go, flai/internal/guard, flai/internal/harness, flaiover/src/lib/settings.ts, flaiover/src/lib/server/agent.ts, design/system/strategic-agents.md, design/system/flai-cli.md, design/adrs, docs/operators/index.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [S-0218, S-0211, S-0223]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 84.52
  by: planner-S-0229
  at: 2026-10-05T05:45:18Z
forecast:
  duration: 1h15m
  delivery: 2026-10-07T01:01:00Z
  basis: "Its own forecast of 1h15m; 7th in the pull order with an in-progress limit of 3, behind S-0296, S-0284, S-0278, S-0223, S-0224 and S-0227."
  by: flai
  at: 2026-10-06T19:35:23Z
---
# S-0229 The strategic agents' settings are edited in the dashboard: permissions, policy, release policy, schedules, and agents

## Goal

The operator decides how much freedom each strategic agent has. The manifest holds the settings; the dashboard should edit them where the agents' pages are, behind the `settings` host action, with each setting explained.

## Acceptance criteria
- [ ] The orchestrator page has a settings panel for `orchestration.permissions` (each with a sentence on what it allows and its risk), `orchestration.policy`, and `orchestration.release` (kind, thresholds, theme, whole epics); the planner page for `planning` (agent, replan, schedule, hour rate, cycle, default duration); the analyzer page for `analysis` (agent, schedule)
- [ ] Saving writes the manifest through a hostapi write (`settings.manifest`) that validates and refuses with the field and reason; a running orchestrator picks up changed permissions on its next decision
- [ ] The panels are read-only with the reason when the `settings` host action is off
- [ ] `design/system/flaiover-dashboard.md`, `project-manifest.md`, and the operator guide describe every setting; tests cover a save, a refusal, and read-only

## Tasks
- T-0952 flai manifest set writes the strategic agents' settings to system-flow.yaml, validating each and refusing with the field and the reason
- T-0954 A running orchestrator is held to changed permissions on its next call, and its prompt says they may change while it runs
- T-0961 The host API's settings.manifest write runs flai manifest set, and settings.get gives the strategic agents' settings with what each means
- T-0963 A StrategicSettings panel edits one block of the manifest through /api/settings, shows a refusal on its field, and is read-only with the reason while settings is off
- T-0964 The orchestrator, planner, and analyzer pages each show their settings panel, and the Settings page points to them
- T-0966 An ADR, the dashboard and manifest design, flai-cli.md, and the operator and user guides describe editing the strategic agents' settings from the dashboard

## Notes

### Planning

Planned by planner-S-0229 on 2026-10-05.

Touches, by where each came from:

- Declared, kept: `flaiover/src/routes`, `flaiover/src/routes/settings`, `flaiover/src/lib/components`, `flai/internal/hostapi`, `flai/internal/manifest`, `design/system/flaiover-dashboard.md`, `design/system/project-manifest.md`, `docs/users/flaiover.md`, `docs/operators/settings.md`.
- Co-change (`flai touches suggest`): `design/system/flai-cli.md` (30%), `docs/users/flai.md` (28%), `docs/operators/index.md` (18%, the `settings` host action's section), `docs/users/flai-reference.md` (13%, regenerated for the new command), `flaiover/src/lib/server/agent.ts` (10%, `REQUIRED_METHODS`, which `contract_test.go` holds equal to flai's methods), `flai/cmd/serve_actions.go` and its test (7%, where `settings.get` is built).
- Design: `design/adrs` (ADR-0039 lists what the `settings` action lets a dashboard change; widening it to the manifest's strategic keys is a new decision), `design/system/strategic-agents.md` ("On the settings page" says the planning keys are set by hand), `flai/internal/guard` and `flai/internal/harness` (the second criterion: a running orchestrator held to changed permissions, S-0218's guard and prompt).
- Layout: `flai/cmd/manifest.go`, `flai/cmd/manifest_test.go`, and `flai/cmd/root.go`, a new command, since ADR-0039 has every dashboard write run a flai command; `flai/cmd/guard_permissions_change_test.go`; `flaiover/src/lib/settings.ts` (`SETTINGS_KINDS`).
- Left out: suggestions from unrelated work (`design/issues/summary.md`, `design/adrs/README.md`, `flai/internal/check/check.go`, `flai/internal/serve/agents.go`, `flai/internal/mcpserver/server.go`). The manifest's new keys and their `flai check` findings are S-0217's, S-0218's, S-0222's, and S-0223's.

Forecast: `flai forecast` gives 1h, 132 s per unit over 14 done large feature stories on claude-opus-5-5, times size 27 (4 criteria, 23 touches at the time). Raised to 1h15m: the closest done stories, which likewise add a flai write and its dashboard form, took 72m (S-0204) and 63m (S-0201), and this one adds a command, a host API write and read, a guard test, a shared panel on three pages, and an ADR. The delivery is flai's play-out of the board with that duration.

Cost of delay: 84.52 USD a week, as `flai cod` gives it: S-0229 has no inputs of its own, so it takes its share of E-0016's 1500 USD a week from the operator's inputs, 1h15m of the 22h11m forecast over the epic's 17 open stories without inputs. It replaces 97.56, which was worked out from the 2h forecast before this plan.

Tasks: six, in five layers: T-0952 and T-0954, then T-0961, T-0963, T-0964, and T-0966. The plan thread on this story gives the assumptions and a proposal for `after`.
