---
id: T-3106
type: task
nature: remediation
title: Describe the per-login budget state, the harness change's clear, and flai serve budget in the design and the docs
status: backlog
parent: S-0720
owner: alex
created: 2026-11-11T08:55:25Z
updated: 2026-11-11T08:55:25Z
transitions: []
stream: S-0720
tags: [flai, flaiover, docs, budget]
touches: [design/system/flai-cli.md, docs/operators/index.md, docs/users/flai.md, docs/users/flaiover.md, design/system/flaiover-dashboard.md]
after: [T-3104, T-3105, T-3096]
---
# T-3106 Describe the per-login budget state, the harness change's clear, and flai serve budget in the design and the docs

## Work

- Rewrite `design/system/flai-cli.md` § The host's budget state: the key (provider, else harness with its resolved program and login), the old shape still read, the clear on a harness change with its why, `flai serve budget` and `flai serve budget clear`, and `budgetHeldWhy`'s new text. Fix the `budgetKey` link, which points at `project-manifest.md`, where the key is not described.
- In `docs/operators/index.md`, the S-0534 section and its index row: say what holds per login, that a harness change lifts the old login's hold, and that `flai serve budget clear` replaces removing `budget.json`.
- Add the two commands to `docs/users/flai.md`, and the banner's command to `docs/users/flaiover.md` and `design/system/flaiover-dashboard.md`.
- Waits for T-3104, T-3105, and T-3096: it describes the commands, the harness change's clear, and the banner as they built them.

## Done when

- The design and the docs describe the per-login key, the harness change's clear, and the manual clear, as built.
- The markdown lint and `flai check --strict` pass on the changed files.

## Notes
