---
id: I-0102
title: TestWaitForWorkAcrossAFolder fails under the close-out's full Go run when the host is loaded
class: defect
status: closed
count: 1
cost: 10m
first_reported: 2026-10-07T03:13:07Z
last_reported: 2026-10-07T03:13:07Z
updated: 2026-10-07T23:39:52Z
---

# I-0102 TestWaitForWorkAcrossAFolder fails under the close-out's full Go run when the host is loaded

## Description
TestWaitForWorkAcrossAFolder fails under the close-out's full Go run when the host is loaded

## Instances

### 2026-10-07T03:13:07Z
Story: S-0269.
S-0269's close-out stopped at the full Go tests: folder_test.go:203 'handed over: map[]'. The mcpserver package took 178s there against 9s alone; the test sleeps 150ms and waits 3s. Three runs alone passed; the branch changes only the tool list the test checks.

## Remediation

Story S-0310 remediates this issue, created from it at 2026-10-07T06:48:45Z.
Closed 2026-10-07T23:39:52Z: S-0310. The failure was a tool error, not a timeout: readyStoryIn in flai/internal/mcpserver/folder_test.go added its criterion with os.WriteFile, which truncates the story file before filling it, and a wait_for_work held across it could poll in between, fail on the parse error, and answer an error the test dropped as map[]. readyStoryIn now writes through atomicfile.WriteFile, as flai writes items; every call in folder_test.go reports the tool's error; and TestAHeldWaitFailsOnATruncatedStory reproduces the truncated read.
