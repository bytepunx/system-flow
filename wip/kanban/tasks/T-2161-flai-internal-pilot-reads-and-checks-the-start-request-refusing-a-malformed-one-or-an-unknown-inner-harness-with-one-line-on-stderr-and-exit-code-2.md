---
id: T-2161
type: task
nature: feature
title: flai/internal/pilot reads and checks the start request, refusing a malformed one or an unknown inner harness with one line on stderr and exit code 2
status: done
parent: S-0507
owner: alex
created: 2026-10-26T04:49:12Z
updated: 2026-10-28T23:00:21Z
transitions:
  - to: ready
    at: 2026-10-28T22:57:33Z
    by: agent-S-0507
  - to: in-progress
    at: 2026-10-28T22:57:34Z
    by: agent-S-0507
  - to: done
    at: 2026-10-28T23:00:21Z
    by: agent-S-0507
stream: S-0507
tags: [cli, docker]
touches: [flai/internal/pilot/request.go, flai/internal/pilot/request_test.go]
after: [T-2154]
usage:
  source: log
  seconds: 167
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 28
      output: 109
      cache_read: 982591
      cache_write: 58017
      cost: 0.3543
---
# T-2161 flai/internal/pilot reads and checks the start request, refusing a malformed one or an unknown inner harness with one line on stderr and exit code 2

## Work

Write `flai/internal/pilot/request.go`: `Read(io.Reader) (Request, error)` decodes one JSON document, `harness.Request` with the container fields T-2154 added, from standard input, and `Check` refuses one that lacks `Inner`, names an inner harness `harness.Adapters` lacks or that is `container`, lacks `Cache` or names a folder that does not exist, or lacks the MCP variable names when the inner harness needs them; the error is one line with no field value that could be a secret. Define `ExitUsage = 2` for the refusals. Test in `request_test.go` each refusal and a valid request, checking the error text never echoes the environment. Waits for T-2154 for the fields.

## Done when

`flai test flai/internal/pilot` passes with a test per refusal and the valid case, and the error text of each refusal is one line.

## Notes

The package stays free of `os/exec` and the harness start; the next task adds them.
