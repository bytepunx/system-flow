---
id: I-0244
title: TestAHeldWaitFailsOnATruncatedStory times out instead of failing on the truncated story when the host is loaded
class: defect
status: open
count: 1
cost: 5m
first_reported: 2026-11-02T00:22:39Z
last_reported: 2026-11-02T00:22:39Z
updated: 2026-11-02T00:22:39Z
---

# I-0244 TestAHeldWaitFailsOnATruncatedStory times out instead of failing on the truncated story when the host is loaded

## Description
TestAHeldWaitFailsOnATruncatedStory times out instead of failing on the truncated story when the host is loaded

## Instances

### 2026-11-02T00:22:39Z
Story: S-0553.
S-0553's flai test on T-2347's paths (harness, roles, serve, and their importers) failed only TestAHeldWaitFailsOnATruncatedStory in flai/internal/mcpserver/folder_test.go:270: wait_for_work answered timed_out true rather than failing on the truncated story file. The test sleeps 150ms, truncates the story, and expects the held wait's next poll to read the empty file; with four other stories' agents and tests on the host, the poll can miss the window. S-0553 changes nothing mcpserver's wait reads. I-0102 (closed) was the same race in TestWaitForWorkAcrossAFolder.

## Remediation
