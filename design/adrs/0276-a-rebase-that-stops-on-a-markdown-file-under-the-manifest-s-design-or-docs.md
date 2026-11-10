---
id: ADR-0276
title: "A rebase that stops on a markdown file under the manifest's design or docs folder merges it by sections, as one inside claims.shared does"
status: accepted
date: 2026-11-10
supersedes: []
superseded_by: []
refines: [ADR-0175]
---

# ADR-0276 A rebase that stops on a markdown file under the manifest's design or docs folder merges it by sections, as one inside claims.shared does

## Context

[ADR-0175](0175-a-rebase-that-stops-on-a-markdown-file-inside-the-manifest-s-shared-paths.md) merges by sections a markdown file a rebase stops on, in `flai stream sync`, `flai task done`'s sync, `flai accept`, and the merge queue's batch, but only when the file is inside the manifest's `claims.shared`. A design document outside that list is left to the agent.

[I-0315](../issues/I-0315-stories-that-each-add-a-section-at-one-spot-in-a-shared-design-document.md) records the cost. The stories of E-0033 each added a "Cause N fixed" section to `design/system/server-performance.md`, at the same spot, after the S-0596 Stories table. That file is not shared. Each acceptance of one of them stopped the rebase of the next, at its sync or its own acceptance, and an acceptance that stops sends the story back to in-progress. S-0598, S-0599, and S-0603 were each sent back twice. Every stop was settled by hand the same way: keep both sections. Each cost a close-out of 10 to 20 minutes and a review cycle.

The section merge already settles that case: a section both sides added appears once per side, the base branch's first. It skipped the file only because the file was not in `claims.shared`. That list is the operator's, and a document shows up as a hub only after it has conflicted.

## Decision

A rebase that stops on a markdown file under the manifest's design folder (`layout.design`) or docs folder (`layout.docs`) merges it by sections, as ADR-0175 merges one inside `claims.shared`, in every rebase flai runs.

- The file must end in `.md`, as before.
- The issue files keep their merge by instances ([ADR-0126](0126-a-rebase-of-a-story-branch-merges-an-issue-file-both-sides-changed-by-its.md)). The generated files, `design/issues/summary.md` and `design/adrs/README.md`, are still written again ([ADR-0228](0228-flai-generates-the-adr-index-from-the-adr-files-and-checks-the-issues-summary.md)). Neither is merged by sections.
- The wip folder stays out. flai writes it in the main checkout, and no story's branch is meant to change it.
- A markdown file inside `claims.shared` outside those two folders, such as `template/CHANGELOG.md`, is still merged as before.
- ADR-0175's cases it cannot settle are unchanged. Two sides that change the same lines of one section, a section one side removes and the other changes, or a front matter key other than `updated` both change differently still leave the stop for the agent.
- The gate is one matcher, `sharedMarkdown` in `flai/internal/issues/generated.go`. `issues.FilesIn` gives it to every rebase flai runs, so no caller changes.

What stays as ADR-0175 has it:

- The trial merge still reports a conflict on these files between two open stories, because a section merge can fail.
- The merge queue's batch still treats a stop it cannot settle by `claims.shared`. A stop left on a design or docs file outside the list is a same-lines conflict, and it is left out or sent back as any other path is.

## Consequences

- Stories that each add a section at one spot in a design or docs document no longer stop at sync or acceptance, and are not sent back for it.
- Sections both stories added land in acceptance order: the base branch's first, then the story's. That is not always the order a person would pick, such as cause order. The story's own review and the markdown lint are the check, as ADR-0175 says.
- The operator need not add a hub document to `claims.shared` to have it merged. That list still decides what never holds a ready story (ADR-0096) and what the merge queue leaves waiting.

## Alternatives considered

- Add each hub document to `claims.shared`. That list is the operator's, and a file shows up as a hub only after it has conflicted. Rejected as the remedy, though the operator can still do it.
- Merge every markdown file in the repository by sections. Work items and narratives under wip are flai's, and a README or a sub-project's markdown has no section convention to rely on. Rejected: too wide.
- A planning rule that gives each story's section its own file. It changes nothing for stories already written, and it scatters one document's content across files. Rejected.
- Make the trial merge run the section merge too, so that two open stories are not told of a conflict a rebase would settle. Deferred to a story of its own: it changes what holds a ready story and what the stories are messaged.
