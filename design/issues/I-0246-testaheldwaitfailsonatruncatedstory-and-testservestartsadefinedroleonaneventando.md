---
id: I-0246
title: TestAHeldWaitFailsOnATruncatedStory and TestServeStartsADefinedRoleOnAnEventAndOnASchedule time out under flai test's Go run when the host's load is near 80
class: defect
status: open
count: 8
cost: 6m
first_reported: 2026-11-02T00:55:15Z
last_reported: 2026-11-07T02:07:15Z
updated: 2026-11-07T02:07:15Z
---

# I-0246 TestAHeldWaitFailsOnATruncatedStory and TestServeStartsADefinedRoleOnAnEventAndOnASchedule time out under flai test's Go run when the host's load is near 80

## Description
TestAHeldWaitFailsOnATruncatedStory and TestServeStartsADefinedRoleOnAnEventAndOnASchedule time out under flai test's Go run when the host's load is near 80

## Instances

### 2026-11-02T00:55:15Z
Story: S-0420.
S-0420's flai test runs on 2026-11-02: TestAHeldWaitFailsOnATruncatedStory (flai/internal/mcpserver/folder_test.go:270) timed out once in T-1763's sub-agent's run and once in mine, with the wait answering timed_out rather than the truncated read's error; TestServeStartsADefinedRoleOnAnEventAndOnASchedule (flai/internal/serve/serve_test.go:1009) timed out twice in waitFor. S-0420 changes neither package's code under test. The host's load average was 60 to 82 at the time; both pass on an earlier run.

### 2026-11-02T02:31:53Z
Story: S-0553.
S-0553's close-out stopped at the integration tier on TestServeStartsADefinedRoleOnAnEventAndOnASchedule (flai/internal/serve/serve_test.go:1009): never happened: the watcher's run for the move ended. The role lab runs the operator's command stub, so S-0553's change to the claude-code adapter is not on its path. TestAHeldWaitFailsOnATruncatedStory, the other test this issue names, failed the same way in S-0553's flai test earlier, recorded first as I-0244.

### 2026-11-02T03:49:08Z
Story: S-0422.
S-0422's close-out at 93b29b8 stopped at integration on TestServeStartsADefinedRoleOnAnEventAndOnASchedule ("never happened: the watcher's run for the move ended", serve_test.go:1009); a T-1777 sub-agent's flai test run had met the same failure once at about 03:12Z, before S-0422 touched serve.go, and passed on the next run.

### 2026-11-02T05:04:45Z
Story: S-0423.
S-0423's close-out at 8ac61fd stopped at integration on TestServeStartsADefinedRoleOnAnEventAndOnASchedule (flai/internal/serve/serve_test.go:1009: never happened: the watcher's run for the move ended). S-0423 does not change flai/internal/serve. The host was running three other stories' agents and close-outs at the time.

### 2026-11-02T06:01:43Z
Story: S-0553.
S-0553's close-out after its rebase onto main stopped at integration twice in a row, at e616012a: first on TestWaitForWorkPassesOverAHeldStory (flai/internal/mcpserver/work_test.go:266, the wait answered timed_out after the held story cleared), then on TestAHeldWaitFailsOnATruncatedStory (folder_test.go:270, timed_out instead of the truncated read's error). Each passes alone three times over. The host's load average was 93 on 24 cores, with three other stories' cmd.test and mcpserver.test runs going at the same time. S-0553 does not change flai/internal/mcpserver.

### 2026-11-02T06:36:55Z
Story: S-0424.
S-0424's `flai test` over the branch's changes failed only on TestServeStartsADefinedRoleOnAnEventAndOnASchedule (serve_test.go:1009, "never happened: the watcher's run for the move ended"), with S-0506 and S-0518 agents running on the host. T-1786's and T-2434's sub-agents also saw it fail once and pass on the rerun.

### 2026-11-02T19:52:44Z
Story: S-0425.
S-0425's close-out at e6040b6 stopped at integration on TestServeStartsADefinedRoleOnAnEventAndOnASchedule (serve_test.go:1009, "never happened: the watcher's run for the move ended"), with S-0510, S-0518, and S-0566 agents running on the host. The same test passed in S-0425's flai test run of flai/internal/serve at 19:20Z, after T-1791 changed serve/activity.go's charge.

### 2026-11-07T02:07:15Z
Story: S-0580.
S-0580's close-out at 27220c8 stopped at integration on TestServeStartsADefinedRoleOnAnEventAndOnASchedule (flai/internal/serve/serve_test.go:1009: never happened: the watcher's run for the move ended). S-0580 does not change flai/internal/serve. The host's load average was 53 over the last five minutes, with S-0578's agent working too.

## Remediation

Story S-0560 remediates this issue, created from it at 2026-11-02T03:55:02Z.
