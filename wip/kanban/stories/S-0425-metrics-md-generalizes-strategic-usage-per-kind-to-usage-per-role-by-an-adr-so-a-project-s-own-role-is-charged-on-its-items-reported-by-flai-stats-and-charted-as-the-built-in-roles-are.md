---
id: S-0425
type: story
nature: feature
title: metrics.md generalizes strategic usage per kind to usage per role by an ADR, so a project's own role is charged on its items, reported by flai stats, and charted as the built-in roles are
status: backlog
parent: E-0024
owner: alex
created: 2026-10-16T05:41:28Z
updated: 2026-10-17T01:56:07Z
transitions: []
tags: [cli, dashboard]
topics: [cli, dashboard, planning, orchestration]
touches: [design/adrs, design/adrs/README.md, design/system/flaiover-dashboard.md, design/system/metrics.md, docs/users/flai.md, docs/users/flaiover.md, flai/internal/metrics/strategic.go, flai/internal/metrics/strategic_test.go, flai/internal/metrics/usage.go, flai/internal/metrics/usage_test.go, flai/internal/serve/activity.go, flai/internal/serve/activity_test.go, flai/internal/statsread/statsread.go, flai/internal/usage/usage.go, flai/internal/usage/usage_test.go, flai/internal/workitem/usage.go, flai/internal/workitem/usage_test.go, flaiover/src/lib/components/RolePage.svelte, flaiover/src/lib/components/RolePage.svelte.test.ts, flaiover/src/lib/strategic.test.ts, flaiover/src/lib/strategic.ts, flaiover/src/lib/usage.test.ts, flaiover/src/lib/usage.ts, flaiover/src/lib/viz/charts.test.ts, flaiover/src/lib/viz/charts.ts, design/system/flai-cli.md]
after: [S-0423, S-0424]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 50m
  delivery: 2026-10-17T22:08:00Z
  basis: "Its own forecast of 50m; 39th in the pull order with an in-progress limit of 5, behind S-0440, S-0441, S-0442, S-0444, S-0443, S-0241, S-0392, S-0393, S-0394, S-0395, S-0396, S-0397, S-0398, S-0399, S-0400, S-0401, S-0402, S-0403, S-0404, S-0405, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422, S-0423 and S-0424."
  by: flai
  at: 2026-10-17T01:56:07Z
---
# S-0425 metrics.md generalizes strategic usage per kind to usage per role by an ADR, so a project's own role is charged on its items, reported by flai stats, and charted as the built-in roles are

## Goal

What a defined role spends counts as the built-in roles' spend does. An ADR generalizes ADR-0083, ADR-0095, and ADR-0100 from three kinds to any role: a role's file says how its activity is charged, to the item its run was started for, evenly to the work items it names, or evenly to the issues it names, and `usage.strategic` takes the role's name as its kind. `flai stats` reports each role's totals, what its items carry, and its project total, and the dashboard's strategic charts draw a series per role.

## Acceptance criteria

- [ ] An ADR, refining ADR-0083, ADR-0095, and ADR-0100, defines strategic usage per role: the kinds are the roles a project defines, each charged by the rule its role file names (`item`, `named`, or `issues`), the planner's, the orchestrator's, and the analyzer's being the rules they have today; `design/system/metrics.md` § Strategic usage and § Strategic agents say the same.
- [ ] A defined role's activity, logged at its run's end or through `activity_log`, is charged by its rule under its name in `usage.strategic` on the items or issues it reaches and every item above them, and what no item carries stays in its project strategic total.
- [ ] `flai stats --json` reports `strategic` for every role a project defines, with its document's totals, `items`, `issues`, and `project`, and strategic use per day by role; the existing stats tests for the three kinds pass unchanged.
- [ ] The dashboard's Strategic Cost and Strategic Use charts, and each role's page, draw a series per role from `flai stats`, shown by a component test with a fixture role.
- [ ] `docs/users/flai.md`, `docs/users/flaiover.md`, and `design/system/flaiover-dashboard.md` § Views say that strategic usage is per role.

## Tasks
- T-1790 An ADR refining ADR-0083, ADR-0095, and ADR-0100 defines strategic usage per role, and metrics.md says the same
- T-1791 A defined role's activity is charged by its file's rule under its name in usage.strategic, on the items or issues it reaches and every item above them
- T-1792 flai stats reports strategic totals, items, issues, project, and use per day for every role a project defines
- T-1793 The dashboard's Strategic Cost and Strategic Use charts and each role page draw a series per role from flai stats
- T-1794 The user guides and the dashboard's design say that strategic usage is charged, reported, and charted per role

## Notes

- Drafted by the planner for E-0024. Last story of the epic, after the baseline formulas, as the epic's order says; it could start as soon as the dashboard's role page is done, and the plan thread proposes that.
- `metrics.md` is the contract between `flai stats` and the dashboard; it changes only with this ADR, as `CLAUDE.md` asks.

### Planning

- **Touches.** The story declared none. Every touch comes from its five tasks.
  - **From layout:** serve's `activity.go`, `workitem/usage.go`, and `usage/usage.go`, each with its test; `metrics/strategic.go` and `usage.go` with their tests; `statsread.go`; and in `flaiover/src/lib`, `strategic.ts`, `usage.ts`, `viz/charts.ts`, and `RolePage.svelte`, each with its test.
  - **From design:** the ADR's index, `metrics.md`, `flaiover-dashboard.md`, and the user guides.
  - **Co-change, added:** `design/system/flai-cli.md` (48%), which describes `flai stats`.
  - **Co-change, not added:** the operator docs (no setting changes) and `agent-narrative.md`, which S-0417 updates.
  - **Folder touch kept:** `design/adrs`, for the ADR's new file; it is shared.
- **Forecast.** Kept at flai's 50m. The forecast's basis gives the figures, and its delivery is flai's play-out moved by any raise times the cycle factor.
- **Cost of delay.** Not set. Neither the story nor E-0024 has inputs, so `flai cod` refuses. The epic's input is asked on TH-0471, question 1, and the value follows its answer.
