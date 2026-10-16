---
id: T-1885
type: task
nature: feature
title: flai/internal/components reads each component's examples under design/components, and flai check reports components.examples for one that does not fit its component
status: backlog
parent: S-0443
owner: alex
created: 2026-10-16T21:38:22Z
updated: 2026-10-16T21:48:30Z
transitions: []
stream: S-0443
tags: [cli]
touches: [flai/internal/components/examples.go, flai/internal/components/examples_test.go, flai/internal/check/check.go, flai/internal/check/components_test.go]
after: [T-1884]
---
# T-1885 flai/internal/components reads each component's examples under design/components, and flai check reports components.examples for one that does not fit its component

## Work

- `examples.go`: `ExamplesPath(m, component)` maps a component to its file under `design/components/` by T-1884's layout, and `Examples(component)` reads it, with each example's line. `Validate(examples, props)` reports an unknown prop, a missing required one, and a value that does not fit the prop's kind.
- `check.go`: a `checker.components` that walks `design/components/` and reports `components.examples`, including a file whose component does not exist. Without `node`, or with no Svelte sub-project, it gives one note and passes.

Waits for T-1884, whose ADR fixes the layout, the format, and the finding.

## Done when

Tests cover a valid file, each finding with its line, a file with no component, and the note without `node`. `flai test flai/internal/components flai/internal/check` passes.

## Notes

None.
