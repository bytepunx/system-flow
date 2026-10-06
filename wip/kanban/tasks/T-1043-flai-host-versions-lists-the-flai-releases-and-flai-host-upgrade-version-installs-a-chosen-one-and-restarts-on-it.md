---
id: T-1043
type: task
nature: improvement
title: flai host versions lists the flai releases and flai host upgrade --version installs a chosen one and restarts on it
status: backlog
parent: S-0298
owner: alex
created: 2026-10-06T21:45:13Z
updated: 2026-10-06T21:45:13Z
transitions: []
stream: S-0298
tags: [cli]
touches: [flai/internal/host/api.go, flai/internal/host/client.go, flai/internal/host/host.go, flai/internal/host/host_test.go, flai/cmd/host.go, flai/cmd/host_integration_test.go]
after: [T-1040]
---
# T-1043 flai host versions lists the flai releases and flai host upgrade --version installs a chosen one and restarts on it

## Work

The host's `POST /upgrade` runs `flai self-upgrade` for the newest release and restarts itself and its processes on what it installed (`hostLauncher.upgrade` in `flai/cmd/host.go`, `host.upgrade` in `flai/internal/host/host.go`). Let it take an optional version in its body, passed as `self-upgrade --version`, so installing an older release restarts the host exactly as installing the newest does. Add `GET /versions`, which runs `self-upgrade --list --json` the way `GET /check` runs `--check`, and the client methods for both in `client.go`. Add `flai host versions` and `flai host upgrade --version <x>` in `flai/cmd/host.go`, saying what was installed, what it replaced, and that the host restarts on it. Waits for T-1040: the launcher runs the `--list` flag and the published-version refusal that task adds.

## Done when

- `flai host upgrade --version` to an older published release installs it and the host restarts on it with its processes, as `flai host upgrade` does for the newest, tested with a stand-in releases API.
- An unpublished version is refused before anything is replaced.
- `flai host versions` prints the list in text and JSON.
- `scripts/flai-test.sh` passes for `flai/internal/host` and `flai/cmd`.

## Notes
