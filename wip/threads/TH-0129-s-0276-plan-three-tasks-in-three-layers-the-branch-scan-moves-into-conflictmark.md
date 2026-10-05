---
id: TH-0129
title: "S-0276 plan: three tasks in three layers; the branch scan moves into conflictmark"
anchor:
  path: wip/kanban/stories/S-0276-flai-accept-dry-run-the-dashboard-s-acceptance-preview-does-not-report-a-conflict-marker-as-a-blocker.md
  item: S-0276
status: resolved
participants: [planner-S-0276, alex]
created: 2026-10-05T04:12:19Z
updated: 2026-10-05T04:28:21Z
---

# TH-0129 S-0276 plan: three tasks in three layers; the branch scan moves into conflictmark

On wip/kanban/stories/S-0276-flai-accept-dry-run-the-dashboard-s-acceptance-preview-does-not-report-a-conflict-marker-as-a-blocker.md.

## Entries

### 2026-10-05T04:12:19Z planner-S-0276
S-0276 is planned: three tasks in the backlog, one per layer.

1. **Layer 1: T-0871**, conflictmark reads the conflict markers a story branch adds or changes, and flai accept refuses by it.
   - It moves `(*app).conflictMarkers` out of `flai/cmd/branch.go` into `flai/internal/conflictmark`, and `refuseConflictMarkers` calls it.
   - Touches: `flai/internal/conflictmark`, `flai/cmd/branch.go`.
2. **Layer 2: T-0872**, flai accept --dry-run lists a blocker naming each conflict marker the story branch carries. It runs after T-0871.
   - `preview.Accept` calls the same scan and adds one blocker naming each `path:line`.
   - Tests go in `flai/cmd/accept_conflict_test.go`: one branch with a marker, one clean.
   - Touches: `flai/internal/preview`, `flai/cmd/accept_conflict_test.go`.
3. **Layer 3: T-0873**, flai-cli.md and the users' guide say the acceptance preview reports conflict markers. It runs after T-0872, so the docs quote the blocker as it is built.
   - Touches: `design/system/flai-cli.md`, `docs/users/flai.md`.

The figures:

- **Touches:** the story now has six. I added `flai/internal/conflictmark` and `flai/cmd/accept_conflict_test.go` and kept the four you declared.
- **Forecast:** 30m, delivery about 05:19Z. flai gave 20m; I raised it because the three tasks run in series and include a cross-package move.
- **Cost of delay:** 12.50 USD/week, from your 5m on TH-0127.

The reasons for each are in the story's Notes under Planning.

Assumptions:

- `preview` cannot import `cmd`, so criterion 2's single check means moving the scan into `conflictmark`. `conflictmark` will then import `execx` and `storygit`. I checked that neither imports `check` or `conflictmark`, so there is no import cycle.
- The dashboard's Accept dialog shows `blockers[]` as they come. Nothing in `flaiover/` or its docs changes.
- The preview checks the branch as it is committed, which may be behind main. The acceptance checks it after its sync. A sync that meets a conflict stops rather than committing markers, so the two can only differ on a sync that would fail anyway.
- I have no tasks to propose splitting, merging, or dropping.

### 2026-10-05T04:28:21Z alex
Resolved.
