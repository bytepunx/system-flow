---
id: T-2801
type: task
nature: improvement
title: Test that no non-test Go file under flai is over 2,000 lines
status: done
parent: S-0631
owner: alex
created: 2026-11-09T22:05:19Z
updated: 2026-11-10T06:45:40Z
transitions:
  - to: ready
    at: 2026-11-10T06:39:17Z
    by: agent-S-0631
  - to: in-progress
    at: 2026-11-10T06:39:19Z
    by: agent-S-0631
  - to: done
    at: 2026-11-10T06:45:40Z
    by: agent-S-0631
stream: S-0631
tags: [cli]
touches: [flai/internal/check/sourcesize_test.go]
after: [T-2794, T-2795, T-2796, T-2797, T-2798, T-2799]
usage:
  source: log
  seconds: 381
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 10
      output: 1923
      cache_read: 1369759
      cache_write: 41010
      cost: 0.5458
---
# T-2801 Test that no non-test Go file under flai is over 2,000 lines

## Work

Write the test that reproduces I-0298 and keeps it fixed. Put it in a new file, `flai/internal/check/sourcesize_test.go`, beside `TestMonorepoIsClean` in `check_test.go`, which already reads the repository from the package's folder. Call it, for example, `TestNoSourceFileIsOver2000Lines`.

- Walk the flai module root, `../..` from the package's folder.
- Leave out `_test.go` files and `testdata` folders.
- Count the lines of each `.go` file.
- Fail with one line per file over 2,000, giving its path and its count, so a story that grows a file past the limit is told which one.
- Keep the limit in one named constant. `design/conventions/code-quality.md` asks for SRP at every level, and 2,000 is the line I-0298 draws.
- Do not skip it under `-short`. It reads files only and takes milliseconds, so `flai test` on the check package runs it, and so does every close-out's integration tier.

It waits for the six splits, T-2794 to T-2799: until each is done, the test fails on that file. It shares no path with them.

Test files are out of its scope. Seven are over 2,000 lines, and the plan's thread proposes them as a story of their own.

## Done when

- The test fails on main as it was before this story, naming the six files and their counts. Check by running it with the splits stashed, or by lowering the constant for a moment.
- It passes on the story's branch.
- `flai test flai/internal/check/sourcesize_test.go` passes.

## Notes
