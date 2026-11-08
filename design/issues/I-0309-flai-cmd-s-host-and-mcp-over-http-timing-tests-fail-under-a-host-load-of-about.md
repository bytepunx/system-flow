---
id: I-0309
title: flai/cmd's host and MCP-over-HTTP timing tests fail under a host load of about 200
class: impression
status: open
count: 1
cost: 10m
first_reported: 2026-11-08T00:27:14Z
last_reported: 2026-11-08T00:27:14Z
updated: 2026-11-08T00:27:14Z
---

# I-0309 flai/cmd's host and MCP-over-HTTP timing tests fail under a host load of about 200

## Description
flai/cmd's host and MCP-over-HTTP timing tests fail under a host load of about 200

## Instances

### 2026-11-08T00:27:14Z
Story: S-0600.
S-0600's flai test on layer 1's paths (workitem/store.go, mdlint/folder.go, check/check.go) failed four flai/cmd tests the change does not reach: TestImportWithServeStoppedSaysToStartIt, TestImportWithAHostRegistersTheProject (two subtests), and TestMCPOverHTTPAnswersAWaitBeforeClaudeCodeGivesUp ("wait_for_work outlives the client's 3s"). The load average was 199 to 222 on 24 cores while the run took 9m19s under -race. They look like timing tests that time out under load, not defects of the change; not verified on main.

## Remediation
