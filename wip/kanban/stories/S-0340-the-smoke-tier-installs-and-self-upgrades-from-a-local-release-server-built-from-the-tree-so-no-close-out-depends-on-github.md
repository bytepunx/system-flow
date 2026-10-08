---
id: S-0340
type: story
nature: remediation
title: The smoke tier installs and self-upgrades from a local release server built from the tree, so no close-out depends on GitHub
status: in-progress
owner: alex
created: 2026-10-08T07:59:11Z
updated: 2026-10-08T08:22:05Z
transitions:
  - to: ready
    at: 2026-10-08T07:59:29Z
    by: system-flow
  - to: in-progress
    at: 2026-10-08T08:14:33Z
    by: alex
tags: [cli]
topics: [ci, testing, release]
touches: [scripts/install-test.sh, scripts/smoke.sh, scripts/flai-snapshot.sh, scripts/release-server.sh, install.sh, ".github/workflows/system-flow-check.yml", design/tech/ci.md, design/system/devex.md, docs/users/flai.md, flai/.goreleaser.yaml, flai/internal/releaseserver/releaseserver.go, flai/internal/releaseserver/releaseserver_test.go, flai/internal/releaseserver/serve/main.go, scripts/install-published-test.sh, ".github/workflows/install-published.yml", Makefile, scripts/README.md, design/system/flai-cli.md, design/issues/I-0086-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1007
  estimated: true
  turns:
    - day: 2026-10-08
      ceremony: 5
      hand_edits: 1
      work: 32
  models:
    - model: claude-opus-5-5
      input: 266
      output: 1202
      cache_read: 11450762
      cache_write: 538997
      cost: 5.3898
  strategic:
    - kind: planner
      seconds: 343
      estimated: true
      models:
        - model: claude-haiku-4-5-20251001
          input: 324
          output: 4777
          cache_read: 1650124
          cache_write: 153320
          cost: 0.3869
        - model: claude-opus-5-5
          input: 162
          output: 28343
          cache_read: 9560039
          cache_write: 390113
          cost: 5.4075
    - kind: orchestrator
      seconds: 567
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 177
          output: 3000
          cache_read: 32800705
          cache_write: 42557
          cost: 8.0924
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5h30m
    by: alex
    at: 2026-10-08T08:22:05Z
  value: 800
  by: planner-S-0340
  at: 2026-10-08T08:07:07Z
forecast:
  duration: 1h30m
  delivery: 2026-10-08T11:20:00Z
  basis: "Its own forecast of 1h30m; 5th in the pull order with an in-progress limit of 3, behind S-0232, S-0321, S-0339, S-0342, S-0322, S-0291 and S-0341."
  by: flai
  at: 2026-10-08T08:08:17Z
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
- T-1348 scripts/flai-snapshot.sh builds a host-only release with a plain version for the local release server
- T-1350 A local release server answers the three GitHub release endpoints from a GoReleaser dist folder
- T-1352 The latest published release's install from GitHub is checked by its own script, run by CI on main and on a schedule
- T-1354 scripts/install-test.sh installs and self-upgrades from the local release server and passes with GitHub unreachable
- T-1356 flai.md, ci.md, devex.md, flai-cli.md, and scripts/README.md say what the smoke tier installs from and what the GitHub check covers
- T-1357 Close I-0086 saying the smoke tier installs from a local release server, and settle S-0291 as the operator decides

## Notes

Measured on 2026-10-08 on the operator's host, `squatchship`, WSL2 under mirrored networking. Alone, `scripts/install-test.sh` took 285 s and passed. In S-0287's verify record the smoke tier took 11 min 25 s against 3 min 6 s for the whole integration tier. In S-0326's five close-outs the listing fetch dropped each time; the agent measured 2 of 5 fetches of `releases?per_page=50` dropping while 3 of 3 at `per_page=5` came back. The GitHub fetches were the slow and the failing part on every run; nothing of the stories' own code was.

The local server needs no new dependency: a small Go program under `flai/cmd` or a script that serves a directory with the JSON the installer reads is enough. Whichever it is, it serves the same shape the real API answers, so that `install.sh` and `selfupgrade` run unchanged.

### Planning

Two findings in the code shape the plan:

- `selfupgrade.List` keeps only releases tagged with a plain `X.Y.Z` (`buildinfo.Bare`), and `install.sh` names the archive from the tag. A GoReleaser snapshot is versioned `0.0.0-<commit>` and archived as `flai_snapshot_<os>_<arch>`, so neither would accept it. T-1348 gives `scripts/flai-snapshot.sh` a local mode with a plain version, which touches `flai/.goreleaser.yaml`.
- The smoke and integration tiers are selected by `flai/**` alone. A story that changed only `scripts/` would not run smoke in its own close-out. The release server is planned as Go under `flai/internal/releaseserver` (T-1350), so the close-out runs smoke and shows criteria 2 and 3.

The layers are T-1348, T-1350, and T-1352 first, T-1354 second, then T-1356 and T-1357. The plan's thread on S-0340 gives the reasons.

Touches: `flai touches suggest S-0340` found 255 of 1434 commits changing the declared paths. No folder touch is kept.

| Touch | Source |
|-------|--------|
| The eleven declared paths | declared, all kept |
| `flai/internal/selfupgrade/selfupgrade.go`, `selfupgrade_test.go` | declared; no task names them, since `FLAI_RELEASES_API` and the hidden `--api` already set the API base (S-0107) |
| `flai/.goreleaser.yaml` | layout: the snapshot's version and archive name (T-1348) |
| `flai/internal/releaseserver/releaseserver.go`, `releaseserver_test.go`, `serve/main.go` | design: the story's note on a small Go program; layout: an internal package with its own `main` (T-1350) |
| `scripts/install-published-test.sh`, `.github/workflows/install-published.yml` | criterion 4: the GitHub check as its own script and workflow (T-1352) |
| `Makefile` | layout: every script has a Make target, and `install-test` has one (T-1352) |
| `scripts/README.md` | layout: it indexes every script (T-1356) |
| `design/system/flai-cli.md` | co-change, 187 of 255 commits; its `flai self-upgrade` row says the smoke test runs against the real latest release (T-1356) |
| `design/issues/I-0086-…md`, `design/issues/summary.md` | criterion 6 (T-1357) |

Co-changed paths left out: `docs/users/flai-reference.md` (73 commits) is generated from the commands, and no command changes. `docs/operators/index.md`, `docs/operators/settings.md`, and the flaiover documents name no part of the smoke test.

S-0232, in progress, holds this story: it touches `scripts/flai-snapshot.sh`, `flai/.goreleaser.yaml`, and `design/tech/ci.md`. It waits on the operator's release key (TH-0324). `design/system/flai-cli.md`, `docs/users/flai.md`, and `design/issues` are shared paths and hold nothing.

Forecast: `flai forecast` gave 21m at size 17 (6 criteria, 11 touches), at 74 s per unit over 19 done remediation stories. It is raised to 1h30m:

- The story has six tasks and a new Go package with tests.
- The GoReleaser mode needs trial builds.
- Criteria 2 and 3 need smoke runs with GitHub cut off, several minutes each.
- Its close-out runs integration and smoke.

The delivery is flai's start, about 08:51Z, plus 1h30m: 10:30Z. It moves later for as long as S-0232 holds the story.

Cost of delay: `flai cod` gave 37.50 USD a week from the operator's input of 15m per 168h cycle. The value is set to 800 USD a week:

- The step costs every close-out that runs smoke between 285 s and 11 min 25 s. A local server should take about a minute. That saves at least 3 minutes a run.
- 136 stories were accepted in the 7 days to 2026-10-08, by `git log`. Taking half as changing `flai/` gives about 70 smoke runs a week, or about 3.5h.
- I-0086's 7 failures in 2 days come to about 24 a week. Each costs a rerun of integration and smoke and its diagnosis, about 2h a week, as S-0291's planning found.
- That is about 5.5h a week at 150 USD an hour, or 825 USD, rounded down to 800.

The input is the operator's and is left as it is. The plan's thread recommends raising it to about 5h30m.
