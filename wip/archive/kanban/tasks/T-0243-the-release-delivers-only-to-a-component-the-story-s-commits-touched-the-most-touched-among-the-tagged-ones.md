---
id: T-0243
type: task
nature: feature
title: The release delivers only to a component the story's commits touched, the most touched among the tagged ones
status: done
parent: S-0047
owner: alex
created: 2026-09-19T10:30:01Z
updated: 2026-09-19T10:32:32Z
transitions:
  - to: ready
    at: 2026-09-19T10:30:51Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T10:30:51Z
    by: system-flow
  - to: done
    at: 2026-09-19T10:32:32Z
    by: system-flow
stream: S-0047
tags: []
---

# T-0243 The release delivers only to a component the story's commits touched, the most touched among the tagged ones

## Work
`deliverTarget` in `flai/internal/release`: a tag (the story's, then its epic's) delivers only to a component the story's commits touched; when the tags name several touched components, the one with the most touched files wins, ties going to the first tag; a tag naming an untouched component never delivers, so that component gets no release at all; `--deliver` still overrides everything, as the operator's explicit word. With no tag deciding, a single touched component still delivers and several still ask for a tag or `--deliver`. Tests: a [dashboard, cli] story that only touched the CLI; tags naming two touched components; a tag naming only an untouched one with one touched; the epic's tags; `--deliver`.

## Done when
- The cases are tested and fail without the change
- `make flai-test` passes

## Notes
