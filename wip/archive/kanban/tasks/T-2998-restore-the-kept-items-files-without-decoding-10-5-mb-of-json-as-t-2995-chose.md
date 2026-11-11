---
id: T-2998
type: task
nature: improvement
title: Restore the kept items files without decoding 10.5 MB of JSON, as T-2995 chose
status: done
parent: S-0676
owner: alex
created: 2026-11-10T10:31:01Z
updated: 2026-11-11T02:17:01Z
transitions:
  - to: ready
    at: 2026-11-11T01:59:43Z
    by: system-flow
  - to: in-progress
    at: 2026-11-11T01:59:44Z
    by: system-flow
  - to: done
    at: 2026-11-11T02:17:01Z
    by: agent-S-0676
stream: S-0676
tags: [flai, performance]
touches: [flai/internal/workitem/kept.go, flai/internal/workitem/kept_test.go, flai/internal/workitem/store.go, flai/internal/workitem/repo.go, flai/internal/workitem/kept_bench_test.go]
after: [T-2995]
usage:
  source: log
  seconds: 1037
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 67
      output: 18866
      cache_read: 9486337
      cache_write: 175721
      cost: 3.3792
---
# T-2998 Restore the kept items files without decoding 10.5 MB of JSON, as T-2995 chose

## Work

Build the remedy T-2995 recorded in the narrative's `## Decisions`. This task waits for T-2995 because its figures choose what to build. For example, encode the kept items file with `encoding/gob` in `flai/internal/workitem/kept.go`. If the figures say so, also restore the six item folders at once in `Repo.List` (`repo.go`), or make `readAs` lazily in `store.go`.

Keep what S-0600's fix promises (`design/system/server-performance.md` § Parses and lint findings kept between processes (S-0600)):

- a restored item equals `ReadItem`'s, its unknown fields and `readAs` included;
- a missing, damaged, foreign, or older file is a miss, never an error;
- the file is written whole through a rename, only into a `.flai-cache` that exists.

Raise `keptVersion`, so that a file the JSON format wrote is a miss. Rename the file in `keptItemsPath` if its extension no longer fits.

Add tests in `kept_test.go`:

- a version-1 JSON file is a miss, and the next list writes it again in the new format;
- the new format restores `TestKeptItemsRestoreAsReadItemReadsThem`'s items and this repository's (`TestKeptItemsRestoreThisRepositorysItems`) exactly.

## Done when

- `flai test` passes on `flai/internal/workitem/kept.go`, `kept_test.go`, `store.go`, and `repo.go`, with the new tests among them.
- T-2995's benchmark shows the fresh list no longer spending most of its time decoding the kept items files, with its figure before and after recorded in the narrative's `## Decisions`.

## Notes
