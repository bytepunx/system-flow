---
title: ShellCheck
updated: 2026-10-26
status: active
topics: [code]
---

# ShellCheck

| | |
|-|-|
| Version | v0.10.0, pinned as `SHELLCHECK_VERSION` in `scripts/install-tools.sh` |
| Install | `make install-tools` downloads the release asset for the host (linux or darwin, x86_64 or aarch64) from GitHub, checks it against the sha256 pinned in the script, because the release publishes no checksums, and puts `shellcheck` in the main checkout's `bin/`, which `scripts/env.sh` puts on `PATH` in every story worktree. A `bin/shellcheck` that already reports the pinned version is kept. Another version needs `SHELLCHECK_SHA256` set to its asset's sha256 |
| Used in | `scripts/shellcheck.sh`, run by the `shellcheck` tier in `system-flow.yaml`, which selects `scripts/**/*.sh`, so `flai test`, `flai verify`, and the merge queue lint a story's shell scripts (S-0503, [I-0199](../issues/I-0199-no-test-tier-selects-scripts-so-a-story-that-changes-only-shell-scripts-and-the-script-tests-it-adds-pass-flai-verify-and-the-merge-queue-without-any-of-them-running.md)) |

## Why

The scripts under `scripts/` are POSIX `sh` under `set -eu` ([tooling.md](../conventions/tooling.md)), and until S-0503 no test tier read them: a story that changed only scripts passed `flai verify` and the merge queue without a check of any kind. ShellCheck is the standard static analyser for shell. It finds the faults those scripts are prone to, such as unquoted variables, bashisms under `sh`, and `cd` without a check, with no runtime and no configuration beyond its flags. It ships as one static binary per platform, so a pinned download gives every host that runs `make install-tools` the same version and the same findings.

## Considered

- The distro package (`apt install shellcheck`, Homebrew): the version is whatever the host's distribution carries, so findings would differ between hosts. `scripts/env.sh` still finds one on `PATH` after `bin/` when `bin/` has none.
- The Docker image `koalaman/shellcheck`: pinned, but starts a container per lint and needs Docker wherever the tier runs, including the merge queue.
- `checkbashisms`: finds only bashisms, which ShellCheck already reports under a `sh` shebang.
- `shfmt`: formats rather than analyses. It does not find the faults above; it can be added beside ShellCheck if formatting drift becomes a problem.
