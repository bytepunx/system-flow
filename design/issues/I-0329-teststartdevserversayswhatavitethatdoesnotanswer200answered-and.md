---
id: I-0329
title: TestStartDevServerSaysWhatAViteThatDoesNotAnswer200Answered and TestAnOperatorsReplyBecomesADigestEntry fail the integration tier under host load
class: defect
status: closed
count: 2
cost: 15m
first_reported: 2026-11-10T05:34:34Z
last_reported: 2026-11-10T06:11:13Z
updated: 2026-11-10T08:10:13Z
---

# I-0329 TestStartDevServerSaysWhatAViteThatDoesNotAnswer200Answered and TestAnOperatorsReplyBecomesADigestEntry fail the integration tier under host load

## Description
TestStartDevServerSaysWhatAViteThatDoesNotAnswer200Answered and TestAnOperatorsReplyBecomesADigestEntry fail the integration tier under host load

## Instances

### 2026-11-10T05:34:34Z
Story: S-0634.
S-0634's close-out (job S-0634-20261110T052307Z, head a30e7bf) stopped at integration after 10m2s, with a load average of about 115 on 24 cores, on two tests in packages the story does not touch. flai/internal/components/devserver_test.go:417 TestStartDevServerSaysWhatAViteThatDoesNotAnswer200Answered: the fake Vite dev server did not answer within 600ms ("VITE ready, want the 404 it answered"). flai/internal/mcpserver/server_test.go:1662 TestAnOperatorsReplyBecomesADigestEntry: thread_reply took 783ms. Both look like fixed deadlines too short under load, the class I-0309 records for flai/cmd's host and MCP-over-HTTP tests.

### 2026-11-10T06:11:13Z
Story: S-0625.
S-0625's close-out (job S-0625-20261110T055311Z, head e3c14c6) stopped at integration after 9m15s on the same two tests, in packages the story does not touch: flai/internal/components/devserver_test.go:417 TestStartDevServerSaysWhatAViteThatDoesNotAnswer200Answered (the fake Vite did not answer within 600ms) and flai/internal/mcpserver/server_test.go:1662 TestAnOperatorsReplyBecomesADigestEntry (thread_reply took 637ms). Five other stories were in progress on the host. Asked for the close-out again.

## Remediation

Story S-0652 remediates this issue, created from it at 2026-11-10T05:44:45Z.
Closed 2026-11-10T08:10:13Z: Fixed by S-0652: both causes were in the tests. T-2841: TestStartDevServerSaysWhatAViteThatDoesNotAnswer200Answered allows 5s, not 600ms, for the fake vite to start and answer 404, and runs the fake with a 1s start delay, which reproduced the missing 404 on every run before the fix. T-2843: ownersServer sets Options.Slow to an hour, so a thread_reply slowed by load logs no slow request event; TestTheOwnersServerLogsNoSlowCall runs with FLAI_SLOW_REQUEST=1ns, which reproduced the log on every run before the fix.
