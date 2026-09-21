---
id: T-0314
type: task
nature: feature
title: A component's pending bump is the highest delivery type among everything merged and unpublished for it since its last tag
status: done
parent: S-0087
owner: alex
created: 2026-09-21T03:32:54Z
updated: 2026-09-21T03:52:19Z
transitions:
  - to: ready
    at: 2026-09-21T03:46:41Z
    by: system-flow
  - to: in-progress
    at: 2026-09-21T03:46:41Z
    by: system-flow
  - to: done
    at: 2026-09-21T03:52:19Z
    by: system-flow
stream: S-0087
tags: []
---
# T-0314 A component's pending bump is the highest delivery type among everything merged and unpublished for it since its last tag

## Work
Extend `flai/internal/release` from "the bump for one item" to "the bump for everything merged and unreleased, per component, since that component's last tag": walk main's history back to each component's current tag (`CurrentVersion` already reads it), find every done story or epic whose commit lies in that range (commit subjects already carry `[S-nnnn]`/`[E-nnnn]`, per git.md), compute each one's level with the existing `LevelFor`, and take the highest per component it touches or delivers to (an epic done outranks a feature story, which outranks a remediation/improvement). A component with nothing pending is left out of the plan entirely.

## Done when
- A function (or a new mode of `Compute`) returns one Plan per component with pending work, each carrying every story ID it bundles and the highest level among them
- Three patch-level stories and one feature story against the same component since its last tag produce one minor bump, not four patches
- A component nothing has touched since its last tag is absent from the plan
- Unit tests: a mixed batch's bump, an empty batch, a batch spanning an epic's completion (major)

## Notes
This is the core rule the operator asked for. Reuses `LevelFor`, `CurrentVersion`, and the delivery-target logic already in `release.go`; the new part is finding everything pending and rolling several items' levels into one.
