---
id: T-2673
type: task
nature: improvement
title: digest.List and Get read digest entries through the kept store
status: done
parent: S-0602
owner: alex
created: 2026-11-07T23:50:13Z
updated: 2026-11-08T03:22:59Z
transitions:
  - to: ready
    at: 2026-11-08T03:08:11Z
    by: agent-S-0602
  - to: in-progress
    at: 2026-11-08T03:08:11Z
    by: agent-S-0602
  - to: done
    at: 2026-11-08T03:22:59Z
    by: agent-S-0602
stream: S-0602
tags: [cli]
touches: [flai/internal/digest/digest.go, flai/internal/digest/digest_test.go]
after: [T-2660]
usage:
  source: log
  seconds: 860
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 17
      output: 137
      cache_read: 2465858
      cache_write: 31291
      cost: 0.9037
---
# T-2673 digest.List and Get read digest entries through the kept store

## Work

Move `List` and `Get` in `flai/internal/digest/digest.go` onto a package-level `workitem.KeptFiles[Entry]`, with `read` as its parse and a copy of an entry, its slices included, as its clone. Keep the folder under `repo.Root` and the order `List` gives. The writes in `write.go` and `decay.go` (`New`, `Edit`, `Absorb`, `Compact`) save through `atomicfile`, so the next read sees them by stat.

It waits for T-2660, whose store it uses.

## Done when

- `digest_test.go` shows that an entry edited, added, and removed is each seen on the next `List`, and that a second `List` with nothing changed parses no file.
- It shows that two `List` calls at once return equal entries.
- `flai test flai/internal/digest/digest.go flai/internal/digest/digest_test.go` passes.

## Notes
