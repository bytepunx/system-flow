---
id: I-0396
title: A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it
class: defect
status: open
count: 1
cost: 2h30m
first_reported: 2026-11-17T07:13:17Z
last_reported: 2026-11-17T07:13:17Z
updated: 2026-11-17T07:13:18Z
---

# I-0396 A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it

## Description
A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it

## Instances

### 2026-11-17T07:13:17Z
Job S-0817-20261117T042837Z started at 04:28:37Z on S-0817 in review. At 04:29:20Z the operator moved S-0817 to in-progress from the dashboard, with the note 'Rebase appears to fail from review', 43 seconds into the job, while its tiers ran. The job ran on for 33m27s, passed every step at c5bb5c4eb, and recorded outcome 'stayed in review' with the reason 'so S-0817 stays in review, ready to accept', and a verified event of the same. S-0817 was in in-progress. settle in flai/internal/serve/resolvejob.go sets OutcomeStayed on a passed verdict without reading the story; sendBack, the failing path, reads it again and says 'no longer in review'.

## Impact
S-0817 sat in in-progress with a passing close-out at its head and a record saying it was in review, from 05:02Z until the operator asked why at 06:54Z

- time_lost_per_cycle: 30m

## Remediation

When a rebase and resolve job passes, `settle` in `flai/internal/serve/resolvejob.go` reads the story again, as `sendBack` does. For a story no longer in review it records an outcome of its own, not `stayed in review`, and a reason that names the state the story is in, who moved it there and when, and the command that moves it back to review, since its close-out passed at its head. The review page shows that outcome as it shows the others. Decide in the story, with the operator, whether a move out of review while the job runs should also stop the job; record the decision in an ADR refining ADR-0313 if it changes what ADR-0313 says. A test in the resolve job's tests moves the story out of review while the job runs and reads the record.

Story S-0822 remediates this issue, created from it at 2026-11-17T07:13:18Z.
