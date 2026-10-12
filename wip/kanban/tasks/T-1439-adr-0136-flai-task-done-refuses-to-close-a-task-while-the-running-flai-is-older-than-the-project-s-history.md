---
id: T-1439
type: task
nature: remediation
title: "ADR-0136: flai task done refuses to close a task while the running flai is older than the project's history"
status: backlog
parent: S-0363
owner: alex
created: 2026-10-12T14:36:17Z
updated: 2026-10-12T14:36:17Z
transitions: []
stream: S-0363
tags: [cli]
touches: [design/adrs/0136-flai-task-done-refuses-to-close-a-task-while-the-running-flai-is-older-than-the.md, design/adrs/README.md]
---
# T-1439 ADR-0136: flai task done refuses to close a task while the running flai is older than the project's history

## Work

Record the remedy for I-0127 as an ADR before it is built. Write it with `flai adr new`, titled "flai task done refuses to close a task while the running flai is older than the project's history", status `proposed`, refining [ADR-0107](../../../design/adrs/0107-flai-task-done-closes-a-task-in-one-call-commit-sync-move-log-widen-touches.md) and [ADR-0128](../../../design/adrs/0128-flai-task-done-commits-the-paths-the-closing-task-covers-and-those-no-other.md). Add its row to `design/adrs/README.md`. If flai gives it another number or slug, `task_done` widens the touches to the file it wrote.

- **Context.** I-0127's instance: S-0338 closed T-1324 through the MCP tool `task_done` on the installed flai 1.39.13, below 1.40.0 in the project's history. It committed T-1325's and T-1326's paths with T-1324 and widened T-1324's touches, which a done task then keeps. The inbox said `flai_outdated` at the time, and the close went ahead: a warning did not stop it.
- **Decision.** `taskdone.Run` refuses before its first step, changing nothing, while `release.FlaiOutdated` reports the running flai older than the newest `flai/vX.Y.Z` tag reachable from the main checkout's HEAD. The refusal is the outdated message, naming the running and newest versions and `buildinfo.UpgradeCommand`. A dev build is never refused, as `FlaiOutdated` already holds. The MCP tool, `flai task done`, and the host's `task.done` all reach `taskdone.Run`, so one check covers the three.
- **Why it differs from ADR-0117's "the minimum warns".** `task_done` commits and moves a task to done in one call. An older flai splits the commit by older rules, and the result cannot be undone through flai once the task is done. Other reads and writes stay warn-only.
- **Alternatives.** Warn only (the instance shows it is not heeded). Refuse only the commit step (the move and log would still run on old rules). Refuse every write tool (wider than the issue, a story of its own if wanted).

Waits for no task: it is the first layer, and the code and docs tasks implement what it decides.

## Done when

- `design/adrs/0136-*.md` exists with Context, Decision, Consequences, and Alternatives considered, and `flai check --strict` passes on it.
- `design/adrs/README.md` lists it as proposed, refining 0107 and 0128.

## Notes

Drafted by the planner for S-0363.
