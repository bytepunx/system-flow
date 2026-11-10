---
id: I-0324
title: TestEachRequestIsTimedBeneathTheTransport reads the request answered event before the channel client has logged it, so it fails under load since S-0611
class: defect
status: open
count: 5
cost: 11m
first_reported: 2026-11-09T22:00:56Z
last_reported: 2026-11-10T05:35:23Z
updated: 2026-11-10T05:35:23Z
---

# I-0324 TestEachRequestIsTimedBeneathTheTransport reads the request answered event before the channel client has logged it, so it fails under load since S-0611

## Description
TestEachRequestIsTimedBeneathTheTransport reads the request answered event before the channel client has logged it, so it fails under load since S-0611

## Instances

### 2026-11-09T22:00:56Z
Story: S-0616.
S-0616's close-out at ae3f2d4 stopped at the integration tier on flai/internal/channel/channel_test.go:397, TestEachRequestIsTimedBeneathTheTransport: no request answered event in the logs, which held only the connection lines. S-0611 (2465e6c, accepted 2026-11-09T21:50) logs the request answered event after the reply's write, so that send_ms counts the write. The test reads the log as soon as the fake dashboard has the reply, before the client goroutine has logged it: a race that the host's load, with several close-outs and agents running, exposes. S-0616 changes only flai/cmd/new_test.go and docs, not flai/internal/channel. The test should wait for the event, for example by polling the log with a deadline, rather than read it once.

### 2026-11-09T22:42:47Z
Story: S-0622.
S-0622's second close-out, job S-0622-20261109T222615Z, stopped at integration after 13m53s at 0f0e231 with flai/internal/channel/channel_test.go:397 TestEachRequestIsTimedBeneathTheTransport: no request answered event in the log. S-0622 changes nothing under flai/internal/channel; only two test files under flai/cmd and flai/internal/serve. The host was loaded at the time with five stories in progress, their close-outs, planners, and the orchestrator. The close-out was asked for again.

### 2026-11-09T23:16:27Z
Story: S-0644.
S-0644's close-out (job S-0644-20261109T230952Z, at d88d346) stopped at integration on TestEachRequestIsTimedBeneathTheTransport (flai/internal/channel/channel_test.go:397, no request answered event), while other stories' close-outs ran beside it. S-0644 changes nothing under flai/internal/channel; the close-out was asked for again.

### 2026-11-10T05:24:22Z
Story: S-0633.
S-0633's close-out job S-0633-20261110T051519Z stopped at integration after 7m26s on TestEachRequestIsTimedBeneathTheTransport (channel_test.go:397, no request answered event in the log), at 524a1bea5. S-0633 changes nothing under flai/internal/channel. Asked for the close-out again.

### 2026-11-10T05:35:23Z
Story: S-0633.
S-0633's second close-out, job S-0633-20261110T052429Z, stopped at integration again after 10m29s at 141140d6f on TestEachRequestIsTimedBeneathTheTransport (channel_test.go:397, no request answered event), with six stories in progress on the host. S-0649, which remediates this issue, is still in backlog. Asked for the close-out a third time.

## Remediation

Story S-0649 remediates this issue, created from it at 2026-11-10T05:12:08Z.
