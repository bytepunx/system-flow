---
id: S-0232
type: story
nature: feature
title: The release key signs flai's checksums.txt in CI and both components carry the public key
status: in-progress
parent: E-0015
owner: arobson
created: 2026-10-02T12:37:23Z
updated: 2026-10-07T22:28:05Z
transitions:
  - to: ready
    at: 2026-10-07T22:14:01Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-07T22:14:07Z
    by: agent-S-0232
tags: [cli, dashboard]
topics: [release, security]
touches: [".github/workflows/release-flai.yml", flai/.goreleaser.yaml, flaiover/src/lib/server/release.ts, design/tech/ci.md, design/system/release-signing.md, flaiover/src/lib/server/release.test.ts, flai/internal/buildinfo/releasekey.go, flai/internal/buildinfo/releasekey_test.go, docs/operators/runbooks/release-key.md, docs/operators/runbooks/README.md, docs/operators/settings.md, scripts/flai-snapshot.sh]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 882
  turns:
    - day: 2026-10-07
      ceremony: 4
      hand_edits: 3
      work: 34
  models:
    - model: claude-haiku-4-5-20251001
      input: 210085
      output: 8754
      cache_read: 0
      cache_write: 0
      cost: 0.2639
    - model: claude-opus-5-5
      input: 274
      output: 95486
      cache_read: 13408615
      cache_write: 462641
      cost: 7.525
  strategic:
    - kind: orchestrator
      seconds: 823
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 94
          output: 1657
          cache_read: 17636947
          cache_write: 35717
          cost: 4.3563
cost_of_delay:
  value: 2.63
  by: planner-E-0015
  at: 2026-10-07T22:13:39Z
forecast:
  duration: 1h
  delivery: 2026-10-07T23:26:00Z
  basis: "Its own forecast of 1h; 2nd in the pull order with an in-progress limit of 3, behind S-0333 and S-0332."
  by: flai
  at: 2026-10-07T22:10:05Z
---
# S-0232 The release key signs flai's checksums.txt in CI and both components carry the public key

## Goal

The first step of ADR-0070: a cosign key pair exists, the two release workflows can use it, `release-flai.yml` signs `checksums.txt` with it, and the public key is a constant in flai's and flaiover's sources, so that every later story has a signature to check and a key to check it with. The finding is `design/system/release-signing.md`.

## Acceptance criteria
- [ ] The operator documentation has a runbook for generating the pair with `cosign generate-key-pair`, storing the encrypted key and its password as two GitHub Actions secrets, publishing the public key, and rotating or replacing a compromised key; the operator has run it once and the fingerprint of the public key is recorded there.
- [x] `flai/.goreleaser.yaml` has a `signs` section that runs `cosign sign-blob` over the checksum file with the key from the environment and `--tlog-upload=false`, and uploads `checksums.txt.sig` as a release asset; `release-flai.yml` provides the secrets and installs cosign, and a release whose key is missing fails.
- [ ] `release-flai.yml` also runs `actions/attest-build-provenance` over the archives and the checksum file.
- [ ] The public key is a PEM constant in a flai package and in a flaiover server module, each with a test that parses it, and the two are identical.
- [x] The claims marked **to check** in `release-signing.md § Signing releases` that bear on cosign and GoReleaser were checked against their current documentation and the document corrected where it was wrong.
- [x] `design/tech/ci.md` lists cosign and the attestation action with their versions and why.

## Tasks
- T-1221 release-signing.md's cosign, GoReleaser, and attestation claims are checked against their current documentation and corrected
- T-1222 The release-key runbook generates, stores, publishes, and rotates the cosign key pair, and records its fingerprint
- T-1223 flai's release signs checksums.txt with cosign, fails without the key, and attests its build provenance
- T-1224 design/tech/ci.md lists cosign and the attestation action with their versions and why
- T-1225 The release public key is a PEM constant in flai's buildinfo and in flaiover's release module, parsed by a test in each and checked identical

## Notes

### Planning

Touches:

- Declared: `.github/workflows/release-flai.yml`, `flai/.goreleaser.yaml`, `flai/internal/buildinfo`, `flaiover/src/lib/server/release.ts`, `design/tech/ci.md`, `docs/operators`.
- Design: `design/system/release-signing.md`, which criterion 5 corrects.
- Layout: `flaiover/src/lib/server/release.test.ts`, the test that parses the key on the flaiover side.
- Folder touches kept, as declared: `flai/internal/buildinfo`, where the key constant may go in a new file beside `buildinfo.go`; `docs/operators`, where the key runbook may be a new file under `runbooks/`. The story's agent narrows both when it writes the tasks.
- Co-change listed `design/system/flai-cli.md` and `docs/users/flai.md` at 60% and 56%. Not added: nothing a user runs changes here.

Forecast: 1h. flai replays the delivery from the pull order whenever it changes.

- `flai forecast` gave 24m from 116 s per unit of size, over only 3 medium-band feature stories.
- 19 done feature stories with 4 to 7 criteria took a median of about 1h of agent time (S-0198 to S-0229, S-0298). This one spans CI, Go, and TypeScript, and checks cosign and GoReleaser documentation, so 1h.
- The first delivery was flai's playout at its cycle factor of 6.85. It does not count the operator's key-generation step in criterion 1, which can delay it.

Cost of delay: 2.63 USD a week, as `flai cod` gives it: this story's 1h share of the 9h30m forecast over E-0015's eight open stories, of the epic's 25 USD a week penalty, which the operator set on TH-0312. Kept as given: each story closes part of one exposure, and that exposure is closed only when the chain is done, so a share by work fits.
