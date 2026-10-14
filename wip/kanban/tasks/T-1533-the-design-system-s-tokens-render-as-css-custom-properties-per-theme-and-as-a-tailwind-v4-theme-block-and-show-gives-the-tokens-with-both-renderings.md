---
id: T-1533
type: task
nature: feature
title: The design system's tokens render as CSS custom properties per theme and as a Tailwind v4 @theme block, and Show gives the tokens with both renderings
status: backlog
parent: S-0381
owner: alex
created: 2026-10-13T17:46:32Z
updated: 2026-10-14T23:53:25Z
transitions: []
stream: S-0381
tags: [cli]
touches: [flai/internal/designsystem/css.go, flai/internal/designsystem/css_test.go]
after: [T-1592, T-1593]
---
# T-1533 The design system's tokens render as CSS custom properties per theme and as a Tailwind v4 @theme block, and Show gives the tokens with both renderings

## Work

In `flai/internal/designsystem/css.go`, `CSS` renders the tokens as custom properties, a `:root` block for light and a `[data-theme='dark']` block for dark, and `Tailwind` renders them as a Tailwind v4 `@theme` block. Both are deterministic, ordered by namespace then key. `Show` returns the shape `flai design show --json` and `design.get` answer: the path, the tokens, the guidance headings, the problems, and both renderings, or that the project has none. Golden tests in `css_test.go`. It waits for T-1532, whose types it renders.

## Done when

- `go test ./internal/designsystem/` passes with golden tests of both renderings and of `Show` with and without a document.

## Notes

- Layer 2, beside T-1534, with which it shares no file.
