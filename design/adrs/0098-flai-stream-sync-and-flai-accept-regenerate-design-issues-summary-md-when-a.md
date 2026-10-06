---
id: ADR-0098
title: "flai stream sync and flai accept regenerate design/issues/summary.md when a rebase stops on it alone, and the trial merge does not report it"
status: accepted
date: 2026-10-06
supersedes: []
superseded_by: []
refines: [ADR-0069]
topics: [cli]
---

# ADR-0098 flai stream sync and flai accept regenerate design/issues/summary.md when a rebase stops on it alone, and the trial merge does not report it

## Context

`flai issue new`, `bump`, and `close` regenerate `design/issues/summary.md` from the issue files, and the file is committed. Each story records its issues on its own branch. When two branches both change the file, they both rewrite its `updated:` line, and the rows they add or remove are often within three lines of each other. Git then cannot merge them. [I-0074](../issues/I-0074-two-stories-that-each-record-or-close-an-issue-always-conflict-in-design-issues-summary-md-whose-updated-line-both-rewrite.md) counts five instances in two days:

- `flai stream sync`'s trial merge opened a thread for the pair over this one file, twice (TH-0126, TH-0130), and once more on TH-0175. Each story's agent then spent a turn reading the thread.
- The sync that rebases a story onto main after the other was accepted stopped on the file, twice. The agent regenerated it with `flai issue summary` and continued the rebase by hand.

[I-0089](../issues/I-0089-design-issues-summary-md-is-generated-and-committed-so-a-story-branch-that-records-an-issue-conflicts-with-any-issue-recorded-on-main-meanwhile.md) is the same conflict against an issue recorded on main meanwhile.

The resolution is always the same and needs no judgement. Every row of the file comes from the issue files, so once those are merged the right `summary.md` is whatever regeneration writes.

## Decision

**When a rebase that `flai stream sync` runs stops, and `design/issues/summary.md` is the only path in conflict, flai regenerates it from the issue files in the worktree, `git add`s it, and continues the rebase.** It repeats this for each commit the rebase stops on, until the rebase finishes or stops on some other path. This refines [ADR-0069](0069-agents-commit-each-task-on-the-story-branch-after-flai-stream-sync-and-never.md). The rebase still never touches uncommitted work, and any other conflict is still left for the agent with each path named.

- `flai accept` merges through the same sync, so acceptance regenerates the file too.
- When other paths conflict as well, the rebase stops as before and lists them all, `summary.md` included. Its rows depend on the issue files that may be among them, so it is regenerated only once nothing else is in conflict.
- The trial merge leaves `summary.md` out of the paths a pair of open branches conflicts on. A pair whose only conflict is that file counts as clean, so no thread is opened for it, and the pair's open thread is resolved as for any clean pair.
- The generated paths are one list in flai, holding `summary.md` alone today, so a later generated file joins it in one place.

## Consequences

- Two stories that each record or close an issue no longer stop each other's sync or acceptance on `summary.md`, and no agent reads a thread about it. The same goes for an issue recorded on main while a story works (I-0089).
- The regenerated `updated:` line carries the time of the sync, not of either side's last change. The line says when the file was last generated, which it was.
- Two branches that bump the same issue still conflict in that issue's file: its count and instances are data both sides add to, not derived. That is [I-0092](../issues/I-0092-two-story-branches-that-each-bump-the-same-issue-conflict-in-its-front-matter-whose-count-last-reported-and-updated-lines-both-rewrite.md), for S-0297.

## Alternatives considered

- A git merge driver for the file in `.gitattributes`. A driver must also be named in each clone's `git config`, which flai does not install, and `git merge-tree` runs no merge driver at all, so the trial merge would still report the file.
- Dropping the `updated:` line. The rows two stories change still meet within git's three lines of context, as in TH-0126.
- Not committing the file and generating it where it is read. The dashboard, the review page, and the prime pack read it, and a fresh clone would lack it. That change is wider than this story.
