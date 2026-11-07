---
id: T-2595
type: task
nature: improvement
title: A change to flai's host API selects a tier that runs flaiover's tests of the real flai
status: done
parent: S-0592
owner: alex
created: 2026-11-07T22:02:45Z
updated: 2026-11-07T22:19:32Z
transitions:
  - to: ready
    at: 2026-11-07T22:06:31Z
    by: agent-S-0592
  - to: in-progress
    at: 2026-11-07T22:06:35Z
    by: agent-S-0592
  - to: done
    at: 2026-11-07T22:19:32Z
    by: agent-S-0592
stream: S-0592
tags: [flai, flaiover, tests]
touches: [system-flow.yaml, flai/internal/verify/select_test.go]
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 7
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 0
          output: 6
          cache_read: 19330
          cache_write: 573
          cost: 0.005
---
# T-2595 A change to flai's host API selects a tier that runs flaiover's tests of the real flai

## Work

`flai test` on `flai/internal/hostapi/` runs none of flaiover's tests that run the built flai, because the `vitest` tier picks tests by the TypeScript import graph (I-0304's first instance, `flaiover/src/lib/server/writes.test.ts`). Every test file that runs the built flai imports `flaiover/src/lib/server/testing.ts`.

Add a tier `flaiover-flai` to `system-flow.yaml` after `vitest`: `dir: flaiover`, command `["../scripts/flaiover-unit.sh", related, "--run", "--reporter=json", "src/lib/server/testing.ts"]`, `format: vitest-json`, paths `flai/internal/hostapi/**/*.go` and `flai/cmd/hostapi.go`. It is a fixed command, as `svelte-check` is. Add it to the `flaiover` tier's `covers`, so `--all` runs those tests once. Check that `flaiover-unit.sh` builds `bin/flai` before it runs them, and that `vitest related` on `testing.ts` selects every test file that imports it.

Add a repository-manifest test to `flai/internal/verify/select_test.go`, beside `TestRepositoryManifestTypeChecksChangedFlaioverFiles`: a change to `flai/internal/hostapi/writes.go` alone selects `flaiover-flai` with `src/lib/server/testing.ts` as its argument, and `--all` leaves it covered.

This task waits for no other: no other task of the first layer touches its files.

## Done when

- The new test fails without the tier and passes with it.
- `flai test flai/internal/hostapi/writes.go` runs `flaiover-flai`, and it runs `writes.test.ts` among its files.
- `flai check --strict` reports nothing on `system-flow.yaml`.

## Notes
