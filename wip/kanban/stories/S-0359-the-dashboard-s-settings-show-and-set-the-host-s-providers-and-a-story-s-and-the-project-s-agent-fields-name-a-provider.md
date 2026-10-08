---
id: S-0359
type: story
nature: feature
title: The dashboard's settings show and set the host's providers, and a story's and the project's agent fields name a provider
status: backlog
parent: E-0019
owner: alex
created: 2026-10-08T08:52:44Z
updated: 2026-10-08T08:52:58Z
transitions: []
tags: [cli, dashboard]
topics: [agents]
touches: [flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/hostapi/settings.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flaiover/src/lib/settings.ts, flaiover/src/lib/components/SettingsPanel.svelte, flaiover/src/lib/components/SettingsPanel.svelte.test.ts, flaiover/src/lib/components/AgentFields.svelte, flaiover/src/lib/agent.ts, flaiover/src/lib/agent.test.ts, design/system/flaiover-dashboard.md]
after: [S-0349]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
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

Drafted by the planner; see the children.

## Notes

- Whether the key's variable is set is a yes or no, so that the page can say why a start would be refused without holding the key.
