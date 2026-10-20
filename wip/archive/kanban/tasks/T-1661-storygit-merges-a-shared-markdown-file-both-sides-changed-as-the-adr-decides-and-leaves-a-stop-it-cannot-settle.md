---
id: T-1661
type: task
nature: feature
title: storygit merges a shared markdown file both sides changed as the ADR decides, and leaves a stop it cannot settle
status: done
parent: S-0401
owner: alex
created: 2026-10-16T05:02:14Z
updated: 2026-10-19T22:03:51Z
transitions:
  - to: ready
    at: 2026-10-19T21:54:31Z
    by: agent-S-0401
  - to: in-progress
    at: 2026-10-19T21:54:31Z
    by: agent-S-0401
  - to: done
    at: 2026-10-19T22:03:51Z
    by: agent-S-0401
stream: S-0401
tags: [cli]
touches: [flai/internal/storygit/mdmerge.go, flai/internal/storygit/mdmerge_test.go]
after: [T-1645]
usage:
  source: log
  seconds: 560
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 32
      output: 16780
      cache_read: 2098169
      cache_write: 64361
      cost: 1.1422
---
# T-1661 storygit merges a shared markdown file both sides changed as the ADR decides, and leaves a stop it cannot settle

## Work

- `mdmerge.go`: merge a markdown file from a stop's three versions, as T-1645's ADR decides, for example section by section by heading. Return no result for a case it cannot settle, such as both sides changing the same section differently.
- Merge only paths inside `claims.shared`.
- Test two sides adding sections in one place, both changing one section, and one side removing a section the other changed.
- Waits for T-1645, whose decision it implements.

## Done when

The merge settles the cases the ADR says it settles, leaves the rest, and the tests pass.

## Notes
