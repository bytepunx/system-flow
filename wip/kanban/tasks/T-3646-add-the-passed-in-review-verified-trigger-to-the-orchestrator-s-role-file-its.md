---
id: T-3646
type: task
nature: improvement
title: Add the passed-in-review verified trigger to the orchestrator's role file, its baseline, and its template copy
status: backlog
parent: S-0854
owner: alex
created: 2026-11-27T23:59:47Z
updated: 2026-11-27T23:59:47Z
transitions: []
stream: S-0854
tags: [flai, template]
touches: [template/root/design/workflows/roles/orchestrator.md, flai/internal/roles/baseline/orchestrator.md, design/workflows/roles/orchestrator.md, flai/internal/harness/harness_test.go]
after: [T-3642]
---
# T-3646 Add the passed-in-review verified trigger to the orchestrator's role file, its baseline, and its template copy

## Work

This task waits for T-3642. A role file that names `verified` or `passed-in-review` fails validation until T-3642's vocabulary is in, and `flai check` would then refuse the project's own role file.

Add `{event: verified, when: passed-in-review}` to the `trigger:` list after `{event: moved, when: {to: review}}`. Make the change in the template copy first, `template/root/design/workflows/roles/orchestrator.md`, then in `flai/internal/roles/baseline/orchestrator.md`, then in the project's `design/workflows/roles/orchestrator.md`, so that the three lists stay the same. Bump each file's `updated`. Change no other key and no prompt text.

Run the tests that read the baseline role files, among them `flai/internal/harness/harness_test.go` and `flai/internal/roles/roles_test.go`. Update an expectation in `harness_test.go` only where it holds the trigger list.

This task runs beside T-3645. The two change no file in common.

## Done when

- The three role-file copies hold the same trigger list, with the new entry.
- `flai check --strict` reports no `roles.*` finding on the project's role file.
- `flai test` on the changed paths passes.

## Notes
