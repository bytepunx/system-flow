---
id: T-1885
type: task
nature: feature
title: flai/internal/components reads each component's examples under design/components, and flai check reports components.examples for one that does not fit its component
status: done
parent: S-0443
owner: alex
created: 2026-10-16T21:38:22Z
updated: 2026-10-17T04:13:18Z
transitions:
  - to: ready
    at: 2026-10-17T04:04:37Z
    by: agent-S-0443
  - to: in-progress
    at: 2026-10-17T04:04:38Z
    by: agent-S-0443
  - to: done
    at: 2026-10-17T04:13:18Z
    by: agent-S-0443
stream: S-0443
tags: [cli]
touches: [flai/internal/components/examples.go, flai/internal/components/examples_test.go, flai/internal/check/check.go, flai/internal/check/components_test.go]
after: [T-1884]
usage:
  source: log
  seconds: 520
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 41
      output: 16373
      cache_read: 2058580
      cache_write: 69995
      cost: 1.1373
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
