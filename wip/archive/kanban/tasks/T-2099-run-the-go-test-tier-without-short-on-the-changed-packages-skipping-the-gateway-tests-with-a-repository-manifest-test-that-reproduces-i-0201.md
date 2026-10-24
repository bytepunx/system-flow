---
id: T-2099
type: task
nature: remediation
title: Run the go-test tier without -short on the changed packages, skipping the gateway tests, with a repository-manifest test that reproduces I-0201
status: done
parent: S-0490
owner: alex
created: 2026-10-23T21:22:16Z
updated: 2026-10-24T00:03:31Z
transitions:
  - to: ready
    at: 2026-10-24T00:02:10Z
    by: agent-S-0490
  - to: in-progress
    at: 2026-10-24T00:02:11Z
    by: agent-S-0490
  - to: done
    at: 2026-10-24T00:03:31Z
    by: agent-S-0490
stream: S-0490
tags: [flai, tests, manifest]
touches: [system-flow.yaml, flai/internal/verify/select_test.go]
usage:
  source: log
  seconds: 80
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 11
      output: 2197
      cache_read: 853741
      cache_write: 17686
      cost: 0.3562
---
# T-2099 Run the go-test tier without -short on the changed packages, skipping the gateway tests, with a repository-manifest test that reproduces I-0201

## Work

In `system-flow.yaml`, change the `go-test` tier's command to:

```yaml
command: [go, test, "-json", "-race", "-count=1", "-skip", "^TestGateway", "{packages}"]
```

Drop `-short`, so that the tests the changed packages hold that skip under it (`realRepo` in `flai/internal/verify`, `TestReferenceIsCurrent` and `TestSettingsIndexIsComplete` in `flai/cmd`, `TestMonorepoIsClean`, `TestRepositoryLintsClean`, and the rest) run in `flai test <paths>` and in the close-out's `go-test` step while `integration` is left to the merge queue (ADR-0174). Add `-skip ^TestGateway`, as `scripts/integration.sh` does: `TestGatewayOpenRouter` and `TestGatewayLiteLLM` in `flai/internal/serve` make a paid call when the key is set, and only `-short` held them back. Change nothing else in the tier: `integration` still covers it, so `flai test --all` and a close-out with the queue off run the whole suite once, as today. `make test`, `scripts/test.sh`, and CI keep `-short`: they are people's.

Add to `flai/internal/verify/select_test.go`, beside `TestRepositoryManifestLintsChangedFlaioverFiles`, a test for I-0201 that reads this repository's manifest with `CheckoutTiers`, selects the tiers for `flai/internal/verify/paths_test.go` with `Select`, and asserts that the `go-test` tier's argv carries no `-short` and carries `-skip ^TestGateway`. Do not gate it under `testing.Short()`, and say in its comment why: it guards the manifest against `-short` coming back, so it must run under `make test` too; skip only when the monorepo is not present, as its neighbour does.

Run `flai test flai/internal/verify system-flow.yaml` and confirm the new test ran rather than skipped, and that the `go-test` step's duration grew (the real-git tests of `verify` add about five seconds). Run `flai check --strict` for `manifest.tests`.

## Done when

- `system-flow.yaml`'s `go-test` tier runs `go test` without `-short` and with `-skip ^TestGateway` on `{packages}`, and nothing else about the tier changed.
- The new test in `flai/internal/verify/select_test.go` fails on the tier's command as it was and passes on the new one, and runs under `-short`.
- `flai test flai/internal/verify` passes with the real-git tests run, and `flai check --strict` reports no `manifest.tests` finding.

## Notes

Measured on 2026-10-23 in the main checkout, `go test -race -count=1 ./...` with every package in parallel: without `-short` against with it, `cmd` 346 s against 310 s, `mcpserver` 317 s against 148 s (`TestStoryStartStartsAReadyStory`), `components` 55 s against 21 s (its two Vite tests), `check` 35 s against 11 s (`TestMonorepoIsClean`), `mdlint` 24 s against 1 s (`TestRepositoryLintsClean`), `mergequeue` 26 s against 7 s, `workitem` 23 s against 9 s, `taskdone` 18 s against 1 s, `verify` 13 s against 8 s; the other packages within 3 s. The gated tests skip on their own without `git`, `node`, or flaiover's `node_modules`.
