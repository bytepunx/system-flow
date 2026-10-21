---
id: I-0158
title: TestASettledReplyStartsTheAgentWaitingOnItOnce in flai/internal/serve fails now and then in the close-out's integration tier on a branch that does not touch serve
class: defect
status: closed
count: 5
cost: 9m
first_reported: 2026-10-17T04:54:11Z
last_reported: 2026-10-21T09:24:20Z
updated: 2026-10-21T12:10:42Z
---

# I-0158 TestASettledReplyStartsTheAgentWaitingOnItOnce in flai/internal/serve fails now and then in the close-out's integration tier on a branch that does not touch serve

## Description
TestASettledReplyStartsTheAgentWaitingOnItOnce in flai/internal/serve fails now and then in the close-out's integration tier on a branch that does not touch serve

## Instances

### 2026-10-17T04:54:11Z
Story: S-0443.
S-0443's close-out stopped at integration: agents_test.go:1309 'started before the reply', the look started the stub agent again for a message on MS-0001 before the settled reply. S-0443 changes nothing under flai/internal/serve; the message is sent at the run's start second, so the test looks timing-dependent.

### 2026-10-17T18:09:01Z
Story: S-0396.
S-0396: failed again in the close-out's integration tier (agents_test.go:1347, "started before the reply"); S-0396 changes nothing under flai/internal/serve. Ran the close-out again.

### 2026-10-21T07:44:07Z
Story: S-0463.
S-0463's close-out stopped at the integration tier (6m19s): flai/internal/serve/agents_test.go:1348 reported "started before the reply", with flai serve having started stub-agent again for S-0001 on a message on MS-0001 before the settled reply was written. S-0463 changes only flai/internal/host/host_test.go and design/issues.

### 2026-10-21T08:49:08Z
Story: S-0456.
S-0456's close-out at 2026-10-21T08:43Z failed integration only in TestASettledReplyStartsTheAgentWaitingOnItOnce (flai/internal/serve/agents_test.go:1348): "started before the reply", the stub agent started again for S-0001 on MS-0001. S-0456 changes flai/internal/serve/gateway_test.go alone in that package, the LiteLLM gateway check, which the integration tier skips. The run before it, at the same change, passed integration.

### 2026-10-21T09:24:20Z
Story: S-0463.
Merge queue batch 20261021-091124 failed with S-0463 at integration (TH-0564): flai/internal/serve/agents_test.go:1348 'started before the reply', flai serve having started stub-agent again for S-0001 on a message on MS-0001. S-0463 changes only flai/internal/host/host_test.go and design/issues. flai serve sent S-0463 back to in-progress and out of the queue, so this flake now costs a story its place in the queue and a full close-out.

## Remediation

Story S-0473 remediates this issue, created from it at 2026-10-21T11:58:53Z.
Closed 2026-10-21T12:10:42Z: S-0473: the failure came from the test, not from flai serve. TestASettledReplyStartsTheAgentWaitingOnItOnce stamped the other story's run and the message from separate time.Now() calls, and messages.ToSince takes what is strictly after a run's start. So when the clock's second turned between the two putEnded calls, the message was new to the other story's agent and look started it, rightly by ADR-0154. The test now puts the waiting story's run, and so its message, a minute before the other run started. It also puts the other run a second earlier, which reproduced the failure on every run before the fix.
