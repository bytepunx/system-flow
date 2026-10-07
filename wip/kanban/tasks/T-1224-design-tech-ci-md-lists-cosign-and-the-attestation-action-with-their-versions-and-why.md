---
id: T-1224
type: task
nature: feature
title: design/tech/ci.md lists cosign and the attestation action with their versions and why
status: backlog
parent: S-0232
owner: alex
created: 2026-10-07T22:15:40Z
updated: 2026-10-07T22:15:47Z
transitions: []
stream: S-0232
tags: []
touches: [design/tech/ci.md]
after: [T-1223]
---
# T-1224 design/tech/ci.md lists cosign and the attestation action with their versions and why

## Work

List cosign (through `sigstore/cosign-installer`) and `actions/attest-build-provenance` in `design/tech/ci.md` with their versions and why, and update the `release-flai.yml` row of the workflows table.

Waits for T-1223, which fixes the versions the workflow pins.

## Done when

- `design/tech/ci.md` names both with the versions `release-flai.yml` uses and the reason for each.

## Notes
