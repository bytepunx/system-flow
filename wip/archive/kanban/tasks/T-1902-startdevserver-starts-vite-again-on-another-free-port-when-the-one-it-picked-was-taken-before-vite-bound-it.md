---
id: T-1902
type: task
nature: remediation
title: StartDevServer starts Vite again on another free port when the one it picked was taken before Vite bound it
status: done
parent: S-0447
owner: alex
created: 2026-10-18T02:41:24Z
updated: 2026-10-18T02:56:44Z
transitions:
  - to: ready
    at: 2026-10-18T02:54:42Z
    by: agent-S-0447
  - to: in-progress
    at: 2026-10-18T02:54:43Z
    by: agent-S-0447
  - to: done
    at: 2026-10-18T02:56:44Z
    by: agent-S-0447
stream: S-0447
tags: [flai, components, tests]
touches: [flai/internal/components/devserver.go, flai/internal/components/devserver_test.go, design/system/component-library.md]
usage:
  source: log
  seconds: 121
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 19
      output: 4324
      cache_read: 1197261
      cache_write: 34773
      cost: 0.6042
---
# T-1902 StartDevServer starts Vite again on another free port when the one it picked was taken before Vite bound it

## Work

`freePort` in `flai/internal/components/devserver.go` listens on `127.0.0.1:0`, closes the listener, and hands the number to Vite, which binds it later with `--strictPort`. Another process can take the port in between, as in I-0162, where another agent's test run took it and `TestComponentsStartsOneDevServerAndAnswersItAgain` in `flai/internal/serve/components_test.go` failed with `bind: address already in use`. The window cannot be closed: Vite takes a port number, not a listener. So make the start survive it.

- In `StartDevServer`, when Vite exits before it answers and its last lines say the port is in use (Vite's `Port <n> is already in use`, or `address already in use`), pick another free port and start Vite again. Try at most three ports. Any other early exit fails at once, as now. After the last attempt, the error names each port tried and carries Vite's last lines.
- Let a test choose the ports: make the port picker a package variable, or a field the tests can set, that defaults to `freePort`.
- Add a test in `flai/internal/components/devserver_test.go` that reproduces I-0162 without a race. Hold a listener on a port, make the picker return that port first and then a free one, and check that `StartDevServer` answers on the second port and that two fake Vites ran. The fake Vite in `runFakeVite` already exits with `address already in use` when its port is taken. Add a second test: a picker that always returns the held port fails after three attempts, with each port named in the error.
- Add the retry to the dev server bullets in `design/system/component-library.md` § The dev server, and update the comment on `freePort`. S-0448, in progress, also changes this document: run `flai stream sync` before you edit it.

This task waits for nothing.

## Done when

- `StartDevServer` starts Vite again on a new port when the one it picked was taken, at most three ports in all.
- The two new tests pass, along with the existing dev server tests and the `flai/internal/serve` components tests.
- `design/system/component-library.md` describes the retry.
- `flai test` passes on the changed paths.

## Notes

Drafted by the planner. The supervisor in `flai/internal/serve/components.go` calls `StartDevServer` and needs no change. The serve test that failed in I-0162 gets the fix through it.
