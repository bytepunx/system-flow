---
id: T-2656
type: task
nature: improvement
title: stats.get answers identical requests in flight with one computation
status: backlog
parent: S-0607
owner: alex
created: 2026-11-07T23:49:29Z
updated: 2026-11-07T23:49:29Z
transitions: []
stream: S-0607
tags: [cli]
touches: [flai/internal/hostapi/reads.go, flai/internal/hostapi/reads_test.go, flai/internal/hostapi/hostapi.go]
---
# T-2656 stats.get answers identical requests in flight with one computation

## Work

Criterion 2. Of the 120 slow `stats.get` calls in the serve log, 73 came within 3 s of the one before, each computed in full. The handler is `"stats.get"` in `flai/internal/hostapi/reads.go`, built by `readMethods`, which `MethodsFor` in `hostapi.go` calls.

- Key a request by the project's root and its parameters after their defaults: `since` (with `metrics.DefaultWindow` for empty), `type`, `by`, and `bucket`.
- While a computation with that key runs, a request with the same key waits for it and gets its answer, the same `Ran` with its events. Once it ends, the key is gone: a request after it computes again, so the answer is never stale.
- Write it as a small in-flight map under a mutex, owned by the methods `readMethods` builds, so each `MethodsFor` has its own. Use no new module: `golang.org/x/sync/singleflight` would make `golang.org/x/sync` a direct requirement in `flai/go.mod`.
- The first request's context must not cancel the computation for the others. Run the shared work with `context.WithoutCancel`, and let each waiting request stop waiting when its own context ends. Record the `stats.compute` phase on each request: a waiting request records the time it waited.
- Add a test to `reads_test.go`: two requests with the same parameters, started together, get one computation, counted through a hook the test sets. Two requests with different parameters get two computations. A request after the first ends computes again.

This task changes other files than the reads and sections tasks, so it waits for none and runs beside T-2637.

## Done when

- `go test -race ./internal/hostapi/ -run Stats` passes with the three cases above.
- The `stats.get` validation table tests in `reads_test.go` still pass.

## Notes
