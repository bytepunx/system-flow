---
id: T-1566
type: task
nature: feature
title: The free route ends a stalled stream so Claude Code retries, and the check's line has no trailing period
status: done
parent: S-0356
owner: alex
created: 2026-10-13T18:39:35Z
updated: 2026-10-13T18:40:20Z
transitions:
  - to: ready
    at: 2026-10-13T18:39:42Z
    by: agent-S-0356
  - to: in-progress
    at: 2026-10-13T18:39:43Z
    by: agent-S-0356
  - to: done
    at: 2026-10-13T18:40:20Z
    by: agent-S-0356
stream: S-0356
tags: []
touches: [scripts/litellm.sh, flai/internal/serve/gateway_test.go, design/system/agent-adapters.md]
usage:
  source: log
  seconds: 37
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 3
      output: 1243
      cache_read: 408557
      cache_write: 8643
      cost: 0.1705
---

# T-1566 The free route ends a stalled stream so Claude Code retries, and the check's line has no trailing period

## Work

- The close-out's smoke run of `TestGatewayLiteLLM` stalled: the free model's last stream took 4m40s, past the check's 5 minutes. `scripts/litellm.sh` gives the `openrouter-messages/*` deployment `stream_timeout: 45` and `timeout: 120`, so a stalled stream ends and Claude Code retries.
- The check's line drops its trailing period, which one free-model run left off.
- `agent-adapters.md` § Checked against LiteLLM says so.

Found by the close-out, after T-1404.

## Done when

- Three runs of the unified check finish within the 5 minutes.
- `flai test` on the three paths passes, and the close-out passes.

## Notes
