---
id: T-1384
type: task
nature: feature
title: A gateway check test runs claude -p through OpenRouter with the guard and permission_prompt on, and skips naming what is missing
status: backlog
parent: S-0351
owner: alex
created: 2026-10-08T08:45:02Z
updated: 2026-10-08T08:45:02Z
transitions: []
stream: S-0351
tags: [cli]
touches: [flai/internal/serve/gateway_test.go, flai/internal/serve/claudecheck.go]
---
# T-1384 A gateway check test runs claude -p through OpenRouter with the guard and permission_prompt on, and skips naming what is missing

## Work

- Factor the scratch project `claudecheck.go` builds (one story in progress, auto-approve on) into a helper the test can call with extra files and environment; `flai serve`'s own check keeps its behaviour.
- `gateway_test.go`, `TestGatewayOpenRouter`: skip unless `OPENROUTER_API_KEY` is set and `claude` is on `PATH`, saying which is missing. Add `.claude/settings.json` registering `flai guard` to the scratch project. Build the start with the `claude-code` adapter and a provider (`api: anthropic-messages`, `base_url: https://openrouter.ai/api`, `key_env: OPENROUTER_API_KEY`, `models.haiku` from `FLAI_TEST_OPENROUTER_HAIKU`, default `anthropic/claude-haiku-4.5`), with `max_budget_usd` from `FLAI_TEST_GATEWAY_BUDGET`, default `0.10`.
- The prompt asks for a protected `Write`, which `permission_prompt` approves; and for a sub-agent that calls `item_edit`, which the guard refuses. The test reads the stream for both outcomes and a `result`.
- After the run, read `/api/v1/generation?id=<id>` for each assistant message ID and `/api/v1/key` before and after, and `t.Log` the gateway's total beside `total_cost_usd`. A generation the gateway does not find is logged, not failed, since what it returns is what the check is for.
- Use a short timeout and skip under `-short`, so the go-test tier never runs it.

First layer.

## Done when

- Without the key, `go test ./internal/serve/ -run TestGatewayOpenRouter -v` skips, naming `OPENROUTER_API_KEY`.
- With the key, it passes and logs both costs.
- `flai test flai/internal/serve/` passes.

## Notes

Layer 1 of S-0351.
