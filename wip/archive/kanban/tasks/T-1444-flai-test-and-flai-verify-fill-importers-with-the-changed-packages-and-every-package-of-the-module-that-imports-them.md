---
id: T-1444
type: task
nature: improvement
title: "flai test and flai verify fill {importers} with the changed packages and every package of the module that imports them"
status: done
parent: S-0364
owner: alex
created: 2026-10-12T14:40:37Z
updated: 2026-10-12T15:35:21Z
transitions:
  - to: ready
    at: 2026-10-12T15:31:33Z
    by: agent-S-0364
  - to: in-progress
    at: 2026-10-12T15:31:34Z
    by: agent-S-0364
  - to: done
    at: 2026-10-12T15:35:21Z
    by: agent-S-0364
stream: S-0364
tags: [cli]
touches: [flai/internal/verify/importers.go, flai/internal/verify/importers_test.go, flai/internal/verify/select.go, flai/internal/verify/select_test.go, flai/internal/verify/verify.go]
after: [T-1443]
usage:
  source: log
  seconds: 227
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 30
      output: 9155
      cache_read: 1363457
      cache_write: 54515
      cost: 0.8101
---
# T-1444 flai test and flai verify fill {importers} with the changed packages and every package of the module that imports them

## Work

In a new `flai/internal/verify/importers.go`, work out from the checkout's `fs.FS` the packages that import the changed ones. Read the module path from the `go.mod` in the tier's `dir`. Read each `.go` file's imports below it with `go/parser` and `parser.ImportsOnly`, test files included, skipping the directories the go tool ignores as `nearestPackage` does. Build the reverse graph of the module's own packages and close it from the packages `{packages}` names. Ignore build constraints, so the set is never smaller than the go tool's. With no `go.mod`, `{importers}` is what `{packages}` is.

In `select.go`, fill `{importers}` in `fill` with that set, sorted as `./x/y`. In `allArgv`, fill it with `./...` as `{packages}` is filled. Update the placeholder comments in `select.go` and `verify.go`.

Reproduce I-0128 in `importers_test.go` with an `fstest.MapFS` module. Package `a` declares a struct. Package `b` imports `a`. Package `c` imports only `b`, and its test file compares two values of `a`'s type that `b` hands on. A change to `a` must select `./a`, `./b`, and `./c` under `{importers}`, and only `./a` under `{packages}`. In `select_test.go`, cover the argv `Select` and `SelectStory` give a tier with `{importers}`, the `--all` case, and a module without `go.mod`.

This task waits for T-1443, whose constant it fills.

## Done when

- A tier with `{importers}` selected by a change to a package runs on that package and on every package of the module that imports it, directly or transitively, through a test file too.
- The S-0338 case is reproduced by a test that fails without the reverse closure.
- `flai test flai/internal/verify` passes.

## Notes

Drafted by planner-S-0364. Parsing imports keeps `Select` pure over the checkout, as `packages` is. Running `go list` there would need the go tool at selection time.
