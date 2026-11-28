---
id: T-3633
type: task
nature: improvement
title: workflow.md and the user guide say a rebase settles two stories' whole-declaration appends to one Go file
status: backlog
parent: S-0857
owner: alex
created: 2026-11-27T23:58:29Z
updated: 2026-11-27T23:58:29Z
transitions: []
stream: S-0857
tags: [docs, git]
touches: [design/system/workflow.md, docs/users/flai.md]
after: [T-3628]
---
# T-3633 workflow.md and the user guide say a rebase settles two stories' whole-declaration appends to one Go file

## Work

In `design/system/workflow.md`, after the paragraphs on the section merge (S-0401, ADR-0175; S-0637, ADR-0276), add one on the Go merge: what it settles, what it leaves to wait, and that every rebase and the trial merge share it, linking the ADR of T-3623 and I-0409. In `docs/users/flai.md`, add the Go merge to the shared documents' entry and to the two lists of what a rebase settles before the job or the batch gives up. It waits for T-3628, so the text describes the behaviour as built.

## Done when

- Both documents describe the Go merge and link its ADR, and `updated` is bumped in `workflow.md`'s front matter.
- The markdown lint and `flai check --strict` pass on both files.

## Notes
