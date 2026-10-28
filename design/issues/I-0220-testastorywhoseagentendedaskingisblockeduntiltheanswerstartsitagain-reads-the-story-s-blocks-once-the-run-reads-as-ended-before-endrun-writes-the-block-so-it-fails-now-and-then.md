---
id: I-0220
title: TestAStoryWhoseAgentEndedAskingIsBlockedUntilTheAnswerStartsItAgain reads the story's blocks once the run reads as ended, before endRun writes the block, so it fails now and then
class: defect
status: closed
count: 1
cost: 20m
first_reported: 2026-10-28T12:34:04Z
last_reported: 2026-10-28T12:34:04Z
updated: 2026-10-28T12:36:09Z
---

# I-0220 TestAStoryWhoseAgentEndedAskingIsBlockedUntilTheAnswerStartsItAgain reads the story's blocks once the run reads as ended, before endRun writes the block, so it fails now and then

## Description
TestAStoryWhoseAgentEndedAskingIsBlockedUntilTheAnswerStartsItAgain reads the story's blocks once the run reads as ended, before endRun writes the block, so it fails now and then

## Instances

### 2026-10-28T12:34:04Z
Story: S-0503.
Merge queue batch 20261028-122130 failed S-0503 at integration (TH-0679): flai/internal/serve/agents_test.go:3123 "the story idles in progress, blocked []". S-0503 changes no Go outside flai/internal/verify. endRun (flai/internal/serve/agents.go) records the end with l.ended(run) and only then calls blockAsked, by design, so the test's waitFor(!live) can return before the block is on disk and its single read of lab.blocks finds none. Fixed on S-0503's branch: the test waits for the block with waitFor instead of reading once.

## Remediation
Closed 2026-10-28T12:36:09Z: S-0503 (T-2197, commit 7545aa6): the test waits for the block with waitFor instead of reading the story's blocks once the run reads as ended, since endRun records the end before blockAsked writes the block by design.
