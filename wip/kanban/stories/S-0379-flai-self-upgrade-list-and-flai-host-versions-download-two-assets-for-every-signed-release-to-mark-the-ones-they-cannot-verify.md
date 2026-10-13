---
id: S-0379
type: story
nature: improvement
title: flai self-upgrade --list and flai host versions download two assets for every signed release to mark the ones they cannot verify
status: in-progress
owner: alex
created: 2026-10-12T17:44:10Z
updated: 2026-10-13T17:31:19Z
transitions:
  - to: ready
    at: 2026-10-13T17:03:48Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-13T17:03:57Z
    by: agent-S-0379
tags: [flai, cli, release]
topics: [cli, release, security]
touches: [flai/internal/selfupgrade/verdictcache.go, flai/internal/selfupgrade/verdictcache_test.go, flai/internal/selfupgrade/signature.go, flai/cmd/selfupgrade.go, flai/cmd/selfupgrade_test.go, design/system/flai-cli.md, design/system/release-signing.md, docs/users/flai.md, docs/operators/settings.md, docs/operators/runbooks/backup.md, design/issues/I-0142-flai-self-upgrade-list-and-flai-host-versions-download-two-assets-for-every-signed-release-to-mark-the-ones-they-cannot-verify.md, design/issues/summary.md, docs/users/flai-reference.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1640
  estimated: true
  turns:
    - day: 2026-10-13
      ceremony: 1
      hand_edits: 1
      work: 32
  models:
    - model: claude-opus-5-5
      input: 180
      output: 975
      cache_read: 8114247
      cache_write: 421059
      cost: 3.9254
  strategic:
    - kind: orchestrator
      seconds: 273
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 48
          output: 758
          cache_read: 4835192
          cache_write: 21521
          cost: 1.2245
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-10-12T17:44:10Z
  value: 12.5
  by: planner-S-0379
  at: 2026-10-13T17:03:20Z
forecast:
  duration: 40m
  delivery: 2026-10-13T17:49:00Z
  basis: "Its own forecast of 40m; 1st in the pull order with an in-progress limit of 5, behind S-0356, S-0375 and S-0377."
  by: flai
  at: 2026-10-13T17:03:51Z
finalized:
  by: orchestrator
  at: 2026-10-13T17:03:46Z
---
# S-0379 flai self-upgrade --list and flai host versions download two assets for every signed release to mark the ones they cannot verify

## Goal

This story remediates [I-0142](../../../design/issues/I-0142-flai-self-upgrade-list-and-flai-host-versions-download-two-assets-for-every-signed-release-to-mark-the-ones-they-cannot-verify.md), "flai self-upgrade --list and flai host versions download two assets for every signed release to mark the ones they cannot verify". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0142 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0142 is closed with `flai issue close I-0142 --reason` saying what fixed it

## Tasks
- T-1522 selfupgrade keeps each release's signature verdict in a cache keyed by its asset URLs and the trusted keys
- T-1523 flai self-upgrade --list and its check read the verdict cache, so a second listing downloads no signed release's assets
- T-1524 The design and the guides say the release lists reuse a cached verdict and installing never does
- T-1525 Close I-0142 with the reason that the release lists cache each signed release's verdict

## Notes

Cost of delay inputs set by flai from I-0142. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-12T17:28:30Z, 0 days before this story; under one cycle counts as one).

### Planning

Proposed remedy, from I-0142's instance: the release lists keep each signed release's verdict in a cache under flai's `cache_dir`. The key is the release's tag, the API URLs of its `checksums.txt` and `checksums.txt.sig` assets, and a fingerprint of the keys the binary trusts. A listing downloads only the releases it has no verdict for. Download errors, missing assets, and verdicts that read `cosign.pub` from the network are never cached. Installing verifies from the network every time and never reads the cache, so ADR-0070's guarantee is unchanged and no new ADR is planned. `flai host versions` runs the host's own `self-upgrade --list`, so it gains the cache with no change of its own.

Touches. The story declared none, and `flai touches suggest` had nothing to start from. Every touch names a file, so no folder touch is kept.

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/selfupgrade/verdictcache.go`, `verdictcache_test.go` | layout | New file for the cache, beside `signature.go` |
| `flai/internal/selfupgrade/signature.go` | layout | `VerifyRelease` and `TrustedKeys`; the key fingerprint |
| `flai/cmd/selfupgrade.go`, `selfupgrade_test.go` | layout | `signatureRefusals`, `signatureRefusal`, `listFlaiReleases`, `printCheck`; the `serveFlaiReleases` fake counts asset requests |
| `design/system/flai-cli.md` | design | § Verifying a release before it is installed, "The lists" |
| `design/system/release-signing.md` | design | § How a release is verified before it is installed |
| `docs/users/flai.md` | design | § What an upgrade verifies, and the `cache_dir` row |
| `docs/operators/settings.md`, `docs/operators/runbooks/backup.md` | layout | What `cache_dir` holds |
| I-0142's file, `design/issues/summary.md` | goal | Criterion 2 closes the issue |

Forecast. `flai forecast` gave 19m: 81 s per unit of size over 68 large improvement stories on claude-opus-5-5, times size 14. Raised to 40m, because the four tasks run one after another in four layers, the cache needs four regression tests, and five documents change. The delivery is flai's, moved later by the 21m added.

Cost of delay. The value is 12.50 USD a week, as `flai cod` gives it from flai's 5m input at 150 USD an hour. It stands: one instance so far, but the cost grows with every release, and an unauthenticated host reaches GitHub's 60-request limit within weeks.
