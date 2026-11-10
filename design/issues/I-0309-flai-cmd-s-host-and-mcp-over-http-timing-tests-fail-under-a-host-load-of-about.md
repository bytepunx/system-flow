---
id: I-0309
title: flai/cmd's host and MCP-over-HTTP timing tests fail under a host load of about 200
class: impression
status: closed
count: 4
cost: 14m
first_reported: 2026-11-08T00:27:14Z
last_reported: 2026-11-09T22:38:42Z
updated: 2026-11-10T07:51:27Z
---

# I-0309 flai/cmd's host and MCP-over-HTTP timing tests fail under a host load of about 200

## Description
flai/cmd's host and MCP-over-HTTP timing tests fail under a host load of about 200

## Instances

### 2026-11-08T00:27:14Z
Story: S-0600.
S-0600's flai test on layer 1's paths (workitem/store.go, mdlint/folder.go, check/check.go) failed four flai/cmd tests the change does not reach: TestImportWithServeStoppedSaysToStartIt, TestImportWithAHostRegistersTheProject (two subtests), and TestMCPOverHTTPAnswersAWaitBeforeClaudeCodeGivesUp ("wait_for_work outlives the client's 3s"). The load average was 199 to 222 on 24 cores while the run took 9m19s under -race. They look like timing tests that time out under load, not defects of the change; not verified on main.

### 2026-11-08T01:32:32Z
Story: S-0600.
S-0600's close-out (job S-0600-20261108T011757Z, 14m1s, head 8c09c09) stopped at integration on TestImportWithAHostRegistersTheProject/import_--commit alone: "missing dashboard: served by the host flai". Cause found: the test's fakeRunning writes a host state.json stamped now, and host.Status counts the host gone once it is older than host.StaleAfter, 10 s (flai/internal/host/host.go). An import --commit that runs longer than 10 s under a load average of 87 to 116 therefore reads no host. A fix would refresh the fake state while the import runs, or stamp it ahead.

### 2026-11-09T08:25:21Z
Story: S-0606.
S-0606's close-out (job S-0606-20261109T081033Z, 13m27s, head 98d62f0) stopped at integration at a load average of about 100. It failed TestImportWithAHostRegistersTheProject/import_ ("missing dashboard: served by the host flai") and TestImportWithServeStoppedSaysToStartIt ("flai host is not running"). S-0606 does not reach that code. S-0606's T-2709 stamps fakeRunning's state an hour ahead, so it outlasts host.StaleAfter, as this issue's cause suggests. That fixes the two import tests. TestMCPOverHTTPAnswersAWaitBeforeClaudeCodeGivesUp, from the first instance, is not touched.

### 2026-11-09T22:38:42Z
Story: S-0634.
S-0634's close-out (job S-0634-20261109T222352Z, head 0473db7) stopped at integration after 12m39s, with a load average of about 160 on 24 cores, on two timing tests in code the story does not touch: flai/cmd/mcp_http_test.go:357 TestMCPOverHTTPAnswersAWaitBeforeClaudeCodeGivesUp (wait_for_work outlives the client's 3s: context deadline exceeded) and flai/internal/channel/channel_test.go:297 TestBacksOffWhenTheDashboardSaysTheProjectIsHeld (an ordinary close waited 252ms). The second is in flai/internal/channel, not flai/cmd.

## Remediation

Story S-0651 remediates this issue, created from it at 2026-11-10T05:44:44Z.
Closed 2026-11-10T07:51:27Z: S-0651. TestMCPOverHTTPAnswersAWaitBeforeClaudeCodeGivesUp gives its client a 12 s deadline against the server's 2.25 s scaled bound, in place of 3 s, so only a wait held for the 600 s asked fails (T-2839). TestBacksOffWhenTheDashboardSaysTheProjectIsHeld uses a 4 s held back-off, at least 2 s after jitter, and allows an ordinary close 1 s, in place of 300 ms against 250 ms (T-2842). The import tests were fixed by S-0606's T-2709. Pinned to one CPU beside 40 busy loops, main's channel test failed 5 of 10 runs and its MCP test 1 of 5; the story's versions passed 10 of 10 and 5 of 5, and both import tests passed too.
