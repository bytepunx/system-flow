---
id: I-0309
title: flai/cmd's host and MCP-over-HTTP timing tests fail under a host load of about 200
class: impression
status: open
count: 2
cost: 13m
first_reported: 2026-11-08T00:27:14Z
last_reported: 2026-11-08T01:32:32Z
updated: 2026-11-08T01:32:32Z
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

## Remediation
