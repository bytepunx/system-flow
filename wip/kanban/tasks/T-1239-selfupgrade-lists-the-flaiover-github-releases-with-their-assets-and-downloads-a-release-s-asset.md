---
id: T-1239
type: task
nature: feature
title: selfupgrade lists the flaiover GitHub releases with their assets and downloads a release's asset
status: backlog
parent: S-0236
owner: alex
created: 2026-10-07T22:53:23Z
updated: 2026-10-07T22:53:23Z
transitions: []
stream: S-0236
tags: [cli]
touches: [flai/internal/selfupgrade/selfupgrade.go, flai/internal/selfupgrade/selfupgrade_test.go]
---
# T-1239 selfupgrade lists the flaiover GitHub releases with their assets and downloads a release's asset

## Work

In `flai/internal/selfupgrade/selfupgrade.go`, give the release client a way to list the GitHub releases of the configured releases repository whose tag is `flaiover/vX.Y.Z`, with no suffix, neither drafts nor prereleases, newest first by semver, across every page, with their assets. S-0234 creates these releases, each carrying `flaiover_<version>.digests` and its signature. Generalise what `List` does for `flai/vX.Y.Z` by tag prefix rather than copying it, and keep `List`, `ListTags`, and their callers as they are.

Add a function that fetches one such release by version, and one that downloads a named asset of a release into memory with the client's token and size limit, as `Download` fetches `checksums.txt`. Extend the GitHub stand-in in `selfupgrade_test.go` to answer `flaiover/v*` releases with assets.

Waits for no task. S-0233's signature work in the same package is a story `after`, so it is done before this starts.

## Done when

- The client lists published `flaiover/vX.Y.Z` releases with their assets, newest first, and fetches one by version, naming the version when there is none.
- A release's asset is downloaded by name, and a missing asset is an error naming the release and the asset.
- `List` and `ListTags` answer as before.
- `flai test flai/internal/selfupgrade` passes.

## Notes
