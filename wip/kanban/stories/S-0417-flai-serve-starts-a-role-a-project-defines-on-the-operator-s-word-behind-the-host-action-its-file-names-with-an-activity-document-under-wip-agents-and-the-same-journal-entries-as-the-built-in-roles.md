---
id: S-0417
type: story
nature: feature
title: flai serve starts a role a project defines on the operator's word, behind the host action its file names, with an activity document under wip/agents and the same journal entries as the built-in roles
status: backlog
parent: E-0024
owner: alex
created: 2026-10-16T05:40:01Z
updated: 2026-10-23T20:04:41Z
transitions: []
tags: [cli, template]
topics: [cli, planning, orchestration]
touches: [design/system/project-manifest.md, design/system/strategic-agents.md, docs/operators/index.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flai-reference.md, flai/cmd/activity.go, flai/cmd/activity_test.go, flai/cmd/role.go, flai/cmd/role_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/check/activity.go, flai/internal/check/activity_test.go, flai/internal/check/check.go, flai/internal/harness/adapters.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/hostapi_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/mcpserver/activity.go, flai/internal/mcpserver/activity_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/role.go, flai/internal/mcpserver/role_test.go, flai/internal/serve/activity.go, flai/internal/serve/agents.go, flai/internal/serve/role.go, flai/internal/serve/role_test.go, flai/internal/workitem/activity.go, flai/internal/workitem/activity_test.go, flai/internal/workitem/narrative.go, template/CHANGELOG.md, design/system/flai-cli.md, design/system/agent-narrative.md, flai/internal/serve/serve.go, flai/cmd/mcp.go, flai/cmd/mcp_http.go, flai/internal/hostapi/settings.go, flai/internal/hostapi/settings_test.go, flai/internal/hostapi/people.go, flaiover/src/lib/server/agent.ts, flai/cmd/guard.go, flai/internal/guard/roles.go, flai/internal/guard/roles_test.go]
after: [S-0415, S-0416]
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
      seconds: 507
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 69
          output: 1100
          cache_read: 10009086
          cache_write: 22137
          cost: 2.2151
draft: true
cost_of_delay:
  value: 79.4
  by: planner-S-0417
  at: 2026-10-21T11:53:36Z
forecast:
  duration: 1h37m
  delivery: 2026-10-24T03:44:00Z
  basis: "Its own forecast of 1h37m; 12th in the pull order with an in-progress limit of 5, behind S-0430, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415 and S-0416."
  by: flai
  at: 2026-10-23T20:04:41Z
---
# S-0417 flai serve starts a role a project defines on the operator's word, behind the host action its file names, with an activity document under wip/agents and the same journal entries as the built-in roles

## Goal

A role a project defines runs. `flai role run <name> [<item>]`, the MCP tool `role_run`, and the host API's `role.run` start it through `flai serve`, as `flai plan` starts the planner: behind the host action its file names, in the main checkout, with its prompt from its file's body, its harness definition, and its own agent. Its run is recorded and settled as the planner's is, its activity is logged in `wip/agents/<name>.md`, and its starts and ends are journal entries and dashboard notifications. A fixture role shows it end to end.

## Acceptance criteria

- [ ] `flai role run <name> [<item>] [--json]`, the MCP tool `role_run` `{name, item?}`, and the host API's `role.run` start a defined role through `serve`, and each is refused, saying why, while the host action its file names is off, for an item its file does not take or one archived, done, or cancelled, while a run of the role for the same item runs, and when its agent names no harness and no command is set on the host. `flai serve enable <action>` and the Settings page accept a defined role's action. The MCP tool is refused to an agent `flai serve` started, as `plan` is.
- [ ] The run starts in the main checkout through the harness adapter with the role's verb as `Role`, `FLAI_AGENT` set to `<name>` or `<name>-<item>`, `FLAI_ROLE`, `FLAI_ITEM` when it has one, `FLAI_SESSION`, and `FLAI_STARTED_BY=flai-serve`; its prompt is its role file's body with the item and the run's inputs; on `claude-code` it runs with `--agent <name>` over the definition flai generates from its role file and passes in `--agents`, writing nothing under `.claude/` (TH-0471, answer 6); its agent is `roles.<name>.agent` in the manifest merged over the project's `agent`, as `planning.agent` is.
- [ ] The run is recorded in `serve/agents.json` under `roles`, by name and item, with its log `<key>-<name>-<start>.log`, and settled `asked`, `failed`, or `worked` as a planner run is; `serve.LogRunEnd` logs its activity in `wip/agents/<name>.md`, which `wip/agents/index.md` lists and `flai check` validates; `activity_log` and `flai activity log` take the role's name as `kind`.
- [ ] Every start, failure, and end is a journal entry with the role's action, `agent.status` carries the runs under `roles` by name, and an `agent` notification names the role.
- [ ] A serve test starts a fixture role, `reviewer`, from a role file and a fake harness, and shows its environment, its prompt, its record, its activity entry, and each refusal.
- [ ] `design/system/strategic-agents.md`, `design/system/project-manifest.md` (`roles.<name>.agent`), `docs/users/flai.md` with `flai-reference.md`, `docs/operators/index.md` and `docs/operators/settings.md`, and `template/CHANGELOG.md` say how a defined role is started and where its activity is logged.

