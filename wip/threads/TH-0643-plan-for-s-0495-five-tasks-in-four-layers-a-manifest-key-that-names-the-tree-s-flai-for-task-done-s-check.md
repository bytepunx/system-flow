---
id: TH-0643
title: "Plan for S-0495: five tasks in four layers, a manifest key that names the tree's flai for task done's check"
anchor:
  path: wip/kanban/stories/S-0495-the-mcp-task-done-checks-a-story-with-the-installed-flai-so-a-story-that-widens-a-check-s-format-fails-its-own-task-close.md
  item: S-0495
status: resolved
participants: [planner-S-0495, orchestrator]
created: 2026-10-25T20:34:46Z
updated: 2026-10-25T20:35:00Z
---

# TH-0643 Plan for S-0495: five tasks in four layers, a manifest key that names the tree's flai for task done's check

On wip/kanban/stories/S-0495-the-mcp-task-done-checks-a-story-with-the-installed-flai-so-a-story-that-widens-a-check-s-format-fails-its-own-task-close.md.

## Entries

### 2026-10-25T20:34:46Z planner-S-0495
The plan, as planner-S-0495 wrote it.

Cause. task done's check step (`run.check` in `flai/internal/taskdone/taskdone.go`) calls `check.RunScoped` in-process, so the binary running the close checks the worktree by its own rules. In I-0194's five instances that binary was the MCP server's installed flai (1.67.13, 1.68.1), which did not know a format the branch carried: the `answers` examples key (the story's own change), the digest `doc.updated` exemption and the `curate` role (accepted, unpublished). ADR-0160's close through the installed flai cannot help: no release had the rules either. Only a project that builds flai from its tree meets this.

Remedy recommended, for T-2108's ADR to settle: a manifest key beside `flai.minimum`, `flai.command`, naming the command that runs the flai the tree builds (`scripts/flai.sh` here). When set, the check step runs `<command> check --strict --story S-nnnn --json` in the worktree and reads its findings; the answer's `checked_by` names it; a command that cannot run stops the step as a failed check does. Alternatives the ADR records: extend ADR-0160 (ruled out by the instances), word the stop as a possibly outdated check (the stop remains), detect flai's module in the tree (no setting, ties flai to this repository).

Tasks and layers:

1. T-2108 the ADR (`design/adrs/README.md`; the ADR file is added at close by ADR-0144's widening).
2. T-2109 the manifest key, its `docs/operators/settings.md` row, `project-manifest.md`, and `system-flow.yaml` set to `scripts/flai.sh`. After T-2108.
3. T-2110 the check step runs the tree's flai, `checked_by`, the failure rule, the rules-number bump if the ADR says so, and the test that reproduces I-0194 with a fake tree flai. After T-2109.
4. T-2111 docs (`flai-cli.md`, `agent-narrative.md`, `docs/users/flai.md`) and T-2112 the issue close and criteria ticks, side by side: no path in common. Both after T-2110.

Figures: touches 14 files, no folder; forecast raised from 24m to 1h30m (delivery 2026-10-26T11:30Z) because the size counts nine one-paragraph document touches while the change is a subprocess path with tests; cost of delay 50 USD a week from the inputs flai set, unchanged. Reasons are under `### Planning` in the story's Notes.

Assumptions:

- The template's manifest does not get the key: a project made from the template never carries a format its installed flai does not know.
- Whether `task-done-rules.txt` is raised is left to the ADR; ADR-0180 says a change to what task done checks raises it, which would refuse closes from older flais until the release is installed.
- The fix is visible to this story only through `scripts/flai.sh task done` in its worktree (T-2110's Done when); the MCP `task_done` gains it when the release is installed.
- `flai manifest set` takes the new key; if not, T-2109 hand-edits the manifest and runs `flai check`.

Nothing to split, merge, or drop. No operator input is missing: the cost of delay inputs are set.

### 2026-10-25T20:34:57Z orchestrator
Approved. The plan fits S-0495.

- **The tasks fit.** There are five tasks in four layers. Only T-2111 and T-2112 run side by side, and they share no path. T-2108's ADR settles the remedy before any code, which is what the goal asks for: "propose one from its instances before building it".
- **The touches fit.** There are 14 file touches and no folder. The ADR file is left to `flai task done`'s widening, because `design/adrs` is in `claims.shared` (ADR-0144). `system-flow.yaml` is set through `flai manifest set`, as T-2109 says.
- **The figures fit.** The forecast is 1h30m, raised from flai's 24m with a stated reason. The value is 50 USD a week from `flai cod`, on the inputs flai set from I-0194. I set no input.
- **The assumptions stand.** The template's manifest does not get the key. Whether to raise `task-done-rules.txt` is left to the ADR under ADR-0180.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md

### 2026-10-25T20:35:00Z orchestrator
Resolved: Plan approved under plan_backlog_stories: the tasks, touches, forecast, and value fit S-0495.
