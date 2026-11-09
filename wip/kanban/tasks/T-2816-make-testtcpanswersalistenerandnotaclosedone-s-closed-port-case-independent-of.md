---
id: T-2816
type: task
nature: remediation
title: Make TestTCPAnswersAListenerAndNotAClosedOne's closed-port case independent of what answers on the host's loopback
status: backlog
parent: S-0645
owner: agent
created: 2026-11-09T22:10:29Z
updated: 2026-11-09T22:10:29Z
transitions: []
stream: S-0645
tags: [flai, testing]
touches: [flai/cmd/dashboard_porthold.go, flai/cmd/dashboard_porthold_test.go]
---
# T-2816 Make TestTCPAnswersAListenerAndNotAClosedOne's closed-port case independent of what answers on the host's loopback

## Work

First, confirm the cause on the operator's host. Run `wslinfo --networking-mode`, or read the WSL configuration, to check for mirrored networking, which D-0444 records. Then run `go test -count=50 -run TestTCPAnswersAListenerAndNotAClosedOne ./cmd` in `flai/`. When the closed port answers, record what held it with `ss -ltnp` on the Linux side, and on the Windows side where it can be read. Under mirrored networking, a loopback port no Linux process listens on can be relayed to a Windows listener on the same port. So a port that was free a moment ago in Linux can answer, some of the time. Write what you found, or that the failure did not come back, in the narrative's `## Decisions`, citing D-0444.

Then fix the test, not the probe. `tcpAnswers` in `flai/cmd/dashboard_porthold.go` is a thin wrapper over `net.DialTimeout`. The real-socket closed-port assertion tests the host's loopback, not flai's code. Give `tcpAnswers` and `publishAnswers` a dial seam, such as a `dialTCP` function or an inner `tcpAnswersWith(dial, addr, timeout)`. The production callers keep `net.DialTimeout` and `portProbeTimeout`. Keep it safe under `t.Parallel()`: no package variable that a test swaps.

In `flai/cmd/dashboard_porthold_test.go`, keep the real-socket assertions that need no free port: a listener the test holds open answers through `tcpAnswers`, and through `publishAnswers` with a wildcard bind. A malformed publish does not answer. Assert the closed cases through a dialer that refuses, a `*net.OpError` wrapping `syscall.ECONNREFUSED`. Also assert, through a dialer that accepts every address the way the relay does, that the answer follows the dial, and that the publish's port is the one dialled. Never dial a port the test assumes is free.

This task waits for no other task: it is the story's first layer.

## Done when

- `go test -count=50 -run TestTCPAnswersAListenerAndNotAClosedOne ./cmd` passes in `flai/` on the operator's host, and `go test -race` passes on the package's porthold tests.
- No test in `flai/cmd/dashboard_porthold_test.go` asserts that a real port refuses a connection.
- `tcpAnswers` and `publishAnswers` still dial with `net.DialTimeout` and `portProbeTimeout` in production.
- The narrative's `## Decisions` records the cause the host showed, or that it did not come back, and why no test reproduces the host's relay itself.
- `flai test flai/cmd/dashboard_porthold.go flai/cmd/dashboard_porthold_test.go` passes.

## Notes
