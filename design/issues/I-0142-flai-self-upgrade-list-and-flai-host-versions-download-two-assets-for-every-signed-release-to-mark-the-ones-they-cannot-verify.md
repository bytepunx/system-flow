---
id: I-0142
title: flai self-upgrade --list and flai host versions download two assets for every signed release to mark the ones they cannot verify
class: efficiency
status: closed
count: 1
cost: 5m
first_reported: 2026-10-12T17:28:30Z
last_reported: 2026-10-12T17:28:30Z
updated: 2026-10-13T17:31:11Z
---

# I-0142 flai self-upgrade --list and flai host versions download two assets for every signed release to mark the ones they cannot verify

## Description
flai self-upgrade --list and flai host versions download two assets for every signed release to mark the ones they cannot verify

## Instances

### 2026-10-12T17:28:30Z
Story: S-0233.
S-0233 marks each listed flai release this flai cannot verify. A release that publishes checksums.txt.sig costs a download of checksums.txt and the signature on every listing; an unsigned one costs nothing. With two signed releases (1.42.0, 1.43.0) `flai self-upgrade --list` took 3.2 s on 2026-10-12. This project publishes several flai releases a day, so within weeks a listing makes hundreds of asset requests, 8 at a time. It slows the dashboard's Updates page, and without a token it reaches GitHub's unauthenticated limit of 60 requests an hour. Remedy: cache each release's result by its checksums.txt and signature asset URLs, which change when an asset is replaced, under flai's cache dir. Download still verifies every install.

## Remediation

Story S-0379 remediates this issue, created from it at 2026-10-12T17:44:10Z.
Closed 2026-10-13T17:31:11Z: S-0379: flai self-upgrade --list and --check, and so flai host versions, keep each signed release's signature verdict in <cache_dir>/releases/verdicts.json (selfupgrade.Verdicts), keyed by its tag, the API URLs of its checksums.txt, checksums.txt.sig, and cosign.pub, and the trusted keys' fingerprint, so a release is downloaded again only when an asset is replaced or the keys change. TestSelfUpgradeListKeepsEachVerdict shows a second listing downloads no signed release's assets. Installing still verifies from the network every time.
