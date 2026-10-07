---
id: I-0079
title: TestRoundTripRepositoryItems reads the live main checkout and fails a close-out when another agent edits a story mid-run
class: defect
status: open
count: 7
cost: 5m
first_reported: 2026-10-05T05:55:58Z
last_reported: 2026-10-07T23:53:05Z
updated: 2026-10-07T23:53:05Z
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

## Remediation

Story S-0290 remediates this issue, created from it at 2026-10-06T09:56:52Z.
