---
id: T-1523
type: task
nature: improvement
title: flai self-upgrade --list and its check read the verdict cache, so a second listing downloads no signed release's assets
status: done
parent: S-0379
owner: alex
created: 2026-10-13T17:02:21Z
updated: 2026-10-13T17:24:09Z
transitions:
  - to: ready
    at: 2026-10-13T17:09:11Z
    by: agent-S-0379
  - to: in-progress
    at: 2026-10-13T17:09:11Z
    by: agent-S-0379
  - to: done
    at: 2026-10-13T17:24:09Z
    by: agent-S-0379
stream: S-0379
tags: [flai, cli, release]
touches: [flai/cmd/selfupgrade.go, flai/cmd/selfupgrade_test.go, docs/users/flai-reference.md]
after: [T-1522]
usage:
  source: log
  seconds: 895
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 62
      output: 239
      cache_read: 2855421
      cache_write: 105392
      cost: 1.361
---
# T-1523 flai self-upgrade --list and its check read the verdict cache, so a second listing downloads no signed release's assets

## Work

Wire the verdict cache from T-1522 into `flai/cmd/selfupgrade.go`. It waits for T-1522 because it calls the cache that task builds.

- `signatureRefusal` looks a release up in the cache first and calls `VerifyRelease` only on a miss; `signatureRefusals` stores each new verdict and saves the file once, after the eight-at-a-time pass ends.
- The cache directory comes from the config's `cache_dir`. Both callers, `listFlaiReleases` (`--list`) and `printCheck`, use it.
- `flai host versions` needs no change: the host runs its own `self-upgrade --list`, which now reads the host's cache.
- Installing a release still verifies its signature every time, from the network, and never reads the cache.

Reproduce I-0142 in `flai/cmd/selfupgrade_test.go`, with `serveFlaiReleases` and its `asked` counter.

- Run `--list` twice against the same signed releases. The second run asks for no `checksums.txt` or `checksums.txt.sig`, and marks the same releases.
- A release whose `.sig` URL changes is downloaded and verified again.
- A run whose signature download fails stores nothing, and the next run downloads again.
- An install after a cached listing still downloads and verifies the signature.

## Done when

- [ ] The tests above pass, and fail without the cache: the second listing's asset count is zero only with it.
- [ ] `TestSelfUpgradeListMarksUnverifiableReleases` and the existing listing tests pass unchanged, each with its own temporary cache directory.
- [ ] `flai test` passes on the changed paths.

## Notes

Point each test's cache at `t.TempDir()`, so that no test reads or writes `~/.flai/cache`.
