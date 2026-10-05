---
id: S-0254
type: story
nature: improvement
title: flai touches with paths replaces an item's touches, and its help does not say so
status: backlog
owner: alex
created: 2026-10-03T20:33:24Z
updated: 2026-10-05T05:50:01Z
transitions: []
tags: []
topics: [cli]
touches: [flai/cmd/touches.go, flai/cmd/touches_test.go, flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, design/system/flai-cli.md, design/issues/I-0067-flai-touches-with-paths-replaces-an-item-s-touches-and-its-help-does-not-say-so.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
forecast:
  duration: 30m
  delivery: 2026-10-05T23:32:00Z
  basis: "flai forecast's 18m (median 86 s per unit over 7 done medium improvement stories on claude-opus-5-5, size 12, 30th in the pull order) raised to 30m for three tasks worked in sequence, each with its own commit, sync, and test run; delivery moved by the same 12m"
  by: planner-S-0254
  at: 2026-10-05T05:50:01Z
---
# S-0254 flai touches with paths replaces an item's touches, and its help does not say so

## Goal

This story remediates [I-0067](../../../design/issues/I-0067-flai-touches-with-paths-replaces-an-item-s-touches-and-its-help-does-not-say-so.md), "flai touches with paths replaces an item's touches, and its help does not say so". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0067 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0067 is closed with `flai issue close I-0067 --reason` saying what fixed it

## Tasks
- T-0976 flai touches adds with --add and removes with --remove, and its help says paths given alone replace the list
- T-0978 The docs and the design say flai touches replaces unless given --add or --remove
- T-0980 I-0067 records its remediation and is closed

## Notes

### Planning

The proposed fix, from I-0067's one instance: `flai touches` gains `--add` and `--remove`, and its help says that paths given alone replace the list. Replacement stays the default because `flai stream sync`'s hint and any script rely on it. The MCP tool `item_edit` and `flai edit --touches` already say that they replace, so they are left as they are.

Touches. The story declared none.

- `flai/cmd/touches.go`, `flai/cmd/touches_test.go`: layout. The command and its tests live there.
- `flai/cmd/stream_sync.go`, `flai/cmd/stream_sync_test.go`: layout. `printOutside` prints the `flai touches` hint that widens a claim, and it changes to `--add`.
- `design/system/flai-cli.md`: co-change (77% of the seed commits) and design. Its command table has the `flai touches` row.
- `docs/users/flai.md`: design. § Touches and § Outside the touches describe the command and the hint.
- `docs/users/flai-reference.md`, `docs/operators/settings.md`: co-change (29% and 15%). `scripts/flai-reference.sh` regenerates both from the help.
- `design/issues/I-0067-…md`, `design/issues/summary.md`: layout and co-change (15%). Closing the issue writes both.
- Left out: the other co-changed files (`flaiover-dashboard.md`, `docs/operators/index.md`, `mcpserver`, `check`). They change with these seeds in general, not for this fix.

Forecast. `flai forecast` gave 18m, from a median of 86 s per unit over 7 done medium improvement stories, times size 12. It is raised to 30m because the three tasks run in sequence, each with its own commit, sync, and test run. Delivery moves by the same 12m.

Cost of delay. Not set: S-0254 has no inputs and no epic to take a share of, and `flai cod` refuses it. The inputs are asked for on TH-0151, with time_lost_per_cycle 2m recommended, the cost I-0067 records.

Topics. `cli` was added, because the fix is to a command's flags and help.
