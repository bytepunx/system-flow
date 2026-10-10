---
id: S-0350
type: story
nature: feature
title: flai serve starts a claude-code agent through its provider, with the key read from the variable the host names and never written down
status: backlog
parent: E-0019
owner: alex
created: 2026-10-08T08:43:50Z
updated: 2026-10-10T18:58:58Z
transitions: []
tags: [cli]
topics: [agents]
touches: [flai/internal/harness/harness.go, flai/internal/harness/adapters.go, flai/internal/harness/harness_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/serve/plan.go, flai/internal/serve/orchestrate.go, flai/internal/serve/analyze.go, flai/cmd/serve_actions.go, docs/users/flai.md, docs/operators/index.md, design/system/flai-cli.md]
after: [S-0349]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 32.08
  by: planner-E-0019
  at: 2026-10-08T09:00:18Z
forecast:
  duration: 29m
  delivery: 2026-10-10T23:51:00Z
  basis: "Its own forecast of 29m; 14th in the pull order with an in-progress limit of 5, behind S-0232, S-0337, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0304, S-0305, S-0306, S-0313 and S-0349."
  by: flai
  at: 2026-10-10T18:58:58Z
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

- T-1376 The claude-code and command adapters derive their environment from a resolved provider and name the key's variable without reading it
- T-1378 flai serve resolves the agent's provider at every start, refuses an unset key variable, and copies the key into the child alone
- T-1380 The user, operator, and design documents say what a provider sets at start, what it refuses, and where its key's variable must be

## Notes

- Only Claude models work through Claude Code over a gateway (agent-adapters.md § Claude Code over a gateway). flai does not check the model: its meaning is the provider's.
- `claudecheck.go` still checks each new `claude` version against Anthropic direct. Running it through a provider is left out; see the plan's thread.

### Planning

Planned by planner-E-0019 on 2026-10-08. Every touch is a file; no folder touch is kept. Three layers, one task each: the adapters (T-1376), then `flai serve`'s start (T-1378), which fills what the adapters read, then the documents (T-1380).

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/harness/harness.go`, `adapters.go`, `harness_test.go` | design, layout | `Request`, `Start`, and the two adapters (T-1376) |
| `flai/internal/serve/agents.go`, `agents_test.go` | layout | `AgentConfig`, the story's start, and `spawn` (T-1378) |
| `flai/internal/serve/plan.go`, `orchestrate.go`, `analyze.go` | layout | Each calls `harness.For` and `adapter.Start` for its own start (T-1378) |
| `flai/cmd/serve_actions.go` | layout | `agentConfig` builds `serve.AgentConfig` from the host's configuration at every look (T-1378) |
| `docs/users/flai.md`, `docs/operators/index.md`, `design/system/flai-cli.md` | design | The start and where the key must be (T-1380) |

`touches suggest` listed `design/system/flaiover-dashboard.md`, `docs/users/flai-reference.md`, `docs/operators/settings.md`, and `flai/internal/hostapi/writes.go`. None is taken: no flag or setting is added, and the dashboard is S-0359.

Forecast: 29m, as `flai forecast` gives it, 100 s per unit over 34 done feature stories on claude-opus-5-5 in the large band, times size 17. It stands.

Cost of delay: 32.08 USD a week, as `flai cod` works it out: 29m of the 7h32m forecast over E-0019's 12 stories, of the operator's 500 USD a week penalty. It stands.
