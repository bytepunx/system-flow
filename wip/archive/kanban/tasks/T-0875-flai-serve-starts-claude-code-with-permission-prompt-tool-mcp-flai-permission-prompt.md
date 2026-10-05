---
id: T-0875
type: task
nature: remediation
title: flai serve starts claude-code with --permission-prompt-tool mcp__flai__permission_prompt
status: done
parent: S-0257
owner: alex
created: 2026-10-05T04:15:09Z
updated: 2026-10-05T04:21:52Z
transitions:
  - to: ready
    at: 2026-10-05T04:15:53Z
    by: agent-S-0257
  - to: in-progress
    at: 2026-10-05T04:15:54Z
    by: agent-S-0257
  - to: done
    at: 2026-10-05T04:21:52Z
    by: agent-S-0257
stream: S-0257
tags: []
touches: [flai/internal/harness]
usage:
  source: log
  seconds: 358
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 13
      output: 3886
      cache_read: 616798
      cache_write: 16119
      cost: 0.301
---
# T-0875 flai serve starts claude-code with --permission-prompt-tool mcp__flai__permission_prompt

## Work

In the claude-code adapter's `Start` (`flai/internal/harness/adapters.go`), add `--permission-prompt-tool mcp__flai__permission_prompt` to the agent's own argv, before the operator's host args. It must apply whether the host args are the defaults or the operator's own. Check what `--mcp-config` names the server, and use the same name, so that the tool resolves.

## Done when

The harness tests assert the flag and its value, with the default host and with custom host args.

## Notes
