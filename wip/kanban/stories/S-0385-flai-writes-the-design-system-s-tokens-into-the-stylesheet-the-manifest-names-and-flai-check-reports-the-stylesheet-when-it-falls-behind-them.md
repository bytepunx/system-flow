---
id: S-0385
type: story
nature: feature
title: flai writes the design system's tokens into the stylesheet the manifest names, and flai check reports the stylesheet when it falls behind them
status: backlog
parent: E-0020
owner: alex
created: 2026-10-13T17:45:14Z
updated: 2026-10-13T17:45:14Z
transitions: []
tags: [cli]
topics: [design]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/designsystem/stylesheet.go, flai/internal/designsystem/stylesheet_test.go, flai/cmd/design.go, flai/cmd/design_test.go, flai/internal/check/design.go, flai/internal/check/design_test.go, design/system/project-manifest.md, design/system/flai-cli.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [S-0381]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0385 flai writes the design system's tokens into the stylesheet the manifest names, and flai check reports the stylesheet when it falls behind them

## Goal

The design system's tokens reach the project's code by one path: the manifest names the stylesheet a sub-project imports its theme from, `flai design css --write` writes the tokens there as a Tailwind v4 `@theme` block, and `flai check` reports when that stylesheet no longer matches the tokens. A decision taken on the Design page then changes the built UI the next time an agent writes the stylesheet, and drift between the two is caught.

## Acceptance criteria

- [ ] The manifest's `design.stylesheet` names the stylesheet, a path under a sub-project, and `flai check` refuses a path outside every sub-project.
- [ ] `flai design css --write` writes the Tailwind v4 `@theme` block, with a generated-file header naming the command, to that stylesheet, and refuses when the manifest names none.
- [ ] `flai check` reports `design.stale` when the stylesheet differs from what `flai design css --write` would write, naming the command that brings it in step, and passes a project whose manifest names no stylesheet.
- [ ] `design/system/project-manifest.md`, `docs/operators/settings.md`, `design/system/flai-cli.md`, `docs/users/flai.md`, and the generated reference say how.

## Tasks

## Notes

- Drafted by the planner for E-0020. It assumes S-0380 keeps the tokens in Tailwind v4's theme namespaces, so the generated block is a project's whole theme and no mapping file is needed.
