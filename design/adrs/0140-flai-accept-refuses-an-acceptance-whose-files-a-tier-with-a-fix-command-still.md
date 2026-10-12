---
id: ADR-0140
title: "flai accept refuses an acceptance whose files a tier with a fix_command still fails on after the fix, naming each file, line, and error"
status: proposed
date: 2026-10-12
supersedes: []
superseded_by: []
refines: [ADR-0139]
topics: [cli, template]
---

# ADR-0140 flai accept refuses an acceptance whose files a tier with a fix_command still fails on after the fix, naming each file, line, and error

## Context

[ADR-0139](0139-a-test-tier-s-fix-command-rewrites-the-files-a-story-changes-before-flai-task.md) runs a test tier's `fix_command` before `flai task done` commits, after each clean sync's rebase, and before `flai accept` commits the main checkout. No fix stops any of them. A finding the fix cannot fix was still reported only by the tier's `command`, in a story's close-out, and an acceptance would merge it. On TH-0413 the operator decided that markdown the fix cannot fix must block the acceptance, with the files, lines, and errors, so that the operator can correct them (S-0367).

## Decision

`flai accept` runs the `command` of each tier that carries a `fix_command` over what it is about to bring to the main branch, once the fix has run, and refuses the acceptance when one fails (`verify.Unfixed`).

1. **At preflight, before anything changes.** It fixes the main checkout's uncommitted files, then checks them. A failure refuses with nothing changed but what the fix wrote.
2. **After the story branch's rebase and fix, before the merge.** It checks the files the branch changes against the main branch in the story's worktree. A failure refuses before `git merge --ff-only`, so the main branch is unchanged and the story stays in review.
3. **The refusal is a rule refusal.** It names each finding as `path:line: message` with its tier, or the tier's last lines of output for a `plain` tier, and says where to correct the files and to accept again.
4. **What does not refuse.** A manifest whose tiers cannot be read, and paths that cannot be listed, are logged at warn. `--dry-run` neither fixes nor checks. A fix command that fails still stops nothing, as ADR-0139 decides.

## Consequences

- Lint that a tier with a `fix_command` reports never reaches the main branch through an acceptance.
- An acceptance takes the time of up to two more runs of the check. For `scripts/lint-md.sh` that is one `npx` start each.
- Uncommitted files in the main checkout that fail the check block every acceptance until they are corrected or committed.
- A check that cannot start, such as `scripts/lint-md.sh` without `npx`, fails and refuses too, since its exit status cannot be told from lint.
- The operator corrects the files by hand until the dashboard has an editor for them (S-0378).

## Alternatives considered

- **Refuse in the close-out only.** It already fails on what the tier reports, but the main checkout's files and what main brought in after the close-out reach main unchecked.
- **Refuse after the merge.** The main branch would already hold the files, and undoing a merge is not flai's to do.
