---
id: I-0321
title: TestTCPAnswersAListenerAndNotAClosedOne fails on every run on the operator's WSL2 host, on main as well, because a closed 127.0.0.1 port still accepts a connection there
class: defect
status: closed
count: 3
cost: 6m
first_reported: 2026-11-09T17:13:14Z
last_reported: 2026-11-09T22:11:53Z
updated: 2026-11-09T22:54:47Z
---

# I-0321 TestTCPAnswersAListenerAndNotAClosedOne fails on every run on the operator's WSL2 host, on main as well, because a closed 127.0.0.1 port still accepts a connection there

## Description
TestTCPAnswersAListenerAndNotAClosedOne fails on every run on the operator's WSL2 host, on main as well, because a closed 127.0.0.1 port still accepts a connection there

## Instances

### 2026-11-09T17:13:14Z
Story: S-0611.
S-0611's `flai test` on `flai/cmd/serve_actions.go` failed in the go-test tier on `flai/cmd/dashboard_porthold_test.go:145`, "a closed listener answers". The test is S-0615's (c97bd655e), and S-0611 changes nothing it tests. `go test -count=5 -run TestTCPAnswersAListenerAndNotAClosedOne ./cmd` failed 5 of 5 in S-0611's worktree, and `-count=3` failed 3 of 3 in the main checkout on 2026-11-09 at about 17:15 UTC. The test listens on 127.0.0.1:0, closes the listener, and expects `net.DialTimeout` to that port to fail, but on this WSL2 host the dial succeeds. A likely cause is WSL's localhost relay accepting on loopback ports. Every story whose `flai test` or close-out selects `./cmd` fails here until the test, or `tcpAnswers`, is fixed.

### 2026-11-09T17:17:17Z
Story: S-0611.
Correction to the first instance: it is intermittent, not constant. About ten minutes later, at about 17:25 UTC, S-0611's `flai test` ran `./cmd` again in the same worktree and TestTCPAnswersAListenerAndNotAClosedOne passed. So on this host a closed loopback port answers only some of the time, which is consistent with a relay that is sometimes listening.

### 2026-11-09T22:11:53Z
Story: S-0632.
flai test on S-0632's flai/internal/serve paths ran the go-test tier over ./cmd as an importer, and TestTCPAnswersAListenerAndNotAClosedOne (flai/cmd/dashboard_porthold_test.go:145) failed with "a closed listener answers"; S-0632 changes nothing in flai/cmd. The tier's 6m run was lost to it.

## Remediation

Story S-0645 remediates this issue, created from it at 2026-11-09T21:50:55Z.
Closed 2026-11-09T22:54:47Z: S-0645 (T-2816, 2b6821420). The host runs WSL2 with networkingMode=mirrored, where a port free in Linux can be a Windows listener's, so a real dial to a closed 127.0.0.1 port connects whenever the test draws such a port. tcpAnswers and publishAnswers gained tcpAnswersWith and publishAnswersWith, which take the dialer; production keeps net.DialTimeout. The test dials for real only a listener it holds open, and checks the refused case and the accepting relay's case through fake dialers, so it no longer depends on what answers on the host's loopback.
