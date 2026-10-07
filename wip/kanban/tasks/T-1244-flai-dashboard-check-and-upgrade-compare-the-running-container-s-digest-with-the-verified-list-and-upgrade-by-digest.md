---
id: T-1244
type: task
nature: feature
title: flai dashboard check and upgrade compare the running container's digest with the verified list and upgrade by digest
status: backlog
parent: S-0236
owner: alex
created: 2026-10-07T22:54:15Z
updated: 2026-10-07T22:54:15Z
transitions: []
stream: S-0236
tags: [cli]
touches: [flai/cmd/dashboard_upgrade.go, flai/cmd/dashboard.go, flai/cmd/dashboard_test.go]
after: [T-1243]
---
# T-1244 flai dashboard check and upgrade compare the running container's digest with the verified list and upgrade by digest

## Work

In `flai/cmd/dashboard_upgrade.go`, `check` and `upgrade` today compare the image ID of the pulled `image:tag` with the container's. Make them resolve `dashboard.tag`, or `--tag X.Y.Z`, with T-1241's resolver, compare the running container's digest with the platform and index digests of the verified list, and report an update when the container's digest is not on the list of the resolved release. `upgrade` pulls and swaps to the by-digest reference, with the same temporary container and `/_health` check as today. A release the resolver refuses is not pulled, and the reason is said.

`--image` naming another image and `--build` keep comparing and upgrading as today, reported unsigned.

Use the by-digest pull T-1243 put in `dashboard.go`; change that file only where `check` and `upgrade` need it.

Waits for T-1243, which changes the pull and run in `dashboard.go` and the tests in `dashboard_test.go` that this task also changes.

## Done when

- `flai dashboard check` reports up to date when the container's digest is on the resolved release's verified list, and an update otherwise.
- `flai dashboard upgrade`, with and without `--tag`, swaps to the by-digest reference, and refuses a release with no verifiable list before pulling.
- `--image` and `--build` behave as before and are reported unsigned.
- The cases are tested with `fakeRunner` and `flai test flai/cmd` passes.

## Notes
