---
id: T-1295
type: task
nature: remediation
title: Wait for the container name to be released before the swap, and start the previous image again when the new one fails
status: done
parent: S-0316
owner: alex
created: 2026-10-07T23:55:15Z
updated: 2026-10-08T00:07:48Z
transitions:
  - to: ready
    at: 2026-10-08T00:00:16Z
    by: agent-S-0316
  - to: in-progress
    at: 2026-10-08T00:00:17Z
    by: agent-S-0316
  - to: done
    at: 2026-10-08T00:07:48Z
    by: agent-S-0316
stream: S-0316
tags: [flai, dashboard]
touches: [flai/cmd/dashboard_upgrade.go, flai/cmd/dashboard_test.go, docs/users/flai-reference.md]
usage:
  source: log
  seconds: 451
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 27
      output: 7456
      cache_read: 1613852
      cache_write: 71830
      cost: 0.9917
---
# T-1295 Wait for the container name to be released before the swap, and start the previous image again when the new one fails

## Work

Reproduce I-0094 first. Give `fakeRunner` in `flai/cmd/dashboard_test.go` a stand-in for a `--rm` container that Docker is still removing: after `docker stop <name>`, the name stays listed for a set number of `docker ps -a` or `docker container inspect` calls, and `docker run --name <name>` fails meanwhile with exit 125 and `Conflict. The container name /<name> is already in use`. Add a test in which `flai dashboard upgrade` swaps to a healthy new image while the name lingers, and see it fail as the issue describes.

Then fix it in `flai/cmd/dashboard_upgrade.go`:

- After `docker stop` of the real container, wait until Docker no longer lists the name, as a bounded count of tries like `upgradeProbeAttempts`, using `a.sleep` so that the test needs no clock. When the limit passes, `docker rm -f` the name and go on.
- When the new container still fails to start, start the previous image again under the real name and port, from the ref `containerInfo` read before the stop, and report the upgrade as failed with the dashboard on the previous image. Report it as down only when that start fails too.
- `runDashboardRestart` stops and starts under the same name in the same way. Use the same wait there.

Add tests for the rollback and for `flai dashboard restart` with a lingering name. When the upgrade command's `Long` text changes, run `make flai-reference` to regenerate `docs/users/flai-reference.md`.

This task waits for none.

## Done when

- The new test fails on the code as it was and passes after the fix
- `flai dashboard upgrade` and `flai dashboard restart` start the new container once a lingering name is released, in tests with the stand-in
- A failed start after the stop leaves the previous image running, and the upgrade exits non-zero saying so
- `flai test flai/cmd` passes

## Notes
