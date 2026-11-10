---
id: ADR-0257
title: "A rebase of a story branch puts the branch's own changelog lines back under the base branch's top Unreleased, writing each CHANGELOG.md as the base branch has it with the lines the branch added"
status: accepted
date: 2026-11-10
supersedes: []
superseded_by: []
refines: [ADR-0228]
---

# ADR-0257 A rebase of a story branch puts the branch's own changelog lines back under the base branch's top Unreleased, writing each CHANGELOG.md as the base branch has it with the lines the branch added

## Context

[ADR-0228](0228-flai-generates-the-adr-index-from-the-adr-files-and-checks-the-issues-summary.md) §4 gives `CHANGELOG.md` `merge=union`, so that two stories' changelog lines never stop a rebase. A story writes its lines under the top `## Unreleased`, and `flai release` names that section the version and puts a fresh, empty `## Unreleased` above it.

Union keeps both sides' lines of every hunk. That is right while both sides only add lines under one `## Unreleased`. It is wrong when main released in between, because the release changed the line the story's hunk sits next to. [I-0251](../issues/I-0251-a-story-branch-rebased-across-a-release-brings-back-the-changelog-s-old.md) records it twice on S-0518, on 2026-11-02:

1. Release 1.1.17 renamed `## Unreleased` to `## 1.1.17` and added an empty `## Unreleased` above it. Rebasing `story/S-0518` across it kept both sides: the old `## Unreleased` heading came back, with S-0518's line and four lines already released under 1.1.17. The markdown lint then failed the merge queue's batch on `MD024/no-duplicate-heading`, and S-0518 was sent back.
2. Release 1.1.18 did the same an hour later. This time the rebase put S-0518's line inside `## 1.1.18`, and the union dropped main's second `## 1.1.3` heading, which the sync reported as `undoes_main` ([ADR-0155](0155-a-rebased-story-branch-that-undoes-main-s-change-to-a-path-is-recorded-at-sync.md)).

Both times the agent fixed it by hand, from main's file plus the story's line under the top `## Unreleased`. That fix is mechanical: what the story adds to a changelog is lines, and where they go is the top `## Unreleased`.

## Decision

**A rebase of a story branch puts the branch's own changelog lines back under the base branch's top `## Unreleased`: each `CHANGELOG.md` the branch changes is written as the base branch has it, with the lines the branch added under that heading, and committed on the branch.** This refines ADR-0228 §4. `merge=union` stays.

1. **The lines.** The branch's own lines in a changelog are the lines its version before the rebase holds beyond its version at the branch's old merge base, compared as a multiset of line text. Blank lines and headings (lines starting with `#`) are left out: a story adds entries, not sections. A line the base branch's top `## Unreleased` already holds is not added again.
2. **The file.** It is the base branch's file with those lines appended, in their order, to the body of its first `## Unreleased` section. A base file with no `## Unreleased` gets one above its first `##` heading, as `flai release` writes it. A changelog the base branch does not have, or one the branch did not change, is left as the rebase made it.
3. **Where.** The step runs after every rebase that ends clean, before the check of what the rebase undoes (ADR-0155) and before the fold (ADR-0126): in `storygit.Rebase`, which `flai stream sync`, the sync `flai task done` runs, and `flai accept` share, and in the merge queue's batch after each story's rebase onto the batch. A file the step changes is committed on the branch, or on the batch's copy, as `docs: [S-nnnn] put the changelog lines under ## Unreleased`. The sync's report and its `--json` name each file, as `changelogs`.
4. **The merge.** The release package holds the merge, `release.SettleChangelog`, beside the code that cuts `## Unreleased` at a release, so the two read the section the same way.

## Consequences

- A story's changelog line stays under the top `## Unreleased` across any number of releases, and no released line or old heading comes back. The batch's markdown lint no longer fails on it, and the sync reports no `undoes_main` on a changelog.
- A story can no longer change or remove a line main has in a changelog: the step writes main's lines as they are. Released entries are history and are not edited; a story that must correct one does it in a story of its own and keeps the removal at acceptance (ADR-0155).
- A story that wrote its line under a version heading has it moved under `## Unreleased`, which is where `git.md` says it goes.
- The step adds at most one commit per sync, and only when a changelog differs from what main and the story's lines make.

## Alternatives considered

- **A merge driver for `CHANGELOG.md`.** Each clone must name a driver in its `git config`, which flai does not install. ADR-0228 rejected it for that reason.
- **Changelog fragments, one file per story, collected at release.** It removes the shared file, but changes the template's layout, `flai release`, and every project's changelog. The fix here needs none of that.
- **Fold only when the file has two `## Unreleased` headings,** as the first instance suggested. The second instance had one: the line had moved inside `## 1.1.18`. Only a rewrite from main's file catches both.
- **Leave it to the agent.** That is what I-0251 counts, 13 minutes each time, and once it reached the batch and sent a story back.
