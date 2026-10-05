---
id: TH-0153
title: Plan for S-0254
anchor:
  path: wip/kanban/stories/S-0254-flai-touches-with-paths-replaces-an-item-s-touches-and-its-help-does-not-say-so.md
  item: S-0254
status: open
participants: [planner-S-0254]
created: 2026-10-05T05:50:09Z
updated: 2026-10-05T05:50:52Z
---

# TH-0153 Plan for S-0254

On wip/kanban/stories/S-0254-flai-touches-with-paths-replaces-an-item-s-touches-and-its-help-does-not-say-so.md.

## Entries

### 2026-10-05T05:50:09Z planner-S-0254
The proposed fix for I-0067: `flai touches` gains `--add` and `--remove`, and its help says that paths given alone replace the list.

Tasks, in three layers, one after another:

1. T-0976 (layer 1, waits for nothing): add `--add` and `--remove`, with `--add`, `--remove`, and `--clear` refusing each other. `Short`, `Long`, and the examples say that paths given alone replace the list. `flai stream sync` hints with `flai touches <id> --add <outside paths>`. Tests reproduce I-0067. Touches `flai/cmd/touches.go`, `flai/cmd/stream_sync.go`, and their tests.
2. T-0978 (layer 2, after T-0976): update `docs/users/flai.md` and `design/system/flai-cli.md`, and regenerate `flai-reference.md` and `settings.md` with `scripts/flai-reference.sh`. It waits because the reference is generated from T-0976's help.
3. T-0980 (layer 3, after T-0978): write I-0067's Remediation and close it with `flai issue close` from the worktree. It waits so that the issue closes only once the fix is built and documented.

Assumptions:
- Paths given alone keep replacing the list rather than switching to adding. Switching would silently change what existing scripts and sync's hint do; the instance shows the trap was the unclear help, not the replacing itself. If you would rather paths alone add, with `--set` to replace, say so and I'll redraft T-0976.
- `flai edit --touches` and the MCP tool `item_edit` already say "replacing", so they are out of scope.
- Touches: 10 paths. Forecast: 30m, up from 18m for three tasks run in sequence. Cost of delay: waiting on your inputs in TH-0151.

### 2026-10-05T05:50:52Z planner-S-0254
An overlap to know about: S-0251, planned meanwhile, also claims `flai/cmd/stream_sync.go`, `flai/cmd/stream_sync_test.go`, `docs/users/flai.md`, `flai-reference.md`, `flai-cli.md`, and `design/issues/summary.md`. Whichever of the two starts second will be held until the first is done, so they cannot conflict. If you would rather they could run together, I can drop the sync hint change, the `stream_sync` paths, from T-0976 and S-0254. The sync hint already repeats every touch, so it does not cause I-0067; changing it to `--add` is only a tidier hint. The docs and `summary.md` would still overlap.
