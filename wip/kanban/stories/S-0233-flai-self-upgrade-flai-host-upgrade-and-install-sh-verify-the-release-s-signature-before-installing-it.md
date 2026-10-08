---
id: S-0233
type: story
nature: feature
title: flai self-upgrade, flai host upgrade, and install.sh verify the release's signature before installing it
status: backlog
parent: E-0015
owner: arobson
created: 2026-10-02T12:37:23Z
updated: 2026-10-08T09:44:32Z
transitions: []
tags: [cli]
topics: [release, security]
touches: [flai/internal/selfupgrade, flai/cmd/selfupgrade.go, install.sh, docs/users/flai.md, docs/operators, ".github/workflows/system-flow-check.yml", flai/cmd/host.go, flai/cmd/selfupgrade_test.go, scripts/install-test.sh, design/system/flai-cli.md, flai/cmd/host_versions_test.go, flaiover/src/lib/components/HostProcesses.svelte, flaiover/src/lib/components/HostProcesses.svelte.test.ts, docs/users/flaiover.md, docs/users/flai-reference.md, docs/operators/runbooks/install.md, docs/operators/runbooks/update.md]
after: [S-0232]
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
    - kind: planner
      seconds: 35
      estimated: true
      models:
        - model: claude-haiku-4-5-20251001
          input: 10
          output: 4
          cache_read: 0
          cache_write: 9948
          cost: 0.002
        - model: claude-opus-5-5
          input: 30
          output: 1046
          cache_read: 2208472
          cache_write: 110135
          cost: 0.6065
    - kind: orchestrator
      seconds: 272
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 34
          output: 472
          cache_read: 2279538
          cache_write: 29521
          cost: 0.569
cost_of_delay:
  value: 3.29
  by: planner-E-0015
  at: 2026-10-07T22:13:40Z
forecast:
  duration: 1h15m
  delivery: 2026-10-08T11:29:00Z
  basis: "Its own forecast of 1h15m; 5th in the pull order with an in-progress limit of 5, behind S-0232, S-0297, S-0334, S-0344, S-0338, S-0342, S-0337 and S-0343."
  by: flai
  at: 2026-10-08T09:44:32Z
---
# S-0233 flai self-upgrade, flai host upgrade, and install.sh verify the release's signature before installing it

## Goal

A release is installed only when its `checksums.txt` was signed by a key the installer knows (ADR-0070). `flai self-upgrade` and `flai host upgrade` verify with the standard library and the public key in the binary; `install.sh` verifies with `openssl`. What cannot be verified is not installed, and the reason is said.

## Acceptance criteria
- [ ] `flai self-upgrade` and `flai host upgrade` download `checksums.txt.sig` beside `checksums.txt`, verify the ECDSA P-256 signature over the file with `crypto/ecdsa` and the embedded key before checking the archive's hash, and refuse with a message naming the release and the reason when the signature is missing or does not verify.
- [ ] A release signed by a key the binary does not know is refused with the version of the newest release the binary can verify, to install first.
- [ ] `install.sh` downloads the signature and verifies it with `openssl dgst -sha256 -verify` before the checksum; it refuses when `openssl` is absent or the signature fails, and says so.
- [ ] Tests cover a good signature, a bad one, a missing one, and an unknown key, with a test key pair under `testdata/`; `system-flow-check.yml`'s install and self-upgrade step passes against the latest release.
- [ ] `docs/users/flai.md` and the operator documentation say what is verified and what a refusal means.
- [ ] `flai self-upgrade --list` and `flai host versions`, with `--json`, mark each flai release whose signature this binary cannot verify, such as one from before signing, and the Updates page shows it as not installable and offers no way to install it.

## Tasks
- T-1226 selfupgrade verifies checksums.txt.sig with crypto/ecdsa and the embedded key before the archive's hash, and refuses what it cannot verify
- T-1227 install.sh verifies checksums.txt.sig with openssl and the embedded public key before the checksum, and refuses without openssl
- T-1228 flai self-upgrade --list and flai host versions mark each release this binary cannot verify, and flai host upgrade reports the refusal
- T-1229 The Updates page shows a flai release this binary cannot verify as not installable and offers no way to install it
- T-1230 The user, operator, and design documentation say what an upgrade and install.sh verify and what a refusal means

## Notes

The first release that carries this verification is itself installed by older flais without a signature check; from then on every upgrade is verified. Say so in the release notes.

### Planning

Touches:

- Declared, kept: `flai/internal/selfupgrade`, `flai/cmd/selfupgrade.go`, `install.sh`, `docs/users/flai.md`, `docs/operators`, `.github/workflows/system-flow-check.yml`.
- Layout: `flai/cmd/host.go`, whose `flai host upgrade` runs `flai self-upgrade` as a child, and whose `flai host versions` prints `self-upgrade --list`'s JSON through `printHostVersions`; `flai/cmd/selfupgrade_test.go`; `scripts/install-test.sh`, which `scripts/smoke.sh` runs for the CI step.
- Co-change and design: `design/system/flai-cli.md`, changed with these paths in 65% of their commits. Its Versions section says what is verified.
- Co-change: `docs/users/flai-reference.md`, at 24%, generated from the help text that T-1228 changes.
- Layout: `docs/operators/runbooks/install.md` and `update.md`, the runbooks the operator documentation criterion changes.
- Criterion 6, added on TH-0313: `flai/cmd/host_versions_test.go`; `flaiover/src/lib/components/HostProcesses.svelte` and its test, which list and install flai releases on the Updates page; `docs/users/flaiover.md`, which describes that page.
- Tasks name files below the folder touches, which narrow them in the claim (ADR-0096):
  - `flai/internal/selfupgrade/selfupgrade.go` and its test
  - a new `signature.go` and its test
  - a new `installsh_test.go`
  - the two runbooks
- Folder touches kept:
  - `flai/internal/selfupgrade` was declared. T-1226 also keeps `flai/internal/selfupgrade/testdata` as a folder, because the key pairs' file names are not known yet.
  - `docs/operators` was declared.
- Not added: `flai/internal/hostapi/writes.go` and `flai/internal/mcpserver/versions.go`. The host API's `host.upgrade` runs `flai host upgrade`, which refuses through self-upgrade. MCP's list of releases is not named by the criteria.

Forecast: 1h15m, kept. flai replays the delivery from the pull order whenever it changes.

- `flai forecast` now gives 36m, from a median of 106 s per unit of size over 32 large-band feature stories, times a size of 20.
- Done feature stories with 4 to 7 criteria took a median of about 1h of agent time. This one has five tasks in three layers, across Go, `install.sh`, Svelte, and docs, and two test key pairs, so 1h15m stands.

Cost of delay: 3.29 USD a week, as `flai cod` gives it, kept. It is this story's 1h15m share of the 9h30m forecast over E-0015's eight open stories, applied to the epic's 25 USD a week penalty, which the operator set on TH-0312. Each story closes part of one exposure, and the exposure closes only when the whole chain is done, so a share by work fits.
