---
id: I-0324
title: TestEachRequestIsTimedBeneathTheTransport reads the request answered event before the channel client has logged it, so it fails under load since S-0611
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-11-09T22:00:56Z
last_reported: 2026-11-09T22:00:56Z
updated: 2026-11-09T22:00:56Z
---

# I-0324 TestEachRequestIsTimedBeneathTheTransport reads the request answered event before the channel client has logged it, so it fails under load since S-0611

## Description
TestEachRequestIsTimedBeneathTheTransport reads the request answered event before the channel client has logged it, so it fails under load since S-0611

## Instances

### 2026-11-09T22:00:56Z
Story: S-0616.
S-0616's close-out at ae3f2d4 stopped at the integration tier on flai/internal/channel/channel_test.go:397, TestEachRequestIsTimedBeneathTheTransport: no request answered event in the logs, which held only the connection lines. S-0611 (2465e6c, accepted 2026-11-09T21:50) logs the request answered event after the reply's write, so that send_ms counts the write. The test reads the log as soon as the fake dashboard has the reply, before the client goroutine has logged it: a race that the host's load, with several close-outs and agents running, exposes. S-0616 changes only flai/cmd/new_test.go and docs, not flai/internal/channel. The test should wait for the event, for example by polling the log with a deadline, rather than read it once.

## Remediation
