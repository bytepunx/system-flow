---
id: S-0381
type: story
nature: feature
title: A project's design system is a document whose front matter holds its tokens, and flai validates it, renders it as CSS, and serves it to the dashboard
status: backlog
parent: E-0020
owner: alex
created: 2026-10-13T17:44:24Z
updated: 2026-10-13T17:46:50Z
transitions: []
tags: [cli]
topics: [design]
touches: [flai/internal/designsystem/tokens.go, flai/internal/designsystem/tokens_test.go, flai/internal/designsystem/css.go, flai/internal/designsystem/css_test.go, flai/internal/designsystem/starter.md, flai/internal/check/design.go, flai/internal/check/design_test.go, flai/internal/check/check.go, flai/cmd/design.go, flai/cmd/design_test.go, flai/internal/hostapi/design.go, flai/internal/hostapi/design_test.go, flai/internal/hostapi/hostapi.go, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [S-0380]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
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
