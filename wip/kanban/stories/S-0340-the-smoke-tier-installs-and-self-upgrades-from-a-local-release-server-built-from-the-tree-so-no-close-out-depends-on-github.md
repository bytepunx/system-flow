---
id: S-0340
type: story
nature: remediation
title: The smoke tier installs and self-upgrades from a local release server built from the tree, so no close-out depends on GitHub
status: ready
owner: alex
created: 2026-10-08T07:59:11Z
updated: 2026-10-08T07:59:29Z
transitions:
  - to: ready
    at: 2026-10-08T07:59:29Z
    by: system-flow
tags: [cli]
topics: [ci, testing]
touches: [scripts/install-test.sh, scripts/smoke.sh, scripts/flai-snapshot.sh, scripts/release-server.sh, install.sh, flai/internal/selfupgrade/selfupgrade.go, flai/internal/selfupgrade/selfupgrade_test.go, ".github/workflows/system-flow-check.yml", design/tech/ci.md, design/system/devex.md, docs/users/flai.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 22
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 6
          output: 103
          cache_read: 925166
          cache_write: 2463
          cost: 0.2286
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: alex
    at: 2026-10-08T07:59:11Z
---
# S-0340 The smoke tier installs and self-upgrades from a local release server built from the tree, so no close-out depends on GitHub

## Goal

This story remediates [I-0086](../../../design/issues/I-0086-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md), "The close-out's install smoke test failed once and passed when run alone, with no cause in its output", now at seven occurrences, and takes over the remediation S-0291 was opened for.

The smoke tier ends with `scripts/install-test.sh`, which installs the latest published flai from GitHub twice through `install.sh` and twice more through `flai self-upgrade`. Every close-out therefore fetches the release listing (about 1 MB at `per_page=50`) and downloads release archives over the network. On the operator's host on 2026-10-08 that step took 285 s when it passed, 11 min 25 s inside S-0287's one close-out, and failed five close-outs of S-0326 in a row with `curl: (92) HTTP/2 stream 1 was not closed cleanly` while every other tier passed. A story's close-out must not depend on GitHub answering.

The behaviour under test is the installer's and self-upgrade's own: resolving a release, downloading its archive, checking it against `checksums.txt`, and landing the binary where it belongs. All of that can be exercised against a release built from this tree and served from this host. `install.sh` already takes `FLAI_API` and `FLAI_REPO`, and `flai self-upgrade` already takes an API base (`selfupgrade.Options.APIBase`), so the change is in how the smoke test is wired, not in what it checks.

## Acceptance criteria

- [ ] `scripts/install-test.sh` builds a release from the tree with `scripts/flai-snapshot.sh` and serves it from a local HTTP server that answers the three GitHub endpoints the installer and self-upgrade use (the release listing, a release by tag, and an asset download) with that build's archives and `checksums.txt`; `install.sh` and `flai self-upgrade` are pointed at it and at a stand-in repository name, and the script makes no request to `api.github.com` or `github.com`.
- [ ] The smoke tier passes with no GitHub token and with GitHub unreachable, shown by a run with `FLAI_API` pointed at the local server and network access to `api.github.com` refused (an unroutable `https_proxy` is enough), and `scripts/smoke.sh` and `system-flow-check.yml` no longer need `GITHUB_TOKEN` for it.
- [ ] Every check the script makes today is kept against the local server: an explicit `FLAI_INSTALL_DIR`, the default `HOME/.flai/bin` install without sudo and with the `PATH` line printed, `self-upgrade --check`, `self-upgrade --dir`, the resolution of the binary's own path outside a project, and the resolution to `HOME/.flai/bin` inside one; the installed binary reports the snapshot's version.
- [ ] The check that the latest published release installs from GitHub is kept as its own script, run by CI on main and on a schedule but by no story's close-out, and a failure of it opens or bumps an issue through `flai issue` rather than failing a story.
- [ ] `docs/users/flai.md`, `design/tech/ci.md`, and `design/system/devex.md` say what the smoke tier installs from, what the GitHub check covers, and how to run each by hand.
- [ ] I-0086 is closed with `flai issue close I-0086 --reason` saying what fixed it, and S-0291 is cancelled or closed as taken over, as the operator decides.

## Tasks

## Notes

Measured on 2026-10-08 on the operator's host, `squatchship`, WSL2 under mirrored networking. Alone, `scripts/install-test.sh` took 285 s and passed. In S-0287's verify record the smoke tier took 11 min 25 s against 3 min 6 s for the whole integration tier. In S-0326's five close-outs the listing fetch dropped each time; the agent measured 2 of 5 fetches of `releases?per_page=50` dropping while 3 of 3 at `per_page=5` came back. The GitHub fetches were the slow and the failing part on every run; nothing of the stories' own code was.

The local server needs no new dependency: a small Go program under `flai/cmd` or a script that serves a directory with the JSON the installer reads is enough. Whichever it is, it serves the same shape the real API answers, so that `install.sh` and `selfupgrade` run unchanged.
