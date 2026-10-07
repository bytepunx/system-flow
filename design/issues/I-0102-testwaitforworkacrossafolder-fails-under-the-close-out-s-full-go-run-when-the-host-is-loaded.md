---
id: I-0102
title: TestWaitForWorkAcrossAFolder fails under the close-out's full Go run when the host is loaded
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-10-07T03:13:07Z
last_reported: 2026-10-07T03:13:07Z
updated: 2026-10-07T06:48:45Z
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
