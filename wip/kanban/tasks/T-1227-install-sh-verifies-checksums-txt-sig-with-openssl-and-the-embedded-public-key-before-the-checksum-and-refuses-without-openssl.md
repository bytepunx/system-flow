---
id: T-1227
type: task
nature: feature
title: install.sh verifies checksums.txt.sig with openssl and the embedded public key before the checksum, and refuses without openssl
status: backlog
parent: S-0233
owner: alex
created: 2026-10-07T22:40:21Z
updated: 2026-10-07T22:40:21Z
transitions: []
stream: S-0233
tags: [cli]
touches: [install.sh, scripts/install-test.sh, flai/internal/selfupgrade/installsh_test.go, ".github/workflows/system-flow-check.yml"]
---
# T-1227 install.sh verifies checksums.txt.sig with openssl and the embedded public key before the checksum, and refuses without openssl

## Work

In `install.sh`, carry the release public key as PEM, download `checksums.txt.sig` beside `checksums.txt`, and verify it with `openssl dgst -sha256 -verify <key> -signature <sig> checksums.txt` before the archive's checksum. Refuse, saying why, when `openssl` is not on `PATH`, when the release has no signature asset, and when the signature does not verify. Write the key to the script's temporary directory, which it already removes on exit.

Add a Go test, `flai/internal/selfupgrade/installsh_test.go`, that reads `install.sh` from the checkout and fails when its PEM differs from the constant in `flai/internal/buildinfo/releasekey.go`.

In `scripts/install-test.sh`, check that `install.sh` refuses with `openssl` left off `PATH`, and that the install against the latest release still passes. Change `.github/workflows/system-flow-check.yml` only if its runner lacks `openssl` or the step needs to say what it checks.

Waits for no task: it shares no path with T-1226, and needs only S-0232's key constant.

## Done when

- `install.sh` installs the latest signed release and says the signature was verified.
- It refuses, with its reason, when `openssl` is absent, the signature is missing, or it does not verify.
- A test fails when `install.sh`'s key and flai's differ.
- `scripts/install-test.sh` passes, as `system-flow-check.yml`'s smoke step runs it.

## Notes
