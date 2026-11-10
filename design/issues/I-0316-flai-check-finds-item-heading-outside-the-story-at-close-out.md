---
id: I-0316
title: "flai check finds `item.heading` outside the story at close-out"
class: efficiency
status: closed
count: 1
first_reported: 2026-11-09T07:36:46Z
last_reported: 2026-11-09T07:36:46Z
updated: 2026-11-10T20:00:49Z
---

# I-0316 flai check finds `item.heading` outside the story at close-out

## Description
flai check finds `item.heading` outside the story at close-out

## Instances

### 2026-11-09T07:36:46Z
Story: S-0599.
flai check found outside the story:
`wip/kanban/stories/S-0636-13-board-get-calls-at-once-answer-under-300-ms-on-every-served-project-measured.md`: body is missing the "## Tasks" section

## Remediation

Story S-0679 remediates this issue, created from it at 2026-11-10T08:34:56Z.
Closed 2026-11-10T20:00:49Z: S-0679: the close-out's record step (recordOutside, behind flai check --story --record-issues and flai verify --record-issues) now confirms each finding outside the story on an active work item under the checked-write lock, running check.RunItemWrite on it again, and records only what still holds, so a checked creation refused and removed, as S-0636's first write was, leaves no issue (T-3014; TestCheckRecordIssuesConfirmsWorkItemFindingsUnderTheCheckedWriteLock reproduces it). A story's containers and scripts/agent-container.sh mount the main checkout's .flai-cache/items-checked.lock, so the confirmation in a container waits for a checked write on the host (T-3017). ADR-0283 records it (T-3018).
