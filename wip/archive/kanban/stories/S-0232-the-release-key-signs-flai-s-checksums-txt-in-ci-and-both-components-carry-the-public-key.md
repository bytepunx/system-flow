---
id: S-0232
type: story
nature: feature
title: The release key signs flai's checksums.txt in CI and both components carry the public key
status: done
parent: E-0015
owner: arobson
created: 2026-10-02T12:37:23Z
updated: 2026-10-12T15:27:44Z
transitions:
  - to: ready
    at: 2026-10-07T22:14:01Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-07T22:14:07Z
    by: agent-S-0232
  - to: review
    at: 2026-10-12T15:25:28Z
    by: agent-S-0232
  - to: done
    at: 2026-10-12T15:27:44Z
    by: orchestrator
tags: [cli, dashboard]
topics: [release, security]
touches: [".github/workflows/release-flai.yml", flai/.goreleaser.yaml, flaiover/src/lib/server/release.ts, design/tech/ci.md, design/system/release-signing.md, flaiover/src/lib/server/release.test.ts, flai/internal/buildinfo/releasekey.go, flai/internal/buildinfo/releasekey_test.go, docs/operators/runbooks/release-key.md, docs/operators/runbooks/README.md, docs/operators/settings.md, design/issues/I-0121-flai-task-done-widens-touches-with-a-path-the-commit-reverted-to-main-re-growing-a-claim-the-story-narrowed.md, design/issues/summary.md, design/issues/I-0124-flai-stream-sync-replays-a-story-s-revert-of-its-own-change-over-an-identical-change-on-main-undoing-main-s.md, design/issues/I-0125-a-settled-conversation-between-two-stories-agents-keeps-waking-each-in-turn-because-every-reply-awaits-the-other.md, design/issues/I-0129-an-empty-narrative-file-under-wip-agents-makes-every-flai-stream-log-in-the-project-fail.md, flai/internal/buildinfo/cosign.pub, design/issues/I-0136-an-acceptance-on-main-during-a-close-out-s-tiers-fails-its-final-sync-check-so-the-whole-close-out-runs-again.md, design/issues/I-0134-flai-check-finds-issues-no-story-outside-the-story-at-close-out.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3022
  turns:
    - day: 2026-10-07
      ceremony: 7
      hand_edits: 3
      work: 44
    - day: 2026-10-08
      ceremony: 5
      hand_edits: 2
      work: 95
    - day: 2026-10-12
      ceremony: 4
      hand_edits: 5
      work: 34
  models:
    - model: claude-haiku-4-5-20251001
      input: 210085
      output: 8754
      cache_read: 0
      cache_write: 0
      cost: 0.2639
    - model: claude-opus-5-5
      input: 646
      output: 169095
      cache_read: 74298891
      cache_write: 1644480
      cost: 30.6314
  strategic:
    - kind: orchestrator
      seconds: 1641
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 151
          output: 2307
          cache_read: 28535803
          cache_write: 67772
          cost: 7.1119
        - model: claude-sonnet-5-5
          input: 14
          output: 80
          cache_read: 235393
          cache_write: 62790
          cost: 0.241
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
- [x] The operator documentation has a runbook for generating the pair with `cosign generate-key-pair`, storing the encrypted key and its password as two GitHub Actions secrets, publishing the public key, and rotating or replacing a compromised key; the operator has run it once and the fingerprint of the public key is recorded there.
- [x] `flai/.goreleaser.yaml` has a `signs` section that runs `cosign sign-blob` over the checksum file with the key from the environment and `--tlog-upload=false`, and uploads `checksums.txt.sig` as a release asset; `release-flai.yml` provides the secrets and installs cosign, and a release whose key is missing fails.
- [x] `release-flai.yml` also runs `actions/attest-build-provenance` over the archives and the checksum file.
- [x] The public key is a PEM constant in a flai package and in a flaiover server module, each with a test that parses it, and the two are identical.
- [x] The claims marked **to check** in `release-signing.md § Signing releases` that bear on cosign and GoReleaser were checked against their current documentation and the document corrected where it was wrong.
- [x] `design/tech/ci.md` lists cosign and the attestation action with their versions and why.

## Tasks
- T-1221 release-signing.md's cosign, GoReleaser, and attestation claims are checked against their current documentation and corrected
- T-1222 The release-key runbook generates, stores, publishes, and rotates the cosign key pair, and records its fingerprint
- T-1223 flai's release signs checksums.txt with cosign, fails without the key, and attests its build provenance
- T-1224 design/tech/ci.md lists cosign and the attestation action with their versions and why
- T-1225 The release public key is a PEM constant in flai's buildinfo and in flaiover's release module, parsed by a test in each and checked identical

## Notes

### Criterion 4

flai's PEM is the file `flai/internal/buildinfo/cosign.pub`, embedded with `go:embed` and read through `buildinfo.ReleaseKeyPEM()`, rather than a Go `const`. That way GoReleaser publishes the same bytes as the release asset `cosign.pub`, and no second copy can drift. flaiover's is the `const` `RELEASE_PUBLIC_KEY` in `src/lib/server/release.ts`. `releasekey_test.go` and `release.test.ts` each parse the key as ECDSA P-256 and check the fingerprint `ee6df723…c5da40`, and `release.test.ts` fails when its PEM differs from flai's file.

### Criterion 3

`release-flai.yml` runs `actions/attest-build-provenance` over the archives and `checksums.txt` with `if: ${{ !github.event.repository.private }}`: on the Team plan artifact attestations exist only for public repositories, and `bytepunx/system-flow` is private, so the step is skipped until the repository is public or the organisation moves to Enterprise Cloud. The operator took this on TH-0324 (2026-10-12); the criterion is ticked on that basis.

### Snapshot builds

`scripts/flai-snapshot.sh` is S-0340's to change (MS-0014, MS-0019): its d3322fda skips `sign` in both its modes. Until S-0340 is accepted, a plain `make flai-snapshot` on a main that holds S-0232's `signs` entry needs cosign and the key; pass `--skip=sign` by hand meanwhile.

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

### Accepted by the orchestrator

- Verified: 339ac9aa153805d52303ec301802070b6bd37e69
- At: 2026-10-12T15:27:44Z

Verdict: pass at 339ac9aa153805d52303ec301802070b6bd37e69; flai verify passed every step at that commit, the verifier matched all six criteria to the diff (criterion 3's step is conditional on a public repository, as the operator took on TH-0324), and no private key or secret value is in it.

- 1: docs/operators/runbooks/release-key.md, docs/operators/runbooks/README.md, docs/operators/settings.md
- 2: flai/.goreleaser.yaml, .github/workflows/release-flai.yml
- 3: .github/workflows/release-flai.yml
- 4: flai/internal/buildinfo/cosign.pub, flai/internal/buildinfo/releasekey.go, flai/internal/buildinfo/releasekey_test.go, flaiover/src/lib/server/release.ts, flaiover/src/lib/server/release.test.ts
- 5: design/system/release-signing.md
- 6: design/tech/ci.md
