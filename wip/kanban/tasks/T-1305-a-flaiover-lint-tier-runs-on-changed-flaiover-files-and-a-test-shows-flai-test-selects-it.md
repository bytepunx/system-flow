---
id: T-1305
type: task
nature: improvement
title: A flaiover-lint tier runs on changed flaiover files, and a test shows flai test selects it
status: in-progress
parent: S-0319
owner: alex
created: 2026-10-08T00:15:08Z
updated: 2026-10-08T05:53:24Z
transitions:
  - to: ready
    at: 2026-10-08T05:53:24Z
    by: agent-S-0319
  - to: in-progress
    at: 2026-10-08T05:53:24Z
    by: agent-S-0319
stream: S-0319
tags: [flai, flaiover, testing, manifest]
touches: [system-flow.yaml, flai/internal/verify/select_test.go]
after: [T-1304]
---
# T-1305 A flaiover-lint tier runs on changed flaiover files, and a test shows flai test selects it

## Work

Add a `flaiover-lint` tier to `system-flow.yaml`'s `tests` with `flai manifest set tests=<JSON>`, which writes the list whole: keep every other tier as it is, and never edit the manifest by hand. The tier runs `../scripts/flaiover-lint.sh {files}` in `dir: flaiover`, `format: plain`, with no `all_only`. It is selected by `flaiover/**` less `node_modules`, `build`, and `.svelte-kit`, as the `flaiover` tier's paths are. Put it before `vitest`, since tiers run cheapest first and prettier on a few files costs less than vitest.

Add a test to `flai/internal/verify/select_test.go` that loads this repository's `system-flow.yaml` and selects tiers, outside `--all`, for a changed `flaiover/src/**/*.svelte` file. It asserts that a tier running `flaiover-lint.sh` is among them. Before the tier exists it fails, which reproduces I-0105. Mark it as reading the monorepo: skip it under `-short` if that is how the package's other tests that read the repository do it.

It waits for T-1304 because the tier runs the script T-1304 writes.

## Done when

- [ ] `flai check --strict` reports no `manifest.tests` finding
- [ ] `flai test flaiover/src/routes/+page.svelte` runs `flaiover-lint` and passes
- [ ] With a prettier fault planted in that file, `flai test` on it fails at `flaiover-lint` and names the file; the fault is reverted after
- [ ] The new test fails without the tier and passes with it

## Notes

Drafted by the planner for S-0319. The test's file is a prediction: a test file of its own beside `select_test.go` is as good, and the touches widen when the task closes.
