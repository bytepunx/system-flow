---
id: I-0186
title: TestWaitForEventsAcrossAFolder can wake on the new story's creation before its move to ready, and finds no event
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-10-21T13:12:46Z
last_reported: 2026-10-21T13:12:46Z
updated: 2026-10-21T13:12:46Z
---

# I-0186 TestWaitForEventsAcrossAFolder can wake on the new story's creation before its move to ready, and finds no event

## Description
TestWaitForEventsAcrossAFolder can wake on the new story's creation before its move to ready, and finds no event

## Instances

### 2026-10-21T13:12:46Z
Story: S-0479.
S-0479's third close-out failed in the integration tier at flai/internal/mcpserver/folder_test.go:344. The wait answered changed [beta/wip/kanban/stories/S-0001-beta-work.md] with no events. readyStoryIn creates the story, writes its criterion, and moves it to ready, all while wait_for_events polls, so a poll between the creation and the move wakes the wait on a backlog story, which is no event. The test passed in the two runs before on the same branch. Repo.Save's read under the items lock widens the gap slightly. S-0479 now writes the story before the wait and moves it during the wait.

## Remediation
