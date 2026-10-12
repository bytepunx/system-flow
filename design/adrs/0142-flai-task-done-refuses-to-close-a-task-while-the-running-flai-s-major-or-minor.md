---
id: ADR-0142
title: "flai task done refuses to close a task while the running flai's major or minor version is behind the newest flai release in the project's history"
status: proposed
date: 2026-10-12
supersedes: []
superseded_by: []
refines: [ADR-0107, ADR-0128]
---

# ADR-0142 flai task done refuses to close a task while the running flai's major or minor version is behind the newest flai release in the project's history

## Context

[ADR-0107](0107-flai-task-done-closes-a-task-in-one-call-commit-sync-move-log-widen-touches.md) made `flai task done` close a task in one call: commit, sync, move, log, widen touches, check, and inbox. [ADR-0128](0128-flai-task-done-commits-the-paths-the-closing-task-covers-and-those-no-other.md) changed its commit step so the tasks of one layer close apart. A flai that predates a rule like ADR-0128 still runs the close, by its own older rules.

[I-0127](../issues/I-0127-an-outdated-mcp-flai-s-task-done-commits-every-changed-path-under-the-closing-task-so-a-layer-s-tasks-cannot-close-apart.md) records it. S-0338 closed T-1324 through the MCP tool `task_done` on the installed flai 1.39.13, while 1.40.0 was tagged in the project's history. It committed T-1325's and T-1326's paths under T-1324 and widened T-1324's touches with them, which a done task keeps. The inbox reported `flai_outdated` at the time, and the close went ahead: the warning was not heeded.

Today an outdated flai only warns: in `inbox`, `project.info`, `flai serve`, and the inbox `flai task done` answers. [ADR-0117](0117-the-host-api-installs-a-release-the-dashboard-names-only-when-flai-lists-it-as.md) decided that the manifest's `flai.minimum` warns and is never refused.

## Decision

`flai task done` refuses to close a task, changing nothing, while the running flai's major or minor version is behind the newest flai release in the project's history.

1. **What is behind.** The newest release is the highest `flai/vX.Y.Z` tag reachable from the main checkout's HEAD, as `release.FlaiOutdated` reads it. The close is refused when the running flai is older than that tag and its major or minor version is below the tag's. The threshold is the operator's, confirmed on TH-0401.
2. **A patch behind warns.** A flai of the same major and minor version but an older patch closes the task as before; the inbox it answers carries `flai_outdated`, as today.
3. **Changing nothing.** The check runs in `taskdone.Run`'s start, before the commit step: no commit, message, sync, move, log, touches, or check. The refusal names the running and newest versions and the upgrade command, `flai host upgrade` (`flai self-upgrade` where no host runs).
4. **One place.** The MCP tool `task_done`, `flai task done`, and the host channel's `task.done`, which runs the CLI, all reach `taskdone.Run`, so the one check covers the three.
5. **Development builds.** A flai whose version is not a release, such as `dev` or the tree's `scripts/flai.sh` build, is never refused, as `FlaiOutdated` already holds.

Every other command keeps ADR-0117's warn-only rule.

## Consequences

- An outdated flai can no longer split a layer's commits by older rules, nor widen a done task's touches with a sibling's paths.
- Once a minor or major release is tagged on main, every story's agent on a host still running the older flai stops at its next task close until the operator runs `flai host upgrade`. The agent asks for it on a thread. Publishing a patch release stops no one.
- A patch behind can still lack a fix, as 1.39.13 lacked S-0322's in 1.39.14. The inbox's `flai_outdated` warning covers that case.

## Alternatives considered

- **Warn only.** I-0127's instance closed with the warning showing.
- **Refuse on any version behind, a patch included.** Every publish would stop every story's agent until the host upgrades; the operator chose the minor threshold on TH-0401.
- **Refuse only the commit step.** The move, log, and touches would still run on older rules.
- **Refuse every write tool.** Wider than the issue; a story of its own if wanted.
