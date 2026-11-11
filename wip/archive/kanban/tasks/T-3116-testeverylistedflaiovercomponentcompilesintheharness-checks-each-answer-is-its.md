---
id: T-3116
type: task
nature: remediation
title: TestEveryListedFlaioverComponentCompilesInTheHarness checks each answer is its own vite's and says so when another server answered
status: done
parent: S-0714
owner: alex
created: 2026-11-11T09:05:39Z
updated: 2026-11-11T19:00:57Z
transitions:
  - to: ready
    at: 2026-11-11T18:54:18Z
    by: agent-S-0714
  - to: in-progress
    at: 2026-11-11T18:54:18Z
    by: agent-S-0714
  - to: done
    at: 2026-11-11T19:00:57Z
    by: agent-S-0714
stream: S-0714
tags: [flai, components, test]
touches: [flai/internal/components/flaiover_test.go]
after: [T-3115]
usage:
  source: log
  seconds: 399
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 21
      output: 7701
      cache_read: 1873951
      cache_write: 96919
      cost: 1.1443
---
# T-3116 TestEveryListedFlaioverComponentCompilesInTheHarness checks each answer is its own vite's and says so when another server answered

## Work

`fetch` in `flai/internal/components/flaiover_test.go` reads only the status and the body, so I-0342's failure said that the setup did not compile when another server had answered. Make `fetch` give the answer's harness header too. Each check of the setup, each component, and the broken control then compares it with the dev server's token first. On a mismatch it fails saying another server answers on vite's port, with whether vite still runs, from `Err` and `Done`, and vite's last lines, rather than reporting a compile failure.

Waits for T-3115, which gives the token and the header this reads.

## Done when

- An answer without the dev server's token fails the test with a message naming another server on vite's port, vite's state, and its last lines, never a compile failure.
- The test passes on this tree with node and flaiover's dependencies installed, run with `flai test integration` or `go test -run TestEveryListedFlaioverComponentCompilesInTheHarness ./internal/components/` in `flai/`.

## Notes
