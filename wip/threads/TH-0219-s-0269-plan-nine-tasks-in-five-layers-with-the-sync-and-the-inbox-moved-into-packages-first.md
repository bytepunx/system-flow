---
id: TH-0219
title: "S-0269 plan: nine tasks in five layers, with the sync and the inbox moved into packages first"
anchor:
  path: wip/kanban/stories/S-0269-one-command-closes-a-task-flai-task-done-commits-syncs-moves-logs-widens-touches-checks-and-answers-the-inbox.md
  item: S-0269
status: open
participants: [planner-S-0269]
created: 2026-10-06T22:55:40Z
updated: 2026-10-06T22:55:40Z
---

# TH-0219 S-0269 plan: nine tasks in five layers, with the sync and the inbox moved into packages first

On wip/kanban/stories/S-0269-one-command-closes-a-task-flai-task-done-commits-syncs-moves-logs-widens-touches-checks-and-answers-the-inbox.md.

## Entries

### 2026-10-06T22:55:40Z planner-S-0269
I've planned S-0269. There are nine tasks, all drafted in the backlog. Nothing here needs an answer unless you disagree with an assumption below.

The tasks run in five layers:

1. T-1080, T-1085 and T-1092 run together:
   - T-1080: an ADR refining ADR-0069, plus the `flai-cli.md` entry. It sets the steps, the stop rule, the answer's fields, and the names `task_done` and `task.done`.
   - T-1085: `flai stream sync`'s logic moves out of cmd's `app` methods into a storygit function that returns a structured result.
   - T-1092: the agent inbox and its cursor move out of `mcpserver` into a new `flai/internal/inbox` package.
2. T-1098: the `taskdone` package runs the seven steps in order. It has tests for the happy path, a refused sync, and a failed check.
3. T-1102 and T-1107 run together:
   - T-1102: `flai task done` on the CLI, as text and `--json`.
   - T-1107: the MCP tool `task_done`, and the server's instructions.
4. T-1113 and T-1116 run together:
   - T-1113: `task.done` on the host channel, a write spec that runs the CLI.
   - T-1116: `git.md`, `work-management.md`, their template copies, the CHANGELOG, and the harness prompt.
5. T-1119: `docs/users/flai.md`, and the regenerated reference.

The plan rests on these assumptions:

- **Why layer 1 moves code:** the MCP server cannot import cmd. So the sync and the inbox move into packages first, as `criteria_tick` already calls `itemedit`, rather than MCP shelling out to the CLI.
- **Inbox cursor:** the CLI's inbox advances the same per-agent cursor as MCP, so a change is reported once. No `flai inbox` command is added.
- **Order and stopping:** the steps run in the goal's order (commit, sync, move, log, touches, check, inbox). Nothing to commit is not a failure. A failed check still leaves the task done and logged, because the move comes before the check.
- **Touches:** both the task's and the story's touches are widened, with only the paths the task's commit changed that they do not cover.
- **Left to the agent:** running the task's tests (S-0273's `flai test`) and ticking criteria (`flai criteria tick`).
- **Dashboard:** the host channel only exposes `task.done`; no dashboard page calls it.

The story's figures:

- **Touches:** all 23 declared touches are kept, and ten files the layout shows are added. The only folder touch kept is `design/adrs`.
- **Forecast:** raised from 55m to 80m, against flai's 38m, for the extraction work across nine tasks. S-0217 took 68m.
- **Cost of delay:** 427 USD a week, unchanged.

The story's `### Planning` notes give the reasons. I'd split, merge or drop nothing.
