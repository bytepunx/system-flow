---
id: T-1855
type: task
nature: feature
title: flai/internal/components finds a sub-project's components and reads their props with an embedded script run on the project's own svelte/compiler
status: done
parent: S-0437
owner: alex
created: 2026-10-16T21:32:45Z
updated: 2026-10-16T22:05:04Z
transitions:
  - to: ready
    at: 2026-10-16T21:54:45Z
    by: agent-S-0437
  - to: in-progress
    at: 2026-10-16T21:54:45Z
    by: agent-S-0437
  - to: done
    at: 2026-10-16T22:05:04Z
    by: agent-S-0437
stream: S-0437
tags: [cli]
touches: [flai/internal/components/components.go, flai/internal/components/components_test.go, flai/internal/components/props.go, flai/internal/components/props_test.go, flai/internal/components/props.mjs, flai/internal/components/testdata]
usage:
  source: log
  seconds: 619
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 69
      output: 28432
      cache_read: 3216772
      cache_write: 115201
      cost: 1.8809
---
# T-1855 flai/internal/components finds a sub-project's components and reads their props with an embedded script run on the project's own svelte/compiler

## Work

- `components.go`: `Find(root, project, include, exclude)` walks the sub-project and returns each component's path and name, sorted. It takes the globs as arguments, so it waits for no manifest change.
- `props.mjs`, embedded with `//go:embed`: it imports `svelte/compiler` from the sub-project's `node_modules`, parses each file with `parse(source, {modern: true})`, finds `$props()` in the instance script, and reads the destructured names and defaults. It resolves the type annotation when it is a local `type`, an `interface`, or an inline literal. An imported type is `unknown` with its text. The script prints JSON on standard output and one error per file it cannot parse.
- `props.go`: runs `node` on the script with the files on standard input, with a timeout, and maps its answer to `Prop{Name, Kind, Type, Default, Required}`. A missing `node`, or a sub-project without `svelte`, is an error naming the fix.
- `testdata/`: fixture `.svelte` files, one per kind, an unparsable one, and a `package.json` that the test links to flaiover's installed `svelte`.

Waits for T-1853, whose ADR fixes the kinds.

## Done when

Tests cover each prop kind, an unparsable component, and both missing-tool errors; those that need `node` skip when it is absent. `flai test flai/internal/components` passes.

## Notes

`testdata` is a folder touch: the fixtures' names are the agent's to choose.
