---
id: I-0007
title: Review column exceeds its WIP limit while acceptance is batched
class: efficiency
status: closed
count: 1
cost: 5m
first_reported: 2026-09-15T22:42:51Z
last_reported: 2026-09-15T22:42:51Z
updated: 2026-10-03T03:26:21Z
---

# I-0007 Review column exceeds its WIP limit while acceptance is batched

## Description
The agent moves stories to review as they finish; the operator accepts them in batches. With a review limit of 3, the fourth finished story breaches the limit and `flai check --strict` fails, which blocks the commit-at-landing rule until someone accepts.

## Instances

### 2026-09-15T22:42:51Z
S-0026 reached review while S-0010, S-0022, and S-0023 were still awaiting acceptance. Committed with the warning logged.

## Remediation
Either raise the review limit in `wip/kanban/board.md` to match the operator's acceptance cadence, or accept stories before pulling the next one. Operator's call; see the S-0026 open questions.
Closed 2026-10-03T03:26:21Z: S-0243 (ADR-0073): a full review holds the pull, so flai serve, inbox, and wait_for_work start or offer no story while review is at or over its limit, and flai check --strict passes over review over its limit, so it no longer stops another story's close-out. Tests: TestReviewOverItsLimitDoesNotFailStrict, TestNoAgentStartsWhileReviewIsFull, TestWaitForWorkWaitsWhileReviewIsFull.
