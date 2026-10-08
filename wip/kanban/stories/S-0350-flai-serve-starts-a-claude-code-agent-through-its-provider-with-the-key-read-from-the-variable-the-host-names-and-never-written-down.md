---
id: S-0350
type: story
nature: feature
title: flai serve starts a claude-code agent through its provider, with the key read from the variable the host names and never written down
status: backlog
parent: E-0019
owner: alex
created: 2026-10-08T08:43:50Z
updated: 2026-10-08T08:51:57Z
transitions: []
tags: [cli]
topics: [agents]
touches: [flai/internal/harness/harness.go, flai/internal/harness/adapters.go, flai/internal/harness/harness_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/serve/plan.go, flai/internal/serve/orchestrate.go, flai/internal/serve/analyze.go, docs/users/flai.md, docs/operators/index.md, design/system/flai-cli.md]
after: [S-0349]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
finalized:
  by: alex
  at: 2026-10-08T08:51:57Z
---
# S-0350 flai serve starts a claude-code agent through its provider, with the key read from the variable the host names and never written down

## Goal

With S-0349 a story's agent names a provider and flai can resolve it, but nothing uses it. This story makes `flai serve` start an agent through it ([ADR-0129](../../../design/adrs/0129-a-story-s-agent-names-a-provider-beside-its-harness-and-model-and-the-providers.md)). The `claude-code` adapter derives Claude Code's gateway variables from the resolved entry. The `command` adapter passes the entry on. `flai serve` reads the key from the variable `key_env` names, puts it only in the child's environment, and refuses a start it cannot make, naming the setting that fixes it. It covers the story's agent, the planner, the orchestrator, and the analyzer, which start through the same adapters.

## Acceptance criteria

- [ ] With a provider whose `api` is `anthropic-messages`, the `claude-code` adapter sets `ANTHROPIC_BASE_URL` to `base_url`, `ANTHROPIC_AUTH_TOKEN` from the variable `key_env` names, `ANTHROPIC_API_KEY` to the empty string, and `ANTHROPIC_DEFAULT_OPUS_MODEL`, `ANTHROPIC_DEFAULT_SONNET_MODEL`, and `ANTHROPIC_DEFAULT_HAIKU_MODEL` from `models`' `opus`, `sonnet`, and `haiku`; with no provider it sets none of them, as today.
- [ ] The `command` adapter gets `FLAI_PROVIDER`, `FLAI_PROVIDER_API`, `FLAI_PROVIDER_BASE_URL`, and `FLAI_PROVIDER_KEY_ENV`, and `{provider}` in its arguments; it reads the key itself.
- [ ] `flai serve` refuses a start whose provider the host lacks, whose `key_env` variable is unset or empty in its environment, or whose `api` the harness cannot speak, and the refusal names the setting or variable that fixes it, for a story's agent and for the planner, the orchestrator, and the analyzer.
- [ ] The key's value is in no log, no work item, no thread, and no start record: the adapter names the variable to copy, and `flai serve` copies the value into the child's environment as it spawns it; a test reads the run log and the start record of a start with a provider and finds no key.
- [ ] `docs/users/flai.md` says what a provider sets and refuses; `docs/operators/index.md` says the key's variable must be in the environment `flai host` or `flai serve` runs in; `design/system/flai-cli.md` describes the start.

## Tasks

Drafted by the planner; see the children.
- T-1376 The claude-code and command adapters derive their environment from a resolved provider and name the key's variable without reading it
- T-1378 flai serve resolves the agent's provider at every start, refuses an unset key variable, and copies the key into the child alone
- T-1380 The user, operator, and design documents say what a provider sets at start, what it refuses, and where its key's variable must be

## Notes

- Only Claude models work through Claude Code over a gateway (agent-adapters.md § Claude Code over a gateway). flai does not check the model: its meaning is the provider's.
- `claudecheck.go` still checks each new `claude` version against Anthropic direct. Running it through a provider is left out; see the plan's thread.
