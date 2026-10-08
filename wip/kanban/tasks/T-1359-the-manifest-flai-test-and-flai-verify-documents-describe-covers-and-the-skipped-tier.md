---
id: T-1359
type: task
nature: improvement
title: The manifest, flai test, and flai verify documents describe covers and the skipped tier
status: backlog
parent: S-0342
owner: alex
created: 2026-10-08T08:07:12Z
updated: 2026-10-08T08:07:12Z
transitions: []
stream: S-0342
tags: [docs]
touches: [design/system/project-manifest.md, design/system/flai-cli.md, docs/users/flai.md]
after: [T-1349]
---
# T-1359 The manifest, flai test, and flai verify documents describe covers and the skipped tier

## Work

Describe the key and the skip where each document already describes the tiers.

- `design/system/project-manifest.md`: the `tests` example block gains a commented `covers` line, and the list of a tier's fields under `tests` gains `covers`: the tiers whose tests this one runs too; when a run selects both, the covered one is skipped, and alone it runs. Name the refusals `flai check` reports as `manifest.tests` (an unknown name, the tier itself, a duplicate, a cycle) and say that this repository's `integration` covers `go-test`.
- `design/system/flai-cli.md` § Commands: the `flai test` row and the `flai verify` row say that a selected tier another selected tier covers is not run and is listed as `skipped: covered by <tier>`, the `skipped` state and `covered_by` in `--json`; the `not-reached` wording stays for tiers after a failure.
- `docs/users/flai.md` § Run the tests for what changed and § Verify a story before review: the same for users, with the text line `skipped go-test: covered by integration`, and that `flai test` on a few paths still runs `go-test`, since `integration` runs only under `--all` and in verify.
- Bump each file's `updated`.

It waits for T-1349, which settles the state's name and the `covered_by` field; it runs beside T-1355 and T-1358, which take the text line from the story's criterion.

## Done when

- The three documents describe `covers`, its validation, the skip, and its report in text and `--json`, and agree with what T-1349 built.
- `flai test` on the three files (the markdown tier) passes.

## Notes

Layer 3 of S-0342's plan. `design/system/flai-cli.md` and `docs/users/flai.md` are also in S-0341's touches. `docs/operators/settings.md` is T-1345's, whose settings test needs the row.
