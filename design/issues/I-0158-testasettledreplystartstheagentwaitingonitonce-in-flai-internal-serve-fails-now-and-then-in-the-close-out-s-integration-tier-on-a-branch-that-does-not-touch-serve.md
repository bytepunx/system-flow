---
id: I-0158
title: TestASettledReplyStartsTheAgentWaitingOnItOnce in flai/internal/serve fails now and then in the close-out's integration tier on a branch that does not touch serve
class: defect
status: open
count: 2
cost: 6m
first_reported: 2026-10-17T04:54:11Z
last_reported: 2026-10-17T18:09:01Z
updated: 2026-10-17T18:09:01Z
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

## Remediation
