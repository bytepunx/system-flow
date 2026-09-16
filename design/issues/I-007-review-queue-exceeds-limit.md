---
id: I-007
title: Review column exceeds its WIP limit while acceptance is batched
class: efficiency
status: open
count: 1
cost: 5m
first_reported: 2026-09-15T22:42:51Z
last_reported: 2026-09-15T22:42:51Z
updated: 2026-09-15T22:42:51Z
---

# I-007 Review column exceeds its WIP limit while acceptance is batched

## Description
The agent moves stories to review as they finish; the operator accepts them in batches. With a review limit of 3, the fourth finished story breaches the limit and `flai check --strict` fails, which blocks the commit-at-landing rule until someone accepts.

## Instances

### 2026-09-15T22:42:51Z
S-026 reached review while S-010, S-022, and S-023 were still awaiting acceptance. Committed with the warning logged.

## Remediation
Either raise the review limit in `wip/kanban/board.md` to match the operator's acceptance cadence, or accept stories before pulling the next one. Operator's call; see the S-026 open questions.
