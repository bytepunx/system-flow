---
id: I-0246
title: TestAHeldWaitFailsOnATruncatedStory and TestServeStartsADefinedRoleOnAnEventAndOnASchedule time out under flai test's Go run when the host's load is near 80
class: defect
status: open
count: 1
cost: 6m
first_reported: 2026-11-02T00:55:15Z
last_reported: 2026-11-02T00:55:15Z
updated: 2026-11-02T00:55:15Z
---

# I-0246 TestAHeldWaitFailsOnATruncatedStory and TestServeStartsADefinedRoleOnAnEventAndOnASchedule time out under flai test's Go run when the host's load is near 80

## Description
TestAHeldWaitFailsOnATruncatedStory and TestServeStartsADefinedRoleOnAnEventAndOnASchedule time out under flai test's Go run when the host's load is near 80

## Instances

### 2026-11-02T00:55:15Z
Story: S-0420.
S-0420's flai test runs on 2026-11-02: TestAHeldWaitFailsOnATruncatedStory (flai/internal/mcpserver/folder_test.go:270) timed out once in T-1763's sub-agent's run and once in mine, with the wait answering timed_out rather than the truncated read's error; TestServeStartsADefinedRoleOnAnEventAndOnASchedule (flai/internal/serve/serve_test.go:1009) timed out twice in waitFor. S-0420 changes neither package's code under test. The host's load average was 60 to 82 at the time; both pass on an earlier run.

## Remediation
