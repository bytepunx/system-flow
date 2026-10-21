---
id: T-1922
type: task
nature: improvement
title: flai release writes the task-done rules number into each flai tag's annotation
status: cancelled
parent: S-0445
owner: alex
created: 2026-10-20T20:41:58Z
updated: 2026-10-20T20:44:45Z
transitions:
  - to: cancelled
    at: 2026-10-20T20:44:45Z
    by: system-flow
stream: S-0445
tags: [flai, release]
touches: [flai/internal/release/release.go, flai/internal/release/release_test.go]
after: [T-1921]
---
# T-1922 flai release writes the task-done rules number into each flai tag's annotation

## Work

- In `flai/internal/release/release.go`, `Tag` and `TagPending` already make annotated tags (`git tag -a <tag> -m <msg>`). For a `flai/vX.Y.Z` tag, add a line `task-done-rules: N` to the message, with N from `buildinfo.TaskDoneRules`. Leave other components' tags as they are.
- Test in `release_test.go` that a flai tag's annotation carries the line and a flaiover tag's does not.

It waits for T-1921 for the constant. It shares no path with T-1923, so the two run together.

## Done when

- A release made by `flai release` or `flai push --pending` annotates each flai tag with `task-done-rules: N`.
- The tests pass with `flai test` on the two paths.

## Notes

The number recorded is that of the flai that tags the release. In this repository that is the tree's build, so it matches the code being released.
- 2026-10-20T20:44:45Z: moved to cancelled: The rules number is read from flai/internal/buildinfo/task-done-rules.txt at the newest flai tag instead of written into its annotation: a release tagged by the installed flai, as the board's Publish does, would record that flai's number, not the released tree's
