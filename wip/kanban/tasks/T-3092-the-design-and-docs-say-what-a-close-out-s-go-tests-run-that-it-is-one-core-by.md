---
id: T-3092
type: task
nature: improvement
title: The design and docs say what a close-out's Go tests run, that it is one core by design, and when to run flai test race
status: backlog
parent: S-0722
owner: alex
created: 2026-11-11T06:06:52Z
updated: 2026-11-11T06:06:52Z
transitions: []
stream: S-0722
tags: []
touches: [design/system/workflow.md, design/system/project-manifest.md, design/conventions/code-quality.md, design/conventions/tooling.md, design/tech/go.md, docs/contributors/index.md, flai/README.md, Makefile]
after: [T-3091]
---
# T-3092 The design and docs say what a close-out's Go tests run, that it is one core by design, and when to run flai test race

## Work

- `design/system/workflow.md`, and where the tiers are described (`design/system/project-manifest.md`, the project additions of `design/conventions/code-quality.md` and `tooling.md`), say what a close-out's Go tests run, that it is one core and no race detector by design, with T-3091's figures, and when to run `flai test race`.
- `design/tech/go.md`, `docs/contributors/index.md`, `flai/README.md`, and the `Makefile`'s help lines stop saying the close-out or `make test` runs `-race`.
- The template's copy of anything changed is changed with it; the template's scripts and baselines run no `-race` today.
- Waits for T-3091: it cites the measured figures and the timeouts T-3091 sets.

## Done when

- No document under `design/`, `docs/`, or `flai/README.md` says a close-out runs `-race`, and the markdown lint passes.

## Notes
