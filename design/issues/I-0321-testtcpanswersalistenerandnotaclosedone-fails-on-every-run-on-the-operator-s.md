---
id: I-0321
title: TestTCPAnswersAListenerAndNotAClosedOne fails on every run on the operator's WSL2 host, on main as well, because a closed 127.0.0.1 port still accepts a connection there
class: defect
status: open
count: 2
cost: 6m
first_reported: 2026-11-09T17:13:14Z
last_reported: 2026-11-09T17:17:17Z
updated: 2026-11-09T21:50:55Z
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

## Remediation

Story S-0645 remediates this issue, created from it at 2026-11-09T21:50:55Z.
