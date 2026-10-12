---
id: TH-0424
title: "S-0374 plan: the close-out syncs and verifies again, reusing the tiers main's new paths do not select (new ADR), five tasks in four layers"
anchor:
  path: wip/kanban/stories/S-0374-an-acceptance-on-main-during-a-close-out-s-tiers-fails-its-final-sync-check-so-the-whole-close-out-runs-again.md
  item: S-0374
status: open
participants: [planner-S-0374]
created: 2026-10-12T17:35:18Z
updated: 2026-10-12T17:35:18Z
---

# TH-0424 S-0374 plan: the close-out syncs and verifies again, reusing the tiers main's new paths do not select (new ADR), five tasks in four layers

On wip/kanban/stories/S-0374-an-acceptance-on-main-during-a-close-out-s-tiers-fails-its-final-sync-check-so-the-whole-close-out-runs-again.md.

## Entries

### 2026-10-12T17:35:18Z planner-S-0374
**Recommendation:** approve the remedy below, which a new ADR refining ADR-0135 and ADR-0110 records.

## The remedy

- **Today:** when an acceptance lands on main during the tiers, the close-out stops at its last check. The agent syncs and runs every tier again, about 6 minutes each time in I-0136's three instances.
- **After:** `scripts/close-out.sh` runs `flai stream sync` itself, then `flai verify` and the sync check again, at most three rounds.
- `flai verify` reuses a tier the last report passed across that clean rebase when the branch's own change is the same (`git patch-id --stable`). The exception is a tier that the paths main changed select, `{importers}` included: it runs again.
- **Why not pass over main's commits outright,** as I-0136's first instance suggests: a main change can break the branch without touching its files, such as a Go package it imports. ADR-0135 rejected that widening for the same reason.

## Tasks

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-1493 the ADR | `design/adrs/README.md` (task done adds the ADR file) |
| 2 | T-1494 verify's reuse across a clean rebase, with a real-repository test reproducing I-0136 | `flai/internal/verify/story.go`, `paths.go`, `story_test.go`, `paths_test.go` |
| 2 | T-1495 the close-out loop, both copies, and the changelog | `scripts/close-out.sh`, `template/root/scripts/close-out.sh`, `template/CHANGELOG.md` |
| 3 | T-1496 docs and the convention | `docs/users/flai.md`, `design/system/workflow.md`, `design/system/flai-cli.md`, `design/conventions/work-management.md` and its template copy |
| 4 | T-1497 close I-0136 | the issue file, `design/issues/summary.md` |

T-1494 and T-1495 share no path, so they run together.

## Assumptions to confirm

1. **Three rounds** bound the loop. After that the close-out stops as it does now.
2. **The baseline rule changes.** T-1496 rewrites one sentence of `work-management.md` § Rules in both copies: "Any other commit the branch lacks, such as an acceptance, stops the close-out: run `flai stream sync` and run it again". The template releases that rule to every project.
3. **Reuse keeps resume's prefix rule.** Once one tier has to run again, every tier after it runs too, as resume works today.

## Figures

- Every touch names a file. No folder touch is kept.
- Forecast is 1h, raised from flai's 23m, for five tasks with an ADR, a Go change, two scripts, and docs.
- Cost of delay is 45 USD a week, from flai's 18m input. That figure is likely low: the three instances came within 1.5 hours of each other.

S-0374 stays a draft.
