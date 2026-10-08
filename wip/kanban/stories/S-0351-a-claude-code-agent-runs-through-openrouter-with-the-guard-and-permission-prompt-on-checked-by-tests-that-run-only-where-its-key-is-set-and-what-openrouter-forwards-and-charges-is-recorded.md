---
id: S-0351
type: story
nature: feature
title: A claude-code agent runs through OpenRouter with the guard and permission_prompt on, checked by tests that run only where its key is set, and what OpenRouter forwards and charges is recorded
status: backlog
parent: E-0019
owner: alex
created: 2026-10-08T08:44:52Z
updated: 2026-10-08T08:51:59Z
transitions: []
tags: [cli]
topics: [agents, testing]
touches: [flai/internal/serve/gateway_test.go, flai/internal/serve/claudecheck.go, scripts/gateway-smoke.sh, scripts/smoke.sh, docs/contributors/index.md, design/system/agent-adapters.md]
after: [S-0350]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
finalized:
  by: alex
  at: 2026-10-08T08:51:59Z
---
# S-0351 A claude-code agent runs through OpenRouter with the guard and permission_prompt on, checked by tests that run only where its key is set, and what OpenRouter forwards and charges is recorded

## Goal

The finding of S-0339 leaves two facts unchecked that the later stories' contracts rest on: whether OpenRouter's `/api/v1/messages` forwards every beta header and body field Claude Code sends, and what OpenRouter's spend log reports beside Claude Code's own `total_cost_usd` ([agent-adapters.md](../../../design/system/agent-adapters.md) § LiteLLM and OpenRouter). This story runs a headless `claude -p` with the adapter's own arguments through OpenRouter, with S-0350's provider environment, `flai guard` registered, and `permission_prompt` as the permission tool, as `claudecheck.go` runs one against Anthropic. It records what it found before the contracts are fixed. The tests run only where an OpenRouter key is configured, and say they were skipped where it is not, as the operator asked on TH-0368.

## Acceptance criteria

- [ ] A Go test in `flai/internal/serve` builds a scratch project with one story in progress, `flai guard` registered in its `.claude/settings.json`, and auto-approve on. It runs `claude -p` through OpenRouter with the `claude-code` adapter's arguments and a provider entry for OpenRouter. It passes when a protected `Write` goes through `permission_prompt`, a sub-agent's work-item write is refused by the guard, the `haiku` alias resolves through `models`, and the run ends with a `result`.
- [ ] The test runs only when `OPENROUTER_API_KEY` is set and `claude` is on `PATH`, caps its spend with `max_budget_usd`, and otherwise skips with a message naming what is missing; the integration tier passes on a host without the key.
- [ ] `scripts/gateway-smoke.sh`, called from `scripts/smoke.sh`, runs the check verbosely when the key is set and prints one line saying it was skipped and why when it is not; the smoke tier passes on a host without the key.
- [ ] The test reads OpenRouter's `/api/v1/generation` for the run's calls and `/api/v1/key` before and after, and logs the gateway's cost beside the run's `total_cost_usd`.
- [ ] `design/system/agent-adapters.md` gains a section with the date, the `claude` version, and the answers: which beta headers and body fields reached the model, whether the aliases resolved, whether the message IDs in the stream match OpenRouter's generation IDs, and the gateway's cost against `total_cost_usd`. Each open item the run settles is no longer marked to check.
- [ ] `docs/contributors/index.md` says how to run the gateway checks with a key and a spend cap.

## Tasks

Drafted by the planner; see the children.
- T-1384 A gateway check test runs claude -p through OpenRouter with the guard and permission_prompt on, and skips naming what is missing
- T-1385 The smoke tier runs the gateway checks where a key is set and prints why it skipped them where none is
- T-1386 Run the OpenRouter check with the operator's key and record what OpenRouter forwarded and charged in agent-adapters.md

## Notes

- Spend: the operator configures an OpenRouter key with a few dollars of credit and a limit on the key (E-0019's notes). The story's agent needs `OPENROUTER_API_KEY` in its environment to run the check. When it is missing, ask on a thread rather than skip criterion 5.
- `claudecheck.go` is touched only to share its scratch-project builder with the test, not to change what `flai serve` checks.
