---
id: T-1445
type: task
nature: improvement
title: "The manifest, settings, and flai test documents describe {importers} and which tier uses it"
status: backlog
parent: S-0364
owner: alex
created: 2026-10-12T14:40:44Z
updated: 2026-10-12T14:40:44Z
transitions: []
stream: S-0364
tags: [cli]
touches: [design/system/project-manifest.md, design/system/flai-cli.md, docs/operators/settings.md, docs/users/flai.md]
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
