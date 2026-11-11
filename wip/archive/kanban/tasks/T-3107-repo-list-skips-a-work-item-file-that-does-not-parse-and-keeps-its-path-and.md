---
id: T-3107
type: task
nature: remediation
title: Repo.List skips a work item file that does not parse and keeps its path and error for the check
status: done
parent: S-0716
owner: alex
created: 2026-11-11T08:56:33Z
updated: 2026-11-11T14:54:46Z
transitions:
  - to: ready
    at: 2026-11-11T14:44:56Z
    by: agent-S-0716
  - to: in-progress
    at: 2026-11-11T14:44:56Z
    by: agent-S-0716
  - to: done
    at: 2026-11-11T14:54:46Z
    by: agent-S-0716
stream: S-0716
tags: [flai, workitem]
touches: [flai/internal/workitem/store.go, flai/internal/workitem/repo.go, flai/internal/workitem/store_test.go]
usage:
  source: log
  seconds: 590
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 38
      output: 9880
      cache_read: 1860942
      cache_write: 74812
      cost: 1.0181
---
# T-3107 Repo.List skips a work item file that does not parse and keeps its path and error for the check

## Work

Today `lookAll` in `flai/internal/workitem/store.go` drops a file whose parse failed and then fails its whole folder with that file's error. `Repo.List` in `repo.go` passes the error on, so one file that does not parse, such as I-0344's 11-byte `PLACEHOLDER` task, fails every caller.

Make `lookAll` and `Repo.List` keep loading. Leave the file out of the list and record its path and parse error where callers can read them, for example `Repo.Unreadable()`. Keep `Get` of that file's ID failing with its parse error. `NextID` globs file names, so the unreadable file's ID is still never given out again: keep it that way.

Replace `TestAListOfADamagedItemFailsWithItsParseError` in `store_test.go`, which pins today's hard failure, with a test that reproduces I-0344: a project whose tasks folder holds a `PLACEHOLDER` file with no front matter. Its list returns every other item and names that file with `no front matter`.

This task waits for nothing. It is the first layer.

## Done when

- `Repo.List` on a project with a file that has no front matter returns the other items with no error, and the file's path and error can be read from the repository.
- `Get` of that ID and `NextID` behave as before.
- The new test passes and the old test is gone: `flai test flai/internal/workitem/store.go` passes.

## Notes
