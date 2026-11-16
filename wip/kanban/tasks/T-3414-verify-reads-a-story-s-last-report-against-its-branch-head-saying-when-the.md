---
id: T-3414
type: task
nature: remediation
title: verify reads a story's last report against its branch head, saying when the report is from before the head moved and which close-out job runs at the new head
status: backlog
parent: S-0799
owner: alex
created: 2026-11-16T01:45:32Z
updated: 2026-11-16T01:45:32Z
transitions: []
stream: S-0799
tags: [flai, verify]
touches: [flai/internal/verify/story.go, flai/internal/verify/story_test.go]
after: [T-3409]
---
# T-3414 verify reads a story's last report against its branch head, saying when the report is from before the head moved and which close-out job runs at the new head

## Work

`flai queue`, the acceptance preview, and the review page each read the story's last report, `verify.LastReport` (`flai/internal/verify/story.go`, line 1242), and none of them compares a failed report's commit with the branch's head. Give them one reading to share, so that the three tasks after this one say the same thing.

- Give `Report` (`story.go`, line 99) two fields that are not stored and are filled only by the new reader: `Moved bool` (JSON `moved`, omitted when false), true when the report's `Commit` is not the branch's head; and `Running` (JSON `running`, omitted when nil), the story's running close-out job as `closeout` answers it from T-3409, with its id, start, head, and step, set only when the job's head is not the report's commit.
- Add the reader beside `LastReport`, such as `LastReportAt(project, story, head string)`, which fills both. Callers that have no head give an empty one, and `Moved` stays false.
- Make the text of `QueueTiersNotRun` (line 595) say, when a job runs at the head, that a close-out runs there (its job and step) in place of only "its last flai verify ran at X, not at its head Y".
- Test it in `story_test.go`: a failed report at commit A with the head at B gives `moved`; with a running job at B it also gives `running` with the job's step; a report at the head gives neither. The first case is I-0384's: before this change nothing said the report was from before the sync.

This task waits for T-3409, whose `closeout` reader it calls.

## Done when

- `LastReportAt` answers `moved` and `running` as above, and `LastReport`'s callers are unchanged.
- `flai test flai/internal/verify` passes.

## Notes
