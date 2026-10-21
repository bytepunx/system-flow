---
id: T-1923
type: task
nature: improvement
title: flai task done refuses only when the newest flai tag's task-done rules number is above the running flai's, and keeps the minor rule for a tag without one
status: done
parent: S-0445
owner: alex
created: 2026-10-20T20:42:08Z
updated: 2026-10-20T20:54:19Z
transitions:
  - to: ready
    at: 2026-10-20T20:52:14Z
    by: agent-S-0445
  - to: in-progress
    at: 2026-10-20T20:52:14Z
    by: agent-S-0445
  - to: done
    at: 2026-10-20T20:54:19Z
    by: agent-S-0445
stream: S-0445
tags: [flai, task-done]
touches: [flai/internal/release/outdated.go, flai/internal/release/outdated_test.go, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go]
after: [T-1921]
usage:
  source: log
  seconds: 125
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 35
      output: 9767
      cache_read: 2711061
      cache_write: 71184
      cost: 1.2581
---
# T-1923 flai task done refuses only when the newest flai tag's task-done rules number is above the running flai's, and keeps the minor rule for a tag without one

## Work

- In `flai/internal/release/outdated.go`, make `FlaiOutdated` also read the newest `flai/vX.Y.Z` tag's annotation (`git tag -l --format=%(contents) <tag>`) and carry its `task-done-rules` number on `Outdated`, absent when the tag has no such line.
- In `flai/internal/taskdone/taskdone.go`, change `refuseOutdated` so that when the newest tag carries a rules number, the close is refused only when that number is above `buildinfo.TaskDoneRules`; when it carries none, `BehindRelease`'s minor rule applies as today. Keep the refusal's text naming both versions, and add both rules numbers.
- Tests that reproduce I-0157: a running flai a minor version behind a newest tag that carries the same rules number closes the task; one whose tag carries a higher number refuses, changing nothing; a tag without the line refuses on the minor rule as `TestRunRefusesAMinorOrMajorVersionBehind` does now. Add the annotation case to `TestFlaiOutdated`.

It waits for T-1921 for the constant. It shares no path with T-1922, so the two run together.

## Done when

- A minor release that does not raise the rules number no longer stops `flai task done`.
- A release that raises it still does, and an old tag without the line behaves as before.
- The tests pass with `flai test` on the four paths.

## Notes

`FlaiOutdated` caches its answer for a minute per project and running version; the rules number rides on the same answer.
