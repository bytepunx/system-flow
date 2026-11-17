---
id: T-3452
type: task
nature: feature
title: The workflow package derives, from a manifest's queue and when, the verify and verify_at.batch lines that stand in their place
status: backlog
parent: S-0813
owner: alex
created: 2026-11-16T23:49:15Z
updated: 2026-11-16T23:49:15Z
transitions: []
stream: S-0813
tags: [flai]
touches: [flai/internal/workflow/migrate.go, flai/internal/workflow/migrate_test.go, flai/internal/workflow/workflow.go, flai/internal/workflow/workflow_test.go]
after: [T-3450]
---
# T-3452 The workflow package derives, from a manifest's queue and when, the verify and verify_at.batch lines that stand in their place

## Work

Criteria 2 and 3. The check's advice, the reading of the keys, and the upgrade all need one answer to "what lines does this manifest's `queue` and `when` stand for". Write it once here. It waits for T-3450, whose ADR fixes the mapping.

- In `flai/internal/workflow/migrate.go`, add a function that takes the tiers' name, suite (S-0805), `queue`, and `when`, and a workflow. It returns the workflow with the move to review's `verify:` and `verify_at.batch` as the ADR maps them, and the YAML lines a project would write. Take a plain tier list, not `manifest.TestTier`, so the workflow package gains no import of the manifest package.
- Add a predicate for "this file names suites for `batch`", the point at which the ADR retires the keys.
- If the ADR gives a `verify:` entry a `when` of its own, parse and validate it in `flai/internal/workflow/workflow.go`, beside S-0806's `verify:`, and refuse it on any other key as the strict parse does.
- A manifest with no `queue` and no `when` returns the workflow unchanged and no lines.

## Done when

- `migrate_test.go` covers a `queue` tier with no `when`, a `queue` tier with `when`, a manifest with neither, this repository's `integration` and `smoke` tiers, and a file that already names `batch`.
- The lines it returns parse back, with the workflow package's loader, into the workflow it returns.
- `flai test` passes for `flai/internal/workflow`.

## Notes
