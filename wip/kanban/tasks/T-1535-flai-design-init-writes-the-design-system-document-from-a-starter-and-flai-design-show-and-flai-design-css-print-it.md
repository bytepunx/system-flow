---
id: T-1535
type: task
nature: feature
title: flai design init writes the design system document from a starter, and flai design show and flai design css print it
status: backlog
parent: S-0381
owner: alex
created: 2026-10-13T17:46:41Z
updated: 2026-10-13T17:50:44Z
transitions: []
stream: S-0381
tags: [cli]
touches: [flai/cmd/design.go, flai/cmd/design_test.go, flai/cmd/root.go, flai/internal/designsystem/starter.md]
after: [T-1533]
---
# T-1535 flai design init writes the design system document from a starter, and flai design show and flai design css print it

## Work

Add the command group `flai design` in `flai/cmd/design.go`. `init` writes the document from `flai/internal/designsystem/starter.md`, embedded, with neutral tokens in every namespace and the guidance headings, and refuses when the document exists. `show` prints a summary, and with `--json` what `designsystem.Show` returns. `css` prints the custom properties, and with `--tailwind` the `@theme` block. Tests run each command against a temporary project. It waits for T-1533, whose `Show` and renderings it prints.

## Done when

- `go test ./cmd/ -run Design` passes for init (new and existing), show (text and JSON, with and without a document), and css (both forms).
- The starter document passes `flai check --strict`.

## Notes

- Layer 3, beside T-1536.
