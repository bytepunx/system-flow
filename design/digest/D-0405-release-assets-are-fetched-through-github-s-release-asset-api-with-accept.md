---
id: D-0405
title: Release assets are fetched through GitHub's release asset API with Accept application/octet-stream, because browser_download_url fails for a private repository
kind: decision
sources: [S-0034]
topics: [release, ci, host]
paths: [install.sh, flai/internal/selfupgrade/selfupgrade.go]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0405 Release assets are fetched through GitHub's release asset API with Accept application/octet-stream, because browser_download_url fails for a private repository

S-0034 chose the asset API because `browser_download_url` "would fail while private". The API form works with a token while the repository is private and without one once it is public. `install.sh` fetches the archive, `checksums.txt`, and the signature that way, and `flai self-upgrade` does the same. The token comes from `GITHUB_TOKEN`, `GH_TOKEN`, then `gh auth token`; none is written anywhere.

Source: S-0034.
