---
id: T-1490
type: task
nature: remediation
title: flai check counts an unreadable narrative as listed in the index and tests narrative.front-matter
status: done
parent: S-0373
owner: alex
created: 2026-10-12T17:28:52Z
updated: 2026-10-12T18:16:39Z
transitions:
  - to: ready
    at: 2026-10-12T18:14:40Z
    by: agent-S-0373
  - to: in-progress
    at: 2026-10-12T18:14:41Z
    by: agent-S-0373
  - to: done
    at: 2026-10-12T18:16:39Z
    by: agent-S-0373
stream: S-0373
tags: [cli]
touches: [flai/internal/check/check.go, flai/internal/check/check_test.go]
usage:
  source: log
  seconds: 118
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 9
      output: 1692
      cache_read: 655251
      cache_write: 19009
      cost: 0.317
---
# T-1490 flai check counts an unreadable narrative as listed in the index and tests narrative.front-matter

## Work

`narratives()` in `flai/internal/check/check.go` reports a narrative that does not parse as `narrative.front-matter` and skips it before it marks the stream as having a narrative. Once the index lists that file as an unreadable row, the index check would add a second, wrong finding: `index.md lists S-nnnn which has no narrative in agents/`.

- Mark the stream from the file's name before parsing, so an unreadable narrative counts as present for the `narrative.index` checks.
- Keep `narrative.front-matter` an Error. Its message keeps the parse error and adds that the index lists the file as unreadable until its front matter is restored.
- Add a test in `flai/internal/check/check_test.go`. Put an empty `S-nnnn.md` in a project whose index lists it as `| [S-nnnn](S-nnnn.md) | unreadable, see flai check | | | |`. Check reports one `narrative.front-matter` on that file and no `narrative.index`. No test covers that rule today.

It waits for no task. It relies only on the row format the story's plan fixes, which the index task writes at the same time, and it shares no path with it.

## Done when

- The new test passes, and fails without the change to `narratives()`.
- `flai test` passes on `flai/internal/check/check.go` and `flai/internal/check/check_test.go`.
- `TestBadFixtureFindings` passes unchanged.

## Notes

The rule is at `check.go` in `narratives()`, about line 701; the index checks follow it.
