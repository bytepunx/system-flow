---
id: T-1239
type: task
nature: feature
title: selfupgrade lists the flaiover GitHub releases with their assets and downloads a release's asset
status: in-progress
parent: S-0236
owner: alex
created: 2026-10-07T22:53:23Z
updated: 2026-10-12T18:17:51Z
transitions:
  - to: ready
    at: 2026-10-12T18:17:50Z
    by: agent-S-0236
  - to: in-progress
    at: 2026-10-12T18:17:51Z
    by: agent-S-0236
stream: S-0236
tags: [cli]
touches: [flai/internal/selfupgrade/selfupgrade.go, flai/internal/selfupgrade/selfupgrade_test.go]
---
# T-1239 selfupgrade lists the flaiover GitHub releases with their assets and downloads a release's asset

## Work

In `flai/internal/selfupgrade`, `List(ctx, opt, prefix)` already lists the published releases of a tag prefix with their assets; `DashboardTagPrefix` is `flaiover/v`. Add what the dashboard needs beside it, keeping `List`, `ListTags`, `Download`, `VerifyRelease`, and their callers as they are:

- `Find(ctx, opt, prefix, version string) (Release, error)`: the published release of that prefix and version (`latest` or empty is the newest), an error naming the version when there is none.
- `FetchAsset(ctx, opt, rel Release, name string) ([]byte, error)`: one named asset of a release, in memory, with the client's token and retries, as `Download` fetches `checksums.txt`; a missing asset is an error naming the release and the asset.
- `FetchSigned(ctx, opt, rel Release, name string) (data, sig []byte, err error)`: the asset and `<name>.sig`, verified with the trusted keys exactly as `verifiedSums` verifies `checksums.txt`, so the error is a `*SignatureError` of kind `missing`, `bad`, or `unknown-key` (telling an unknown key from a bad signature by the release's own `cosign.pub`, as now). Refactor `verifiedSums` onto it rather than copying it.
- `VerifyBlob(data, sig []byte, keys []*ecdsa.PublicKey) bool` and an exported `TrustedKeys(opt) ([]*ecdsa.PublicKey, error)`: the base64 cosign signature checked against the keys, for a cache that verifies again what it reads.

Extend the GitHub stand-in in `selfupgrade_test.go` to answer `flaiover/v*` releases with assets, and test each function with the test key pair.

Waits for no task; T-1240 waits for it.

## Done when

- `Find` answers the newest and a named `flaiover/vX.Y.Z` release with its assets, and names a version there is none of.
- `FetchAsset` downloads an asset by name; a missing one is an error naming the release and the asset.
- `FetchSigned` answers a good signature, and `missing`, `bad`, and `unknown-key` as `*SignatureError`; `VerifyBlob` checks bytes against keys.
- `List`, `ListTags`, `Download`, and `VerifyRelease` answer as before.
- `flai test flai/internal/selfupgrade` passes.

## Notes
