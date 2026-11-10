---
id: I-0329
title: TestStartDevServerSaysWhatAViteThatDoesNotAnswer200Answered and TestAnOperatorsReplyBecomesADigestEntry fail the integration tier under host load
class: defect
status: open
count: 1
cost: 11m
first_reported: 2026-11-10T05:34:34Z
last_reported: 2026-11-10T05:34:34Z
updated: 2026-11-10T05:34:34Z
---

# I-0329 TestStartDevServerSaysWhatAViteThatDoesNotAnswer200Answered and TestAnOperatorsReplyBecomesADigestEntry fail the integration tier under host load

## Description
TestStartDevServerSaysWhatAViteThatDoesNotAnswer200Answered and TestAnOperatorsReplyBecomesADigestEntry fail the integration tier under host load

## Instances

### 2026-11-10T05:34:34Z
Story: S-0634.
S-0634's close-out (job S-0634-20261110T052307Z, head a30e7bf) stopped at integration after 10m2s, with a load average of about 115 on 24 cores, on two tests in packages the story does not touch. flai/internal/components/devserver_test.go:417 TestStartDevServerSaysWhatAViteThatDoesNotAnswer200Answered: the fake Vite dev server did not answer within 600ms ("VITE ready, want the 404 it answered"). flai/internal/mcpserver/server_test.go:1662 TestAnOperatorsReplyBecomesADigestEntry: thread_reply took 783ms. Both look like fixed deadlines too short under load, the class I-0309 records for flai/cmd's host and MCP-over-HTTP tests.

## Remediation
