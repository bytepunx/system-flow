---
id: T-3116
type: task
nature: remediation
title: TestEveryListedFlaioverComponentCompilesInTheHarness checks each answer is its own vite's and says so when another server answered
status: backlog
parent: S-0714
owner: alex
created: 2026-11-11T09:05:39Z
updated: 2026-11-11T09:05:39Z
transitions: []
stream: S-0714
tags: [flai, components, test]
touches: [flai/internal/components/flaiover_test.go]
after: [T-3115]
---
# T-3116 TestEveryListedFlaioverComponentCompilesInTheHarness checks each answer is its own vite's and says so when another server answered

## Work

`fetch` in `flai/internal/components/flaiover_test.go` reads only the status and the body, so I-0342's failure said that the setup did not compile when another server had answered. Make `fetch` give the answer's harness header too. Each check of the setup, each component, and the broken control then compares it with the dev server's token first. On a mismatch it fails saying another server answers on vite's port, with whether vite still runs, from `Err` and `Done`, and vite's last lines, rather than reporting a compile failure.

Waits for T-3115, which gives the token and the header this reads.

## Done when

- An answer without the dev server's token fails the test with a message naming another server on vite's port, vite's state, and its last lines, never a compile failure.
- The test passes on this tree with node and flaiover's dependencies installed, run with `flai test integration` or `go test -run TestEveryListedFlaioverComponentCompilesInTheHarness ./internal/components/` in `flai/`.

## Notes
