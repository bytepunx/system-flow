---
id: S-0355
type: story
nature: feature
title: A harness that cannot run the guard or hold a call to ask runs a story's agent but none of its roles, unless the host sets guard none for it, and each refusal names what lifts it
status: backlog
parent: E-0019
owner: alex
created: 2026-10-08T08:48:54Z
updated: 2026-10-08T08:52:24Z
transitions: []
tags: [cli, dashboard]
topics: [agents]
touches: [flai/internal/config/config.go, flai/internal/config/config_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/hostapi/settings.go, flai/internal/hostapi/writes_test.go, flaiover/src/lib/components/SettingsPanel.svelte, flaiover/src/lib/components/SettingsPanel.svelte.test.ts, docs/operators/settings.md, docs/users/flai.md, design/system/flai-cli.md]
after: [S-0352]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
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

Drafted by the planner; see the children.
- T-1398 The host's harness entry takes guard none and deny_protected, set with flai serve agent harness
- T-1399 flai serve reads the adapter's capabilities before a story's start and refuses roles without the guard and a start without asking, naming what lifts each
- T-1400 settings.harness and the Settings page show and set a harness's guard and deny_protected
- T-1401 The user and design documents say what a harness without the guard or asking loses and what lifts each refusal

## Notes

- `deny_protected` is the host's say that its arguments for the harness already deny the protected paths in the harness's own vocabulary (ADR-0131's last rule). flai cannot read another harness's arguments, so the operator states it; see the plan's thread.
- The planner, the orchestrator, and the analyzer are outside ADR-0131's decision and start as today.
