---
id: I-0158
title: TestASettledReplyStartsTheAgentWaitingOnItOnce in flai/internal/serve fails now and then in the close-out's integration tier on a branch that does not touch serve
class: defect
status: open
count: 3
cost: 6m
first_reported: 2026-10-17T04:54:11Z
last_reported: 2026-10-21T08:49:08Z
updated: 2026-10-21T08:49:08Z
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

### 2026-10-21T08:49:08Z
Story: S-0456.
S-0456's close-out at 2026-10-21T08:43Z failed integration only in TestASettledReplyStartsTheAgentWaitingOnItOnce (flai/internal/serve/agents_test.go:1348): "started before the reply", the stub agent started again for S-0001 on MS-0001. S-0456 changes flai/internal/serve/gateway_test.go alone in that package, the LiteLLM gateway check, which the integration tier skips. The run before it, at the same change, passed integration.

## Remediation
