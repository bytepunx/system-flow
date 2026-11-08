---
id: T-2685
type: task
nature: improvement
title: The markdown lint keeps each file's findings under .flai-cache, keyed by its stat and the lint configuration's hash
status: done
parent: S-0600
owner: alex
created: 2026-11-07T23:50:53Z
updated: 2026-11-08T00:52:16Z
transitions:
  - to: ready
    at: 2026-11-08T00:27:28Z
    by: agent-S-0600
  - to: in-progress
    at: 2026-11-08T00:27:33Z
    by: agent-S-0600
  - to: done
    at: 2026-11-08T00:52:16Z
    by: agent-S-0600
stream: S-0600
tags: [cli]
touches: [flai/internal/mdlint/config.go, flai/internal/mdlint/kept.go, flai/internal/mdlint/kept_test.go, flai/internal/mdlint/folder.go, flai/internal/check/check.go]
after: [T-2678]
usage:
  source: log
  seconds: 1438
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 45
      output: 258
      cache_read: 6043215
      cache_write: 670922
      cost: 2.4132
---
# T-2685 The markdown lint keeps each file's findings under .flai-cache, keyed by its stat and the lint configuration's hash

## Work

Give `Config` a hash in `flai/internal/mdlint/config.go`: of the configuration file's bytes, as `Load` reads them, and of the rule set `Rules()` names with a format version, so that a flai whose rules changed misses. Keep the findings of each file the folder lint lints in a file under `.flai-cache`, in a new `flai/internal/mdlint/kept.go`: the root, the hash, and per path the size, modification time, and findings. The folder lint in `folder.go` takes the cache's folder, lints only the files whose key differs, and writes the file again, to a temporary file renamed over the old one, when something changed. A missing, damaged, or foreign file, or one with another hash, is a miss; a failed write fails nothing; the `.flai-cache` folder is never made. `markdown()` in `flai/internal/check/check.go` passes `c.repo.CacheDir()`.

Second layer: it waits for T-2678, which writes the folder lint in `folder.go` and the call in `check.go` that this extends.

## Done when

- [ ] A second `flai check` lints no file under `wip` that did not change, and a test shows it.
- [ ] A test changes the lint configuration and sees every file linted again; another changes one file and sees only it linted again.
- [ ] Tests give a missing, a truncated, and another root's cache file and get the same findings as with none, with no error.
- [ ] A project without `.flai-cache` checks as before and has no cache folder afterwards.
- [ ] `go test -race ./internal/mdlint/ ./internal/check/` passes.

## Notes
