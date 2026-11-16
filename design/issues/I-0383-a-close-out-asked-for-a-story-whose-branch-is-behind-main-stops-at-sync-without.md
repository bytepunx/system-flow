---
id: I-0383
title: A close-out asked for a story whose branch is behind main stops at sync without syncing, and the failed report makes a story the merge queue could batch not eligible
class: defect
status: closed
count: 1
cost: 2h
first_reported: 2026-11-16T00:41:30Z
last_reported: 2026-11-16T00:41:30Z
updated: 2026-11-16T02:18:53Z
---

# I-0383 A close-out asked for a story whose branch is behind main stops at sync without syncing, and the failed report makes a story the merge queue could batch not eligible

## Description
A close-out asked for a story whose branch is behind main stops at sync without syncing, and the failed report makes a story the merge queue could batch not eligible

## Instances

### 2026-11-16T00:41:30Z
ADR-0291 §4 and the verify_behind_main blocker's words say that asking for the close-out again, with the MCP tool verify or flai verify <S-nnnn> --job, syncs the story's branch onto main and runs again only the tiers main's new paths send back. scripts/close-out.sh syncs only when main moved during its run, at its last sync check (ADR-0156); when the first flai verify stops at sync it ends. Jobs S-0739-20261115T210803Z, S-0742-20261115T210805Z, and S-0747-20261115T210806Z stopped so, 37 to 51 commits behind main. Each story's last report was then a failed one, so the queue, which had listed all three as eligible on their earlier passing verify, refused them as verify_failed.

## Impact
On 2026-11-15 the orchestrator's close-outs of S-0739, S-0742, and S-0747 each ended in under a second at sync, and flai queue then refused all three as verify_failed (TH-1226, TH-1227, TH-1228) until each was synced by hand and closed out again

- time_lost_per_cycle: 2h

## Remediation

The close-out a story is asked for when its branch is behind main syncs first: `scripts/close-out.sh` runs `flai stream sync` when its first `flai verify` stops at sync, as it already does when main moves during the run (ADR-0156), within the same three rounds, and stops with the conflicts when the sync does. Until then a close-out that stops at sync before any tier ran should not replace a passing report as the story's last, or the queue's `verify_failed` should tell that case apart. Bring ADR-0291 §4's words and the `verify_behind_main` blocker's in step with what the close-out does.

Story S-0798 remediates this issue, created from it at 2026-11-16T00:41:49Z.
Closed 2026-11-16T02:18:53Z: S-0798: the close-out syncs when a round's first flai verify stops at sync, committing what is outstanding first, within ADR-0156's three rounds (ADR-0310), in scripts/close-out.sh, the template's, and a container story's close-out; flai/internal/closeout/script_test.go and the verifyjob tests reproduce the stop and pass.
