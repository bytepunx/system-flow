---
id: S-0385
type: story
nature: feature
title: flai writes the design system's tokens into the stylesheet the manifest names, and flai check reports the stylesheet when it falls behind them
status: backlog
parent: E-0020
owner: alex
created: 2026-10-13T17:45:14Z
updated: 2026-10-14T18:27:22Z
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
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 55
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 18
          output: 328
          cache_read: 2888091
          cache_write: 4910
          cost: 0.7293
draft: true
cost_of_delay:
  value: 103.85
  by: planner-E-0020
  at: 2026-10-13T17:52:16Z
forecast:
  duration: 27m
  delivery: 2026-10-14T20:05:00Z
  basis: "Its own forecast of 27m; 12th in the pull order with an in-progress limit of 5, behind S-0380, S-0241, S-0304, S-0305, S-0306, S-0371, S-0372, S-0374, S-0381, S-0382, S-0383 and S-0384."
  by: flai
  at: 2026-10-14T18:27:22Z
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
- T-1556 The manifest's design.stylesheet names the stylesheet the design system's tokens are written to
- T-1557 flai/internal/designsystem renders the stylesheet's whole content with a generated-file header and says when a stylesheet differs from it
- T-1558 flai design css --write writes the stylesheet the manifest names, and refuses when it names none
- T-1559 flai check reports design.stale when the stylesheet falls behind the tokens, and a design.stylesheet outside every sub-project
- T-1560 The manifest, settings, flai design, and user docs and the reference describe design.stylesheet, flai design css --write, and design.stale

## Notes

- Drafted by the planner for E-0020. It assumes S-0380 keeps the tokens in Tailwind v4's theme namespaces, so the generated block is a project's whole theme and no mapping file is needed.

### Planning

- **Touches.** None declared. All are the planner's, every one a file:
  - `flai/internal/manifest/manifest.go` with its test: layout, where the manifest's keys are. Criterion 1.
  - `flai/internal/designsystem/stylesheet.go` with its test: new, beside S-0381's renderers. Criteria 2 and 3.
  - `flai/cmd/design.go` and `flai/internal/check/design.go`, with their tests: S-0381's command and rule file, extended. Criteria 1 to 3.
  - `design/system/project-manifest.md` and `docs/operators/settings.md`: design, and `documentation.md`'s project addition, whose test fails on a manifest key without a row. With `design/system/flai-cli.md`, `docs/users/flai.md`, and the reference. Criterion 4.
- **Co-change.** `flai touches suggest` lists `docs/operators/index.md`, the dashboard docs, and `flai/cmd/root.go`. None is added: the settings page indexes every key, not the operators' index; the story adds no dashboard view; and `flai design` is registered in `root.go` by S-0381.
- **No folder touch.**
- **Layers.** T-1556 and T-1557 together; then T-1558 and T-1559 together; then T-1560.
- **Forecast: 27m, delivery 2026-10-13T22:11:00Z.** `flai forecast` gave 27m: 95 s per unit of size over 50 done feature stories on claude-opus-5-5 in the large band, times size 17. It stands: each task extends a file S-0381 starts, in five tasks in three layers.
- **Cost of delay: 103.85 USD a week.** `flai cod`: 27m of the 5h12m forecast over E-0020's seven open stories, of the epic's 1200 USD a week. It stands.
