---
id: T-1464
type: task
nature: feature
title: The local release server serves checksums.txt.sig and cosign.pub, and the smoke test verifies its test-signed releases with a test build of flai
status: in-progress
parent: S-0233
owner: alex
created: 2026-10-12T17:04:34Z
updated: 2026-10-12T17:17:12Z
transitions:
  - to: ready
    at: 2026-10-12T17:17:12Z
    by: agent-S-0233
  - to: in-progress
    at: 2026-10-12T17:17:12Z
    by: agent-S-0233
stream: S-0233
tags: []
touches: [flai/internal/releaseserver/releaseserver.go, flai/internal/releaseserver/releaseserver_test.go, flai/internal/buildinfo/releasekey.go, flai/internal/buildinfo/releasekey_test.go, flai/internal/selfupgrade/signature.go, scripts/install-test.sh, flai/cmd/installsh_test.go]
after: [T-1228]
---
# T-1464 The local release server serves checksums.txt.sig and cosign.pub, and the smoke test verifies its test-signed releases with a test build of flai

## Work

T-1226 and T-1227 left the smoke tier failing: `scripts/install-test.sh` signs its local releases with a test key, but `flai/internal/releaseserver` serves only the archives and `checksums.txt`, and `bin/flai` trusts only the embedded release key, so its self-upgrade steps refuse the local releases.

- `flai/internal/releaseserver` serves `checksums.txt.sig` and `cosign.pub` when the dist folder has them.
- `flai/internal/buildinfo/releasekey.go` takes one more trusted key at link time only, through `-ldflags -X`, which no release build sets. `selfupgrade`'s trusted keys add it. The smoke test builds a flai with its test key this way and upgrades with it.
- `flai/cmd/installsh_test.go`'s stand-in serves a test signature and runs a copy of `install.sh` carrying the test key.

Waits for T-1228: both change `flai/internal/selfupgrade/signature.go`.

## Done when

- `scripts/install-test.sh` passes, with each install and upgrade verifying a test-signed release.
- A release build trusts only the embedded key, and a test says so.
- `flai test` passes on the paths changed.

## Notes
