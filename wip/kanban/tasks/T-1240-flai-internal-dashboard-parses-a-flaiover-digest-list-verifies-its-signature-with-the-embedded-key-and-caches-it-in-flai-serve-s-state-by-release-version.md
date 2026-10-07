---
id: T-1240
type: task
nature: feature
title: flai/internal/dashboard parses a flaiover digest list, verifies its signature with the embedded key, and caches it in flai serve's state by release version
status: backlog
parent: S-0236
owner: alex
created: 2026-10-07T22:53:30Z
updated: 2026-10-07T22:53:30Z
transitions: []
stream: S-0236
tags: [cli]
touches: [flai/internal/dashboard/digests.go, flai/internal/dashboard/digests_test.go, flai/internal/dashboard/cache.go, flai/internal/dashboard/cache_test.go]
---
# T-1240 flai/internal/dashboard parses a flaiover digest list, verifies its signature with the embedded key, and caches it in flai serve's state by release version

## Work

Create the package `flai/internal/dashboard`. In `digests.go`, parse a digest list in the format S-0234 recorded in `design/system/release-signing.md` (the index digest and each platform's manifest digest, each with the reference it was pushed under), and verify its signature with the verifier S-0233 added to `flai/internal/selfupgrade`, which takes the trusted keys through its options, so that tests pass the test key pair under `flai/internal/selfupgrade/testdata/`. A list that is malformed, unsigned, or signed by an unknown key is an error that says which.

In `cache.go`, keep each verified list, with its signature and its release version, beside `flai serve`'s state: under the directory `serve.DirFor` gives, in a file or folder named for the dashboard's digest lists. Write it atomically. Reading the cache verifies the list again, so that the dial-time check of S-0238 trusts nothing it did not verify. Give the cache a lookup by version and a lookup by digest, which S-0238 needs to ask whether a running container's digest is on any verified list.

Waits for no task: it needs only S-0233's verifier and S-0234's format, which the story's `after` waits for.

## Done when

- A good list parses and verifies; a malformed list, a bad signature, a missing signature, and an unknown key are each refused with a message that says which.
- A verified list is cached with its version, read back by version and by digest, and a cached list that no longer verifies is not trusted.
- `flai test flai/internal/dashboard` passes.

## Notes
