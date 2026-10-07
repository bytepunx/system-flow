---
id: S-0232
type: story
nature: feature
title: The release key signs flai's checksums.txt in CI and both components carry the public key
status: backlog
parent: E-0015
owner: arobson
created: 2026-10-02T12:37:23Z
updated: 2026-10-02T12:37:23Z
transitions: []
tags: [cli, dashboard]
topics: [release, security]
touches: [".github/workflows/release-flai.yml", flai/.goreleaser.yaml, flai/internal/buildinfo, flaiover/src/lib/server/release.ts, design/tech/ci.md, docs/operators]
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
