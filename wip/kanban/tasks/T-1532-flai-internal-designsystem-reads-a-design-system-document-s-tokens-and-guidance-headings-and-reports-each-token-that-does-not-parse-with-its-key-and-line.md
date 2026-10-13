---
id: T-1532
type: task
nature: feature
title: flai/internal/designsystem reads a design system document's tokens and guidance headings and reports each token that does not parse, with its key and line
status: backlog
parent: S-0381
owner: alex
created: 2026-10-13T17:46:26Z
updated: 2026-10-13T17:46:26Z
transitions: []
stream: S-0381
tags: [cli]
touches: [flai/internal/designsystem/tokens.go, flai/internal/designsystem/tokens_test.go]
---
# T-1532 flai/internal/designsystem reads a design system document's tokens and guidance headings and reports each token that does not parse, with its key and line

## Work

Add the package `flai/internal/designsystem`. `Read` parses the document's front matter `tokens` block, in Tailwind v4's theme namespaces (colours with a light and an optional dark value, fonts, a type scale, a spacing step, radii, shadows), and lists the body's guidance headings. `Validate` returns one problem per token that does not parse: a colour that is not a CSS colour, a length without a unit, a dark colour with no light one, an unknown namespace; each with its key path and front matter line. Table tests cover each case and a valid document. It waits for nothing.

## Done when

- `go test ./internal/designsystem/` passes with a test per problem and a valid document.

## Notes

- Layer 1. The format follows S-0380's decision; the path the package reads is a constant here, `design/system/design-system.md`, unless that decision names another.
