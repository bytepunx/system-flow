---
id: T-1319
type: task
nature: remediation
title: An ADR refining ADR-0085 and ADR-0098 records the remedy for I-0112, proposed from its instance
status: backlog
parent: S-0326
owner: alex
created: 2026-10-08T00:33:37Z
updated: 2026-10-08T00:33:37Z
transitions: []
stream: S-0326
tags: [adr, design]
touches: [design/adrs, design/adrs/README.md]
---
# T-1319 An ADR refining ADR-0085 and ADR-0098 records the remedy for I-0112, proposed from its instance

## Work

Propose the remedy for I-0112 from its one instance and record it with `flai adr new`, refining ADR-0085 (a close-out records findings outside the story in one issue per rule) and ADR-0098 (a sync regenerates `design/issues/summary.md`).

The instance: S-0275's close-out opened I-0112 on its branch for its `narrative.state` notes, while S-0265's branch had opened I-0111 under the same title. Once S-0265 was accepted and S-0275 synced, both files stood on the branch, the next close-out bumped I-0111 (`issues.FindOpenByTitle` takes the lowest ID), and I-0112 was left a duplicate, removed by hand.

`issues.RecordOnce` in `flai/internal/issues/record.go` matches titles in the story's own checkout only, so it cannot see an issue another story branch opened. Looking there is no remedy: it cannot bump a file its branch does not have.

The proposal to weigh: when a rebase of a story branch onto the main branch completes (`flai stream sync`, the sync `flai task done` runs, and `flai accept`'s rebase through `syncStoryBranch`), flai folds each open issue the branch added into the open issue of the same title that the main branch has. It moves the branch issue's instances into the main branch's issue as bumps (count, cost, `last_reported`), deletes the branch's file, regenerates the summary, and commits the fold on the branch. Say what happens to a story already made from the folded issue, and why the fold is not done at record time.

It goes first because the story asks for the remedy to be proposed before it is built, and the tasks after it build what it decides. The ADR's file name is not known until `flai adr new` numbers it, so this task touches the folder.

## Done when

- An accepted ADR names I-0112's cause, the remedy, where flai applies it, and what it does to a story made from the folded issue
- It links ADR-0085 and ADR-0098 as the ADRs it refines
- `flai check --strict` passes

## Notes
