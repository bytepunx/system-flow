---
id: T-1076
type: task
nature: improvement
title: verify over MCP and verify.run and verify.status on the host channel answer the verification result
status: done
parent: S-0270
owner: alex
created: 2026-10-06T22:52:32Z
updated: 2026-10-07T03:49:36Z
transitions:
  - to: ready
    at: 2026-10-07T03:36:08Z
    by: agent-S-0270
  - to: in-progress
    at: 2026-10-07T03:36:09Z
    by: agent-S-0270
  - to: done
    at: 2026-10-07T03:49:36Z
    by: agent-S-0270
stream: S-0270
tags: [flai, mcp]
touches: [flai/internal/mcpserver/folder.go, flai/internal/mcpserver/verify_test.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/mcpserver/verify.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/folder_test.go, flaiover/src/lib/server/agent.ts]
after: [T-1071]
usage:
  source: log
  seconds: 807
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 80
      output: 33378
      cache_read: 5778089
      cache_write: 167706
      cost: 2.8701
---
# T-1076 verify over MCP and verify.run and verify.status on the host channel answer the verification result

## Work

Criterion 2: the same operation from MCP and the host channel.

- Add the MCP tool `verify` in `flai/internal/mcpserver/folder.go`, beside `prime`. It takes a story and answers T-1071's result as JSON. Its description says that it runs what the close-out runs, that a finding outside the story is a note, and that the sub-agent is kept for the review against the criteria and the conventions.
- Add `verify.run` to the host channel's write methods in `flai/internal/hostapi/writes.go`. Gate it behind the `checks` host action that gates `checks.run` today, with progress while it runs, as `checks.run` does. Add the read method `verify.status`, which answers the stored last result from `.flai-cache/verify/<story>.json`, or nothing when there is none.
- Tests: the tool in `flai/internal/mcpserver/verify_test.go`, and the methods, their gate, and the stored result in `flai/internal/hostapi/writes_test.go`.
- Waits for T-1071. Runs alongside T-1075, with no path in common.

## Done when

- [ ] `verify` over MCP and `verify.run` answer the same result as `flai verify --json`.
- [ ] `verify.status` answers the last result, and `verify.run` is refused while the `checks` host action is off.
- [ ] `scripts/flai-test.sh` passes.

## Notes

Drafted by the planner. Assumes the `checks` host action gates `verify.run`, rather than a new action. This is on the plan thread.
