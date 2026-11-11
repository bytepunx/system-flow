---
id: I-0267
title: "TestARepositoryIsOfferedImportedAndThenServed fails the close-out's integration tier with \"failed to read frame header: EOF\" under host load, on a branch that does not touch flai/internal/serve"
class: defect
status: closed
count: 1
cost: 9m
first_reported: 2026-11-07T02:16:58Z
last_reported: 2026-11-07T02:16:58Z
updated: 2026-11-11T23:19:34Z
---

# I-0267 TestARepositoryIsOfferedImportedAndThenServed fails the close-out's integration tier with "failed to read frame header: EOF" under host load, on a branch that does not touch flai/internal/serve

## Description
TestARepositoryIsOfferedImportedAndThenServed fails the close-out's integration tier with "failed to read frame header: EOF" under host load, on a branch that does not touch flai/internal/serve

## Instances

### 2026-11-07T02:16:58Z
Story: S-0580.
S-0580's second close-out, at 13346d1, stopped at integration on TestARepositoryIsOfferedImportedAndThenServed (flai/internal/serve/candidates_test.go:125: failed to get reader: failed to read frame header: EOF). The first close-out, at 27220c8, ran the same code (the two commits differ only in design/issues) and passed this test, stopping on I-0246's test instead. S-0580 does not change flai/internal/serve. The host's load average was 56 over five minutes, with S-0578's agent and close-out running too.

## Remediation

Story S-0703 remediates this issue, created from it at 2026-11-10T08:55:40Z.
Closed 2026-11-11T23:19:34Z: S-0703: the EOF came from the test dashboard, not flai serve. channeltest.Dashboard read a proven connection only inside Ask, Next, and Last, and coder/websocket answers a ping only while a read runs, so the client's 4 s ping went unanswered whenever a loaded host left a connection unread past its tick and timeout; the client then closed the connection, and the next Ask read 'failed to read frame header: EOF'. The dashboard now reads each connection for its whole life and hands its messages to Ask, Next, and Last (T-3044). TestACandidateLeftQuietPastItsPingsAnswersImportRun pings every 300 ms and leaves the candidate quiet for four pings before import.run: it failed 3 of 3 runs against the old dashboard with this error and passes with the fix (T-3042). The other path to the EOF, serve halting a candidate under its import.run answer, was closed earlier by S-0560 (6920a061c).
