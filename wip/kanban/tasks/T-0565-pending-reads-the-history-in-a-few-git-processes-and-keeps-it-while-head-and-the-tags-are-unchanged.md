---
id: T-0565
type: task
nature: feature
title: Pending reads the history in a few git processes and keeps it while HEAD and the tags are unchanged
status: done
parent: S-0157
owner: alex
created: 2026-09-29T19:23:06Z
updated: 2026-09-29T19:31:50Z
transitions:
  - to: ready
    at: 2026-09-29T19:23:30Z
    by: agent-S-0157
  - to: in-progress
    at: 2026-09-29T19:23:31Z
    by: agent-S-0157
  - to: done
    at: 2026-09-29T19:31:50Z
    by: agent-S-0157
stream: S-0157
tags: []
touches: [flai/internal/release]
usage:
  source: log
  seconds: 499
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 54
      output: 28379
      cache_read: 4416020
      cache_write: 67438
      cost: 1.9905
---

# T-0565 Pending reads the history in a few git processes and keeps it while HEAD and the tags are unchanged

## Work

- `release.Pending` reads what it needs from git through a history kept per repository root: `git show-ref --head --tags -d` gives HEAD and the tags, and is the key; on a change, one `git log` of the commits not yet known (with their files), a batched `git diff-tree --stdin` for item commits whose files are not known, and the template's `git log -S` only when its version or `template.yaml` changed.
- Accepted IDs per component range, item commits, and touched files are worked out in memory from that history; the plan itself is recomputed each call from the work items.
- `Compute` for a single item keeps its own path; the pure part is shared.

## Done when

- A behaviour test compares `PendingIDs` with today's per-process implementation for pending items in several components, none, and a release just cut, on the same repository as it changes.
- A test counts the git processes: one while HEAD and the tags are unchanged, at most three after an acceptance.
- `make test` and the linter pass.

## Notes
