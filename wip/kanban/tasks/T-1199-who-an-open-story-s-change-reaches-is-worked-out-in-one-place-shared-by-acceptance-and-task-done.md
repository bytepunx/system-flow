---
id: T-1199
type: task
nature: improvement
title: Who an open story's change reaches is worked out in one place, shared by acceptance and task done
status: done
parent: S-0333
owner: alex
created: 2026-10-07T20:15:31Z
updated: 2026-10-07T22:13:16Z
transitions:
  - to: ready
    at: 2026-10-07T22:10:54Z
    by: agent-S-0333
  - to: in-progress
    at: 2026-10-07T22:10:54Z
    by: agent-S-0333
  - to: done
    at: 2026-10-07T22:13:16Z
    by: agent-S-0333
stream: S-0333
tags: [flai]
touches: [flai/internal/itemedit/covers.go, flai/internal/itemedit/covers_test.go, flai/cmd/accept_overlap.go, flai/cmd/accept_overlap_test.go]
usage:
  source: log
  seconds: 142
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 18
      output: 157
      cache_read: 1209155
      cache_write: 10689
      cost: 0.5468
---
# T-1199 Who an open story's change reaches is worked out in one place, shared by acceptance and task done

## Work

Move `overlapsOf` and `covered` out of `flai/cmd/accept_overlap.go` into `flai/internal/itemedit`, so `flai/internal/taskdone` can call them. It waits for nothing: it changes no behaviour.

- One function takes the items, the projects, the story that changed, and the paths it changed, and returns each other open story whose claim covers one of them, shared paths included, with the paths; a story with an empty claim gets every path.
- `tellOverlaps` calls it, and its behaviour at acceptance stays the same.

## Done when

- The tests of `accept_overlap_test.go` pass unchanged, and `covers_test.go` covers the function on its own.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
