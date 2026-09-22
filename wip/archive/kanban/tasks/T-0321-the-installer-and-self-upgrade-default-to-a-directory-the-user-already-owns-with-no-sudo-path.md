---
id: T-0321
type: task
nature: feature
title: The installer and self-upgrade default to a directory the user already owns, with no sudo path
status: done
parent: S-0091
owner: alex
created: 2026-09-22T21:02:07Z
updated: 2026-09-22T21:06:19Z
transitions:
  - to: ready
    at: 2026-09-22T21:02:11Z
    by: system-flow
  - to: in-progress
    at: 2026-09-22T21:02:11Z
    by: system-flow
  - to: review
    at: 2026-09-22T21:06:19Z
    by: system-flow
  - to: done
    at: 2026-09-22T21:06:19Z
    by: system-flow
stream: S-0091
tags: []
---
# T-0321 The installer and self-upgrade default to a directory the user already owns, with no sudo path

## Work
`install.sh` currently defaults `INSTALL_DIR` to `/usr/local/bin` and falls back to `sudo install` when that is not writable, which almost every host requires. Change the default to `$HOME/.flai/bin`, matching the config and cache locations already under `~/.flai` (ADR-0008). Keep `FLAI_INSTALL_DIR` as the override. Remove the `sudo` fallback entirely: `mkdir -p` the directory and install into it; if that genuinely fails (a read-only home, an `FLAI_INSTALL_DIR` pointed somewhere unwritable), fail with a message naming the path and suggesting `FLAI_INSTALL_DIR`, not a `sudo` retry. Always print where it installed and a `PATH` line to add if `command -v flai` doesn't already resolve to it (today the PATH warning only fires when nothing shadows it; always show the export line so it's there to copy either way). `flai self-upgrade` already replaces whatever `os.Executable()` returns, so once installed under `~/.flai/bin` a bare `flai self-upgrade` needs no `--dir` and no sudo; confirm this in the smoke test rather than assuming it. Update `docs/users/flai.md`'s install section and `design/system/flai-cli.md`'s self-upgrade row to describe the new default and drop the sudo language.

## Done when
`scripts/install-test.sh` covers the default path (with `HOME` pointed at a scratch directory, not the real one) as well as the existing explicit-`FLAI_INSTALL_DIR` and `--dir` cases, and passes; `scripts/flai-test.sh` is clean; docs match the new behavior; no sudo invocation remains in `install.sh`.

## Notes
No ADR: this changes a default path and a script's behavior, not a structure, technology, or contract between parts (`decisions.md`). Recorded as a living-design/docs edit plus a narrative note.
