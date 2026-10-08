---
id: I-0079
title: TestRoundTripRepositoryItems reads the live main checkout and fails a close-out when another agent edits a story mid-run
class: defect
status: closed
count: 11
cost: 5m
first_reported: 2026-10-05T05:55:58Z
last_reported: 2026-10-08T08:27:55Z
updated: 2026-10-08T09:02:45Z
---

# I-0079 TestRoundTripRepositoryItems reads the live main checkout and fails a close-out when another agent edits a story mid-run

## Description
TestRoundTripRepositoryItems reads the live main checkout and fails a close-out when another agent edits a story mid-run

## Instances

### 2026-10-05T05:55:58Z
Story: S-0276.
S-0276's close-out stopped at the full Go tests: TestRoundTripRepositoryItems in flai/internal/workitem read wip/kanban/stories/S-0251-*.md from the main checkout while a planner was writing it, and reported only an updated timestamp and a Tasks line. Two reruns passed. The close-out had to be finished step by step.

### 2026-10-06T22:58:33Z
Story: S-0299.
S-0299's close-out: TestRoundTripRepositoryItems read the main checkout's wip/ while a planner rewrote draft stories (S-0294, then S-0274) and failed; it passed on a rerun with -count=1.

### 2026-10-07T03:19:14Z
Story: S-0269.
S-0269's second close-out stopped at the full Go tests: TestRoundTripRepositoryItems read S-0246, S-0265, and S-0279 in the main checkout while flai's forecast replan rewrote them (forecast and timestamps differ); internal/workitem passed alone.

### 2026-10-07T08:25:53Z
Story: S-0213.
S-0213's close-out stopped at integration in flai/internal/workitem while S-0214's agent was writing a task file in the main checkout; the package passed alone on a rerun a minute later.

### 2026-10-07T09:09:01Z
Story: S-0246.
S-0246's close-out stopped at integration: TestRoundTripRepositoryItems read S-0215's T-0934 in the main checkout while S-0215's agent was writing its usage block, so its marshal differed; S-0246 changes no work-item code. The close-out was run again.

### 2026-10-07T23:37:11Z
Story: S-0280.
S-0280's close-out failed its integration tier in flai/internal/workitem at 2026-10-07T23:35Z while planner-S-0310 was editing S-0310 and T-1280 in the main checkout (23:30:30Z to 23:31:09Z). The kept output was only the tail, which named the package but not the test (I-0102, S-0313). go test ./internal/workitem passed when run again a minute later.

### 2026-10-07T23:53:05Z
Story: S-0310.
S-0310's second close-out stopped at the integration tier: flai/internal/workitem FAILed, its output showing S-0313's story body, which planner-S-0313 was editing on main at 23:38 while the run read it. S-0310 changes nothing under internal/workitem; the run before it passed every tier, and stopped only because main moved.

### 2026-10-08T00:22:00Z
Story: S-0317.
S-0317's close-out stopped at the integration tier: flai/internal/workitem failed while the planner wrote the planning notes of the S-0320 story (I-0106's remedy) in the main checkout during the run; the tier's output was that story's body. Nothing S-0317 changed is in internal/workitem. The close-out was run again.

### 2026-10-08T04:24:20Z
Story: S-0316.
S-0316's close-out at 04:19Z failed integration in flai/internal/workitem, which S-0316 does not touch. The kept output was the body of S-0326's story, the planning notes S-0326's planner was writing in the main checkout during the run. The close-out showed only the last lines, so the failing test is not named (I-0113). Run again.

### 2026-10-08T08:17:30Z
Story: S-0321.
S-0321's close-out failed its integration tier in flai/internal/workitem while agent-S-0291 moved T-1340 to ready and in-progress in the main checkout (08:14:03Z to 08:14:04Z). The test passed on a rerun with the same branch head.

### 2026-10-08T08:27:55Z
Story: S-0321.
S-0321's next close-out failed integration in flai/internal/workitem again, while agent-S-0291 cancelled T-1340 in the main checkout at 08:24:52Z. The package passed run alone with -race straight after (go test -race ./internal/workitem, 7.3s).

## Remediation

Story S-0290 remediates this issue, created from it at 2026-10-06T09:56:52Z.
Closed 2026-10-08T09:02:45Z: Fixed by S-0290 (T-1367): TestRoundTripRepositoryItems now parses, validates, and compares one read of each listed item file, so an item another agent rewrites in the main checkout during a close-out is compared with itself; Save writes atomically, so one read is one whole version. A file gone since the list is skipped. TestRoundTripItemsRewrittenSinceTheList reproduces the race, and fails against the old two-read check.
