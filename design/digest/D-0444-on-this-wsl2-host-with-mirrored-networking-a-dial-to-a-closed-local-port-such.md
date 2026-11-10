---
id: D-0444
title: "On this WSL2 host with mirrored networking a dial to a closed local port such as 127.0.0.1:1 hangs to the connect timeout, so a test that needs an unreachable URL dials a listener it opened on 127.0.0.1:0 and closed"
kind: cause
sources: [S-0548, T-2365, T-2362, S-0645, T-2816, I-0321]
topics: [testing, go, code]
paths: [flai/cmd/dashboard_test.go, flai/internal/selfupgrade/selfupgrade.go, flai/internal/serve/mcp_test.go, flai/cmd/dashboard_porthold.go, flai/cmd/dashboard_porthold_test.go]
date: 2026-10-31
status: active
absorbed_by: ""
asked_again: []
cited: [S-0645, T-2816, TH-0958]
---

# D-0444 On this WSL2 host with mirrored networking a dial to a closed local port such as 127.0.0.1:1 hangs to the connect timeout, so a test that needs an unreachable URL dials a listener it opened on 127.0.0.1:0 and closed

S-0548's T-2365 found why `TestDashboardRunsTheCachedListWhenTheReleasesCannotBeReached` (`flai/cmd/dashboard_test.go`) took 184 s, 40 percent of `cmd`'s time. The planner's brief says: "On this WSL2 host with mirrored networking a closed local port hangs to the 60 s connect timeout, and `selfupgrade`'s request retries three times."

The fix has two parts. The test dials the address of a `net.Listener` it opened on `127.0.0.1:0` and then closed, so the dial is refused at once. A closed listener alone still left the test at 4 s, "since `retryPause` is a package variable `cmd` cannot reach". So `selfupgrade.Options` gained `Attempts`, where zero means the default of three. The test went from 184 s to 0.10 s (wip/archive/agents/S-0548.md, Decisions).

The audit found the same in other tests. "Three `127.0.0.1:1` uses in serve's `mcp_test.go` and one in cmd's `dashboard_test.go` were dialed and now reach a closed listener; the seven others are stored, never dialed."

So a later test that needs a URL nothing answers uses the closed-listener address, never `127.0.0.1:1` or another port it assumes is free. A `127.0.0.1:1` that is only stored and never dialed can stay.

## A closed port can also answer (S-0645, I-0321)

Mirrored networking has a second effect: a port free in Linux can be one a Windows process listens on, so a real dial to a "closed" 127.0.0.1 port sometimes connects. `TestTCPAnswersAListenerAndNotAClosedOne` (`flai/cmd/dashboard_porthold_test.go`) failed intermittently on this host for that reason, main included (I-0321). From S-0645's Decisions (`wip/archive/agents/S-0645.md`):

```text
Cause confirmed in part: `.wslconfig` sets `networkingMode=mirrored`. 300 bind-close-connect probes on 127.0.0.1 found no closed port answering at 22:45 UTC, which matches I-0321's second instance (intermittent). Under mirrored networking a Windows listener is reachable on Linux loopback, so a port the kernel hands out as free in Linux can be one a Windows process listens on; whether the test fails depends on which ephemeral port it gets.
```

```text
Fix the test, not the probe: `tcpAnswersWith` and `publishAnswersWith` take a `dialFunc`; `tcpAnswers` and `publishAnswers` keep `net.DialTimeout`. Only a listener the test holds open is dialled for real. The closed case dials a fake that refuses; a second fake that accepts every address shows the relay's case and that the probe closes what it opened.
```

So the closed-listener address is enough when a test only needs a dial to fail fast, but a test that asserts a closed port does not answer injects the dialer and fakes the refusal, dialling for real only a listener it holds open. The production probe's exposure on a mirrored host (ADR-0267's port-held state) was left out of S-0645's scope.

Sources: S-0548 (T-2365, T-2362), S-0645 (T-2816), I-0321.
