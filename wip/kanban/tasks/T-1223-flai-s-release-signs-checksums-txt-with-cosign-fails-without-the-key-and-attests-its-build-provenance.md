---
id: T-1223
type: task
nature: feature
title: flai's release signs checksums.txt with cosign, fails without the key, and attests its build provenance
status: done
parent: S-0232
owner: alex
created: 2026-10-07T22:15:33Z
updated: 2026-10-08T08:25:00Z
transitions:
  - to: ready
    at: 2026-10-07T22:24:04Z
    by: agent-S-0232
  - to: in-progress
    at: 2026-10-07T22:24:05Z
    by: agent-S-0232
  - to: done
    at: 2026-10-07T22:27:39Z
    by: agent-S-0232
stream: S-0232
tags: []
touches: [flai/.goreleaser.yaml, ".github/workflows/release-flai.yml"]
after: [T-1221]
usage:
  source: log
  seconds: 214
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 16
      output: 4364
      cache_read: 1536508
      cache_write: 26003
      cost: 0.579
---
# T-1223 flai's release signs checksums.txt with cosign, fails without the key, and attests its build provenance

## Work

Add a `signs` section to `flai/.goreleaser.yaml` that runs `cosign sign-blob` over the checksum file with the key from the environment and `--tlog-upload=false`, writing `checksums.txt.sig` as a release asset. In `release-flai.yml`, install cosign, pass `COSIGN_PRIVATE_KEY` and `COSIGN_PASSWORD`, fail before GoReleaser when the key is missing, and run `actions/attest-build-provenance` over the archives and `checksums.txt`.

Waits for T-1221, whose check fixes the flags, versions, and whether the attestation runs on this repository.

## Done when

- A release with the key signs `checksums.txt` and uploads the signature; one without the key fails.
- The attestation step covers the archives and the checksum file.

## Notes
