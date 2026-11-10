---
id: T-2933
type: task
nature: improvement
title: The user guide says what ready-empty counts and that a move to ready warns of an after hold, and I-0285 is closed
status: done
parent: S-0673
owner: alex
created: 2026-11-10T08:49:31Z
updated: 2026-11-10T10:31:07Z
transitions:
  - to: ready
    at: 2026-11-10T10:30:51Z
    by: agent-S-0673
  - to: in-progress
    at: 2026-11-10T10:30:51Z
    by: agent-S-0673
  - to: done
    at: 2026-11-10T10:31:07Z
    by: agent-S-0673
stream: S-0673
tags: [docs, issues]
touches: [docs/users/flai.md, design/issues/I-0285-after-chains-on-planned-stories-hold-ready-stories-for-days-while-in-progress.md, design/issues/summary.md]
after: [T-2930, T-2931]
usage:
  source: log
  seconds: 16
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 2
      output: 431
      cache_read: 321151
      cache_write: 15790
      cost: 0.1631
---
# T-2933 The user guide says what ready-empty counts and that a move to ready warns of an after hold, and I-0285 is closed

## Work

In `docs/users/flai.md`, say in the `inbox` and `wait_for_events` rows and the trigger events that `ready-empty` counts the ready stories not held on `after`, and in `flai move`'s section that a move to ready warns of an `after` hold. Use the words T-2930 and T-2931 wrote.

Then close I-0285 with `flai issue close I-0285 --reason`, saying what fixed it: the planner's rule from ADR-0230 for the chains, the new ADR's `ready-empty` count, and the move's warning, naming the tests. `flai` regenerates `design/issues/summary.md`.

It waits for T-2930 and T-2931: it documents what they built, and the reason names their tests.

## Done when

- `docs/users/flai.md` says both, and the markdown lint passes on it.
- I-0285 is closed with its reason, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` is clean.

## Notes
