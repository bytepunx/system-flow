---
id: S-0034
type: story
nature: feature
title: flai installer script and self-upgrade command
status: done
parent: E-0001
owner: alex
created: 2026-09-17T07:21:06Z
updated: 2026-09-17T07:45:39Z
transitions:
  - to: ready
    at: 2026-09-17T07:22:04Z
    by: alex
  - to: in-progress
    at: 2026-09-17T07:22:04Z
    by: alex
  - to: review
    at: 2026-09-17T07:27:40Z
    by: alex
  - to: done
    at: 2026-09-17T07:45:39Z
    by: alex
tags: [cli, docs]
---

# S-0034 flai installer script and self-upgrade command

## Goal
Installing and upgrading flai is one command: a script at the repository root that a README one-liner pipes into a shell, and `flai self-upgrade`, which does the same from an installed binary. Both find the latest `flai/v*` release, download the binary for the platform, verify it against `checksums.txt`, and install it, working against the private repository through a GitHub token or `gh`.

## Acceptance criteria
- [x] `install.sh` at the repository root detects OS and architecture, resolves the latest flai release (or `FLAI_VERSION`), downloads the archive and `checksums.txt`, verifies the SHA-256, installs to `FLAI_INSTALL_DIR` (default `/usr/local/bin`, sudo when needed), and warns when the directory is not on `PATH`
- [x] The script authenticates with `GITHUB_TOKEN`, `GH_TOKEN`, or `gh auth token` so it works while the repository is private, and without any of them once it is public
- [x] `flai self-upgrade` performs the same operation from the binary: `--check` reports current and latest, `--version` pins, `--dir` installs elsewhere; the running binary is replaced atomically; unit tests cover resolution, verification, and replacement
- [x] README highlights the one-liner near the top; docs/users/flai.md and design/system/flai-cli.md describe both paths
- [x] An install test runs the script against the real latest release in `make smoke` and in CI

## Tasks
- T-0109 install.sh at the repository root: platform detection, latest flai release, checksum verification, install dir, PATH advice, token or gh for private repos
- T-0110 flai self-upgrade: resolve release, download asset and checksums, verify, replace the running binary; --version, --check, --dir; tests
- T-0111 README highlights the one-liner; docs/users install section; flai-cli design; install test in smoke and CI

## Notes
- Modelled on bytepunx/kluster's install.sh (operator's pointer, 2026-09-17). Differences: monorepo tags (`flai/vX.Y.Z`), tar.gz archives, and a private repository for now.
