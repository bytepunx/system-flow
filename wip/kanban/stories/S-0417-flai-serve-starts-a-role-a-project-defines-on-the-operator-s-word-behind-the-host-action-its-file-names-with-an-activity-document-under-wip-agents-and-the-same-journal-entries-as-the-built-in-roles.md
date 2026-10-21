---
id: S-0417
type: story
nature: feature
title: flai serve starts a role a project defines on the operator's word, behind the host action its file names, with an activity document under wip/agents and the same journal entries as the built-in roles
status: backlog
parent: E-0024
owner: alex
created: 2026-10-16T05:40:01Z
updated: 2026-10-21T08:32:26Z
transitions: []
tags: [cli, template]
topics: [cli, planning, orchestration]
touches: [design/system/project-manifest.md, design/system/strategic-agents.md, docs/operators/index.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flai-reference.md, flai/cmd/activity.go, flai/cmd/activity_test.go, flai/cmd/role.go, flai/cmd/role_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/check/activity.go, flai/internal/check/activity_test.go, flai/internal/check/check.go, flai/internal/harness/adapters.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/hostapi_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/mcpserver/activity.go, flai/internal/mcpserver/activity_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/role.go, flai/internal/mcpserver/role_test.go, flai/internal/serve/activity.go, flai/internal/serve/agents.go, flai/internal/serve/role.go, flai/internal/serve/role_test.go, flai/internal/workitem/activity.go, flai/internal/workitem/activity_test.go, flai/internal/workitem/narrative.go, template/CHANGELOG.md, design/system/flai-cli.md, design/system/agent-narrative.md]
after: [S-0415, S-0416]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 1h15m
  delivery: 2026-10-21T16:00:00Z
  basis: "Its own forecast of 1h15m; 16th in the pull order with an in-progress limit of 5, behind S-0456, S-0463, S-0465, S-0459, S-0461, S-0462, S-0460, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415 and S-0416."
  by: flai
  at: 2026-10-21T08:32:26Z
---
# S-0417 flai serve starts a role a project defines on the operator's word, behind the host action its file names, with an activity document under wip/agents and the same journal entries as the built-in roles

## Goal

A role a project defines runs. `flai role run <name> [<item>]`, the MCP tool `role_run`, and the host API's `role.run` start it through `flai serve`, as `flai plan` starts the planner: behind the host action its file names, in the main checkout, with its prompt from its file's body, its harness definition, and its own agent. Its run is recorded and settled as the planner's is, its activity is logged in `wip/agents/<name>.md`, and its starts and ends are journal entries and dashboard notifications. A fixture role shows it end to end.

## Acceptance criteria

- [ ] `flai role run <name> [<item>] [--json]`, the MCP tool `role_run` `{name, item?}`, and the host API's `role.run` start a defined role through `serve`, and each is refused, saying why, while the host action its file names is off, for an item its file does not take or one archived, done, or cancelled, while a run of the role for the same item runs, and when its agent names no harness and no command is set on the host. `flai serve enable <action>` and the Settings page accept a defined role's action. The MCP tool is refused to an agent `flai serve` started, as `plan` is.
- [ ] The run starts in the main checkout through the harness adapter with the role's verb as `Role`, `FLAI_AGENT` set to `<name>` or `<name>-<item>`, `FLAI_ROLE`, `FLAI_ITEM` when it has one, `FLAI_SESSION`, and `FLAI_STARTED_BY=flai-serve`; its prompt is its role file's body with the item and the run's inputs; on `claude-code` it runs with `--agent` over the definition its file names; its agent is `roles.<name>.agent` in the manifest merged over the project's `agent`, as `planning.agent` is.
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

## Notes

- Drafted by the planner for E-0024. Waits for the prime story and the guard story: a role is started only once its pack and its guard read its file, so it never runs unguarded.
- A defined role's activity charges no item until the metrics story generalizes `usage.strategic` per role: it stays in the role's project strategic total, which `flai stats` already works out per kind.
- Starting on an event, on a schedule, or for as long as the action is on is the next story's.

### Planning

- **Touches.** The story declared none. Every touch comes from its six tasks.
  - **From layout:** the manifest and its settings table; the harness and the claude-code adapter; the activity documents in `workitem`, `check`, `mcpserver`, and `cmd`; a new `flai/internal/serve/role.go`, with `agents.go` and `activity.go`; a new MCP tool `role.go` registered in `folder.go`; the host API's `writes.go` and `hostapi.go`; and `flai/cmd/role.go` and `serve_actions.go`. Each has its tests.
  - **From design:** `strategic-agents.md`, `project-manifest.md`, the operator and user docs, `docs/operators/settings.md` for the new keys, and `template/CHANGELOG.md`.
  - **Co-change, added:** `design/system/flai-cli.md` (39%), and `design/system/agent-narrative.md`, whose § Strategic agents' activity documents lists the three documents by name.
  - **Co-change, not added:** `flaiover/src/lib/server/agent.ts` and the dashboard's design. The dashboard reads `roles` from S-0420 on. Also `mcpserver/server.go`: the new tool registers in `folder.go`, as `plan` does.
  - No folder touch.
- **Forecast.** Kept at flai's 1h15m. Most touches are entry points and tests modelled on the planner's. The forecast's basis gives the figures, and its delivery is flai's play-out moved by any raise times the cycle factor.
- **Cost of delay.** Not set. Neither the story nor E-0024 has inputs, so `flai cod` refuses. The epic's input is asked on TH-0471, question 1, and the value follows its answer.
