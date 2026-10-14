---
id: S-0381
type: story
nature: feature
title: A project's design system is a document whose front matter holds its tokens, and flai validates it, renders it as CSS, and serves it to the dashboard
status: backlog
parent: E-0020
owner: alex
created: 2026-10-13T17:44:24Z
updated: 2026-10-14T18:27:22Z
transitions: []
tags: [cli]
topics: [design]
touches: [flai/internal/designsystem/tokens.go, flai/internal/designsystem/tokens_test.go, flai/internal/designsystem/css.go, flai/internal/designsystem/css_test.go, flai/internal/designsystem/starter.md, flai/internal/check/design.go, flai/internal/check/design_test.go, flai/internal/check/check.go, flai/cmd/design.go, flai/cmd/design_test.go, flai/cmd/root.go, flai/internal/hostapi/design.go, flai/internal/hostapi/design_test.go, flai/internal/hostapi/hostapi.go, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [S-0380]
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
          input: 19
          output: 329
          cache_read: 2888092
          cache_write: 4910
          cost: 0.7293
draft: true
cost_of_delay:
  value: 173.08
  by: planner-E-0020
  at: 2026-10-13T17:51:23Z
forecast:
  duration: 45m
  delivery: 2026-10-14T19:34:00Z
  basis: "Its own forecast of 45m; 8th in the pull order with an in-progress limit of 5, behind S-0380, S-0241, S-0304, S-0305, S-0306, S-0371, S-0372 and S-0374."
  by: flai
  at: 2026-10-14T18:27:22Z
---
# S-0381 A project's design system is a document whose front matter holds its tokens, and flai validates it, renders it as CSS, and serves it to the dashboard

## Goal

A project keeps its design system in one markdown document, `design/system/design-system.md` unless S-0380's decision names another, whose front matter holds its tokens and whose body holds its guidance. flai owns the format: it writes a starter, validates the tokens, renders them as CSS custom properties and as a Tailwind v4 `@theme` block, and answers them to the dashboard, so the dashboard and agents read the same tokens.

## Acceptance criteria

- [ ] `flai design init` writes the design system document from a starter: tokens for colours in light and dark, font families, a type scale, a spacing step, radii, and shadows, named after Tailwind v4's theme namespaces, and guidance headings in its body; it refuses when the document exists.
- [ ] `flai check` reports each token that does not parse, such as a colour that is not a CSS colour, a length without a unit, or a dark colour with no light one, as `design.tokens` with its key and line, and passes a project with no design system document.
- [ ] `flai design css` prints the tokens as CSS custom properties, one block per theme, and with `--tailwind` as a Tailwind v4 `@theme` block; `flai design show --json` gives the tokens, the guidance headings, and both renderings.
- [ ] The host method `design.get` answers what `flai design show --json` prints, and says when the project has no design system document.
- [ ] `design/system/flai-cli.md`, `docs/users/flai.md`, and the generated reference say how.

## Tasks
- T-1532 flai/internal/designsystem reads a design system document's tokens and guidance headings and reports each token that does not parse, with its key and line
- T-1533 The design system's tokens render as CSS custom properties per theme and as a Tailwind v4 @theme block, and Show gives the tokens with both renderings
- T-1534 flai check reports a design system token that does not parse as design.tokens with its key and line
- T-1535 flai design init writes the design system document from a starter, and flai design show and flai design css print it
- T-1536 The host method design.get answers the dashboard what flai design show --json prints
- T-1537 The flai design and user docs and the reference describe the design system document, flai design, design.tokens, and design.get

## Notes

- Drafted by the planner for E-0020. It assumes S-0380 recommends tokens in the front matter of one markdown document under `design/system`; the story is checked against S-0380's decision before it is finalized.

### Planning

- **Touches.** None declared. All are the planner's, every one a file:
  - `flai/internal/designsystem/` (`tokens.go`, `css.go`, `starter.md`, with tests): layout, a new package beside `docedit` and `adr`. Criteria 1 and 3.
  - `flai/internal/check/design.go` with its test, and `check.go`, whose `Run` calls each rule: layout. Criterion 2.
  - `flai/cmd/design.go` with its test, and `flai/cmd/root.go`, which registers each command group: layout, and co-change (`root.go` changes with new commands). Criteria 1 and 3.
  - `flai/internal/hostapi/design.go` with its test, and `hostapi.go`, whose method table merges each group's methods: layout. Criterion 4.
  - `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/users/flai-reference.md`: design, as `documentation.md` asks for a new command. Criterion 5.
- **Co-change.** `flai touches suggest` also lists `docs/operators/index.md`, `docs/operators/settings.md`, the dashboard docs, and `flai/internal/hostapi/writes.go`. None is added: the story adds no setting, no dashboard view, and no write.
- **No folder touch.**
- **Layers.** T-1532; then T-1533 and T-1534 together; then T-1535 and T-1536 together; then T-1537.
- **Forecast: 45m, delivery 2026-10-13T21:38:00Z.** `flai forecast` gave 34m: 95 s per unit of size over 50 done feature stories on claude-opus-5-5 in the large band, times size 21. Raised to 45m: the story starts a package from nothing, with a parser, two renderers with golden tests, a check rule, a command group, and a host method, across six tasks in four layers, which a rate per touch does not weigh. The delivery adds the 11m to flai's 21:27Z.
- **Cost of delay: 173.08 USD a week.** `flai cod`: 45m of the 5h12m forecast over E-0020's seven open stories, of the epic's 1200 USD a week. It stands.
