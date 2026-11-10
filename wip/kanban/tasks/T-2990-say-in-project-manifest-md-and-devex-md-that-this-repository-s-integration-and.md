---
id: T-2990
type: task
nature: remediation
title: "Say in project-manifest.md and devex.md that this repository's integration and smoke paths select template/**"
status: backlog
parent: S-0696
owner: alex
created: 2026-11-10T10:29:54Z
updated: 2026-11-10T10:29:54Z
transitions: []
stream: S-0696
tags: [docs, testing]
touches: [design/system/project-manifest.md, design/system/devex.md]
after: [T-2987]
---
# T-2990 Say in project-manifest.md and devex.md that this repository's integration and smoke paths select template/**

## Work

- In `design/system/project-manifest.md`, in the `tests` section, the paragraph on `all_only` tiers says "this repository gives `integration` and `smoke` `flai/**`". Make it say `flai/**` and `template/**`, and that a template-only change therefore runs both in its close-out and in the merge queue's batch (I-0270). In the next paragraph, check that what it says of `when: ["flai/**", "template/**"]` still reads true.
- In `design/system/devex.md`, in the Test tiers row, change "integration, and smoke (both `flai/**`)" to name `template/**` too.
- Read the passages in `design/conventions/code-quality.md`, `design/conventions/tooling.md`, and `docs/users/flai.md` that say a story that touches or changes `template/` runs both tiers. They should now be true as written. If one is not, fix it and add its path to this task's touches.
- Waits for T-2987, because it describes the paths that task writes. If T-2987 takes another fix, these words change with it.

## Done when

- No document under `design/` or `docs/` says that this repository's `integration` or `smoke` paths are `flai/**` alone.
- `flai test design/system/project-manifest.md design/system/devex.md` passes.

## Notes
