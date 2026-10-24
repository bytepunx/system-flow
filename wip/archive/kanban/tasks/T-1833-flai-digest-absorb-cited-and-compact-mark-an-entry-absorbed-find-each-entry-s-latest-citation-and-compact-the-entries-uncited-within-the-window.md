---
id: T-1833
type: task
nature: feature
title: flai digest absorb, cited, and compact mark an entry absorbed, find each entry's latest citation, and compact the entries uncited within the window
status: done
parent: S-0432
owner: alex
created: 2026-10-16T06:06:53Z
updated: 2026-10-24T14:07:24Z
transitions:
  - to: ready
    at: 2026-10-24T13:42:16Z
    by: agent-S-0432
  - to: in-progress
    at: 2026-10-24T13:42:17Z
    by: agent-S-0432
  - to: done
    at: 2026-10-24T14:06:44Z
    by: agent-S-0432
stream: S-0432
tags: [cli]
touches: [flai/internal/digest/decay.go, flai/internal/digest/decay_test.go, flai/internal/digest/digest.go, flai/internal/digest/digest_test.go, flai/cmd/digest.go, flai/cmd/digest_test.go, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-1832]
usage:
  source: log
  seconds: 1467
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 85
      output: 33891
      cache_read: 5592543
      cache_write: 231844
      cost: 3.0981
---
# T-1833 flai digest absorb, cited, and compact mark an entry absorbed, find each entry's latest citation, and compact the entries uncited within the window

## Work

- `flai/internal/digest/decay.go`: `Absorb(id, path, heading)` sets `status: absorbed` and `absorbed_by`, refusing a path or heading that does not exist; `Cited(root)` scans narratives, threads, and work items, archived ones included, for each entry's ID and gives its latest citation and date; `Compact(root, window, dryRun)` turns each active entry with no citation within the window, counted from its own date when it has none, into `compacted`, its body reduced to its title, and returns what it changed.
- `flai/cmd/digest.go`: `flai digest absorb <D-nnnn> --by <path> [--heading <h>] [--json]`, `flai digest cited [--json]`, and `flai digest compact [--dry-run] [--json]`, the window from `digest.compact_after`, or the default while S-0432's settings task is not in.
- Regenerate `docs/users/flai-reference.md`.

Waits for T-1831. Runs beside the settings task.

## Done when

`decay_test.go` covers an absorption and a refused one, a citation in an archived narrative, a cited entry kept, an uncited one compacted, a dry run that changes nothing, and a second run that compacts nothing more; `flai/cmd/digest_test.go` covers the three commands; `flai test flai/internal/digest flai/cmd` passes.

## Notes

S-0430's pack already leaves out absorbed entries and loads compacted ones as their titles.
