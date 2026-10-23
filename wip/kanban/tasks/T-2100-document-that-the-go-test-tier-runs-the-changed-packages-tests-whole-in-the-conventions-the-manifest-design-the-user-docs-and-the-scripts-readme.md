---
id: T-2100
type: task
nature: remediation
title: Document that the go-test tier runs the changed packages' tests whole, in the conventions, the manifest design, the user docs, and the scripts' README
status: backlog
parent: S-0490
owner: alex
created: 2026-10-23T21:22:45Z
updated: 2026-10-23T21:22:45Z
transitions: []
stream: S-0490
tags: [flai, tests, docs]
touches: [design/conventions/tooling.md, design/conventions/code-quality.md, design/system/project-manifest.md, docs/users/flai.md, scripts/README.md, scripts/flai-test.sh]
after: [T-2099]
---
# T-2100 Document that the go-test tier runs the changed packages' tests whole, in the conventions, the manifest design, the user docs, and the scripts' README

## Work

Waits for T-2099, so that the command and the figures it records are what the documents describe.

Say in each document that this repository's `go-test` tier runs the Go tests of the changed packages whole, the real-git and monorepo tests that skip under `-short` included and the gateway tests left out, while `make test`, `scripts/test.sh`, and CI keep `-short`:

- `design/conventions/tooling.md`, the project addition that lists what `flai test` runs: replace "the short Go tests" with the Go tests of the changed packages, real-git ones included, and say that a test that skips under `-short` therefore runs between tasks and in the close-out, not first in the merge queue's batch (I-0201).
- `design/conventions/code-quality.md`, the project addition "Tiers here": keep what it says of `make test` and `make integration`, and add that `flai test`'s `go-test` tier runs the changed packages without `-short`, so `-short` selects the tiers for people and CI only.
- `design/system/project-manifest.md`, the `tests` bullet: where it says a close-out that changes Go files under `flai/` runs `go-test` and leaves `integration` and `smoke` to the batch, add that `go-test` runs without `-short` and skips `^TestGateway`, why (I-0201, linked), and what it costs by the figures in T-2099's Notes. Bump `updated`.
- `docs/users/flai.md`, § Run the tests for what changed: one sentence after "integration covers go-test", saying `go-test` runs the changed packages' tests whole, real-git ones included. Bump `updated`.
- `scripts/README.md`, the first paragraph and the `flai-test.sh` row: "the short Go tests" is now the changed packages' Go tests; `test.sh`'s row stays, since the script keeps `-short`.
- `scripts/flai-test.sh`, the comment that says `go-test` is `test.sh`'s: say that `go-test` runs the changed packages without `-short` and the full run here holds every test once.

Run `flai test` on the files for the markdown lint.

## Done when

- No document under `design/` or `docs/`, and no script comment, says that `flai test` runs the short Go tests.
- Each file above names the `go-test` tier as running without `-short` and skipping the gateway tests, and `project-manifest.md` links I-0201 and gives the cost.
- `flai test` on the changed files passes.

## Notes

The template is not touched: its manifest has no `go-test` tier, and its `tooling.md` baseline does not say "short".
