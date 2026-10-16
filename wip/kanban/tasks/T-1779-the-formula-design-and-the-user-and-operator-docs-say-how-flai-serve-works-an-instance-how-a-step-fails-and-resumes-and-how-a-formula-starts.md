---
id: T-1779
type: task
nature: feature
title: The formula design and the user and operator docs say how flai serve works an instance, how a step fails and resumes, and how a formula starts
status: backlog
parent: S-0422
owner: alex
created: 2026-10-16T05:48:10Z
updated: 2026-10-16T05:50:22Z
transitions: []
stream: S-0422
tags: [cli]
touches: [design/system/strategic-agents.md, docs/users/flai.md, docs/operators/index.md, docs/operators/settings.md, design/system/flai-cli.md]
after: [T-1775, T-1776, T-1777, T-1778]
---
# T-1779 The formula design and the user and operator docs say how flai serve works an instance, how a step fails and resumes, and how a formula starts

## Work

- In the formula section of `design/system/strategic-agents.md`, or its own design document if the formula story made one (then widen the touches to it), describe:
  - the `formula` host action;
  - the step runner, its environment, and the role steps;
  - the failure and `flai formula resume`;
  - an instance's acceptance with no branch;
  - the starters.
- Describe `flai formula resume` in `docs/users/flai.md`.
- Add how to turn the action on and read a blocked instance to `docs/operators/index.md`.
- Add the action's row to `docs/operators/settings.md`.

Last layer.

## Done when

The markdown lint, the settings-index test, and `flai check --strict` pass.

## Notes
