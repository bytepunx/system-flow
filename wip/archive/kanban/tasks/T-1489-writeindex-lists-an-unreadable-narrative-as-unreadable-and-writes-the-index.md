---
id: T-1489
type: task
nature: remediation
title: WriteIndex lists an unreadable narrative as unreadable and writes the index
status: done
parent: S-0373
owner: alex
created: 2026-10-12T17:28:46Z
updated: 2026-10-12T18:16:34Z
transitions:
  - to: ready
    at: 2026-10-12T18:14:38Z
    by: agent-S-0373
  - to: in-progress
    at: 2026-10-12T18:14:39Z
    by: agent-S-0373
  - to: done
    at: 2026-10-12T18:16:34Z
    by: agent-S-0373
stream: S-0373
tags: [cli]
touches: [flai/internal/workitem/narrative.go, flai/internal/workitem/narrative_test.go, flai/internal/workitem/activity_test.go]
usage:
  source: log
  seconds: 115
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 9
      output: 1692
      cache_read: 655251
      cache_write: 19009
      cost: 0.317
---
# T-1489 WriteIndex lists an unreadable narrative as unreadable and writes the index

## Work

This is the cause of I-0129. `Repo.ActiveNarratives` in `flai/internal/workitem/narrative.go` returns at the first `S-*.md` that `ReadNarrative` refuses, such as an empty file (`no front matter`). `WriteIndex` then fails, and so does every command that rewrites `wip/agents/index.md`: stream open, log, and state, move, edit, archive, accept, story start, and task done. Each fails after its own change is saved, so a retry writes it twice, as the duplicate log entry in S-0232's narrative shows.

Change it so that one story's narrative never fails another's write:

- `WriteIndex` lists an `S-*.md` that does not parse as one row, `| [S-nnnn](S-nnnn.md) | unreadable, see flai check | | | |`, and writes the index. This is how it already lists an unreadable activity document under `## Strategic agents`.
- The readable narratives are listed as before, in ID order, with the unreadable rows in that order among them.
- A story's own narrative that does not parse still fails that story's `flai stream log` and `flai stream state`, as before. Only the index stops depending on the other narratives.

Add a test in `flai/internal/workitem/narrative_test.go`, beside `TestSetStreamStateReplacesBothSections`. Put one good narrative, one empty `S-nnnn.md`, and one with no front matter in a project. `LogStream` on the good story and `WriteIndex` both succeed, and the index lists all three, the two bad ones as unreadable. This reproduces I-0129: it fails before the change.

It waits for no task. The row format above is fixed by the plan, so the check task of this story can rely on it at the same time.

## Done when

- The new test fails on the main branch and passes with the change.
- `flai test` passes on `flai/internal/workitem/narrative.go` and `flai/internal/workitem/narrative_test.go`.
- `TestStreamsAndIndex` and `TestIndexListsAnUnreadableActivityDocument` pass unchanged.

## Notes

The precedent is the activity-document loop in `WriteIndex` and its test in `flai/internal/workitem/activity_test.go`. `ActiveNarratives` has no other caller, so its signature may change.
