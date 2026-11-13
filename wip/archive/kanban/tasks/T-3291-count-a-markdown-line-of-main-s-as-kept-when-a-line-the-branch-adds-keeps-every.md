---
id: T-3291
type: task
nature: remediation
title: Count a markdown line of main's as kept when a line the branch adds keeps every word of it in order
status: done
parent: S-0767
owner: alex
created: 2026-11-13T11:48:15Z
updated: 2026-11-13T12:39:24Z
transitions:
  - to: ready
    at: 2026-11-13T12:05:34Z
    by: agent-S-0767
  - to: in-progress
    at: 2026-11-13T12:05:35Z
    by: agent-S-0767
  - to: done
    at: 2026-11-13T12:14:28Z
    by: agent-S-0767
stream: S-0767
tags: [flai, git]
touches: [flai/internal/storygit/undone.go, flai/internal/storygit/undone_test.go, design/adrs/0297-a-markdown-line-of-main-s-that-a-rebased-branch-removes-is-not-undone-when-a.md]
usage:
  source: log
  seconds: 533
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 19
      output: 4804
      cache_read: 1564864
      cache_write: 37979
      cost: 0.6705
  strategic:
    - kind: orchestrator
      seconds: 4
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 8
          cache_read: 20200
          cache_write: 757
          cost: 0.0061
---
# T-3291 Count a markdown line of main's as kept when a line the branch adds keeps every word of it in order

## Work

In `flai/internal/storygit/undone.go`, change the undone check so that, in a path ending `.md`, a line main added that the branch's net change removes is not undone when a line the branch adds in the same path keeps every word of it in order: the words of main's line, split with `strings.Fields`, are a subsequence of the added line's words. Each added line answers for one of main's lines at most. Other paths keep the multiset of whole lines (ADR-0155).

- Give `undoneLines` the branch's added lines and the path, or add a function beside it, and keep its multiset behaviour for lines that match exactly.
- `FindUndone` passes the added lines of `git diff onto HEAD -- <path>`, which `diffLines` already returns and the code drops.
- `RecordedUndone` passes the added lines of `git diff <main> story/<id> -- <path>`, so a recorded path the branch now extends drops out of the record.
- Update the comment at the top of the file and the doc comments of the functions you change.

In `flai/internal/storygit/undone_test.go`, add cases: an extended table row and an extended prose line in a `.md` path are not undone; a line that drops or reorders one of main's words is; the same extended line in a `.go` path is; two of main's lines both kept in one added line count one of them only. Add a `FindUndone` case with scripted diffs for the extension. This task waits for none.

## Done when

- [ ] `undoneLines`, or its replacement, keeps a markdown line of main's that a branch's added line extends, and the unit tests above pass with `flai test flai/internal/storygit/undone.go flai/internal/storygit/undone_test.go`
- [ ] `TestFindUndone` and `TestRecordedUndoneDropsAPathTheBranchNoLongerRemoves` still pass

## Notes
