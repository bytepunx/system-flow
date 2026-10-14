---
id: T-1574
type: task
nature: remediation
title: flai stream restore S-nnnn --path commits main's version of a path the branch undoes onto the story branch, and clears it from the record
status: done
parent: S-0371
owner: alex
created: 2026-10-14T20:58:37Z
updated: 2026-10-14T22:09:18Z
transitions:
  - to: ready
    at: 2026-10-14T21:43:55Z
    by: agent-S-0371
  - to: in-progress
    at: 2026-10-14T21:43:56Z
    by: agent-S-0371
  - to: done
    at: 2026-10-14T22:09:18Z
    by: agent-S-0371
stream: S-0371
tags: [cli, go]
touches: [flai/cmd/stream_restore.go, flai/cmd/stream_restore_test.go, flai/cmd/stream.go]
after: [T-1477]
usage:
  source: log
  seconds: 1522
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 88
      output: 32887
      cache_read: 5188616
      cache_write: 243199
      cost: 3.0311
---
# T-1574 flai stream restore S-nnnn --path commits main's version of a path the branch undoes onto the story branch, and clears it from the record

## Work

It waits for T-1477, whose record names the paths a branch undoes. It shares no path with T-1479, so the two run together.

- In a new `flai/cmd/stream_restore.go`, `flai stream restore <story-id> --path <path>` (repeatable), after the pattern of `stream_diff.go`. In the story's worktree it checks out main's version of each path (`git checkout <main> -- <path>`, or `git rm` when main has no such file), commits it on `story/S-nnnn` with a message naming the path and I-0124's check, and answers what it restored. It refuses, changing nothing, while the worktree has uncommitted changes, while a rebase is unfinished, for a story that is not in progress or in review, and for a path the record does not hold, naming the paths it does. `--json` gives the commit and the paths.
- After the commit, the record's reader no longer lists the path, since the branch no longer removes main's lines.
- In `flai/cmd/stream.go`, register the command in `newStreamCmd`.
- In a new `flai/cmd/stream_restore_test.go`, with the sync tests' helpers: I-0124's sequence, a sync that records the path, then `flai stream restore` commits main's lines back and the path leaves the record; and the refusals for uncommitted changes and for a path the record does not hold.

## Done when

- `flai stream restore` on I-0124's sequence commits main's version of the path on the story branch and the record no longer lists it.
- Its refusals change nothing.
- `flai test` on the three files passes.

## Notes

`docs/users/flai-reference.md`, generated from the help, is regenerated in the docs task, once, after every help changed.
