---
id: S-0425
type: story
nature: feature
title: metrics.md generalizes strategic usage per kind to usage per role by an ADR, so a project's own role is charged on its items, reported by flai stats, and charted as the built-in roles are
status: backlog
parent: E-0024
owner: alex
created: 2026-10-16T05:41:28Z
updated: 2026-10-29T01:57:03Z
transitions: []
tags: [cli, dashboard]
topics: [cli, dashboard, planning, orchestration, template]
touches: [design/adrs, design/adrs/README.md, design/system/flaiover-dashboard.md, design/system/metrics.md, design/system/work-hierarchy.md, design/system/strategic-agents.md, design/system/continuous-improvement.md, design/system/flai-cli.md, design/tech/charts.md, docs/users/flai.md, docs/users/flaiover.md, design/workflows/roles/planner.md, design/workflows/roles/orchestrator.md, design/workflows/roles/analyzer.md, template/root/design/workflows/roles/planner.md, template/root/design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/analyzer.md, template/CHANGELOG.md, flai/cmd/activity.go, flai/cmd/activity_test.go, flai/cmd/stats.go, flai/cmd/check_stats_test.go, flai/internal/check/check.go, flai/internal/check/roles_test.go, flai/internal/issues/usage.go, flai/internal/issues/usage_test.go, flai/internal/mcpserver/activity.go, flai/internal/mcpserver/activity_test.go, flai/internal/metrics/strategic.go, flai/internal/metrics/strategic_test.go, flai/internal/metrics/usage.go, flai/internal/metrics/usage_test.go, flai/internal/roles/roles.go, flai/internal/roles/roles_test.go, flai/internal/serve/activity.go, flai/internal/serve/activity_test.go, flai/internal/statsread/statsread.go, flai/internal/usage/usage.go, flai/internal/usage/usage_test.go, flai/internal/workitem/usage.go, flai/internal/workitem/usage_test.go, flaiover/src/lib/components/RolePage.svelte, flaiover/src/lib/components/RolePage.svelte.test.ts, flaiover/src/lib/strategic.test.ts, flaiover/src/lib/strategic.ts, flaiover/src/lib/usage.test.ts, flaiover/src/lib/usage.ts, flaiover/src/lib/viz/charts.test.ts, flaiover/src/lib/viz/charts.ts, "flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts"]
after: [S-0423, S-0424]
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
      seconds: 2244
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 67
          output: 735
          cache_read: 4706101
          cache_write: 100574
          cost: 1.0572
draft: true
cost_of_delay:
  value: 62.36
  by: planner-S-0425
  at: 2026-10-23T15:35:46Z
forecast:
  duration: 1h30m
  delivery: 2026-10-29T12:29:00Z
  basis: "Its own forecast of 1h30m; 14th in the pull order with an in-progress limit of 5, behind S-0519, S-0532, S-0537, S-0534, S-0536, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412 and S-0424."
  by: flai
  at: 2026-10-29T01:57:03Z
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
- T-2073 A role file names its charge rule, item, named, or issues, the baseline role files carry theirs, and flai check reports an unknown rule

## Notes

- Drafted by the planner for E-0024. Last story of the epic, after the baseline formulas, as the epic's order says; it could start as soon as the dashboard's role page is done, and the plan thread proposes that.
- `metrics.md` is the contract between `flai stats` and the dashboard; it changes only with this ADR, as `CLAUDE.md` asks.

### Planning

- **Plan thread.** Replanned 2026-10-23 on the story's own thread, after the code was read. One task added, T-2073: the role file S-0414 defines has no field for the charge rule criterion 1 names, so the field, the three baseline role files, and its check are a task of their own.
- **Touches.** 51, all files but `design/adrs`.
  - **Declared, kept:** the 26 earlier touches, from the five earlier tasks.
  - **From layout, added (code read 2026-10-23):**
    - `flai/internal/issues/usage.go` and its test on T-1791: `ChargeStrategic` there refuses any kind but the three, as `workitem/usage.go` does.
    - `flai/cmd/activity.go` and its test on T-1791: `chargedTo` says what an activity charged with a branch per kind.
    - `flai/internal/mcpserver/activity.go` and its test on T-1791: `activity_log`'s description states the three rules by kind.
    - `flai/cmd/stats.go` and `flai/cmd/check_stats_test.go` on T-1792: the printed labels, and the stats tests criterion 3 keeps unchanged, beside a new one with a fixture role.
    - `flaiover/src/routes/charts/[kind]/+page.svelte` and its test on T-1793: the table view's columns come from `STRATEGIC_AGENTS` in `viz/charts.ts`, which also filters a role outside the list out of the bars.
    - `flai/internal/roles/roles.go`, `flai/internal/check/check.go`, and their tests on T-2073, for the `charge` field.
  - **From design, added:** `design/system/work-hierarchy.md` on T-1790, which defines the `strategic` entry with `metrics.md` (ADR-0083); the six baseline role files and `template/CHANGELOG.md` on T-2073; `design/system/strategic-agents.md` on T-1794, whose sections on what each run costs describe the charges; `design/system/continuous-improvement.md` § Issue file, which says one entry per kind; `design/tech/charts.md` § Palette and rules, for the slot a role beyond the three takes.
  - **Co-change, kept:** `design/system/flai-cli.md` (48%), shared, so it holds nothing.
  - **Co-change, not added:** `docs/operators/index.md` and `settings.md` (no setting changes), `docs/users/flai-reference.md` (no flag changes), `workflow.md`, `work-management.md`, `project-manifest.md`, and `agent-narrative.md`, which S-0417 updates for defined roles' activity documents.
  - **Left to earlier stories:** `workitem/activity.go`'s kind gate, `serve/agents.go`, and `hostapi/people.go`, which S-0417 (T-1737) opens to any role's name; the dashboard's role page and `strategic.ts`'s run types, which S-0420 builds; the three roles' start path, which S-0419 moves.
  - **Folder touch kept:** `design/adrs`, because the ADR's file name is known only once `flai adr new` makes it. The folder is in `claims.shared`, so it holds no story.
- **Forecast.** 1h30m, cut from flai's 1h41m (107 s a unit over 85 done feature stories on this model in the large band, times size 56). Fourteen of the touches are one-line edits: a `charge:` line in each of six role files, a CHANGELOG line, and a sentence each in seven documents. The Go charge dispatch and the dashboard's series per role are the work, and they are the size of S-0417's or S-0420's. The delivery is flai's play-out with the new duration: S-0424's 08:22Z plus 1h30m times the cycle factor 1.12, 2026-10-24T10:02Z.
- **Cost of delay.** 62.36 USD a week, worked out by hand from TH-0471's answer 1, as the planners of S-0414 to S-0420 did, because `flai cod` refuses while E-0024 has no inputs.
  - 4h lost per cycle at `planning.hour_rate` 150 over the default 168h cycle gives E-0024 600 USD a week.
  - S-0425 takes the share its 1h30m is of the 14h26m forecast for E-0024's twelve open stories: 90 / 866 × 600.
  - The input the operator chose is still not on E-0024; the plan thread asks for it. Once it is set, `flai cod S-0425` gives the same figure while the forecasts stand.
