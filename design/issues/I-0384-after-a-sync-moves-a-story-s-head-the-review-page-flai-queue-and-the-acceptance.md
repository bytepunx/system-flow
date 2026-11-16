---
id: I-0384
title: After a sync moves a story's head, the review page, flai queue, and the acceptance preview report the verify failure from before the sync while a new close-out runs at the new head
class: defect
status: closed
count: 1
cost: 20m
first_reported: 2026-11-16T00:41:30Z
last_reported: 2026-11-16T00:41:30Z
updated: 2026-11-16T03:40:17Z
---

# I-0384 After a sync moves a story's head, the review page, flai queue, and the acceptance preview report the verify failure from before the sync while a new close-out runs at the new head

## Description
After a sync moves a story's head, the review page, flai queue, and the acceptance preview report the verify failure from before the sync while a new close-out runs at the new head

## Instances

### 2026-11-16T00:41:30Z
On 2026-11-16 at 00:35Z S-0739 and S-0747 had been synced onto main (0 commits behind) and their close-out jobs S-0739-20261116T002837Z and S-0747-20261116T003049Z were running the integration tier. flai queue still listed both as verify_failed at its step sync, naming the commits da74987 and be8fafb from before the sync, and the operator read that as the sync having done nothing. The last report under .flai-cache/verify stays the finished one until the running job writes its own; nothing says a close-out is under way at another head.

## Remediation

While a close-out job runs for a story at a head other than its last report's, `flai queue`, the acceptance preview, and the review page say so, naming the job and its step, in place of or beside the last finished report's failure: the job's state is in `.flai-cache/verify/jobs/<S-nnnn>.job.json`. A `verify_failed` whose report's commit is not the branch's head says the report is from before the head moved.

Story S-0799 remediates this issue, created from it at 2026-11-16T00:41:49Z.
Closed 2026-11-16T03:40:17Z: S-0799 fixed it. A running close-out job now records the head it runs at (T-3409), and verify.LastReportAt reads the last report against the branch head, giving moved and the running job with its step (T-3414). flai queue lists such a story as verify_running with the job, and a failed report from before the head moved says so (T-3415). The acceptance preview carries verify_running beside unchanged blockers (T-3416). flai verify --last and verify.status say the same (T-3417), and so do the review page and the merge queue panel (T-3418).
