---
id: T-2367
type: task
nature: remediation
title: The cmd tests whose fixtures are per test run with t.Parallel, and the package globals that stop them are made per test or named
status: backlog
parent: S-0548
owner: alex
created: 2026-10-31T07:25:25Z
updated: 2026-10-31T07:26:18Z
transitions: []
stream: S-0548
tags: [cli]
touches: [flai/cmd/dashboard_test.go, flai/cmd/dashboard.go, flai/cmd/root.go, flai/cmd/forge.go, flai/cmd/components.go, flai/cmd/serve_mcp.go, flai/cmd/selfupgrade.go, flai/cmd/issue.go, flai/cmd/accept_remote.go]
after: [T-2365]
---
# T-2367 The cmd tests whose fixtures are per test run with t.Parallel, and the package globals that stop them are made per test or named

## Work

`flai/cmd` has 744 tests at about 0.6 s each and none calls `t.Parallel`, so the package's time is their sum, 450 s under `-race` in batch 20261030-190545, and a close-out that touches `cmd` runs go-test for 394 s at the median. Add `t.Parallel()` to every test, and subtest, whose fixture is its own (`t.TempDir`, its own `fakeRunner` and `app` through `runWith` and `runWithApp` in `dashboard_test.go`).

The package globals tests override stop that: `dashboardReleases` (`dashboard.go`, set by `setDashboardReleases`), `apps` (`root.go`), `githubAPIBase` (`forge.go`), `componentsRunner` and `componentsInstaller` (`components.go`), `exitWithProcessEvery` (`serve_mcp.go`), `releaseKeys` and `executablePath` (`selfupgrade.go`), and the vars in `issue.go` and `accept_remote.go`. Make each one a field of the `app` or its options so a test sets its own, or leave the tests that set it serial and name the global and the tests in the story's notes as the reason. `t.Setenv` and `t.Chdir` tests stay serial; Go refuses `t.Parallel` with them.

Leave `test.go` and `test_test.go` to T-2368, which changes them in the same layer; add `t.Parallel` there only once T-2368 is done. The task widens its touches to each test file it changes when it closes.

Second layer: it waits for T-2365, which changes `dashboard_test.go` and takes the 184 s dashboard test out of the package first, so this task's timing measures parallelism alone.

## Done when

`go test -race -count=1 ./cmd` passes three times in a row and finishes under two minutes on the host, its time and five slowest tests (`go test -json` elapsed) recorded in the narrative. Every test left serial is serial for a reason the story's notes name. Story criteria 3 and 4 for `cmd`.

## Notes

Drafted by the planner of S-0548 on 2026-10-31. It first drafted this task for `mcpserver`, which the story's agent's T-2361 covers, and rewrote it for `cmd`.
