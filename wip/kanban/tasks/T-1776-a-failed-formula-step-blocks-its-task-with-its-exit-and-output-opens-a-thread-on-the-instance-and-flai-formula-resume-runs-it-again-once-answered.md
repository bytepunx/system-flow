---
id: T-1776
type: task
nature: feature
title: A failed formula step blocks its task with its exit and output, opens a thread on the instance, and flai formula resume runs it again once answered
status: backlog
parent: S-0422
owner: alex
created: 2026-10-16T05:47:50Z
updated: 2026-10-16T05:47:50Z
transitions: []
stream: S-0422
tags: [cli]
touches: [flai/internal/serve/formula_fail.go, flai/internal/serve/formula_fail_test.go, flai/cmd/formula.go, flai/cmd/formula_test.go, docs/users/flai-reference.md]
after: [T-1774]
---
# T-1776 A failed formula step blocks its task with its exit and output, opens a thread on the instance, and flai formula resume runs it again once answered

## Work

On a step that exits non-zero:

- block its task with `flai block`, the exit code, and the last 40 lines of its output in the task's Notes;
- open a thread on the story naming the step;
- run nothing after it.

Then:

- `flai formula resume <S-nnnn>` runs the blocked step again and goes on, once the thread is answered or with `--force`.
- Test a failure, a resume that passes, and a resume refused while the thread is open.
- Regenerate `docs/users/flai-reference.md`.

Waits for the step runner. Runs beside the role step task, since they share no path.

## Done when

`go test ./internal/serve/... ./cmd/...` passes.

## Notes
