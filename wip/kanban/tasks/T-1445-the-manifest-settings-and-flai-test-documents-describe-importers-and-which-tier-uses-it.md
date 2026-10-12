---
id: T-1445
type: task
nature: improvement
title: "The manifest, settings, and flai test documents describe {importers} and which tier uses it"
status: done
parent: S-0364
owner: alex
created: 2026-10-12T14:40:44Z
updated: 2026-10-12T15:31:22Z
transitions:
  - to: ready
    at: 2026-10-12T15:26:58Z
    by: agent-S-0364
  - to: in-progress
    at: 2026-10-12T15:26:59Z
    by: agent-S-0364
  - to: done
    at: 2026-10-12T15:31:22Z
    by: agent-S-0364
stream: S-0364
tags: [cli]
touches: [design/system/project-manifest.md, design/system/flai-cli.md, docs/operators/settings.md, docs/users/flai.md]
usage:
  source: log
  seconds: 263
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 29
      output: 8778
      cache_read: 1307267
      cache_write: 52268
      cost: 0.7767
---
# T-1445 The manifest, settings, and flai test documents describe {importers} and which tier uses it

## Work

Describe `{importers}` where `{packages}` and `{files}` are described:

- `design/system/project-manifest.md`, the `tests[].command` key and its example.
- `docs/operators/settings.md`, the `tests[].command` row.
- `design/system/flai-cli.md`, the `flai test` row and its tier selection.
- `docs/users/flai.md`, the `flai test` section.

Say what it stands for: the changed packages and every package of the module that imports one of them, through test files too. Say why this repository's `vet` tier uses it: a change that breaks a package importing the changed one is found by `flai test`, not first by the integration tier. Say that `go-test` and `golangci-lint` keep `{packages}`, and why.

This task waits for none: it touches only documents, so it runs beside T-1443.

## Done when

- Each of the four documents names `{importers}` beside `{packages}` and `{files}`, and says what it selects.
- `flai test` on the four documents passes the markdown tier.

## Notes

Drafted by planner-S-0364.
