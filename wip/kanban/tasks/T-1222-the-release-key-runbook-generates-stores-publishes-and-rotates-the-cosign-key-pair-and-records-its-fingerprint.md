---
id: T-1222
type: task
nature: feature
title: The release-key runbook generates, stores, publishes, and rotates the cosign key pair, and records its fingerprint
status: done
parent: S-0232
owner: alex
created: 2026-10-07T22:15:32Z
updated: 2026-10-07T22:27:34Z
transitions:
  - to: ready
    at: 2026-10-07T22:24:04Z
    by: agent-S-0232
  - to: in-progress
    at: 2026-10-07T22:24:04Z
    by: agent-S-0232
  - to: done
    at: 2026-10-07T22:27:34Z
    by: agent-S-0232
stream: S-0232
tags: []
touches: [docs/operators/runbooks/release-key.md, docs/operators/runbooks/README.md, docs/operators/settings.md]
after: [T-1221]
usage:
  source: log
  seconds: 210
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 25
      output: 8346
      cache_read: 1313524
      cache_write: 40871
      cost: 0.6914
---
# T-1222 The release-key runbook generates, stores, publishes, and rotates the cosign key pair, and records its fingerprint

## Work

Write `docs/operators/runbooks/release-key.md`: generate the pair with `cosign generate-key-pair`, store the encrypted key and its password as two GitHub Actions secrets of `bytepunx/system-flow`, publish the public key, record its SHA-256 fingerprint, and rotate or replace a compromised key, as `release-signing.md` § Keys says. List it in the runbooks index and the two secrets in `docs/operators/settings.md`. Once the operator has run it, record the fingerprint they report.

Waits for T-1221, whose check fixes the commands the runbook gives.

## Done when

- The runbook covers generation, the two secrets, publishing, rotation, and compromise.
- The fingerprint of the operator's public key is recorded in it.

## Notes
