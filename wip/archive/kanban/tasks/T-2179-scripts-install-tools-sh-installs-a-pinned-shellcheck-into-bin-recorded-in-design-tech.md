---
id: T-2179
type: task
nature: improvement
title: scripts/install-tools.sh installs a pinned ShellCheck into bin/, recorded in design/tech
status: done
parent: S-0503
owner: alex
created: 2026-10-26T05:22:29Z
updated: 2026-10-26T05:31:11Z
transitions:
  - to: ready
    at: 2026-10-26T05:28:34Z
    by: agent-S-0503
  - to: in-progress
    at: 2026-10-26T05:28:35Z
    by: agent-S-0503
  - to: done
    at: 2026-10-26T05:31:11Z
    by: agent-S-0503
stream: S-0503
tags: [scripts, tooling]
touches: [scripts/install-tools.sh, design/tech/shellcheck.md, design/tech/README.md, scripts/README.md, Makefile]
usage:
  source: log
  seconds: 156
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 21
      output: 6643
      cache_read: 1210252
      cache_write: 36207
      cost: 0.6137
---
# T-2179 scripts/install-tools.sh installs a pinned ShellCheck into bin/, recorded in design/tech

## Work

- Extend `scripts/install-tools.sh` to install ShellCheck beside golangci-lint and GoReleaser: download the pinned release (`SHELLCHECK_VERSION`, default `v0.10.0`, the version on this host) for the host's OS and architecture from its GitHub release archive, check it against the checksum the release publishes, and place `shellcheck` in `$CACHE_ROOT/bin`, which `scripts/env.sh` puts on PATH in the main checkout and every story worktree. Skip the download when that version is already there. Print its version at the end as the other tools' versions are printed. ShellCheck is neither a Go module nor an npm package, so `go install` and `npm install` do not fit.
- Add `design/tech/shellcheck.md` with the version, where it is used (the `shellcheck` tier and `scripts/shellcheck.sh`), why, and what was considered instead, and its row in the table of `design/tech/README.md`.
- Update the `install-tools` row of `scripts/README.md` and the `install-tools` help line of the `Makefile` to name ShellCheck.

Waits for nothing: first layer.

## Done when

- On a tree without `bin/shellcheck`, `scripts/install-tools.sh` leaves an executable `bin/shellcheck` whose `--version` prints the pinned version, and a second run does not download it again.
- `design/tech/shellcheck.md` exists and `design/tech/README.md` names it.
- `flai check --strict` and the markdown lint are clean on the files changed.

## Notes

Drafted by planner-S-0503. The host has ShellCheck 0.10.0 at `/usr/bin/shellcheck` already; the pinned copy in `bin/` makes the tier reproducible on a host without it, as `scripts/install-tools.sh` does for golangci-lint (tooling convention: extend the install script rather than rely on what happens to be on PATH).
