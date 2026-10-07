---
id: T-1230
type: task
nature: feature
title: The user, operator, and design documentation say what an upgrade and install.sh verify and what a refusal means
status: backlog
parent: S-0233
owner: alex
created: 2026-10-07T22:40:45Z
updated: 2026-10-07T22:40:45Z
transitions: []
stream: S-0233
tags: [cli]
touches: [docs/users/flai.md, docs/users/flai-reference.md, docs/operators/runbooks/install.md, docs/operators/runbooks/update.md, design/system/flai-cli.md]
after: [T-1227, T-1228]
---
# T-1230 The user, operator, and design documentation say what an upgrade and install.sh verify and what a refusal means

## Work

Write what is verified, with what key, and what each refusal means, and what to do about it:

- the signature is missing
- the signature does not verify
- the key is unknown, so install the release named first
- `openssl` is absent from the host `install.sh` runs on
- the version lists mark a release that cannot be verified

Put it where each reader looks:

- `docs/users/flai.md`, under `flai self-upgrade`, `flai host upgrade`, and `flai host versions`.
- `docs/operators/runbooks/install.md` and `docs/operators/runbooks/update.md`. Link the key's fingerprint in S-0232's `docs/operators/runbooks/release-key.md` rather than copying it.
- `design/system/flai-cli.md`, under Versions: the host's flai and the tree.

State that the first release that carries the check is installed by older flais without one, and that every upgrade after it is verified. Regenerate `docs/users/flai-reference.md` with `make flai-reference` after T-1228's help text changes.

Waits for T-1227 and T-1228: it describes what they built, and the reference is generated from T-1228's help text.

## Done when

- The user and operator documentation say what is verified and what each refusal means.
- `design/system/flai-cli.md` says what self-upgrade and host upgrade verify.
- `docs/users/flai-reference.md` matches the commands' help.
- The markdown lint passes on the files changed.

## Notes
