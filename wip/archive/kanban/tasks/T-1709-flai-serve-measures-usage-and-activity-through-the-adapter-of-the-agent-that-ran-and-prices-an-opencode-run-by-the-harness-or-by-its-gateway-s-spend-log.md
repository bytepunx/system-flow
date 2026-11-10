---
id: T-1709
type: task
nature: feature
title: flai serve measures usage and activity through the adapter of the agent that ran, and prices an OpenCode run by the harness or by its gateway's spend log
status: done
parent: S-0410
owner: alex
created: 2026-10-16T05:19:01Z
updated: 2026-11-07T21:31:13Z
transitions:
  - to: ready
    at: 2026-11-07T21:10:44Z
    by: agent-S-0410
  - to: in-progress
    at: 2026-11-07T21:10:48Z
    by: agent-S-0410
  - to: done
    at: 2026-11-07T21:31:13Z
    by: agent-S-0410
stream: S-0410
tags: [cli]
touches: [flai/internal/serve/usage.go, flai/internal/serve/usage_test.go, flai/internal/serve/activity.go, flai/internal/serve/activity_test.go]
after: [T-1708]
usage:
  source: log
  seconds: 1225
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 62
      output: 21595
      cache_read: 9666045
      cache_write: 169139
      cost: 3.4278
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
