---
id: ADR-0139
title: "A test tier's fix_command rewrites the files a story changes before flai task done commits them, after each clean sync's rebase, and before flai accept commits the main checkout, and never stops any of them"
status: proposed
date: 2026-10-12
supersedes: []
superseded_by: []
refines: [ADR-0069, ADR-0107]
topics: [cli, template]
---

# ADR-0139 A test tier's fix_command rewrites the files a story changes before flai task done commits them, after each clean sync's rebase, and before flai accept commits the main checkout, and never stops any of them

## Context

S-0367 asks that bad markdown not reach the main branch. Three things let it through.

1. **Nothing fixes markdown before it is committed.** The `markdown` test tier runs `scripts/lint-md.sh` over the files a change selects, and reports what it finds. An agent fixes each finding by hand, or a close-out finds it later. `markdownlint-cli2 --fix` fixes most of the rules on its own.
2. **A rebase brings old findings into a story's files.** A file the branch changes may carry a bad line that another story committed to the main branch. After the rebase, the story's close-out reports it, and the story's agent fixes a line it did not write.
3. **The main checkout is committed without a story.** `flai accept` writes work items, narratives, threads, and conversations under `wip/` in the main checkout and commits them with `git add -A`. No story's lint runs over them.

A test tier already names the command that checks a set of files. The question is where the command that fixes them is named, and when flai runs it.

## Decision

A test tier in `system-flow.yaml` may carry `fix_command`, and flai runs it at three points, over the files each point is about to commit. A fix that fails never stops any of them.

1. **`fix_command` is an argument list beside `command`.** It runs as `command` runs: never through a shell, in the tier's `dir`, with the same placeholders, `{files}`, `{packages}`, and `{importers}`, filled for the paths the tier's `paths` select. It runs only over files that exist. It rewrites them in place. Its exit status is ignored: a command that cannot start or exits non-zero is logged at warn, and what it changed before it stopped is kept. The tier's `command` still reports what is left. A tier with no `fix_command` is not run. `verify.Fix` runs the tiers in list order and answers the paths whose content changed.
2. **`flai task done` fixes before it commits.** Over the paths it is about to commit, before the commit step, so the commit holds the fixed content. The MCP tool `task_done` and the host channel's `task.done` do the same.
3. **A clean rebase is followed by a fix.** `flai stream sync`, and the rebase `flai accept` runs on the story's branch, run the fix over every file the branch changes against the main branch, once the rebase is clean and the issue folds are committed. What it changed is committed on the branch as `chore: [S-nnnn] fix the lint of the files the branch changes, after the sync`, and the result lists those paths as `fixed`. Every clean sync does it, the close-out's included, so what the branch changes is fixed before review.
4. **`flai accept` fixes the main checkout before it commits.** Over the paths it has not committed, the work items, narratives, threads, and conversations it wrote, just before its `git add -A` and commit.

This repository's `markdown` tier carries `fix_command: [scripts/lint-md.sh, --fix, "{files}"]`. This refines [ADR-0107](0107-flai-task-done-closes-a-task-in-one-call-commit-sync-move-log-widen-touches.md), whose commit step now fixes first, and [ADR-0069](0069-agents-commit-each-task-on-the-story-branch-after-flai-stream-sync-and-never.md), whose sync may now add a commit of its own to the branch.

## Consequences

- What a fix command can fix never reaches a commit flai makes, on a story's branch or on the main branch. Findings it cannot fix are still reported by the tier's `command`.
- A story's agent no longer fixes by hand a bad line a rebase brought in. The sync's fix commit holds it, so the story's diff can include lines it did not write.
- A sync, a task's close, and an acceptance each take as long as the fix commands take. For `scripts/lint-md.sh` that is one `npx` start and the files it is given.
- A fix command must be safe to run over any file its tier selects, and to run again. One that rewrites a file wrongly is committed without review.
- A project with no `fix_command` sees no change. The template's manifest carries none.

## Alternatives considered

- **Run the fix in the close-out only.** It catches every file before review, but a task's commits keep the bad content, and the main checkout's commits get none.
- **Fix every file in the repository after a rebase.** It would put changes to files the story never touched in the story's diff, and overlap other stories' claims.
- **A git pre-commit hook.** It runs for commits flai does not make and is skipped by others, a hook lives outside the repository's files, and a hook that rewrites files it is committing is fragile.
- **Stop the step when the fix fails.** A missing `npx` or a network fault would then stop a task's close, a sync, or an acceptance, for a finding the tier's `command` reports anyway.
