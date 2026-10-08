---
id: S-0355
type: story
nature: feature
title: A harness that cannot run the guard or hold a call to ask runs a story's agent but none of its roles, unless the host sets guard none for it, and each refusal names what lifts it
status: backlog
parent: E-0019
owner: alex
created: 2026-10-08T08:48:54Z
updated: 2026-10-08T09:09:50Z
transitions: []
tags: [cli, dashboard]
topics: [agents]
touches: [flai/internal/config/config.go, flai/internal/config/config_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/hostapi/settings.go, flai/internal/hostapi/writes_test.go, flaiover/src/lib/components/SettingsPanel.svelte, flaiover/src/lib/components/SettingsPanel.svelte.test.ts, docs/operators/settings.md, docs/users/flai.md, docs/users/flai-reference.md, design/system/flai-cli.md]
after: [S-0352, S-0359]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 37.61
  by: planner-E-0019
  at: 2026-10-08T09:00:31Z
forecast:
  duration: 34m
  delivery: 2026-10-08T21:56:00Z
  basis: "Its own forecast of 34m; 29th in the pull order with an in-progress limit of 5, behind S-0232, S-0322, S-0341, S-0344, S-0345, S-0348, S-0338, S-0346, S-0342, S-0337, S-0334, S-0343, S-0297, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0289, S-0304, S-0305, S-0306, S-0313, S-0347, S-0349, S-0350, S-0351, S-0352, S-0353 and S-0354."
  by: flai
  at: 2026-10-08T09:09:50Z
finalized:
  by: alex
  at: 2026-10-08T08:52:24Z
---
# S-0355 A harness that cannot run the guard or hold a call to ask runs a story's agent but none of its roles, unless the host sets guard none for it, and each refusal names what lifts it

## Goal

With S-0352 each adapter states whether its harness runs the guard and can hold a call to ask. [ADR-0131](../../../design/adrs/0131-a-harness-without-a-guard-hook-runs-a-story-s-agent-but-none-of-its-roles.md) says what `flai serve` does with that before a start. On a harness without the guard, the story's agent may run but a story whose `roles` name a sub-agent is refused, unless the operator sets `agent.harnesses.<name>.guard: none` on the host. On a harness that cannot hold a call to ask, a start is refused unless the project's auto-approve is on or the host's arguments deny the protected paths. Every refusal names the setting that lifts it. Today the `command` harness is the one without either: this story is where its roles and its protected writes are first held back.

## Acceptance criteria

- [ ] The host's `agent.harnesses.<name>` takes `guard: none`, set and cleared with `flai serve agent harness <name> --guard none` and `--guard ""`, and through the dashboard's `settings.harness`; any other value is refused.
- [ ] Before it starts a story's agent, `flai serve` reads the adapter's capabilities. Without `Guard`, a story whose agent has `roles` is refused, naming the harness and `flai serve agent harness <name> --guard none`; with `guard: none` set it starts, and the start log says the roles run unguarded.
- [ ] Without `Ask`, a start is refused unless the project's `auto-approve` host action is on or the host's `deny_protected` for the harness is set, and the refusal names both; tests cover each case for the `command` harness and none of them changes a `claude-code` start.
- [ ] The refusal reaches the story's agent status as any refused start does, so that the dashboard's story page shows it; a test reads it there.
- [ ] The Settings page shows and sets a harness's guard, with a line saying what `none` gives up.
- [ ] `docs/operators/settings.md` indexes `agent.harnesses.<name>.guard` and `deny_protected` and `TestSettingsIndexIsComplete` passes; `docs/users/flai.md` says what a harness without a guard loses; `design/system/flai-cli.md` describes the checks.

## Tasks

- T-1398 The host's harness entry takes guard none and deny_protected, set with flai serve agent harness
- T-1399 flai serve reads the adapter's capabilities before a story's start and refuses roles without the guard and a start without asking, naming what lifts each
- T-1400 settings.harness and the Settings page show and set a harness's guard and deny_protected
- T-1401 The user and design documents say what a harness without the guard or asking loses and what lifts each refusal

## Notes

- `deny_protected` is the host's say that its arguments for the harness already deny the protected paths in the harness's own vocabulary (ADR-0131's last rule). flai cannot read another harness's arguments, so the operator states it; see the plan's thread.
- The planner, the orchestrator, and the analyzer are outside ADR-0131's decision and start as today.

### Planning

Planned by planner-E-0019 on 2026-10-08. Every touch is a file; no folder touch is kept. Three layers: the settings (T-1398); then the refusals (T-1399) and the dashboard (T-1400) together, with no file in common; then the documents (T-1401).

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/config/config.go`, `config_test.go` | design, layout | `HarnessHost` gains `Guard` and `DenyProtected` (T-1398) |
| `flai/cmd/serve_actions.go`, `serve_actions_test.go` | layout | `flai serve agent harness` (T-1398), and `agentConfig`, which hands them to `flai serve` (T-1399) |
| `docs/operators/settings.md`, `docs/users/flai-reference.md` | design | ADR-0131's consequences; the new flags regenerate both (T-1398) |
| `flai/internal/serve/agents.go`, `agents_test.go` | layout | The story's start, where the checks go (T-1399) |
| `flai/internal/hostapi/settings.go`, `writes_test.go`, `flaiover/src/lib/components/SettingsPanel.svelte`, `SettingsPanel.svelte.test.ts` | layout | `settings.harness` and each harness's box (T-1400) |
| `docs/users/flai.md`, `design/system/flai-cli.md` | design | ADR-0131's consequences (T-1401) |

`touches suggest` listed `flai/internal/hostapi/writes.go`, `flaiover/src/lib/server/agent.ts`, and `flai/internal/mcpserver/server.go`. None is taken: `settings.harness` is already a method the dashboard may call, and no MCP tool changes. `StoryAgent.svelte` already shows why a start failed, so the story page needs no change.

Forecast: 34m, as `flai forecast` gives it, 100 s per unit over 34 done feature stories on claude-opus-5-5 in the large band, times size 20. It stands.

Cost of delay: 37.61 USD a week, as `flai cod` works it out: 34m of the 7h32m forecast over E-0019's 12 stories, of the operator's 500 USD a week penalty. It stands.
