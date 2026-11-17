---
id: T-3455
type: task
nature: feature
title: flai check advises a manifest that sets queue or when of the workflow file's lines to write in their place
status: backlog
parent: S-0813
owner: alex
created: 2026-11-16T23:49:30Z
updated: 2026-11-16T23:49:30Z
transitions: []
stream: S-0813
tags: [flai]
touches: [flai/internal/check/check_layout.go, flai/internal/check/tests_test.go]
after: [T-3452]
---
# T-3455 flai check advises a manifest that sets queue or when of the workflow file's lines to write in their place

## Work

Criterion 2. It waits for T-3452, whose function gives the lines.

- In `flai/internal/check/check_layout.go`, beside the `manifest.tests` findings, add an advisory finding, such as `manifest.tests-queue`, on each tier that sets `queue` or `when`. While the project's file names no suites for `batch`, it gives the lines T-3452 derives and the file to put them in, the project's `design/workflows/board.md` or the one to create. Once the file names them, it says the keys are ignored and may be removed.
- Advisory, as `touches.folder` is: `flai check --strict` passes, and the close-out does not stop on it.
- Keep the `manifest.tests` error for `when` without `queue: true` as it is.

## Done when

- `tests_test.go` covers a manifest with `queue` alone, with `queue` and `when`, with neither (no finding), and one whose file already names `batch`.
- The finding's text names the lines and the file.
- `flai test` passes for `flai/internal/check`.

## Notes
