---
id: ADR-0201
title: "When the manifest names flai.command, task done's check step runs that command's flai check in the story's worktree in place of the running flai's own check, and stops when it gives no result"
status: accepted
date: 2026-10-25
supersedes: []
superseded_by: []
refines: [ADR-0107, ADR-0160, ADR-0180]
topics: [cli, mcp]
---

# ADR-0201 When the manifest names flai.command, task done's check step runs that command's flai check in the story's worktree in place of the running flai's own check, and stops when it gives no result

## Context

`flai task done`, the MCP tool `task_done`, and the host channel's `task.done` close a task in one call (ADR-0107). The check step ran `flai check --strict` scoped to the story in-process, so the rules of the flai running the close checked the worktree.

I-0194 records five closes in three stories where that flai did not know a format the worktree carried, and stopped at the check after committing, syncing, and moving the task. S-0476's branch added the `answers` key to examples files (ADR-0187); S-0429's worktree carried S-0428's exemption of digest entries from `doc.updated`, accepted and not yet published; S-0432's carried S-0431's `curate` role, accepted and not yet published. Each time the MCP server ran the installed flai (1.67.13, then 1.68.1), and each time the tree's `scripts/flai.sh check --strict --story` passed. The close was finished by hand with the tree's CLI.

ADR-0160's close through the installed flai does not reach this: no release yet knew the format. The cause is only met by a project whose tree builds the flai that reads it, as this repository's does.

## Decision

When the story's worktree's manifest names `flai.command`, task done's check step runs that command's `flai check` in the worktree in place of the running flai's own check, and stops when it gives no result.

- The manifest's `flai` block takes `command` beside `minimum`: a path relative to the repository root, of a program that runs the flai the tree builds, such as `scripts/flai.sh`. `Load` refuses an absolute path, or one that leaves the repository, naming the key. Only a project that builds flai sets it; the template's manifest does not.
- The step runs `<worktree>/<command> check --strict --story <S-nnnn> --json` in the worktree, through the runner, and reads the result it prints in place of the in-process check: its findings inside the story stop the close, and its notes outside are listed as ADR-0198 lists them. The step passes when the command exits 0. The result's `check.checked_by` names the command and the version its `version --json` prints, as `closed_by` names the installed flai (ADR-0160).
- A command that cannot run, or whose output holds no check result, stops the step as a failed check does, with the command's output in the error: a branch whose flai does not build closes no task. Without the key the step is unchanged.
- The change raises the task-done rules number in `flai/internal/buildinfo/task-done-rules.txt` from 1 to 2 (ADR-0180): what task done checks a branch with has changed, and a flai older than this one would check with its own rules again. Once this release is published, an older flai refuses the close, and the MCP `task_done` closes with the installed flai once it has caught up, as ADR-0160 does.

## Consequences

- A story that widens a format the check reads closes its tasks through the MCP `task_done` once the installed flai carries this change, since the branch's own flai checks it.
- Each close with the key set runs the tree's flai twice more: `scripts/flai.sh` builds `bin/flai` in the worktree when a source is newer, so the first close after a Go change waits for the build.
- The close-out was already right: `scripts/close-out.sh` runs the tree's `flai verify`.

## Alternatives considered

- Run the check with the installed flai at the server's path, extending ADR-0160. The instances rule it out: no release knew the format.
- Say in the stop that the finding may come from a flai older than the branch's format. The stop and the hand-finished close remain.
- Detect flai's own Go module in the tree and build it. It ties flai to this repository with no setting, and guesses how to build it.
- Leave the rules number at 1. An older flai would keep checking with its own rules after this release, the cause itself.
