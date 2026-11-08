---
id: T-2682
type: task
nature: improvement
title: A flai process keeps parsed items under .flai-cache and reads again only the files whose key differs
status: done
parent: S-0600
owner: alex
created: 2026-11-07T23:50:45Z
updated: 2026-11-08T00:51:32Z
transitions:
  - to: ready
    at: 2026-11-08T00:27:20Z
    by: agent-S-0600
  - to: in-progress
    at: 2026-11-08T00:27:24Z
    by: agent-S-0600
  - to: done
    at: 2026-11-08T00:51:32Z
    by: agent-S-0600
stream: S-0600
tags: [cli]
touches: [flai/internal/workitem/kept.go, flai/internal/workitem/kept_test.go, flai/internal/workitem/store.go, flai/internal/workitem/transition.go, flai/cmd/promote_test.go, flai/cmd/verify_test.go, flai/internal/check/check_test.go, flai/internal/workitem/repo.go]
after: [T-2674]
usage:
  source: log
  seconds: 1447
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 59
      output: 229
      cache_read: 8778983
      cache_write: 400134
      cost: 3.2991
---
# T-2682 A flai process keeps parsed items under .flai-cache and reads again only the files whose key differs

## Work

Keep the parses `store.list` makes in a file under `(*Repo).CacheDir()`, as `flai/internal/storygit/committed.go` keeps the stories' commits: the file holds the repository root, a format version, and, per folder and file name, the key (size and modification time) and the item as `ParseItem` left it. A process with nothing in memory loads it once, reuses an entry whose key matches the file's stat, parses the rest on the pool, and writes the file again only when something changed, to a temporary file beside it renamed over the old one. A missing, damaged, or foreign file (another root or another version) is a miss; a failed write fails nothing; the `.flai-cache` folder is never made.

What the cache must not change, from the store as it is:

- `Path` and `Archived` come from where the file is, not from the cache.
- `readAs` is set again from the restored item (`itemlock.go` merges on it, I-0175), and `Unknown` fields are kept so `WarnUnknown` still warns once per process.
- An entry is trusted only when the file was stat'ed at least `racyWindow` after its modification time, as the in-memory `kept` is.

In `flai/internal/workitem/transition.go`, `TransitionAll` lists with the archive at line 70 and lists again without it at line 130 for `WriteIndex`: take the second from what changed in the first rather than listing again, if the index then stays byte for byte the same.

Second layer: it waits for T-2674, which reshapes `store.list` in `store.go`, the function this hooks into.

## Done when

- [ ] A second process lists the items without parsing any file that did not change, and a test shows it by the parse counter.
- [ ] A test changes one file's size or modification time and sees only that file parsed again.
- [ ] Tests give the loader a missing, a truncated, and another root's cache file and get the same items as with none, with no error.
- [ ] A project without `.flai-cache` lists as before and has no cache folder afterwards.
- [ ] `go test -race ./internal/workitem/` passes.

## Notes
