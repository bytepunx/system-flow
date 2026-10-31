---
id: T-2365
type: task
nature: remediation
title: The dashboard releases test reaches a listener it opened and closed, and no cmd test dials a port nothing listens on
status: in-progress
parent: S-0548
owner: alex
created: 2026-10-31T07:25:11Z
updated: 2026-10-31T07:32:15Z
transitions:
  - to: ready
    at: 2026-10-31T07:32:13Z
    by: agent-S-0548
  - to: in-progress
    at: 2026-10-31T07:32:15Z
    by: agent-S-0548
stream: S-0548
tags: [cli]
touches: [flai/cmd/dashboard_test.go, flai/cmd/dashboard_versions_test.go, flai/cmd/selfupgrade_test.go, flai/cmd/host_integration_test.go, flai/cmd/serve_mcp_test.go, flai/cmd/mcp_http_test.go, flai/internal/selfupgrade/selfupgrade.go, flai/internal/selfupgrade/selfupgrade_test.go]
---
# T-2365 The dashboard releases test reaches a listener it opened and closed, and no cmd test dials a port nothing listens on

## Work

`TestDashboardRunsTheCachedListWhenTheReleasesCannotBeReached` (`flai/cmd/dashboard_test.go`) points the release client at `http://127.0.0.1:1` through `setDashboardReleases`. On this WSL2 host with mirrored networking a closed local port hangs to the 60 s connect timeout, and `selfupgrade`'s request retries three times (`attempts` in `flai/internal/selfupgrade/selfupgrade.go`), so the test takes 184 s, 40 percent of `cmd`.

Give the test a URL that fails at once: open a `net.Listener` on `127.0.0.1:0`, take its address, close it, and point `dashboardReleasesAt` there; or give the `selfupgrade.Options` client a connect timeout under a second for the test. Keep `retryPause` at zero in tests.

Then audit the other tests that pass `127.0.0.1:1` as a URL: `--api http://127.0.0.1:1` in `flai/cmd/dashboard_versions_test.go`, the registry entries in `flai/cmd/selfupgrade_test.go` and `flai/cmd/host_integration_test.go`, and the one each in `flai/cmd/serve_mcp_test.go` and `flai/cmd/mcp_http_test.go`. Replace each one that is dialed with the same closed-listener address; leave an entry that is only stored and never dialed, and say so in the final message. The `127.0.0.1:1` entries in `flai/internal/serve` tests are T-2362's.

First layer: it waits for nothing. The `cmd` parallel task waits for it, since both change `flai/cmd/dashboard_test.go`.

## Done when

`go test -race -run TestDashboardRunsTheCachedListWhenTheReleasesCannotBeReached ./cmd` takes under two seconds on the host, and every `cmd` test that dials the release or dashboard API reaches a closed listener or a short timeout, never a port nothing listens on. Story criterion 1.

## Notes

Drafted by the planner of S-0548 on 2026-10-31.
