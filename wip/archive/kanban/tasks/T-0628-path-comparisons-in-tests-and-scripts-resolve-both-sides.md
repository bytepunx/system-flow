---
id: T-0628
type: task
nature: feature
title: Path comparisons in tests and scripts resolve both sides
status: done
parent: S-0180
owner: arobson
created: 2026-10-01T08:22:54Z
updated: 2026-10-01T08:23:04Z
transitions:
  - to: ready
    at: 2026-10-01T08:23:04Z
    by: agent-S-0180
  - to: in-progress
    at: 2026-10-01T08:23:04Z
    by: agent-S-0180
  - to: done
    at: 2026-10-01T08:23:04Z
    by: agent-S-0180
stream: S-0180
tags: []
usage:
  source: log
  seconds: 0
  models: []
---

# T-0628 Path comparisons in tests and scripts resolve both sides

## Work

`checks_test.go` reads the check's directory with `pwd -P`, so both sides of the comparison are resolved (`EvalSymlinks` on the other); `install-test.sh` resolves the `mktemp -d` path it greps for with `cd … && pwd -P`, as flai resolves its own path.

## Done when

`TestRunChecksSubstitutesStoryAndRootAndRunsInTheWorktree` and `scripts/install-test.sh` pass on macOS.

## Notes
