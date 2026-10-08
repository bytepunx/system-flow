---
id: S-0359
type: story
nature: feature
title: The dashboard's settings show and set the host's providers, and a story's and the project's agent fields name a provider
status: backlog
parent: E-0019
owner: alex
created: 2026-10-08T08:52:44Z
updated: 2026-10-08T09:08:56Z
transitions: []
tags: [cli, dashboard]
topics: [agents]
touches: [flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/hostapi/settings.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flaiover/src/lib/settings.ts, flaiover/src/lib/server/agent.ts, flaiover/src/lib/components/SettingsPanel.svelte, flaiover/src/lib/components/SettingsPanel.svelte.test.ts, flaiover/src/lib/components/AgentFields.svelte, flaiover/src/lib/components/ItemEditor.svelte, flaiover/src/lib/components/NewItemForm.svelte, flaiover/src/lib/agent.ts, flaiover/src/lib/agent.test.ts, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [S-0349]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 40.93
  by: planner-E-0019
  at: 2026-10-08T09:00:27Z
forecast:
  duration: 37m
  delivery: 2026-10-08T22:15:00Z
  basis: "Its own forecast of 37m; 33rd in the pull order with an in-progress limit of 5, behind S-0232, S-0322, S-0341, S-0344, S-0345, S-0338, S-0346, S-0342, S-0337, S-0334, S-0343, S-0297, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0289, S-0304, S-0305, S-0306, S-0313, S-0347, S-0348, S-0349, S-0350, S-0351, S-0352, S-0353, S-0354, S-0355, S-0356, S-0357 and S-0358."
  by: flai
  at: 2026-10-08T09:08:56Z
finalized:
  by: alex
  at: 2026-10-08T08:52:58Z
---
# S-0359 The dashboard's settings show and set the host's providers, and a story's and the project's agent fields name a provider

## Goal

S-0349 puts providers on the host and `provider` on an agent, set from a shell. [ADR-0129](../../../design/adrs/0129-a-story-s-agent-names-a-provider-beside-its-harness-and-model-and-the-providers.md) and E-0019's first criterion also have the dashboard's settings set them, under the `settings` host action as a harness's program is set today ([ADR-0039](../../../design/adrs/0039-a-settings-host-action-turned-on-only-in-a-shell-lets-the-dashboard-change-the.md)). This story adds a providers section to the Settings page through a `settings.provider` write. It also adds a provider field to `AgentFields`, used for the project's default agent, a new story, and a story's agent. A key is never shown or entered: the page shows and sets the name of its variable.

## Acceptance criteria

- [ ] `settings.get`'s `host` carries `agent.providers`, each entry's `api`, `base_url`, `key_env`, `models`, and whether `key_env` is set in `flai serve`'s environment, never its value.
- [ ] `settings.provider` sets or resets an entry through `flai serve agent provider`, refuses what that command refuses, and is host-wide, as `settings.harness` is.
- [ ] `settings.default_agent`, `settings.manifest`, and the item writes that take an agent accept `provider`, and their refusal text names it among the agent's keys.
- [ ] The Settings page has a Providers section: each entry with its fields, a line saying whether its key's variable is set, Save and Reset, and a form to add one; read-only with the enabling command while `settings` is off.
- [ ] `AgentFields` has a provider field with the project's default as its placeholder, used wherever the agent is edited, and sends it only when typed.
- [ ] Tests cover the read, the write and a refusal in `writes_test.go`, and the section and the field in the flaiover tests; `design/system/flaiover-dashboard.md` describes both.

## Tasks

- T-1415 settings.get carries the host's providers, settings.provider sets them, and the agent writes accept provider
- T-1416 The Settings page has a Providers section that shows, adds, saves, and resets the host's providers
- T-1417 AgentFields has a provider field, with the project's default as its placeholder, sent only when typed
- T-1418 The dashboard's design document describes the Providers section and the provider field

## Notes

- Whether the key's variable is set is a yes or no, so that the page can say why a start would be refused without holding the key.

### Planning

Planned by planner-E-0019 on 2026-10-08. Every touch is a file; no folder touch is kept. Three layers: flai's read and writes (T-1415); then the Providers section (T-1416) and the agent field (T-1417) together, with no file in common; then the documents (T-1418).

| Touch | Source | Why |
|-------|--------|-----|
| `flai/cmd/serve_actions.go`, `serve_actions_test.go` | layout | `hostSettings` builds `settings.get`'s `host` (T-1415) |
| `flai/internal/hostapi/settings.go`, `writes.go`, `writes_test.go` | layout | `settings.harness` and the agent writes live here; `writes.go` co-changes 11% with S-0349's files (T-1415) |
| `flaiover/src/lib/settings.ts`, `flaiover/src/lib/server/agent.ts`, `flaiover/src/lib/components/SettingsPanel.svelte`, `SettingsPanel.svelte.test.ts` | layout, co-change (14%) | The settings types, the methods the dashboard may call, and the page (T-1416) |
| `flaiover/src/lib/components/AgentFields.svelte`, `ItemEditor.svelte`, `NewItemForm.svelte`, `flaiover/src/lib/agent.ts`, `agent.test.ts` | layout | The agent field and the two other forms that bind it (T-1417) |
| `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md` | design, co-change (48%) | The pages' design and their users' guide (T-1418) |

`touches suggest` also listed `design/system/flai-cli.md` and `docs/users/flai.md`. Neither is taken: S-0349 describes the commands.

Forecast: 37m, as `flai forecast` gives it, 100 s per unit over 34 done feature stories on claude-opus-5-5 in the large band, times size 22. It stands.

Cost of delay: 40.93 USD a week, as `flai cod` works it out: 37m of the 7h32m forecast over E-0019's 12 stories, of the operator's 500 USD a week penalty. It stands.
