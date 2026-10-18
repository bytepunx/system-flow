---
id: T-1637
type: task
nature: feature
title: The workflow, the git convention and the template's, and the CLI and operator docs describe acceptance through the remote beside the local merge
status: done
parent: S-0397
owner: alex
created: 2026-10-16T04:40:23Z
updated: 2026-10-18T03:18:06Z
transitions:
  - to: ready
    at: 2026-10-18T03:12:29Z
    by: agent-S-0397
  - to: in-progress
    at: 2026-10-18T03:12:29Z
    by: agent-S-0397
  - to: done
    at: 2026-10-18T03:18:06Z
    by: agent-S-0397
stream: S-0397
tags: [cli, template]
touches: [design/system/workflow.md, design/system/flai-cli.md, design/conventions/git.md, template/root/design/conventions/git.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/index.md, template/CHANGELOG.md]
after: [T-1633, T-1635, T-1636]
usage:
  source: log
  seconds: 337
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 74
      output: 34019
      cache_read: 4825443
      cache_write: 189224
      cost: 2.7117
---
# T-1637 The workflow, the git convention and the template's, and the CLI and operator docs describe acceptance through the remote beside the local merge

## Work

- `design/system/workflow.md`, the review to done row: under pull-request acceptance, `flai accept` merges through the remote by the project's method, linking T-1633's ADR.
- `design/conventions/git.md` and `template/root/design/conventions/git.md`, baseline: "Accepting an item merges it …" now covers both modes. Change the template first, then this repository's.
- `design/system/flai-cli.md` and `docs/users/flai.md`: `flai accept` under pull-request acceptance, its dry run, and its refusals.
- `docs/operators/index.md`: the two modes side by side. Cover choosing one, the token, the host action, what a protected main branch means for publishing, and a runbook for switching a project from one mode to the other.
- Add a `template/CHANGELOG.md` entry.
- Waits for T-1633 (the ADR), T-1635, and T-1636: the docs describe what they built.

## Done when

Every document above describes both modes, and the markdown lint is clean.

## Notes
