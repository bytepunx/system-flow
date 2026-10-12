---
id: T-1241
type: task
nature: feature
title: flai/internal/dashboard resolves dashboard.tag to a flaiover release and its verified image digest, and refuses a release with no list it can verify
status: in-progress
parent: S-0236
owner: alex
created: 2026-10-07T22:53:42Z
updated: 2026-10-12T18:26:45Z
transitions:
  - to: ready
    at: 2026-10-12T18:26:44Z
    by: agent-S-0236
  - to: in-progress
    at: 2026-10-12T18:26:45Z
    by: agent-S-0236
stream: S-0236
tags: [cli]
touches: [flai/internal/dashboard/resolve.go, flai/internal/dashboard/resolve_test.go]
after: [T-1239, T-1240]
---
# T-1241 flai/internal/dashboard resolves dashboard.tag to a flaiover release and its verified image digest, and refuses a release with no list it can verify

## Work

In `flai/internal/dashboard/resolve.go`, resolve a configured image repository and `dashboard.tag` to the image to run. `latest` (or empty) is the newest published flaiover release and `X.Y.Z` (with or without a leading `v`) is that release, both through `selfupgrade.Find` with `DashboardTagPrefix` (T-1239). Fetch the release's `flaiover_<version>.digests` and its signature with `selfupgrade.FetchSigned`, parse it with T-1240's parser, check that the list's repository is the configured image repository, cache it with T-1240's cache, and answer the release version, the list, and the reference `<repository>@sha256:<index digest>`.

When the release does not exist, the tag is neither `latest` nor a version, the release has no digest list, the list does not verify or parse, or it names another repository, answer an error naming the release and the reason, which the commands print as why the container was not started; nothing is cached for it. When the release's assets cannot be fetched (a network error, not a refusal) and the cache holds a verified list for that version, use it and say so. Whether the configured repository is the signed one is the caller's: `--image` naming another repository and `--build` stay outside the resolver, and the commands run them as they are and report them unsigned. Take the release client's options and the cache directory as inputs, so tests pass a GitHub stand-in, the test key pair, and a temporary directory.

Waits for T-1239, whose client it calls, and T-1240, whose parser and cache it calls.

## Done when

- `latest` and a version each resolve to the release's version and a by-digest reference of the configured image.
- A missing release, a tag that is no version, a missing list, a list that does not verify, and a list for another repository each answer an error naming the release and the reason, and nothing is cached for them.
- A verified list is cached with its version, and a cached list is used when the release cannot be reached.
- `flai test flai/internal/dashboard` passes, with a GitHub stand-in and the test key pair.

## Notes
