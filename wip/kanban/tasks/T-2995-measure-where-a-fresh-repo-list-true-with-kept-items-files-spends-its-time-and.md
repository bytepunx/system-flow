---
id: T-2995
type: task
nature: research
title: Measure where a fresh Repo.List(true) with kept items files spends its time, and choose the remedy
status: backlog
parent: S-0676
owner: alex
created: 2026-11-10T10:30:49Z
updated: 2026-11-10T10:30:49Z
transitions: []
stream: S-0676
tags: [flai, performance]
touches: [flai/internal/workitem/kept_bench_test.go]
---
# T-2995 Measure where a fresh Repo.List(true) with kept items files spends its time, and choose the remedy

## Work

Write `BenchmarkFreshListWithKeptItems` in `flai/internal/workitem/kept_bench_test.go`. It copies a corpus the size of this repository's items into `b.TempDir()`, about 3,300 items with the archive at 6.0 MB of tasks and 3.8 MB of stories in the kept files, as `TestKeptItemsRestoreThisRepositorysItems` reads this repository. It lists once so the kept items files are written. Then, for each round, it drops the in-process store and times `Repo.List(true)` as a fresh process would.

Split the time by phase with sub-benchmarks or `-cpuprofile`:

- reading the files;
- `readKeptItems`'s `json.Unmarshal` of the whole file into `keptRecord`s;
- `restoreItem`'s `json.Unmarshal` of each item and its `it.Marshal()` for `readAs`;
- the stats in `lookup`;
- the clone and sort in `List`.

Time the candidates the issue's instance names:

- `encoding/gob` (or another binary encoding from the standard library) in place of JSON;
- decoding only the entries whose key matches;
- restoring the six folders at once rather than one after another in `Repo.List`;
- making `readAs` lazily.

Do not list the archive less in `flai move`'s `Follow`: that cuts a move, not a list, and the plan's thread proposes it apart.

This task waits for nothing.

## Done when

- `go test -run '^$' -bench BenchmarkFreshListWithKeptItems ./internal/workitem` runs and reports each phase.
- The narrative's `## Decisions` records the phase figures, each candidate's figure, and the remedy chosen, with its reason.

## Notes
