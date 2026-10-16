---
id: T-1709
type: task
nature: feature
title: flai serve measures usage and activity through the adapter of the agent that ran, and prices an OpenCode run by the harness or by its gateway's spend log
status: backlog
parent: S-0410
owner: alex
created: 2026-10-16T05:19:01Z
updated: 2026-10-16T05:19:01Z
transitions: []
stream: S-0410
tags: [cli]
touches: [flai/internal/serve/usage.go, flai/internal/serve/usage_test.go, flai/internal/serve/activity.go, flai/internal/serve/activity_test.go]
after: [T-1708]
---
# T-1709 flai serve measures usage and activity through the adapter of the agent that ran, and prices an OpenCode run by the harness or by its gateway's spend log

## Work

- In `flai/internal/serve/usage.go` (`ReadRates`) and `activity.go`, read through the adapter of the run's agent, from its start record, in place of `harness.Adapters[ClaudeCode].Reader()`.
- For a run on `opencode`: with no provider, keep `step_finish`'s cost and set `priced_by: harness`; with one, join the run to the gateway's spend log by what S-0406 found `step_finish` carries, set `priced_by: gateway`, and keep OpenCode's figure as tokens only, as `usage.For` does for Claude Code.
- Test a Claude Code run unchanged, an OpenCode run with and without a provider against the fake OpenRouter spend server E-0019's tests use, and a strategic agent's activity on OpenCode.
- Waits for T-1708: it reads through the adapter's reader.

## Done when

Usage and activity come through the run's adapter, OpenCode runs carry the right `priced_by`, Claude Code runs are as they were, and the serve tests pass.

## Notes
