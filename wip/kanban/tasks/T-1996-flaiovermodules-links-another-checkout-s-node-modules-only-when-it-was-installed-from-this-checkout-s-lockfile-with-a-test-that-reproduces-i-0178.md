---
id: T-1996
type: task
nature: remediation
title: flaioverModules links another checkout's node_modules only when it was installed from this checkout's lockfile, with a test that reproduces I-0178
status: done
parent: S-0482
owner: alex
created: 2026-10-21T12:08:38Z
updated: 2026-10-21T12:22:13Z
transitions:
  - to: ready
    at: 2026-10-21T12:20:19Z
    by: agent-S-0482
  - to: in-progress
    at: 2026-10-21T12:20:20Z
    by: agent-S-0482
  - to: done
    at: 2026-10-21T12:22:13Z
    by: agent-S-0482
stream: S-0482
tags: [flai, tests, components]
touches: [flai/internal/components/integration_test.go]
usage:
  source: log
  seconds: 113
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 13
      output: 99
      cache_read: 928160
      cache_write: 11444
      cost: 0.3183
---
# T-1996 flaioverModules links another checkout's node_modules only when it was installed from this checkout's lockfile, with a test that reproduces I-0178

## Work

`flaioverModules` in `flai/internal/components/integration_test.go` takes the first `flaiover/node_modules` that has vite, this checkout's or the main checkout's, without asking whether it was installed from this checkout's `flaiover/pnpm-lock.yaml`. A stale one compiles the components against missing packages, and the test fails in `flaiover_test.go` with a bare 500 instead of naming the cause.

- Factor the choice into a function of the candidates and this checkout's lockfile, so a test can call it on temporary folders.
- Accept a candidate only when it has vite and the lockfile pnpm recorded in it, `node_modules/.pnpm/lock.yaml`, matches this checkout's `flaiover/pnpm-lock.yaml`.
- When a candidate has vite but its lockfile differs, fail rather than skip, naming the folder and `scripts/flaiover-install.sh --if-needed`, so a stale install is never mistaken for a broken component. Skip only when no candidate has vite, as now.
- Add a test that reproduces I-0178: this checkout has no `node_modules`, and the other has vite and a `lock.yaml` that differs from this checkout's lockfile, as the main checkout's did once S-0458 added `@lucide/svelte`. The choice refuses it and names the install.
- Add a test that this checkout's own install wins over the other's.

It waits for no task: it shares no path with the script task. Leave `props_test.go`'s `installed()` as it is: it reads only svelte's compiler, which a new dependency does not change.

## Done when

- The two new tests pass, and the reproducing one fails against the old choice.
- `TestEveryListedFlaioverComponentCompilesInTheHarness` and `TestServeCarriesHotUpdatesOfTheComponentAndTheStylesheet` still pass with a current install.
- `flai test flai/internal/components/integration_test.go` passes.

## Notes
