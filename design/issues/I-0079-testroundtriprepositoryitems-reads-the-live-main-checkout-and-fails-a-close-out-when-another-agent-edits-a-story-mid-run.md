---
id: I-0079
title: TestRoundTripRepositoryItems reads the live main checkout and fails a close-out when another agent edits a story mid-run
class: defect
status: open
count: 3
cost: 6m
first_reported: 2026-10-05T05:55:58Z
last_reported: 2026-10-07T03:19:14Z
updated: 2026-10-07T03:19:14Z
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

## Remediation

Story S-0290 remediates this issue, created from it at 2026-10-06T09:56:52Z.
