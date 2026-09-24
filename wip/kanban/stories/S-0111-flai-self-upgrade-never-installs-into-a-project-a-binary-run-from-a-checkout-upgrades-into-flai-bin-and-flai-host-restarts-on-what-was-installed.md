---
id: S-0111
type: story
nature: remediation
title: "flai self-upgrade never installs into a project: a binary run from a checkout upgrades into ~/.flai/bin, and flai host restarts on what was installed"
status: backlog
owner: alex
created: 2026-09-24T05:09:05Z
updated: 2026-09-24T05:09:05Z
transitions: []
tags: [cli]
touches: [flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0111 flai self-upgrade never installs into a project: a binary run from a checkout upgrades into ~/.flai/bin, and flai host restarts on what was installed

## Goal

Upgrading flai never writes a release into a repository. The installed flai lives under the home folder, where `install.sh` puts it, whatever binary happened to run the upgrade.

## Acceptance criteria

- [ ] `flai self-upgrade` with no `--dir`, run by a flai that sits inside a system-flow project (a checkout's `bin/flai`), installs to `$FLAI_INSTALL_DIR` or `~/.flai/bin` instead of over itself, says where, and says to put that folder on PATH ahead of the checkout
- [ ] its up-to-date answer is about the flai in that folder, so a checkout's build never stops the home install from being made or upgraded
- [ ] `flai host upgrade` restarts the host on the binary that was installed, so the host and its children run the installed flai after an upgrade
- [ ] an installed flai outside any project (`~/.flai/bin`, `/usr/local/bin`) is still replaced in place, as before

## Tasks

## Notes

Found on the operator's machine on 2026-09-24. Their shell puts the repository's `bin` on PATH ahead of `/usr/local/bin`, so `flai self-upgrade` wrote 1.15.1 into `system-flow/bin/flai`. That is the file `scripts/flai.sh` rebuilds from the tree, so the next build replaces the installed release with a development build, and on other machines that folder does not exist.
