---
id: D-0292
title: flai self-upgrade --list and flai host versions download two assets for every signed release to mark the ones they cannot verify
kind: cause
sources: [I-0142]
topics: [cli, release, security, flai, go, code, all]
paths: [flai/internal/selfupgrade/verdictcache.go, flai/internal/selfupgrade/verdictcache_test.go, flai/internal/selfupgrade/signature.go, flai/cmd/selfupgrade.go, flai/cmd/selfupgrade_test.go, design/system/flai-cli.md, design/system/release-signing.md, docs/users/flai.md, docs/operators/settings.md, docs/operators/runbooks/backup.md, design/issues/I-0142-flai-self-upgrade-list-and-flai-host-versions-download-two-assets-for-every-signed-release-to-mark-the-ones-they-cannot-verify.md, design/issues/summary.md, docs/users/flai-reference.md, design/issues/I-0136-an-acceptance-on-main-during-a-close-out-s-tiers-fails-its-final-sync-check-so-the-whole-close-out-runs-again.md]
date: 2026-10-13
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0292 flai self-upgrade --list and flai host versions download two assets for every signed release to mark the ones they cannot verify

I-0142 was closed on 2026-10-13.

The cause:

```text
flai self-upgrade --list and flai host versions download two assets for every signed release to mark the ones they cannot verify
```

The fix:

```text
S-0379: flai self-upgrade --list and --check, and so flai host versions, keep each signed release's signature verdict in <cache_dir>/releases/verdicts.json (selfupgrade.Verdicts), keyed by its tag, the API URLs of its checksums.txt, checksums.txt.sig, and cosign.pub, and the trusted keys' fingerprint, so a release is downloaded again only when an asset is replaced or the keys change. TestSelfUpgradeListKeepsEachVerdict shows a second listing downloads no signed release's assets. Installing still verifies from the network every time.
```

Source: I-0142.
