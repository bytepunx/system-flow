---
id: T-1854
type: task
nature: feature
title: The manifest's sub-projects take components with include and exclude globs and a setup module, checked, written by flai, and documented
status: done
parent: S-0437
owner: alex
created: 2026-10-16T21:32:38Z
updated: 2026-10-16T22:04:54Z
transitions:
  - to: ready
    at: 2026-10-16T21:54:44Z
    by: agent-S-0437
  - to: in-progress
    at: 2026-10-16T21:54:44Z
    by: agent-S-0437
  - to: done
    at: 2026-10-16T22:04:54Z
    by: agent-S-0437
stream: S-0437
tags: [cli]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, design/system/project-manifest.md, docs/operators/settings.md, flai/internal/check/check.go, flai/internal/check/components_test.go]
usage:
  source: log
  seconds: 610
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 80
      output: 32901
      cache_read: 3722385
      cache_write: 133308
      cost: 2.1765
---
# T-1854 The manifest's sub-projects take components with include and exclude globs and a setup module, checked, written by flai, and documented

## Work

- `manifest.go`: add `Components *Components` to `Project` (`Include`, `Exclude`, `Setup`), and an accessor that returns the defaults for a `sveltekit` or `svelte` sub-project without the key: include `src/lib/**/*.svelte`. Validate it with the manifest's other rules: a glob that does not parse, a `setup` that names no file under the sub-project, and the key on another kind.
- `settings.go`: a writer, `SetComponents(project, components)`, that sets or removes one sub-project's `components` and keeps every other key, the order, and the comments, as `flai manifest set`'s writer does. `flai manifest set` reaches no sub-project key today, and `CLAUDE.md` gives flai the rest of the manifest, so T-1856's command writes through this writer.
- `design/system/project-manifest.md` and `docs/operators/settings.md`: the key, its defaults, and its rules. The settings test fails while a manifest key has no row.

Waits for T-1853, whose ADR fixes the key's shape.

## Done when

Manifest tests cover the defaults, each refusal, a sub-project with the key, and the writer keeping comments and order. `flai test flai/internal/manifest docs/operators/settings.md` passes.

## Notes

None.
