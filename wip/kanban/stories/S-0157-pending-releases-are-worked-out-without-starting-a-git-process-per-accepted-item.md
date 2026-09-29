---
id: S-0157
type: story
nature: improvement
title: Pending releases are worked out without starting a git process per accepted item
status: review
parent: E-0012
owner: alex
created: 2026-09-29T07:00:28Z
updated: 2026-09-29T19:37:45Z
transitions:
  - to: ready
    at: 2026-09-29T19:18:25Z
    by: alex
  - to: in-progress
    at: 2026-09-29T19:18:42Z
    by: agent-S-0157
  - to: review
    at: 2026-09-29T19:37:45Z
    by: agent-S-0157
tags: [cli]
topics: [server-side, back-end]
touches: [flai/internal/release, design/system/server-performance.md, design/system/flai-cli.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1170
  models:
    - model: claude-opus-5-5
      input: 178
      output: 92833
      cache_read: 14445640
      cache_write: 220602
      cost: 6.5113
---
# S-0157 Pending releases are worked out without starting a git process per accepted item

## Goal

Every board, from the dashboard or MCP, works out which accepted items are not yet published by starting 25 git processes one after another (13 `diff-tree`, 6 `log`, 6 `tag` on this repository): 75 ms each time, and more as items are accepted between releases. Cause 2 of `design/system/server-performance.md`. Work it out with a few git processes, or once per change of `HEAD` and the tags, and keep the answer.

## Acceptance criteria
- [x] `release.pending` in a `board.get` event takes under 15 ms on this repository, with at most 3 `exec.git` steps, or none while `HEAD` and the tags are unchanged.
- [x] `release.PendingIDs` answers exactly what it answers today, for a repository with pending items in several components, none, and a release just cut: a behaviour test compares them.

## Tasks
- T-0565 Pending reads the history in a few git processes and keeps it while HEAD and the tags are unchanged
- T-0566 Measure release.pending on this repository and record it in the design

## Notes

Measured by S-0152: `release.pending=75.4 exec.git.diff-tree=34.7x13 exec.git.log=27.8x6 exec.git.tag=11.3x6`. Proposed: one `git log --name-only` over the range since the oldest component tag instead of a `diff-tree` per commit, and a cache keyed by `git rev-parse HEAD` and the tag list.

Delivered (S-0157): `release.Pending` reads git through a history kept per repository root (`flai/internal/release/history.go`); figures in `design/system/server-performance.md` under cause 2, design in `design/system/flai-cli.md` § Internal structure.

How the first criterion was verified: `board.get` was called six times in one process through `hostapi.Methods`, the table `flai serve` answers with, timed by the same `perf` recorder that writes its `request answered` event: `release.pending` took 57 ms with four `exec.git` steps the first time, then 8.2 to 9.5 ms with one, `show-ref`. Also `flai mcp` built from the story branch, over stdio on this repository, asked `board` six times in one process; its `request answered` events showed `release.pending` 7.9 to 10.2 ms (one look at 25.7 ms while the host was loaded and `repo.list` also tripled) with one `exec.git` step, `show-ref`. The MCP `board` and the host API's `board.get` call the same `release.PendingIDs` on a long-lived process. A second `flai serve` was not started beside the operator's, so no `board.get` event was read from a running `flai serve`. In a clone, the look after an acceptance took 10.8 ms with two steps, and after a publish tag 8.3 ms with one. The first look in a process reads the whole history: four steps, about 60 ms, against 242 to 259 ms and 69 steps before. The criterion's "none while HEAD and the tags are unchanged" is met as "at most 3": the check that they are unchanged is itself the one step.

How the second was verified: `TestPendingFromTheKeptHistoryAnswersAsTheProcessesDid` compares `Pending` and `PendingIDs` with the per-process implementation, kept in `history_test.go`, as one repository changes: nothing accepted; items in several components (one refused for touching two with no tag); research and an epic; a release just cut; accepted after the release; HEAD moved back, then forward again; an amended commit; a tag off HEAD's history. `TestPendingOnThisRepositoryAnswersAsTheProcessesDid` (integration tier) does the same on this repository, and it was also run by hand on clones of it at five earlier release points, one with items pending in all three components. `TestPendingStartsOneGitProcessWhileNothingChanged` counts the processes.
