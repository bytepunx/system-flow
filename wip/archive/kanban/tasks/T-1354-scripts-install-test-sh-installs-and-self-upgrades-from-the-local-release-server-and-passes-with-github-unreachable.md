---
id: T-1354
type: task
nature: remediation
title: scripts/install-test.sh installs and self-upgrades from the local release server and passes with GitHub unreachable
status: done
parent: S-0340
owner: alex
created: 2026-10-08T08:06:02Z
updated: 2026-10-08T08:31:59Z
transitions:
  - to: ready
    at: 2026-10-08T08:25:08Z
    by: agent-S-0340
  - to: in-progress
    at: 2026-10-08T08:25:08Z
    by: agent-S-0340
  - to: done
    at: 2026-10-08T08:31:59Z
    by: agent-S-0340
stream: S-0340
tags: [cli]
touches: [scripts/install-test.sh, scripts/smoke.sh, install.sh, ".github/workflows/system-flow-check.yml"]
after: [T-1348, T-1350]
usage:
  source: log
  seconds: 411
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 35
      output: 14273
      cache_read: 2083972
      cache_write: 62103
      cost: 1.0975
---
# T-1354 scripts/install-test.sh installs and self-upgrades from the local release server and passes with GitHub unreachable

## Work

Rewrite `scripts/install-test.sh` to build the release with T-1348's local mode of `scripts/flai-snapshot.sh` and serve it with T-1350's `scripts/release-server.sh`. Point every install at that server under a stand-in repository name: `FLAI_API` and `FLAI_REPO` for `install.sh`, and `FLAI_RELEASES_API` and `--repo` for `flai self-upgrade`, passed through each `env -i` as well.

- Keep every check the script makes today: no `sudo` in `install.sh`; an explicit `FLAI_INSTALL_DIR`; the default `HOME/.flai/bin` install, without sudo and with the `PATH` line printed; `self-upgrade --check`; `self-upgrade --dir`; the binary's own path outside a project; and `HOME/.flai/bin` inside one.
- Check that each installed binary reports the snapshot's version.
- Drop the token requirement and the `gh auth token` lookup.
- Print `install.sh`'s output when any install fails. I-0086's last instance lost it to `set -e`.
- Change `install.sh` only where the local server shows it must, keeping its defaults for users.
- In `.github/workflows/system-flow-check.yml`, drop `GITHUB_TOKEN` from the smoke step and install GoReleaser for the snapshot, such as `goreleaser/goreleaser-action` with `install-only`. Change `scripts/smoke.sh` only if its wording names GitHub.

`self-upgrade --dir` over a binary already at the latest version answers "already the latest release" and downloads nothing, as it does today. Serving a second, newer release would make that step download; that is optional.

Waits for T-1348 and T-1350: it runs the local build and the server.

## Done when

- `scripts/install-test.sh` passes with `GITHUB_TOKEN` and `GH_TOKEN` unset and with `https_proxy` and `HTTPS_PROXY` at an address that does not answer, so that any request to `api.github.com` or `github.com` fails. The run is recorded in the narrative.
- `scripts/smoke.sh` passes the same way.
- The script requests nothing from `api.github.com` or `github.com`.
- The smoke step of `system-flow-check.yml` is given no `GITHUB_TOKEN`.

## Notes

Drafted by the planner from S-0340's criteria 1 to 3. `flai self-upgrade` already reads its API base from `FLAI_RELEASES_API` or the hidden `--api` (S-0107), and `install.sh` from `FLAI_API`, so neither is expected to change.
