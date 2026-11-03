---
id: T-2446
type: task
nature: feature
title: The preflight test's markdown tier fixes a file without sed -i, so the fix runs under BSD sed on macOS as under GNU sed
status: done
parent: S-0572
owner: alex
created: 2026-11-02T23:59:52Z
updated: 2026-11-03T00:05:14Z
transitions:
  - to: ready
    at: 2026-11-03T00:00:31Z
    by: system-flow
  - to: in-progress
    at: 2026-11-03T00:00:32Z
    by: system-flow
  - to: done
    at: 2026-11-03T00:05:14Z
    by: system-flow
stream: S-0572
tags: []
touches: [flai/internal/mergequeue/preflight_test.go]
---
# T-2446 The preflight test's markdown tier fixes a file without sed -i, so the fix runs under BSD sed on macOS as under GNU sed

## Work

In `flai/internal/mergequeue/preflight_test.go`, the `markdownTier` fixture's `fix_command` is `sed -i s/FIXABLE/fixed/ "$@"`. BSD sed takes `-i`'s argument as a backup suffix, so on macOS the script is read as the suffix and the first file as the script: nothing is rewritten and `the fix clears it` fails. Replace it with a POSIX loop that writes `sed s/FIXABLE/fixed/` to a temporary file and moves it over the original, which works the same with every sed. Waits for nothing: it changes a file no other task touches.

## Done when

- `go test -race -run TestThePreflightConfirmsALintFindingInAScratchWorktree ./internal/mergequeue/` passes on this host
- the fix command uses no `sed -i`, `perl -i`, or other flag whose meaning differs between GNU and BSD

## Notes
