---
id: T-1226
type: task
nature: feature
title: selfupgrade verifies checksums.txt.sig with crypto/ecdsa and the embedded key before the archive's hash, and refuses what it cannot verify
status: done
parent: S-0233
owner: alex
created: 2026-10-07T22:40:03Z
updated: 2026-10-12T17:04:02Z
transitions:
  - to: ready
    at: 2026-10-12T16:54:52Z
    by: agent-S-0233
  - to: in-progress
    at: 2026-10-12T16:54:53Z
    by: agent-S-0233
  - to: done
    at: 2026-10-12T17:04:02Z
    by: agent-S-0233
stream: S-0233
tags: [cli]
touches: [flai/internal/selfupgrade/selfupgrade.go, flai/internal/selfupgrade/selfupgrade_test.go, flai/internal/selfupgrade/signature.go, flai/internal/selfupgrade/signature_test.go, flai/internal/selfupgrade/testdata]
usage:
  source: log
  seconds: 549
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 50
      output: 210
      cache_read: 2268283
      cache_write: 119429
      cost: 1.0776
---
# T-1226 selfupgrade verifies checksums.txt.sig with crypto/ecdsa and the embedded key before the archive's hash, and refuses what it cannot verify

## Work

In `flai/internal/selfupgrade`, make `Download` fetch `checksums.txt.sig` beside `checksums.txt` and verify the ECDSA P-256 signature over the file with `crypto/ecdsa` and the release public key in `flai/internal/buildinfo/releasekey.go` (S-0232), before `Verify` checks the archive's hash. Put the verification in a new `signature.go`, with the keys it trusts given through `Options`, so that tests can pass a test key.

Refuse with an error that names the release and the reason in each case:

- `checksums.txt.sig` is not a release asset.
- The signature does not verify over `checksums.txt`.
- The release is signed by a key the binary does not know. Name the newest published release the binary can verify, to install first.

Generate a test key pair, and a second one that stands for an unknown key, under a new `testdata/`, and extend `fakeGitHub` to serve a signature asset. Cover a good signature, a bad one, a missing one, and an unknown key.

Waits for no task. It needs S-0232's key constant, which the story's `after` already waits for.

## Done when

- `Download` refuses a release whose signature is missing, bad, or made with an unknown key, and names the release and the reason.
- For an unknown key, the error names the newest release the binary can verify.
- A release signed with the trusted key installs as before.
- The four cases are tested with the key pairs under `flai/internal/selfupgrade/testdata/`, and `flai test` passes on the package.

## Notes
