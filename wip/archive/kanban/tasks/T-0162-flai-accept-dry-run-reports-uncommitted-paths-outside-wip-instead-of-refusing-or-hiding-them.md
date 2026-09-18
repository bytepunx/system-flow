---
id: T-0162
type: task
nature: improvement
title: flai accept --dry-run reports uncommitted paths outside wip instead of refusing or hiding them
status: done
parent: S-0051
owner: alex
created: 2026-09-18T21:00:23Z
updated: 2026-09-18T21:01:32Z
transitions:
  - to: ready
    at: 2026-09-18T21:00:24Z
    by: alex
  - to: in-progress
    at: 2026-09-18T21:00:25Z
    by: alex
  - to: done
    at: 2026-09-18T21:01:32Z
    by: alex
stream: S-0051
tags: []
touches: [flai/cmd]
---

# T-0162 flai accept --dry-run reports uncommitted paths outside wip instead of refusing or hiding them

## Work
In `flai/cmd/accept.go`, add `Uncommitted []string` (`json:"uncommitted,omitempty"`) to the accept result, filled from `dirtyOutsideWip` whenever git is in use. A dry run never refuses for uncommitted paths, with or without `--yes`: it reports them, and the text output prints them with what the real run will do. The real run is unchanged: refused without `--yes`, included with it. Tests with real git: dry run with a stray file reports it and exits 0 without `--yes`; the real run refuses and leaves the story in review; with `--yes` it accepts.

## Done when
- The new tests fail without the change and pass with it
- `go test -race ./cmd/...` passes

## Notes
