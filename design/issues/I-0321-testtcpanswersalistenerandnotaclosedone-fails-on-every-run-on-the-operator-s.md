---
id: I-0321
title: TestTCPAnswersAListenerAndNotAClosedOne fails on every run on the operator's WSL2 host, on main as well, because a closed 127.0.0.1 port still accepts a connection there
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-11-09T17:13:14Z
last_reported: 2026-11-09T17:13:14Z
updated: 2026-11-09T17:13:14Z
---

# I-0321 TestTCPAnswersAListenerAndNotAClosedOne fails on every run on the operator's WSL2 host, on main as well, because a closed 127.0.0.1 port still accepts a connection there

## Description
TestTCPAnswersAListenerAndNotAClosedOne fails on every run on the operator's WSL2 host, on main as well, because a closed 127.0.0.1 port still accepts a connection there

## Instances

### 2026-11-09T17:13:14Z
Story: S-0611.
S-0611's `flai test` on `flai/cmd/serve_actions.go` failed in the go-test tier on `flai/cmd/dashboard_porthold_test.go:145`, "a closed listener answers". The test is S-0615's (c97bd655e), and S-0611 changes nothing it tests. `go test -count=5 -run TestTCPAnswersAListenerAndNotAClosedOne ./cmd` failed 5 of 5 in S-0611's worktree, and `-count=3` failed 3 of 3 in the main checkout on 2026-11-09 at about 17:15 UTC. The test listens on 127.0.0.1:0, closes the listener, and expects `net.DialTimeout` to that port to fail, but on this WSL2 host the dial succeeds. A likely cause is WSL's localhost relay accepting on loopback ports. Every story whose `flai test` or close-out selects `./cmd` fails here until the test, or `tcpAnswers`, is fixed.

## Remediation
