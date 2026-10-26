---
id: T-2136
type: task
nature: improvement
title: storygit.Committed keeps the stories' commits per repository root for the life of the process, keyed by the branch tips, and reads the log again only when a tip moves
status: done
parent: S-0500
owner: alex
created: 2026-10-26T04:26:38Z
updated: 2026-10-26T04:34:27Z
transitions:
  - to: ready
    at: 2026-10-26T04:29:12Z
    by: agent-S-0500
  - to: in-progress
    at: 2026-10-26T04:29:16Z
    by: agent-S-0500
  - to: done
    at: 2026-10-26T04:34:27Z
    by: agent-S-0500
stream: S-0500
tags: [flai]
touches: [flai/internal/storygit/committed.go, flai/internal/storygit/committed_test.go, flai/internal/statsread/statsread_test.go, flai/internal/hostapi/timing_test.go]
usage:
  source: log
  seconds: 311
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 40
      output: 11380
      cache_read: 1909605
      cache_write: 71036
      cost: 1.0853
---
# T-2136 storygit.Committed keeps the stories' commits per repository root for the life of the process, keyed by the branch tips, and reads the log again only when a tip moves

## Work

In `flai/internal/storygit/committed.go`, keep what `Committed` reads per repository root between calls, as `release.readHistory` keeps the release history (`flai/internal/release/history.go`, S-0157): a package-level map from root to the kept result under a mutex, each entry holding its key and the per-story files. The key is the output of `git show-ref --heads` in the main checkout, the tips of the main branch and every story branch, read at every call. While the key equals the kept one, `Committed` returns a deep copy of the kept map and starts no `git log`; when it differs, or nothing is kept, it runs the whole log as today, parses it, and keeps the result under the new key. Keep the lookups of git and the main branch as they are, so the error messages stand; a read that fails keeps nothing. Keep the `wip` folder prefix out of the kept files as today.

Tests in `flai/internal/storygit/committed_test.go`, with a counting runner like `release/history_test.go`'s that records each git command: a second call while nothing changed starts `show-ref` and no `log`, and answers the same map; a commit on a story branch or on main changes the answer; a rebase that moves a story tip to a commit the old tips do not reach is read whole; two roots are kept apart; the caller's copy may be changed without changing the next answer. A test in `flai/internal/statsread/statsread_test.go` that `Read` called twice in one process starts one `git log`, beside `TestReadInARepositoryReadsTheCommits`. In `flai/internal/hostapi/timing_test.go`, add `stats.get` to the read methods with the phases it marks if that test can hold it, so a warm `stats.get` is seen to mark `exec.git.show-ref` and not `exec.git.log`.

This is the first task: it waits for nothing.

## Done when

- A second `Committed` on the same root with unchanged branch tips starts `git show-ref --heads` and no `git log`, and answers the same files per story, in a test with a counting runner
- A commit on a story branch, a commit on the main branch, and a story branch rebased onto a new base each make the next call read the log again and answer the new files, in tests
- `statsread.Read` twice in one process starts one `git log`, in a test
- Existing tests in `storygit`, `statsread`, `hostapi`, and `cmd` pass unchanged in what they assert, and `flai test` on the changed paths passes

## Notes

The memo lives in `storygit`, not in `statsread`, so `flai accept`'s call after its acceptance commit gains it too; that call sees a changed key and reads again, as it must. A whole read on change is chosen over reading only the commits added since the kept tips: `flai stream sync` rebases story branches, which moves their tips to commits the kept history does not reach, and the whole read costs 0.67 s on this repository only when a tip moved.
