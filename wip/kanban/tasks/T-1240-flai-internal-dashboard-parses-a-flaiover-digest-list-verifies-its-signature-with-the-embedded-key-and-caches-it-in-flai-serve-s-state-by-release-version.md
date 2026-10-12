---
id: T-1240
type: task
nature: feature
title: flai/internal/dashboard parses a flaiover digest list, verifies its signature with the embedded key, and caches it in flai serve's state by release version
status: backlog
parent: S-0236
owner: alex
created: 2026-10-07T22:53:30Z
updated: 2026-10-12T18:17:19Z
transitions: []
stream: S-0236
tags: [cli]
touches: [flai/internal/dashboard/digests.go, flai/internal/dashboard/digests_test.go, flai/internal/dashboard/cache.go, flai/internal/dashboard/cache_test.go]
after: [T-1239]
---
# T-1240 flai/internal/dashboard parses a flaiover digest list, verifies its signature with the embedded key, and caches it in flai serve's state by release version

## Work

Create the package `flai/internal/dashboard`.

In `digests.go`, parse a digest list in the format `design/system/release-signing.md` § The digest list gives: lines of exactly three space-separated fields, a `sha256:` digest of 64 lower-case hex characters, the reference `<repository>:<version>`, and the kind (`index` first, then `<os>/<arch>[/<variant>]`); no blank lines, one final newline. Any other shape is an error naming the line. The parsed list answers its version, its repository, its index digest, every digest it lists, and whether a `RepoDigests` entry (`<repository>@sha256:…`) matches it, by splitting at `@` as that section says. Verify a list's signature with `selfupgrade.VerifyBlob` and `selfupgrade.TrustedKeys` from T-1239; `selfupgrade.FetchSigned` already answers `missing`, `bad`, and `unknown-key` when fetching.

In `cache.go`, keep each verified list with its signature, by release version, beside `flai serve`'s state: a folder `dashboard-digests` under the directory `serve.DirFor` gives (take the directory as a parameter, so `flai/cmd` passes `a.serveDir()`), one `<version>.digests` and one `<version>.digests.sig`, written with `atomicfile.WriteFile`. Reading the cache verifies each list again with the trusted keys, so the dial-time check of S-0238 trusts nothing it did not verify; a cached list that no longer verifies, or does not parse, is skipped and named. Give the cache a lookup by version and a lookup by any digest (which version's list names it), which S-0238 needs.

Waits for T-1239, whose exported verifier it calls.

## Done when

- A good list parses; each malformed shape is refused naming the line; a `RepoDigests` entry matches by repository and digest.
- A verified list is cached with its version, read back by version and by digest, and a cached list that no longer verifies is not trusted.
- `flai test flai/internal/dashboard` passes.

## Notes
