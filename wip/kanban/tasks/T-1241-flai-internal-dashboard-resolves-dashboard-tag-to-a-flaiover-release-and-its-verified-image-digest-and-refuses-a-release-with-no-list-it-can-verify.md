---
id: T-1241
type: task
nature: feature
title: flai/internal/dashboard resolves dashboard.tag to a flaiover release and its verified image digest, and refuses a release with no list it can verify
status: backlog
parent: S-0236
owner: alex
created: 2026-10-07T22:53:42Z
updated: 2026-10-07T22:53:42Z
transitions: []
stream: S-0236
tags: [cli]
touches: [flai/internal/dashboard/resolve.go, flai/internal/dashboard/resolve_test.go]
after: [T-1239, T-1240]
---
# T-1241 flai/internal/dashboard resolves dashboard.tag to a flaiover release and its verified image digest, and refuses a release with no list it can verify

## Work

In `flai/internal/dashboard/resolve.go`, resolve a configured `dashboard.tag` to the image to run. `latest` is the newest published flaiover release and `X.Y.Z` is that release, both through the release client T-1239 extended. Download the release's digest list and its signature, verify and cache them with T-1240's code, and answer the release version and the reference `ghcr.io/bytepunx/flaiover@sha256:<index digest>`, with the platform digests the list carries for comparing later.

When the release has no digest list, the list does not verify, or the release does not exist, answer an error that the commands print as the reason the container was not started. Use the cached list when the release's assets cannot be fetched and the cache holds the version, and say so. Leave `--image` naming another image and `--build` outside the resolver: the commands run those as they are and report them as unsigned.

Waits for T-1239, whose client it calls, and T-1240, whose verification and cache it calls.

## Done when

- `latest` and a version each resolve to the release's version and a by-digest reference of the configured image.
- A missing release, a missing list, and a list that does not verify each answer an error naming the release and the reason, and nothing is cached for them.
- A verified list is cached with its version.
- `flai test flai/internal/dashboard` passes, with a GitHub stand-in and the test key pair.

## Notes
