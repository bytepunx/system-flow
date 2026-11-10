---
id: T-2815
type: task
nature: improvement
title: Close I-0315 with what fixed it
status: done
parent: S-0637
owner: alex
created: 2026-11-09T22:09:44Z
updated: 2026-11-10T08:13:19Z
transitions:
  - to: ready
    at: 2026-11-10T08:13:06Z
    by: agent-S-0637
  - to: in-progress
    at: 2026-11-10T08:13:07Z
    by: agent-S-0637
  - to: done
    at: 2026-11-10T08:13:19Z
    by: agent-S-0637
stream: S-0637
tags: [issues]
touches: [design/issues/I-0315-stories-that-each-add-a-section-at-one-spot-in-a-shared-design-document.md, design/issues/summary.md]
after: [T-2813, T-2814]
usage:
  source: log
  seconds: 12
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 2
      output: 313
      cache_read: 320721
      cache_write: 2563
      cost: 0.0895
---
# T-2815 Close I-0315 with what fixed it

## Work

Run `flai issue close I-0315 --reason "<reason>"`. The reason names the new ADR, the matcher in `flai/internal/issues/generated.go` that now merges by sections a stopped markdown file under the design or docs folder, and the two tests T-2814 added. Let flai write `design/issues/summary.md` again.

Waits for T-2813 and T-2814: the issue closes once the fix, its tests, and its documents are in.

## Done when

- I-0315 is closed with that reason, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` passes.

## Notes
