---
id: T-0659
type: task
nature: feature
title: The tooling convention has a shell rule for zsh and for gating each step of a chain
status: done
parent: S-0187
owner: arobson
created: 2026-10-01T09:26:25Z
updated: 2026-10-01T09:27:07Z
transitions:
  - to: ready
    at: 2026-10-01T09:26:41Z
    by: agent-S-0187
  - to: in-progress
    at: 2026-10-01T09:26:41Z
    by: agent-S-0187
  - to: done
    at: 2026-10-01T09:27:07Z
    by: agent-S-0187
stream: S-0187
tags: []
usage:
  source: log
  seconds: 26
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 11
      output: 3249
      cache_read: 449474
      cache_write: 11040
      cost: 0.2318
---

# T-0659 The tooling convention has a shell rule for zsh and for gating each step of a chain

## Work

Add two baseline rules to `tooling.md`, here and in `template/root`: the host shell may be zsh (quote globs, URLs, and variables; `[ a = b ]`; a sequence of more than a few commands goes in a script under `scripts/`), and each step of a chain that must succeed is joined with `&&` or runs under `set -e`, never `;` or a pipe that hides its exit code, with what a commit records checked before the commit. Update the `tooling.md` row in `design/system/conventions.md`, and the template's version and changelog.

## Done when

- Both copies of `tooling.md` carry the same two rules above the marker.
- `design/system/conventions.md` describes them, and `template/template.yaml` and `template/CHANGELOG.md` record the patch.

## Notes
