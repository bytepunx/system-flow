---
id: ADR-0311
title: "The merge queue's preflight checks main again with a flai built from main when the running flai is behind main's flai sub-project, and otherwise holds naming the installed flai and the release to install"
status: accepted
date: 2026-11-16
supersedes: []
superseded_by: []
refines: [ADR-0201, ADR-0234]
---

# ADR-0311 The merge queue's preflight checks main again with a flai built from main when the running flai is behind main's flai sub-project, and otherwise holds naming the installed flai and the release to install

## Context

[ADR-0234](0234-before-it-builds-a-batch-the-merge-queue-runs-the-landing-s-preflight-on-the.md) has the merge queue run `flai check` on the main checkout before it builds a batch, and hold itself on each error-level finding with one thread to the operator. It runs the check in flai serve's own process, with the installed flai.

[I-0382](../issues/I-0382-the-merge-queue-s-preflight-holds-the-queue-on-main-s-files-when-the-installed.md) records what that costs when main is ahead of the installed flai. S-0748 changed the role files' form and the code that reads it in one acceptance. flai serve ran flai 1.99.0, built before S-0748, so its check reported `roles.file` on two role files that were right, held the queue for about three hours, and told the operator to fix them. A flai built from the tree found nothing. The hold cleared only once 1.100.0 was installed, and the story that would publish it waited in the held queue.

[ADR-0201](0201-when-the-manifest-names-flai-command-task-done-s-check-step-runs-that-command-s.md) already lets a project that builds flai name `flai.command`, such as `scripts/flai.sh`, so that a branch's own rules check the branch at task done.

## Decision

When the preflight's check finds an error-level finding and the running flai is behind main's flai sub-project, the preflight checks main again with the manifest's `flai.command`, and holds only on what that check finds. Where it cannot, it holds as before, and the hold names the cause and the remedy.

- **Behind.** The running flai is behind when main has commits under the flai sub-project's path, the manifest's project named `flai`, or the whole tree where there is none, that the running flai was not built from. The point compared from is the running flai's build commit when the repository has it; else the tag of its own release, `flai/v<version>`; else the newest `flai/v*` tag merged in main. The count is `git rev-list --count <point>..<main> -- <path>`. A flai that git cannot place is not behind, and the preflight logs why.
- **The second check.** It runs only when both hold: the in-process check found an error-level finding, and the running flai is behind. It runs `<main checkout>/<flai.command> check --json` in the main checkout, unscoped, and reads the result it prints, as task done reads it (ADR-0201). Its answer counts only when `<flai.command> version --json` names a flai that is not behind main itself, so that a wrapper that fell back to the installed flai cannot confirm the installed flai's finding. Its error-level findings, but those on the files a landing leaves (ADR-0178), replace the in-process check's. The lint part of the preflight is unchanged.
- **A hold with its cause.** When the running flai is behind and the second check could not settle it, because the manifest names no `flai.command`, the command gave no result, or its flai is behind too, the hold records the cause beside its findings: the installed flai's version and commit, the point compared from, how many of main's commits under the path it lacks, the newest flai release, and why the second check did not run or count. The thread to the operator and `flai queue` say that the findings may be the installed flai's rather than the files', and that the remedy is to publish the release with `flai release --pending` and install it with `flai host upgrade`. They no longer ask the operator to fix files that may be right. The hold clears as ADR-0234 says, on the first look that finds main clean, which follows the install.
- **What stays.** A running flai that is not behind checks as before, once. A finding the second check confirms holds as ADR-0234 says, with no cause. The `preflight_held` and `preflight_cleared` events keep their fields.

This refines ADR-0234 and ADR-0201.

## Consequences

- In this repository, which names `scripts/flai.sh`, a format change that the installed flai cannot read no longer holds the queue: the tree's flai checks main and finds it clean.
- A project that does not build flai never runs the second check: the running flai is behind only when main's history holds commits the running flai lacks, and such a project has no flai sub-project to be ahead.
- The second check builds flai from the main checkout when its sources changed, which can take a minute, and runs only on a look whose in-process check found an error.
- Where the hold stands, the operator reads the cause and the remedy, publish and install, instead of a list of files to fix.

## Alternatives considered

- **Always check main with `flai.command`.** Lost: it builds and runs a second flai on every look that would build a batch, for the rare case of a stale installed flai.
- **Hold with the cause alone, never checking again.** Lost: the hold would stand until a release is published, and the publish story would wait in the held queue, as I-0382 records.
- **Hold an acceptance that changes a format the running flai reads, such as the role files, the workflow file, or the manifest, until a flai that reads it is installed** (the issue's other option). Lost for now: it holds every such acceptance, even where the second check would settle it, and it needs a list of formats kept in step with the code. The cause on the hold covers a project with no `flai.command`.
- **Trust whatever `flai.command` answers.** Lost: `scripts/flai.sh` falls back to the installed flai when go is missing or the build fails, and would then confirm the very finding the second check exists to test.
