---
id: T-1522
type: task
nature: improvement
title: selfupgrade keeps each release's signature verdict in a cache keyed by its asset URLs and the trusted keys
status: done
parent: S-0379
owner: alex
created: 2026-10-13T17:02:12Z
updated: 2026-10-13T17:09:02Z
transitions:
  - to: ready
    at: 2026-10-13T17:05:17Z
    by: agent-S-0379
  - to: in-progress
    at: 2026-10-13T17:05:18Z
    by: agent-S-0379
  - to: done
    at: 2026-10-13T17:09:02Z
    by: agent-S-0379
stream: S-0379
tags: [flai, cli, release]
touches: [flai/internal/selfupgrade/verdictcache.go, flai/internal/selfupgrade/verdictcache_test.go, flai/internal/selfupgrade/signature.go]
usage:
  source: log
  seconds: 224
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 34
      output: 147
      cache_read: 1131753
      cache_write: 72893
      cost: 0.5538
---
# T-1522 selfupgrade keeps each release's signature verdict in a cache keyed by its asset URLs and the trusted keys

## Work

Add a verdict cache to `flai/internal/selfupgrade`, in a new `verdictcache.go`. It waits for nothing: it is the base the listing builds on.

- One small JSON file under flai's `cache_dir` (`config.CacheDir`, default `~/.flai/cache`), such as `releases/verdicts.json`, written atomically (temp file and rename, as `flai/internal/verify/story.go` does with `atomicfile.WriteFile`).
- Each entry's key is the release's tag, the API URL of its `checksums.txt` asset, the API URL of its `checksums.txt.sig` asset, the API URL of its `cosign.pub` asset (empty when it publishes none), and a fingerprint of the keys the binary trusts (`TrustedKeys` in `signature.go`). A replaced asset gets a new URL, and a build that trusts another key gets a new fingerprint, so either one misses the cache.
- The value is the verdict: verified, or the refusal (`SignatureError`: unknown key, bad signature) with its kind and reason.
- Never cache a download error or a missing asset (a missing asset costs no download). Each of these is checked again on the next listing. A verdict that read `cosign.pub` is cached like any other, since that asset's URL is in the key.
- An unreadable or corrupt cache file is treated as empty, never as an error.
- Add the fingerprint function to `signature.go` if no exported one exists.

## Done when

- [ ] Unit tests in `verdictcache_test.go` show that a stored verdict is returned for the same tag, URLs, and keys, and is missed when either URL or the key fingerprint changes.
- [ ] Unit tests show that a download error is never stored, and that a corrupt file loads as empty.
- [ ] `flai test` passes on the changed paths.

## Notes

The asset struct decodes only `name` and `url` (`selfupgrade.go`, `ghAsset`). The API URL carries the asset's ID, so it changes when an asset is replaced. Widen `ghAsset` with `updated_at` only if the URL turns out not to.
