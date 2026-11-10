---
id: D-0444
title: "On this WSL2 host with mirrored networking a dial to a closed local port such as 127.0.0.1:1 hangs to the connect timeout, so a test that needs an unreachable URL dials a listener it opened on 127.0.0.1:0 and closed"
kind: cause
sources: [S-0548, T-2365, T-2362]
topics: [testing, go, code]
paths: [flai/cmd/dashboard_test.go, flai/internal/selfupgrade/selfupgrade.go, flai/internal/serve/mcp_test.go]
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

Sources: S-0548 (T-2365, T-2362).
