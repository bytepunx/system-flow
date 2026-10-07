---
id: S-0232
type: story
nature: feature
title: The release key signs flai's checksums.txt in CI and both components carry the public key
status: backlog
parent: E-0015
owner: arobson
created: 2026-10-02T12:37:23Z
updated: 2026-10-07T19:34:50Z
transitions: []
tags: [cli, dashboard]
topics: [release, security]
touches: [".github/workflows/release-flai.yml", flai/.goreleaser.yaml, flai/internal/buildinfo, flaiover/src/lib/server/release.ts, design/tech/ci.md, docs/operators, design/system/release-signing.md, flaiover/src/lib/server/release.test.ts]
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
      seconds: 4
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 0
          output: 1
          cache_read: 133363
          cache_write: 347
          cost: 0.0349
forecast:
  duration: 1h
  delivery: 2026-10-08T02:25:00Z
  basis: "flai's 116 s per unit of size rests on 3 medium-band stories and gave 24m; raised to 1h, the median agent time of 19 done feature stories of 4 to 7 criteria, for CI, two languages, and checking cosign and GoReleaser documentation; delivery at flai's cycle factor of 6.85, before the operator's key-generation step"
  by: planner-E-0015
  at: 2026-10-07T19:33:51Z
---
# S-0232 The release key signs flai's checksums.txt in CI and both components carry the public key

## Goal

The first step of ADR-0070: a cosign key pair exists, the two release workflows can use it, `release-flai.yml` signs `checksums.txt` with it, and the public key is a constant in flai's and flaiover's sources, so that every later story has a signature to check and a key to check it with. The finding is `design/system/release-signing.md`.

## Acceptance criteria
- [ ] The operator documentation has a runbook for generating the pair with `cosign generate-key-pair`, storing the encrypted key and its password as two GitHub Actions secrets, publishing the public key, and rotating or replacing a compromised key; the operator has run it once and the fingerprint of the public key is recorded there.
- [ ] `flai/.goreleaser.yaml` has a `signs` section that runs `cosign sign-blob` over the checksum file with the key from the environment and `--tlog-upload=false`, and uploads `checksums.txt.sig` as a release asset; `release-flai.yml` provides the secrets and installs cosign, and a release whose key is missing fails.
- [ ] `release-flai.yml` also runs `actions/attest-build-provenance` over the archives and the checksum file.
- [ ] The public key is a PEM constant in a flai package and in a flaiover server module, each with a test that parses it, and the two are identical.
- [ ] The claims marked **to check** in `release-signing.md § Signing releases` that bear on cosign and GoReleaser were checked against their current documentation and the document corrected where it was wrong.
- [ ] `design/tech/ci.md` lists cosign and the attestation action with their versions and why.

## Tasks

## Notes

### Planning

Touches:

- Declared: `.github/workflows/release-flai.yml`, `flai/.goreleaser.yaml`, `flai/internal/buildinfo`, `flaiover/src/lib/server/release.ts`, `design/tech/ci.md`, `docs/operators`.
- Design: `design/system/release-signing.md`, which criterion 5 corrects.
- Layout: `flaiover/src/lib/server/release.test.ts`, the test that parses the key on the flaiover side.
- Folder touches kept, as declared: `flai/internal/buildinfo`, where the key constant may go in a new file beside `buildinfo.go`; `docs/operators`, where the key runbook may be a new file under `runbooks/`. The story's agent narrows both when it writes the tasks.
- Co-change listed `design/system/flai-cli.md` and `docs/users/flai.md` at 60% and 56%. Not added: nothing a user runs changes here.

Forecast: 1h, delivery 2026-10-08T02:25Z.

- `flai forecast` gave 24m from 116 s per unit of size, over only 3 medium-band feature stories.
- 19 done feature stories with 4 to 7 criteria took a median of about 1h of agent time (S-0198 to S-0229, S-0298). This one spans CI, Go, and TypeScript, and checks cosign and GoReleaser documentation, so 1h.
- Delivery is flai's playout at its cycle factor of 6.85. It does not count the operator's key-generation step in criterion 1, which can delay it.

Cost of delay: no value yet. E-0015 and its stories have no inputs; TH-0312 asks the operator for them.
