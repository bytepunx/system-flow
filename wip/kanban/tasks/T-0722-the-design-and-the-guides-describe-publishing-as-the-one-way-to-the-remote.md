---
id: T-0722
type: task
nature: improvement
title: The design and the guides describe publishing as the one way to the remote
status: ready
parent: S-0195
owner: arobson
created: 2026-10-02T23:31:09Z
updated: 2026-10-02T23:31:26Z
transitions:
  - to: ready
    at: 2026-10-02T23:31:26Z
    by: agent-S-0195
stream: S-0195
tags: []
touches: [design/system, docs, design/conventions/git.md]
after: [T-0718, T-0719, T-0720, T-0721]
---
# T-0722 The design and the guides describe publishing as the one way to the remote

## Work

Describe publishing as the one way accepted work reaches the remote: `design/system/pushing-from-the-board.md`, `flai-cli.md`, `flaiover-dashboard.md`, `workflow.md`, `dashboard-host-channel.md`, and the operator and user guides (`docs/operators/index.md`, `docs/operators/settings.md`, `docs/users/flai.md`, `docs/users/flaiover.md`, the generated `docs/users/flai-reference.md`), and the git convention's project addition. `flai push --pending` and `auto-publish` are named as the operator's shell tools, outside the workflow. Waits for T-0718, T-0719, T-0720, and T-0721, whose behaviour it describes.

## Done when

- Every listed document describes publishing (`git fetch`, then `flai release --pending`, or Publish) as the one way to the remote, with links to ADR-0067
- No document tells an agent to push an unpushed acceptance

## Notes
