---
id: T-2529
type: task
nature: remediation
title: TestHostRunsServeAndTheMCPServersAndTakesThemWithIt waits for the MCP server's state within a bound that outlasts its staleness check on a loaded host
status: done
parent: S-0569
owner: alex
created: 2026-11-07T03:47:55Z
updated: 2026-11-07T04:40:48Z
transitions:
  - to: ready
    at: 2026-11-07T04:17:19Z
    by: agent-S-0569
  - to: in-progress
    at: 2026-11-07T04:17:22Z
    by: agent-S-0569
  - to: done
    at: 2026-11-07T04:40:48Z
    by: agent-S-0569
stream: S-0569
tags: [flai, tests]
touches: [flai/cmd/host_integration_test.go]
usage:
  source: log
  seconds: 1405
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 76
      output: 474
      cache_read: 8980939
      cache_write: 60262
      cost: 3.2322
---
# T-2529 TestHostRunsServeAndTheMCPServersAndTakesThemWithIt waits for the MCP server's state within a bound that outlasts its staleness check on a loaded host

## Work

I-0247's fourth instance: `TestHostRunsServeAndTheMCPServersAndTakesThemWithIt` (`flai/cmd/host_integration_test.go`, integration tier only, skipped under `-short`) failed at line 112, "never the MCP server's own state". Its `waitUntil` allows 20 seconds for the MCP server's state file to hold the child's PID. `readMCPState` (`flai/cmd/mcp_http.go`) also wants the file's `Updated` within `MCPStaleAfter`, 10 seconds (`flai/internal/serve/agents.go`), and the PID alive. On a loaded host the state may be written late, or read as stale between two refreshes.

First reproduce it: run `GOMAXPROCS=2 go test -race -count=5 -run TestHostRunsServeAndTheMCPServersAndTakesThemWithIt ./cmd` beside a CPU load, and record whether the state never came or came stale. Then make the wait hold on a loaded host: a longer bound in the test, or a shorter refresh the test sets through the host's options. Do not change `MCPStaleAfter` for running hosts to make a test pass. If the state file's refresh is genuinely late under load, that is a defect in flai: say so in the narrative and on a thread on the story before you change it. Waits for no other task.

## Done when

- [ ] The narrative says whether the state never came or came stale
- [ ] `GOMAXPROCS=2 go test -race -count=5 -run TestHostRunsServeAndTheMCPServersAndTakesThemWithIt ./cmd` passes beside a CPU load, recorded in the narrative
- [ ] The test asserts no less than before

## Notes
