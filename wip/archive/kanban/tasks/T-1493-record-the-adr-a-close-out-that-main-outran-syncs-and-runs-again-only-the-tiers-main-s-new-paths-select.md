---
id: T-1493
type: task
nature: improvement
title: "Record the ADR: a close-out that main outran syncs and runs again only the tiers main's new paths select"
status: done
parent: S-0374
owner: alex
created: 2026-10-12T17:33:48Z
updated: 2026-10-14T21:45:14Z
transitions:
  - to: ready
    at: 2026-10-14T21:44:24Z
    by: agent-S-0374
  - to: in-progress
    at: 2026-10-14T21:44:24Z
    by: agent-S-0374
  - to: done
    at: 2026-10-14T21:45:14Z
    by: agent-S-0374
stream: S-0374
tags: [cli]
touches: [design/adrs/README.md]
usage:
  source: log
  seconds: 50
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 12
      output: 4970
      cache_read: 792707
      cache_write: 18076
      cost: 0.3763
---
# T-1493 Record the ADR: a close-out that main outran syncs and runs again only the tiers main's new paths select

## Work

Write a new ADR with `flai adr new` that refines ADR-0135 and ADR-0110, from I-0136's three instances (S-0232 twice, S-0337 twice more): each time every tier passed, main gained non-wip commits during the run, the final `flai verify --sync-only` failed, and the branch then rebased cleanly.

Decide, as the plan thread proposes unless its answer says otherwise:

1. **Reuse across a clean rebase.** `flai verify` reuses a tier the story's last report passed when the head and base moved but the branch's own change is the same (its diff from the merge base has the same `git patch-id --stable`) and the fingerprint is unchanged, unless the paths main changed between the report's base commit and the new one select that tier (`Select`, `{importers}` included). Those tiers run again. The `rebase`, `sync`, `narrative`, and `check` steps always run.
2. **The close-out syncs and verifies again.** When `flai verify --sync-only` fails at its `sync` step, `scripts/close-out.sh` runs `flai stream sync` and, when it ends clean, `flai verify --record-issues` and the sync check again, at most three rounds, then stops as now with the hint. A sync that stops on a conflict stops the close-out as now.
3. **Why not pass over main's commits.** Main's non-wip change can break the branch even when the two touch different files, such as a Go package the branch imports; reusing only what main's paths cannot reach keeps the run's verdict true of the head that goes to review.

Record the alternatives: pass over commits whose paths are disjoint from the branch's (ADR-0135's rejected widening), and hold acceptances while a close-out runs. Add it to `design/adrs/README.md`. `flai task done` widens the touches with the new ADR file, whose number is not known yet.

## Done when

- The ADR is accepted, refines ADR-0135 and ADR-0110, and names the three decisions and the alternatives.
- `design/adrs/README.md` lists it.
- `flai check --strict` is clean.

## Notes

Layer 1: waits for no task. Every other task builds what this ADR decides.