## Tasks
- T-1735 The manifest takes roles.name.agent and roles.name.permissions, merged and validated as planning.agent and orchestration.permissions are, and flai manifest set writes them
- T-1736 The harness starts a defined role with its prompt from its role file's body, its definition, and its environment, on claude-code with --agent over that definition
- T-1737 Activity documents, activity_log, and flai activity log take any defined role's name as their kind, and wip/agents/index.md lists each and flai check validates it
- T-1738 serve starts, records, and settles a defined role's run under roles in agents.json and logs its activity in wip/agents/name.md, shown end to end with a fixture role
- T-1739 flai role run, the MCP tool role_run, and the host API's role.run start a defined role, flai serve enable takes its action, and agent.status carries roles
- T-1740 The design, the manifest's design, and the user and operator docs say how a defined role is started, as whom, and where its activity is logged
- T-1983 flai guard gates a defined role's permission-bound calls on roles.name.permissions in the manifest, read at each call

## Notes

- Drafted by the planner for E-0024. Waits for the prime story and the guard story: a role is started only once its pack and its guard read its file, so it never runs unguarded.
- A defined role's activity charges no item until the metrics story generalizes `usage.strategic` per role: it stays in the role's project strategic total, which `flai stats` already works out per kind.
- Starting on an event, on a schedule, or for as long as the action is on is the next story's.

### Planning

- **TH-0471's answers (2026-10-21).** The plan follows them.
  - **Answer 4** confirms `roles.<name>.agent`, so T-1735 keeps it.
  - **Answer 6** has flai generate a role's harness definition from its role file. Criterion 2 now says the run takes the definition flai generates and passes in `--agents`, writing nothing under `.claude/`. T-1736's Work still reads `.claude/agents/<name>.md`; its rewrite is proposed on TH-0570.
  - **Answer 2's folder**, `design/workflows/roles/`, is read through S-0414's `roles.Load`, so no touch here names it.
  - **Answer 1's** 4h is not yet on E-0024 (TH-0567, question 3).
- **Touches.** 48, all files; no folder touch.
  - **Declared, kept:** the 41 earlier touches, from the six tasks.
  - **From layout, added (code read 2026-10-21):**
    - `flai/internal/serve/serve.go` on T-1738: the `agent` notification names only an item, the orchestrator, and the analyzer.
    - `flai/internal/serve/activity.go`, `flai/internal/serve/agents.go`, and `flai/internal/hostapi/people.go` on T-1737: they call `workitem.IsActivityKind`, which a defined role's kind must pass, and `openLog` retries a log name only for an activity kind.
    - `flai/cmd/mcp.go` and `flai/cmd/mcp_http.go` on T-1739: they wire each MCP tool's dependency, as `Plans: a.mcpPlan` does.
    - `flai/internal/hostapi/settings.go` and its test on T-1739: `settings.action` accepts only the actions in `hostapi.Actions`.
    - `flaiover/src/lib/server/agent.ts` on T-1739: `TestTheDashboardAsksForExactlyWhatFlaiOffers` fails both ways, so `role.run` joins `REQUIRED_METHODS`.
    - `docs/operators/settings.md` on T-1735: `TestSettingsIndexIsComplete` fails on a manifest key with no row, and T-1737 runs `./cmd/...` in the same layer.
  - **Co-change, kept:** `design/system/flai-cli.md` (shared, holds nothing) and `design/system/agent-narrative.md`, whose § Strategic agents' activity documents names the three documents.
  - **Co-change, not added:** `flaiover-dashboard.md`, `docs/users/flaiover.md`, and the dashboard's role pages, since the dashboard shows defined roles from S-0420. `mcpserver/server.go`, since tools register in `folder.go`. `design/adrs/README.md`, since no ADR is planned.
  - **Left to later stories:** `flai/internal/usage` and `flai/internal/metrics`, whose `StrategicKinds` S-0425 generalizes, and `flai/internal/guard`, which S-0416 changes. Wiring `roles.<name>.permissions` into the guard is proposed on TH-0570.
- **Forecast.** flai's 1h37m is kept: 107 s a unit times size 54, six criteria and 48 touches. The seven touches added are one-line wirings, which overstates the size. Generating the definition in memory (answer 6) adds work that offsets it. The delivery is flai's play-out, 2026-10-21T19:28Z.
- **Cost of delay.** 79.40 USD a week, worked out by hand from answer 1, because `flai cod` refuses while E-0024 has no inputs.
  - 4h lost per cycle at `planning.hour_rate` 150 over the default 168h cycle gives E-0024 600 USD a week.
  - S-0417 takes the share its 1h37m is of the 12h13m forecast for E-0024's twelve open stories: 97 / 733 × 600.
  - Once the input is on E-0024, `flai cod S-0417` gives the same figure while the forecasts stand.
