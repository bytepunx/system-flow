---
id: S-0349
type: story
nature: feature
title: A story's agent names a provider from the host's providers map, and a project's manifest may override the entry's api, base_url, and models
status: backlog
parent: E-0019
owner: alex
created: 2026-10-08T08:42:40Z
updated: 2026-10-08T09:00:17Z
transitions: []
tags: [cli, template]
topics: [agents]
touches: [flai/internal/manifest/provider.go, flai/internal/manifest/provider_test.go, flai/internal/manifest/agent.go, flai/internal/manifest/agent_test.go, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/config/config.go, flai/internal/config/config_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/cmd/agent.go, flai/cmd/agent_test.go, flai/cmd/edit.go, flai/cmd/edit_test.go, flai/internal/itemedit/itemedit.go, flai/internal/mcpserver/items_write.go, flai/internal/mcpserver/items_write_test.go, docs/operators/settings.md, docs/users/flai.md, docs/users/flai-reference.md, design/system/project-manifest.md, design/system/work-hierarchy.md, design/system/flai-cli.md, template/root/system-flow.yaml.tmpl, template/CHANGELOG.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 60.84
  by: planner-E-0019
  at: 2026-10-08T09:00:10Z
forecast:
  duration: 55m
  delivery: 2026-10-08T15:02:00Z
  basis: "Its own forecast of 55m; 26th in the pull order with an in-progress limit of 5, behind S-0232, S-0288, S-0322, S-0341, S-0290, S-0344, S-0345, S-0338, S-0346, S-0342, S-0337, S-0334, S-0343, S-0297, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0289, S-0304, S-0305, S-0306, S-0313, S-0347 and S-0348."
  by: flai
  at: 2026-10-08T09:00:17Z
finalized:
  by: alex
  at: 2026-10-08T08:51:54Z
---
# S-0349 A story's agent names a provider from the host's providers map, and a project's manifest may override the entry's api, base_url, and models

## Goal

A story's `agent` names the harness and the model today, and where the model's calls go is nowhere in flai ([ADR-0129](../../../design/adrs/0129-a-story-s-agent-names-a-provider-beside-its-harness-and-model-and-the-providers.md)). This story adds the name and the map it resolves in, and nothing that starts an agent yet: `agent.provider` beside `harness` and `model`, in the manifest's default agent, under `planning`, `orchestration`, and `analysis`, and in a story's front matter; a `providers` map on the host, each entry with `api`, `base_url`, `key_env`, and `models`; and a `providers` map in the manifest that overrides an entry's `api`, `base_url`, and `models`, never `key_env`. One function resolves a provider name against the two, so that the next story's start reads one value.

## Acceptance criteria

- [ ] `agent.provider` is accepted in the manifest's default agent, in `planning.agent`, `orchestration.agent`, and `analysis.agent`, and in a story's front matter, validated by the same name pattern as `harness`, and copied into a new story with the rest of the agent (ADR-0037).
- [ ] The host's configuration holds `agent.providers.<name>` with `api` (`anthropic-messages`, `openai-chat`, or `openai-responses`), `base_url`, `key_env`, and optional `models` (`<alias>: <name on the provider>`), set and reset with `flai serve agent provider <name>`; an unknown `api`, a `base_url` that is not an http or https URL, or a `key_env` that is not a variable name is refused, and no flag takes a key.
- [ ] The manifest's `providers.<name>` may set `api`, `base_url`, and `models`; `key_env` there is refused by the manifest's validation, naming the host setting that holds it.
- [ ] One resolver merges the host's entry with the manifest's override and refuses a provider the host has no entry for, naming `flai serve agent provider <name>`; tests cover the merge, each refusal, and an agent with no provider, which resolves to none.
- [ ] `flai agent`, `flai edit`, `flai manifest set`, and the MCP `item_new` and `item_edit` agent take `provider`, and `flai check --strict` passes on a story that names one.
- [ ] `docs/operators/settings.md` indexes every new key and `TestSettingsIndexIsComplete` passes; `docs/users/flai.md`, `design/system/project-manifest.md`, `design/system/work-hierarchy.md`, and `design/system/flai-cli.md` describe the provider; the template's `system-flow.yaml.tmpl` shows `provider` and `providers` commented out, with a `template/CHANGELOG.md` entry.

## Tasks

- T-1371 The manifest's agent takes provider, the manifest's providers map overrides api, base_url, and models, and one resolver merges it with the host's entry
- T-1372 The host's configuration holds agent.providers, set and reset with flai serve agent provider
- T-1373 flai agent, flai edit, flai manifest set, and the MCP item_new and item_edit agent take provider
- T-1374 The user and design documents describe the provider, and the template's manifest shows provider and providers

## Notes

- Nothing here starts an agent through a provider: the next story derives the harness's environment from the resolved entry.
- The dashboard's settings for providers are their own story, after this one.

### Planning

Planned by planner-E-0019 on 2026-10-08. Every touch is a file; no folder touch is kept. The tasks run in four layers, one each: T-1371, then T-1372, then T-1373, then T-1374. They are chained because each of the first three adds rows to `docs/operators/settings.md`, whose index test must pass in every task.

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/manifest/provider.go`, `provider_test.go` | design | ADR-0129's entry, its validation, and the resolver, a new file (T-1371) |
| `flai/internal/manifest/agent.go`, `agent_test.go`, `manifest.go`, `manifest_test.go` | design, layout | `Agent.Provider` and `Manifest.Providers` (T-1371) |
| `flai/internal/config/config.go`, `config_test.go`, `flai/cmd/serve_actions.go`, `serve_actions_test.go` | layout | `AgentStart` holds `Harnesses` beside which `Providers` goes; `flai serve agent` lives in `serve_actions.go` (T-1372) |
| `flai/cmd/agent.go`, `agent_test.go`, `edit.go`, `edit_test.go`, `flai/internal/itemedit/itemedit.go`, `flai/internal/manifest/settings.go`, `settings_test.go`, `flai/internal/mcpserver/items_write.go`, `items_write_test.go` | layout | Every way an agent is set today (T-1373) |
| `docs/operators/settings.md` | design | The index test fails on a key without a row (T-1371 to T-1373) |
| `docs/users/flai-reference.md` | co-change (25%) | The new flags regenerate it with `make flai-reference` (T-1372, T-1373) |
| `docs/users/flai.md`, `design/system/project-manifest.md`, `design/system/work-hierarchy.md`, `design/system/flai-cli.md`, `template/root/system-flow.yaml.tmpl`, `template/CHANGELOG.md` | design | ADR-0129's consequences and the epic's last criterion (T-1374) |

`touches suggest` also listed `docs/operators/index.md`, `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`, and `flai/internal/hostapi/writes.go`. None is taken: the operator's guide changes with S-0350 and the dashboard with S-0359.

Forecast: 55m, as `flai forecast` gives it, 100 s per unit over 34 done feature stories on claude-opus-5-5 in the large band, times size 33. It stands: four tasks of plain schema and command work.

Cost of delay: 63.66 USD a week, as `flai cod` works it out: this story's 55m of the 7h12m forecast over E-0019's 11 stories, of the operator's 500 USD a week penalty. It stands. The value of the epic arrives with this story, S-0350, and S-0351, which their `after` already put first.
