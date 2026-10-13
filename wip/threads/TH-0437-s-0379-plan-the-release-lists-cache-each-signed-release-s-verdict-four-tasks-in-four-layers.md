---
id: TH-0437
title: "S-0379 plan: the release lists cache each signed release's verdict, four tasks in four layers"
anchor:
  path: wip/kanban/stories/S-0379-flai-self-upgrade-list-and-flai-host-versions-download-two-assets-for-every-signed-release-to-mark-the-ones-they-cannot-verify.md
  item: S-0379
status: resolved
participants: [planner-S-0379, orchestrator]
created: 2026-10-13T17:03:31Z
updated: 2026-10-13T17:03:41Z
---

# TH-0437 S-0379 plan: the release lists cache each signed release's verdict, four tasks in four layers

On wip/kanban/stories/S-0379-flai-self-upgrade-list-and-flai-host-versions-download-two-assets-for-every-signed-release-to-mark-the-ones-they-cannot-verify.md.

## Entries

### 2026-10-13T17:03:31Z planner-S-0379
**Plan:** the release lists cache each signed release's signature verdict under flai's `cache_dir`, so a listing downloads only the releases it has no verdict for. Installing still verifies from the network every time.

## Tasks and layers

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1522 Verdict cache in `flai/internal/selfupgrade` | `verdictcache.go`, `verdictcache_test.go`, `signature.go` | none |
| 2 | T-1523 `--list` and the upgrade check read the cache; the test reproducing I-0142 | `flai/cmd/selfupgrade.go`, `selfupgrade_test.go` | T-1522 |
| 3 | T-1524 Design and guides | `flai-cli.md`, `release-signing.md`, `docs/users/flai.md`, `docs/operators/settings.md`, `runbooks/backup.md` | T-1523 |
| 4 | T-1525 Close I-0142 | I-0142's file, `design/issues/summary.md` | T-1523, T-1524 |

Each task waits on the one before, so nothing runs together. Every touch names a file.

## Figures

- **Forecast:** 40m. flai gave 19m. I raised it for four sequential layers, four regression tests, and five documents.
- **Cost of delay:** 12.50 USD a week, as `flai cod` gives it from flai's 5m input.

## Assumptions

- **The key** is the tag, the API URLs of `checksums.txt` and its `.sig`, and a fingerprint of the trusted keys. The API URL carries the asset's ID, so it changes when an asset is replaced. A build that trusts another key misses the cache.
- **Never cached:** download errors, missing assets, and verdicts that read `cosign.pub` from the network.
- **No new ADR.** ADR-0070 governs installs, and installs never read the cache. The lists only mark releases.
- **`flai host versions`** runs the host's own `self-upgrade --list`, so it gains the cache with no change of its own.
- **`flai dashboard versions`** verifies flaiover digest lists through its own cache, `dashboard-digests`. It is out of scope.

No split, merge, or drop proposed.

### 2026-10-13T17:03:40Z orchestrator
Approved. ADR-0070 § Decision binds `flai self-upgrade` and `flai host upgrade` to verify before they install. It says nothing about how the lists mark releases, and installs never read the cache, so no recorded decision changes and no ADR is needed. T-1523 reproduces I-0142 with the fake that counts asset requests (criterion 1), and T-1525 closes it (criterion 2). Every touch names a file, and the four layers follow each task's after. The figures are flai's: 40m (19m raised, with a reason) and 12.50 USD/week from flai's 5m input. Resolving.

Source: design/adrs/0070-releases-are-signed-with-a-cosign-key-pair-verified-before-they-are-installed.md § Decision

### 2026-10-13T17:03:41Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories.
