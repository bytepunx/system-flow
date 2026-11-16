---
id: T-3259
type: task
nature: feature
title: flai check reports a project's workflow file it refuses, judges the items against the file it runs, and gives no advisory to a project without one
status: done
parent: S-0755
owner: alex
created: 2026-11-12T16:10:11Z
updated: 2026-11-15T22:28:31Z
transitions:
  - to: ready
    at: 2026-11-15T22:14:34Z
    by: agent-S-0755
  - to: in-progress
    at: 2026-11-15T22:14:34Z
    by: agent-S-0755
  - to: done
    at: 2026-11-15T22:28:31Z
    by: agent-S-0755
stream: S-0755
tags: [cli]
touches: [flai/internal/check/check_workflow.go, flai/internal/check/check_items.go, flai/internal/check/workflow_test.go]
after: [T-3253]
usage:
  source: log
  seconds: 837
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 37
      output: 12757
      cache_read: 2630672
      cache_write: 77188
      cost: 1.2357
---
# T-3259 flai check reports a project's workflow file it refuses, judges the items against the file it runs, and gives no advisory to a project without one

## Work

Waits for T-3253: the check reads the refusal that task's open keeps. ADR-0307 §2 and §5.

- With the project's file valid, `workflow.item` and the other per-item judgments read the project's workflow (`c.repo.Flow()` is the loaded file).
- With the file refused, `flai check` still runs, reports the file's finding as today, and leaves out `workflow.item` and the per-item judgments that read the workflow.
- No advisory for a project without the file (ADR-0307 §5): nothing to add; a test pins that a project without the file has no `workflow.` finding.

## Done when

`workflow_test.go` covers a project's valid file with a state of its own that an item is in (no `workflow.item`), a refused file (one finding, no `workflow.item`), and a project without the file (no `workflow.` finding); `flai test` passes on `flai/internal/check`.

## Notes
