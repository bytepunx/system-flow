---
id: TH-1176
title: "Plan for S-0777: resolve the kluster tests' temp roots, then check the whole suite under a symlinked TMPDIR"
anchor:
  path: wip/kanban/stories/S-0777-flai-s-kluster-tests-pass-on-macos-where-the-temp-dir-is-a-symlink.md
  item: S-0777
status: resolved
participants: [planner-S-0777, orchestrator]
created: 2026-11-13T20:05:06Z
updated: 2026-11-13T20:06:20Z
---

# TH-1176 Plan for S-0777: resolve the kluster tests' temp roots, then check the whole suite under a symlinked TMPDIR

On wip/kanban/stories/S-0777-flai-s-kluster-tests-pass-on-macos-where-the-temp-dir-is-a-symlink.md.

## Entries

### 2026-11-13T20:05:06Z planner-S-0777
Plan for S-0777, all tasks in the backlog.

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-3323 kluster package fixtures (`setup`, `deploySetup`) resolve their temp folders | `flai/internal/kluster/kluster_test.go`, `deploy_test.go`, `agentconfig_test.go` |
| 1 | T-3324 `klusterProject` in `flai/cmd` resolves its root | `flai/cmd/kluster_test.go` |
| 2 | T-3325 whole suite under a symlinked `TMPDIR`, `scripts/flai-test.sh`, and the narrative's account of the macOS check; after T-3323 and T-3324 | none unless the run finds another test |

T-3323 and T-3324 run together: they share no path.

Figures set:

- Touches: the four declared, kept. The co-changed production files `flai touches suggest` listed are left out, since flai's behaviour does not change.
- Forecast: 19m, delivery 2026-11-14T02:53:00Z, as `flai forecast` gave it.
- Cost of delay: 75 USD a week, from the operator's `time_lost_per_cycle` 30m (TH-1175).
- After: none.

Assumptions:

- Setting `TMPDIR` to a symlink on Linux reproduces macOS's `/var` to `/private/var` link. The story's branch is not pushed before acceptance, so the `test (macos-latest)` job cannot run on it. Criterion 3 is then met by the narrative's account.
- `tempProject` in `flai/cmd` stays as it is. Every `flai/cmd` test shares it, so only `klusterProject` resolves its root.
- My grep of every test under `flai/` that logs `pwd` found no other file to fix. The others resolve their roots already, compare a base name, or rely on the `PWD` flai sets. T-3325 checks this.

Nothing to split, merge, or drop.

### 2026-11-13T20:06:07Z orchestrator
**Approved:** the plan fits S-0777 as drafted.

- **Tasks:** T-3323 and T-3324 fix the two fixtures that cover the eight tests (criterion 1). They share no path, so they run together. T-3325 runs after both. It covers criterion 2 (the sweep under a symlinked `TMPDIR`), criterion 3 (the narrative's account, since the branch is not pushed before acceptance), and criterion 4 (`scripts/flai-test.sh` on Linux).
- **Touches:** the four test files are right for a change that alters no flai behaviour. If T-3325 finds another test, its agent widens the touches with `task_done`.
- **Figures:** the 19m forecast is flai's figure. The 75 USD a week value comes from your own `time_lost_per_cycle` of 30m (TH-1175). No `after` is needed.

I set nothing. The story is planned.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-13T20:06:20Z orchestrator
Resolved: Plan approved under plan_backlog_stories: the tasks, touches, forecast (19m), and value (75 USD a week, from the operator's input on TH-1175) fit S-0777.
