---
id: T-2800
type: task
nature: improvement
title: Give the nolint directive in flai/internal/watch/watch.go its reason
status: done
parent: S-0631
owner: alex
created: 2026-11-09T22:05:07Z
updated: 2026-11-10T06:38:58Z
transitions:
  - to: ready
    at: 2026-11-10T06:16:51Z
    by: agent-S-0631
  - to: in-progress
    at: 2026-11-10T06:16:52Z
    by: agent-S-0631
  - to: done
    at: 2026-11-10T06:38:58Z
    by: agent-S-0631
stream: S-0631
tags: [cli]
touches: [flai/internal/watch/watch.go]
usage:
  source: log
  seconds: 1326
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 4
      output: 785
      cache_read: 559433
      cache_write: 16749
      cost: 0.2229
---
# T-2800 Give the nolint directive in flai/internal/watch/watch.go its reason

## Work

I-0298 counts 26 `//nolint` directives in flai and one without a reason, which it places at `flai/internal/watch/watch.go:54`. The directive is now at line 94: `return nil //nolint:nilerr`, inside the `filepath.WalkDir` callback's `err != nil` branch. The reason is in the comment on line 93: "A folder that vanished or cannot be read is skipped, not fatal."

Move the reason onto the directive, as every other nolint in flai has it, for example `//nolint:nilerr // a folder that vanished or cannot be read is skipped, not fatal`. Then drop or shorten the comment above so the reason is said once.

Check with `grep -rn 'nolint' flai --include='*.go'` that no other directive lacks a reason, and give any you find one too.

This task waits for nothing and shares no path with the splits, so it runs with them.

## Done when

- Every `//nolint` directive in `flai/` carries a reason after `//`.
- `flai test flai/internal/watch/watch.go` passes, golangci-lint included.

## Notes
