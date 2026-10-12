---
id: T-1244
type: task
nature: feature
title: flai dashboard check and upgrade compare the running container's digest with the verified list and upgrade by digest
status: done
parent: S-0236
owner: alex
created: 2026-10-07T22:54:15Z
updated: 2026-10-12T19:16:29Z
transitions:
  - to: ready
    at: 2026-10-12T18:58:44Z
    by: agent-S-0236
  - to: in-progress
    at: 2026-10-12T18:58:45Z
    by: agent-S-0236
  - to: done
    at: 2026-10-12T19:16:29Z
    by: agent-S-0236
stream: S-0236
tags: [cli]
touches: [flai/cmd/dashboard_upgrade.go, flai/cmd/dashboard.go, flai/cmd/dashboard_test.go, flai/cmd/dashboard_versions_test.go, flai/cmd/dashboard_watch.go, flai/cmd/dashboard_watch_test.go]
after: [T-1243]
usage:
  source: log
  seconds: 1064
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 81
      output: 43402
      cache_read: 5898543
      cache_write: 251857
      cost: 3.4562
---
# T-1244 flai dashboard check and upgrade compare the running container's digest with the verified list and upgrade by digest

## Work

In `flai/cmd/dashboard_upgrade.go`, `check` and `upgrade` today compare the image ID of the pulled `image:tag` with the container's. For the signed repository (`signedDashboardImage`), make them resolve `dashboard.tag`, or `--tag X.Y.Z`, with `a.chooseDashboardImage` (T-1243, which runs T-1241's resolver), and compare the running container's `RepoDigests` with the resolved release's verified list (`List.Matches`): up to date when one matches, an update otherwise. `upgrade` pulls and swaps to the by-digest reference, with the same temporary container and `/_health` check as today, and the record keeps that reference. A release the resolver refuses is not pulled, and the reason is said. `--published` keeps refusing an unverifiable release before anything (T-1242).

`--image` naming another repository and `--build` keep comparing and upgrading as today, reported unsigned.

Also from T-1243: `runDashboardRestart` with no container running starts `s.ref()` by tag without resolving; make it use the record's by-digest reference or resolve. Update `TestDashboardRecordsWhatFlaiHostRestarts`, which expects an upgrade to record `…:latest`, and `TestDashboardUpgradePublishedDeploysAPublishedTag` in `dashboard_versions_test.go`, which expects the by-tag start and upgrade, failing since T-1243.

Waits for T-1243, which changes the pull and run in `dashboard.go` and the tests in `dashboard_test.go` that this task also changes.

## Done when

- `flai dashboard check` reports up to date when the container's digest is on the resolved release's verified list, and an update otherwise.
- `flai dashboard upgrade`, with and without `--tag`, swaps to the by-digest reference, and refuses a release with no verifiable list before pulling.
- `flai dashboard restart` with no container running runs a release by digest.
- `--image` and `--build` behave as before and are reported unsigned.
- The cases are tested with `fakeRunner` and `flai test flai/cmd` passes.

## Notes
