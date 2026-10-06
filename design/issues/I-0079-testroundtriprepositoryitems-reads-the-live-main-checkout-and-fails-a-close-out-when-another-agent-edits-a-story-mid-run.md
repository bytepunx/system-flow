---
id: I-0079
title: TestRoundTripRepositoryItems reads the live main checkout and fails a close-out when another agent edits a story mid-run
class: defect
status: open
count: 1
cost: 6m
first_reported: 2026-10-05T05:55:58Z
last_reported: 2026-10-05T05:55:58Z
updated: 2026-10-06T09:56:52Z
---

# I-0079 TestRoundTripRepositoryItems reads the live main checkout and fails a close-out when another agent edits a story mid-run

## Description
TestRoundTripRepositoryItems reads the live main checkout and fails a close-out when another agent edits a story mid-run

## Instances

### 2026-10-05T05:55:58Z
Story: S-0276.
S-0276's close-out stopped at the full Go tests: TestRoundTripRepositoryItems in flai/internal/workitem read wip/kanban/stories/S-0251-*.md from the main checkout while a planner was writing it, and reported only an updated timestamp and a Tasks line. Two reruns passed. The close-out had to be finished step by step.

## Remediation

Story S-0290 remediates this issue, created from it at 2026-10-06T09:56:52Z.
