---
id: T-1856
type: task
nature: feature
title: flai components list prints each Svelte sub-project's components and their props, and flai components config writes a sub-project's components key
status: done
parent: S-0437
owner: alex
created: 2026-10-16T21:32:53Z
updated: 2026-10-16T22:26:25Z
transitions:
  - to: ready
    at: 2026-10-16T22:05:17Z
    by: agent-S-0437
  - to: in-progress
    at: 2026-10-16T22:05:17Z
    by: agent-S-0437
  - to: done
    at: 2026-10-16T22:26:25Z
    by: agent-S-0437
stream: S-0437
tags: [cli]
touches: [flai/cmd/components.go, flai/cmd/components_test.go, flai/cmd/root.go, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-1854, T-1855]
usage:
  source: log
  seconds: 1268
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 83
      output: 34105
      cache_read: 3858560
      cache_write: 138185
      cost: 2.2562
---
# T-1856 flai components list prints each Svelte sub-project's components and their props, and flai components config writes a sub-project's components key

## Work

- `flai/cmd/components.go`: the `components` group with `list [--project <name>] [--json]`. It loads the manifest, takes each supported sub-project's globs from T-1854's accessor, and calls T-1855's `Find` and props reader. An unsupported sub-project is listed as unsupported with its kind. Register the group in `root.go`.
- `config <sub-project> [--include <glob>]... [--exclude <glob>]... [--setup <path>] [--unset]`: writes the key with T-1854's `SetComponents` and runs the manifest's check, refusing and writing nothing on a finding.
- Regenerate `docs/users/flai-reference.md` with `make flai-reference`.
- Run `list` on this repository and check flaiover's `Chart.svelte` and `Gallery.svelte` against their sources.

Waits for T-1854 for the manifest's globs and writer, and T-1855 for the reader.

## Done when

Command tests cover `list` as text and JSON, `--project`, an unsupported sub-project, the missing-tool error, and `config` writing, unsetting, and refusing. The run on this repository lists flaiover's components with the props criterion 4 names, and `flai test flai/cmd` passes.

## Notes

None.
