---
id: S-0233
type: story
nature: feature
title: flai self-upgrade, flai host upgrade, and install.sh verify the release's signature before installing it
status: backlog
parent: E-0015
owner: arobson
created: 2026-10-02T12:37:23Z
updated: 2026-10-07T20:28:28Z
transitions: []
tags: [cli]
topics: [release, security]
touches: [flai/internal/selfupgrade, flai/cmd/selfupgrade.go, install.sh, docs/users/flai.md, docs/operators, ".github/workflows/system-flow-check.yml", flai/cmd/host.go, flai/cmd/selfupgrade_test.go, scripts/install-test.sh, design/system/flai-cli.md]
after: [S-0232]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
forecast:
  duration: 1h15m
  delivery: 2026-10-07T23:32:00Z
  basis: "Its own forecast of 1h15m; 2nd in the pull order with an in-progress limit of 3, behind S-0329, S-0330 and S-0232."
  by: flai
  at: 2026-10-07T20:28:28Z
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

## Tasks

## Notes

The first release that carries this verification is itself installed by older flais without a signature check; from then on every upgrade is verified. Say so in the release notes.

### Planning

Touches:

- Declared: `flai/internal/selfupgrade`, `flai/cmd/selfupgrade.go`, `install.sh`, `docs/users/flai.md`, `docs/operators`, `.github/workflows/system-flow-check.yml`.
- Layout: `flai/cmd/host.go`, where `flai host upgrade` lives; `flai/cmd/selfupgrade_test.go`; `scripts/install-test.sh`, which tests `install.sh`.
- Co-change and design: `design/system/flai-cli.md`, changed with these paths in 65% of their commits, whose self-upgrade section says what is verified.
- Folder touches kept, as declared: `flai/internal/selfupgrade`, where the verification may be a new file beside `selfupgrade.go`, and the test key pair goes under a new `testdata/`; `docs/operators`, whose install and update runbooks change.

Forecast: 1h15m, delivery 2026-10-08T10:58Z.

- `flai forecast` gave 22m from 116 s per unit of size, over only 3 medium-band feature stories.
- Done feature stories of this size took a median of about 1h of agent time. This one adds signature checks in Go and in `install.sh`, two commands, and four test cases, so 1h15m.
- Delivery is played out after S-0232 at flai's cycle factor of 6.85.

Cost of delay: no value yet. E-0015 and its stories have no inputs; TH-0312 asks the operator for them.
