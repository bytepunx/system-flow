---
id: T-3253
type: task
nature: feature
title: flai reads the project's own workflow file when it has one and the baseline when it has none, everywhere a repository is opened
status: done
parent: S-0755
owner: alex
created: 2026-11-12T16:09:46Z
updated: 2026-11-15T22:14:19Z
transitions:
  - to: ready
    at: 2026-11-15T21:53:51Z
    by: agent-S-0755
  - to: in-progress
    at: 2026-11-15T21:53:51Z
    by: agent-S-0755
  - to: done
    at: 2026-11-15T22:14:19Z
    by: agent-S-0755
stream: S-0755
tags: [cli]
touches: [flai/internal/workflow/load.go, flai/internal/workflow/load_test.go, flai/internal/workflow/validate.go, flai/internal/check/check_workflow.go, flai/internal/workitem/repo.go, flai/internal/workitem/repo_flow_test.go, flai/cmd/check.go, flai/cmd/project.go, flai/internal/check/concurrent_test.go]
usage:
  source: log
  seconds: 1228
  estimated: true
  priced_by: harness
  models:
    - model: claude-haiku-5-5
      input: 34
      output: 24870
      cache_read: 991725
      cache_write: 98044
      cost: 0.0346
    - model: claude-opus-5-5
      input: 60
      output: 20405
      cache_read: 4208022
      cache_write: 123470
      cost: 1.9766
---
# T-3253 flai reads the project's own workflow file when it has one and the baseline when it has none, everywhere a repository is opened

## Work

Waits for nothing: it is the load every other task reads. ADR-0307 §1 and §2.

- Move the rules that judge a workflow, `flowCheck` in `flai/internal/check/check_workflow.go` (counts, reach, effect, vocabulary), into the workflow package (`validate.go`), returning findings with a level, a code, a line, and a message; `check_workflow.go` maps them to its findings, word for word as today, so `check/workflow_test.go` passes unchanged.
- Add `load.go`: read `<layout.design>/workflows/board.md` under a root, parse it, and judge it. Absent: no workflow, the repository runs the baseline. A parse problem or an error-level finding: an error naming the code, the path and line, and the message. Warnings stop nothing.
- `workitem.Open` loads it from the repository's root and fails on that error. `Repo.Flow()` returns, in order, a test's `Workflow`, the loaded file, the baseline. Give the repository a way to read the file again when its size or modification time changed (for the MCP server, a later task), safe for concurrent `Flow()` calls.
- A second open for `flai check` alone that never fails on the file and keeps the refusal for the check to read; `flai/cmd/check.go` (and `cmd/project.go` where check gets its repository) uses it.

## Done when

`load_test.go` covers a project with a valid file (Open runs it), with none (the baseline), with a file that does not parse, and with one that parses but has an error finding (Open refuses naming it), and with a warning only (Open runs it); `flai test` passes on the changed packages, and `check/workflow_test.go` passes unchanged.

## Notes
