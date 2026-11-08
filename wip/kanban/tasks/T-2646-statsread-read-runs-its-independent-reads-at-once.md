---
id: T-2646
type: task
nature: improvement
title: statsread.Read runs its independent reads at once
status: backlog
parent: S-0607
owner: alex
created: 2026-11-07T23:49:07Z
updated: 2026-11-07T23:49:07Z
transitions: []
stream: S-0607
tags: [cli]
touches: [flai/internal/statsread/statsread.go, flai/internal/statsread/statsread_test.go]
after: [T-2637]
---
# T-2646 statsread.Read runs its independent reads at once

## Work

Cause 17: `Read` in `flai/internal/statsread/statsread.go` runs 11 reads one after another, 230 to 300 ms warm. Run them at once, so it takes about as long as the slowest read, about 90 ms.

- Run each read in its own goroutine with `sync.WaitGroup`: the items list, the activity documents, the roles, the threads, the issues, the board, the stories' commits, the conversations, the conflict log, the merge queue events, the digest, and the forge events when acceptance is by pull request. Use no new module: `errgroup` would make `golang.org/x/sync` a direct requirement in `flai/go.mod`.
- Each goroutine writes only its own variables; fill `opt` once all have ended. The manifest fields (`Fallback`, `Projects`, `Shared`, `Owner`) stay where they are, since they read memory only.
- Keep which error wins. When several reads fail, return the error of the read that comes first in today's order, with today's wording. Keep each warning `log` gets, with its `detail`. `slog` handlers are safe to call at once, but the order of the warnings may change. Note that in the function's comment.
- Check that what the reads share is safe to read at once: `repo.List` and the item store's folder locks, `repo.CacheDir`, `repo.Manifest`, and the runner `r`, which `execx.Timed` records phases on through `perf.Recorder`, which takes its own lock.
- Add a test to `statsread_test.go` that two failing reads still return the first one's error.

It waits for T-2637, whose golden test shows the report did not change.

## Done when

- `go test -race ./internal/statsread/` passes, the golden test of T-2637 among it.
- `Read` starts its reads at once, and its errors and warnings are those it gave before.

## Notes
