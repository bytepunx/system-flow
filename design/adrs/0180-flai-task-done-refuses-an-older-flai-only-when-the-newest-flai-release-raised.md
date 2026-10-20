---
id: ADR-0180
title: "flai task done refuses an older flai only when the newest flai release raised the task-done rules number its source carries, and keeps the minor-version rule for a release without one"
status: proposed
date: 2026-10-20
supersedes: []
superseded_by: []
refines: [ADR-0142, ADR-0160]
---

# ADR-0180 flai task done refuses an older flai only when the newest flai release raised the task-done rules number its source carries, and keeps the minor-version rule for a release without one

## Context

ADR-0142 refuses `flai task done` while the running flai's major or minor version is behind the newest `flai/vX.Y.Z` tag in the project's history, and ADR-0160 closes the MCP `task_done` through the installed flai once that has caught up by the same rule. Every feature story publishes a minor flai release. Ten were tagged between 2026-10-16T22:45Z and 2026-10-18T02:19Z, with `auto-upgrade` off, so the host stayed behind and `task_done` refused nearly every close (I-0157): each agent fell back to the tree's CLI, and `flai serve` opened one upgrade thread per release. Most of those releases did not change how a task closes. The question is what a release must change for an older flai to be refused.

## Decision

`flai task done` refuses only when the newest flai release changed task done's rules, by a task-done rules number the flai source carries.

1. **The number.** `flai/internal/buildinfo/task-done-rules.txt` holds one positive integer, starting at 1. `buildinfo.TaskDoneRules` is that file, embedded in the binary, and `flai version --json` prints it as `task_done_rules`. A change to what task done commits, messages, syncs, moves, logs, adds to the touches, or checks raises it by one in the same change.
2. **Where a release records it.** A `flai/vX.Y.Z` tag records the number in the file at the commit it tags, which `git show <tag>:flai/internal/buildinfo/task-done-rules.txt` reads. Nothing is written into the tag: the number is the released tree's, whichever flai tagged it.
3. **The comparison.** When the running flai is older than the newest tag, `taskdone.Run`'s start refuses the close, changing nothing, while the tag's number is above the running flai's, or the tag's major version is above the running flai's. A minor or patch release that keeps the number closes the task; the inbox's `flai_outdated` still warns.
4. **The fallback.** A tag whose commit has no such file, such as every tag before this ADR, keeps ADR-0142's rule: a major or minor version behind refuses.
5. **ADR-0160's caught-up check.** The installed flai has caught up when its `task_done_rules` is not below the tag's number and its major version is not below the tag's; when the tag or the installed flai has no number, by ADR-0142's rule as before.
6. **Unchanged.** A build that is not a release, such as `dev`, is never refused (ADR-0142 point 5). The check stays in one place, `taskdone.Run`, for the CLI, the MCP tool, and the host channel.

## Consequences

- A minor release that does not raise the number no longer stops any agent's close, with `auto-upgrade` on or off, and including the minutes between a tag and its publication.
- Whoever changes task done's steps must raise the number, or an older flai will close tasks by rules the project has outgrown. `flai-cli.md` and the comment on `TaskDoneRules` say so.
- The check reads one more git object, at the newest tag, inside `FlaiOutdated`'s minute cache.
- The path of the file is fixed: moving it makes every later tag fall back to the minor rule.
- Upgrade threads still open per release (ADR-0158, ADR-0161); fewer of them block work.

## Alternatives considered

- Turn on `auto-upgrade` for the project: builds nothing, but each release still refuses closes until it is published and installed, 10 to 15 minutes each.
- Write `task-done-rules: N` into each tag's annotation at release: the number would be the tagging flai's, and the board's Publish runs the installed flai, which can be behind the tree it tags.
- Resolve upgrade threads a newer release supersedes: fewer threads, same refusals.
- Have the orchestrator batch releases: fewer refusals, but each batch still refuses until the host is upgraded.
