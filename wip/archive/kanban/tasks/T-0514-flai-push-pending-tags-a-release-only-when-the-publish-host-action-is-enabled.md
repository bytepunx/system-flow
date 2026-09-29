---
id: T-0514
type: task
nature: remediation
title: flai push --pending tags a release only when the publish host action is enabled
status: done
parent: S-0144
owner: alex
created: 2026-09-29T01:14:19Z
updated: 2026-09-29T01:16:48Z
transitions:
  - to: ready
    at: 2026-09-29T01:14:25Z
    by: agent-S-0144
  - to: in-progress
    at: 2026-09-29T01:14:25Z
    by: agent-S-0144
  - to: done
    at: 2026-09-29T01:16:48Z
    by: agent-S-0144
stream: S-0144
tags: []
touches: [flai/cmd/push.go, flai/cmd/push_test.go, flai/internal/hostapi/writes.go]
---
# T-0514 flai push --pending tags a release only when the publish host action is enabled

## Work

- Add the `publish` host action to `flai/internal/hostapi` (`ActionPush`'s sibling), off by default like every host action, and reword `push`'s meaning so it no longer promises a release.
- `flai push --pending` computes, applies, and tags the pending release before it pushes only when `publish` is enabled for the project in the configuration flai reads; otherwise it pushes the merged commits and any tags already made, and tags nothing. `--dry-run` shows a release plan only when it would tag one. The command's help says so.
- Tests: with `publish` off, a push after an acceptance creates no tag and no version bump; with it on, it does as before.

## Done when

- `go test` for `flai/cmd` and `flai/internal/hostapi` passes with the new tests, and `flai serve actions` lists `publish`.

## Notes
