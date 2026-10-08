---
id: TH-0369
title: "S-0341 plan: five tasks in three layers; resume keyed on head, base, and a fingerprint of the tiers"
anchor:
  path: wip/kanban/stories/S-0341-flai-verify-resumes-at-the-tier-that-failed-when-the-branch-head-and-its-base-are-unchanged-so-a-retry-re-runs-the-failure-and-not-the-tiers-that-passed.md
  item: S-0341
status: open
participants: [planner-S-0341]
created: 2026-10-08T08:06:01Z
updated: 2026-10-08T08:06:01Z
---

# TH-0369 S-0341 plan: five tasks in three layers; resume keyed on head, base, and a fingerprint of the tiers

On wip/kanban/stories/S-0341-flai-verify-resumes-at-the-tier-that-failed-when-the-branch-head-and-its-base-are-unchanged-so-a-retry-re-runs-the-failure-and-not-the-tiers-that-passed.md.

## Entries

### 2026-10-08T08:06:01Z planner-S-0341
The plan for S-0341 is written. Five tasks in three layers. Touches widened from 9 files to 15, all files, no folders. Forecast 1h30m. Cost of delay 12.50 USD a week.

## Tasks and layers

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1343 flai verify resumes a story's run from its last record | `flai/internal/verify/{verify,story,story_test,run,run_test}.go` | none |
| 2 | T-1344 `--fresh`, reused tier lines, and a last line that says it resumed | `flai/cmd/verify.go`, `verify_test.go`, `docs/users/flai-reference.md` | T-1343 |
| 2 | T-1346 The review page shows a reused tier with its origin time | `flaiover/src/lib/review.ts`, `review.test.ts`, `components/Review.svelte` | T-1343 |
| 3 | T-1347 The close-out's last line says when tiers were reused | `scripts/close-out.sh` | T-1344 |
| 3 | T-1351 The docs say when a run resumes and how to force a full one | `docs/users/flai.md`, `design/system/flai-cli.md`, `docs/users/flaiover.md` | T-1344, T-1346 |

No two tasks in the same layer share a path.

## Assumptions

1. **What forces a full run.** A run resumes only when the record's `commit` and `base` match and a new fingerprint matches. The fingerprint is a hash of the manifest's `tests` plus each selected tier's name, command, and directory. A record without one, written by an older flai, runs in full.
2. **A reused tier's time.** It carries `reused_from`, the time of the run that actually ran it. That time carries through a chain of resumed runs.
3. **A record that passed every tier.** The next run at the same head reuses every tier. It still runs rebase, sync, narrative, and check. A record that stopped before the tiers runs every tier.
4. **Only `flai verify` resumes.** `flai test` never does. The MCP `verify` tool and the host's `verify.run` call the same story run, so they resume too. They get no `fresh` parameter in this story.
5. **The dashboard is in scope (T-1346).** Criterion 2 says the review page must show a whole result. Today a `reused` step would render grey with no duration, like a step not reached. This adds three flaiover files to the story's touches. Nothing open touches them.
6. **Forecast raised from 19m to 1h30m.** The last two `flai verify` stories took 81 minutes (S-0326) and 100 minutes (S-0270).

## Proposals

None to split, merge, or drop. If you would rather leave the dashboard out, cancel T-1346 and drop its three touches. The criterion is then met only by `--last` and the raw state name on the page.
