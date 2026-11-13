---
id: TH-1175
title: "S-0777 needs cost of delay inputs: it has none and no epic to take a share of"
anchor:
  path: wip/kanban/stories/S-0777-flai-s-kluster-tests-pass-on-macos-where-the-temp-dir-is-a-symlink.md
  item: S-0777
status: resolved
participants: [planner-S-0777, orchestrator]
created: 2026-11-13T20:01:21Z
updated: 2026-11-13T20:03:24Z
---

# TH-1175 S-0777 needs cost of delay inputs: it has none and no epic to take a share of

On wip/kanban/stories/S-0777-flai-s-kluster-tests-pass-on-macos-where-the-temp-dir-is-a-symlink.md.

## Entries

### 2026-11-13T20:01:21Z planner-S-0777
`flai cod S-0777` cannot work out a value: the story has no cost of delay inputs and no epic.

Recommended: `time_lost_per_cycle` 10m, and no revenue or penalty. The red `test (macos-latest)` job costs a reader a few minutes on each push to `main` to see that it is the known failure. It also hides any real macOS defect. It brings no revenue and carries no penalty.

Alternatives:

- `penalty_per_week` 0 and `time_lost_per_cycle` 0: value 0, so the story ranks last.
- `time_lost_per_cycle` 30m, if the red job is checked by hand on every release.

Give the inputs with `flai edit S-0777 --time-lost-per-cycle 10m` or answer here. I will work out the value once they are set. The touches, forecast, and tasks are planned meanwhile.

### 2026-11-13T20:03:13Z orchestrator
Settled by the operator: alex set `time_lost_per_cycle: 30m` on S-0777 at 2026-11-13T20:02:24Z, the planner's third alternative. I set nothing over it: the inputs are the operator's.

The planner still runs for S-0777 (pid 245780) and works out the value and the forecast duration from that input. `flai promote --candidates` lists S-0777 as lacking only those two.

The two tasks, T-3323 and T-3324, in one layer, and the touches fit the story's criteria. I resolve this thread.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-13T20:03:24Z orchestrator
Resolved: The operator set time_lost_per_cycle 30m on S-0777; the planner works out the value on its run.
