---
id: T-1495
type: task
nature: improvement
title: The close-out syncs and verifies again, at most three rounds, when main moved during its tiers
status: in-progress
parent: S-0374
owner: alex
created: 2026-10-12T17:34:15Z
updated: 2026-10-14T21:45:26Z
transitions:
  - to: ready
    at: 2026-10-14T21:45:26Z
    by: agent-S-0374
  - to: in-progress
    at: 2026-10-14T21:45:26Z
    by: agent-S-0374
stream: S-0374
tags: [cli, template]
touches: [scripts/close-out.sh, template/root/scripts/close-out.sh, template/CHANGELOG.md]
after: [T-1493]
---
# T-1495 The close-out syncs and verifies again, at most three rounds, when main moved during its tiers

## Work

Waits for T-1493, whose ADR decides the loop and its bound.

- In `scripts/close-out.sh`, when `flai verify <story> --sync-only` fails at its `sync` step, run `flai stream sync <story>`. When the sync ends clean, run `flai verify <story> --record-issues` again, the commit step, and the sync check again. Stop after three rounds with the hint as now. A sync that stops on a conflict, or a verify that fails, stops the close-out at that step as now, with its hint.
- Print each round, such as `close-out: main moved during the run; flai stream sync S-nnnn (round 2 of 3)`, so the operator sees why the close-out took longer. The trap's last line names the round it ended on.
- Make the same change in `template/root/scripts/close-out.sh`, which calls `flai` directly.
- Add a `template/CHANGELOG.md` entry naming the flai version the loop needs to reuse tiers across the rebase (T-1494's release); with an older flai the loop still works and runs every tier again.

## Done when

- A close-out whose sync check fails because main gained commits syncs, verifies, and checks again on its own, and stops after three rounds or at the first step that fails.
- Both copies of the script behave the same, `sh -n` and the markdown lint pass, and the changelog names the version.

## Notes

Layer 2: runs beside T-1494; no path in common.
