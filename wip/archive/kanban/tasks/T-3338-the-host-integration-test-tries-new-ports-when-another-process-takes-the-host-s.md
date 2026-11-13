---
id: T-3338
type: task
nature: remediation
title: The host integration test tries new ports when another process takes the host's or the project's
status: done
parent: S-0779
owner: alex
created: 2026-11-13T21:45:58Z
updated: 2026-11-13T21:51:11Z
transitions:
  - to: ready
    at: 2026-11-13T21:46:18Z
    by: agent-S-0779
  - to: in-progress
    at: 2026-11-13T21:46:19Z
    by: agent-S-0779
  - to: done
    at: 2026-11-13T21:51:11Z
    by: agent-S-0779
stream: S-0779
tags: []
touches: [flai/cmd/host_integration_test.go]
usage:
  source: log
  seconds: 292
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 1966
      cache_read: 851077
      cache_write: 10448
      cost: 0.2931
---

# T-3338 The host integration test tries new ports when another process takes the host's or the project's

## Work

T-3333's second load run found it. `TestHostRunsServeAndTheMCPServersAndTakesThemWithIt` (`flai/cmd/host_integration_test.go`) failed in the fifty-count run: `flai host already runs on this machine for /tmp/TestHostUpgradeToAVersion.../cfg.json`. The test asks the kernel for two ports, closes them, and gives them to the host's child processes as the host's address and the project's remembered MCP address. A host that the whole-package run beside it started took the first port.

Run the scenario as a try on fresh ports, and try again, a bounded number of times, when a try loses a port. A try loses its host's port when `flai host start` fails because another host answers there, or because the host's log says its address is in use. It loses the project's port when the host's log says the MCP address moved because another process holds it. Build flai once, outside the tries. Any other failure fails the test as before.

Waits for nothing in the story: no other task touches the file.

## Done when

- The test no longer fails when another process takes either port it gave away, and still fails on anything else.
- `flai test` on `flai/cmd/host_integration_test.go` passes.

## Notes
