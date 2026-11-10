---
id: T-2920
type: task
nature: improvement
title: flai serve's close-out job sets CLOSE_OUT_FRESH and CLOSE_OUT_ALL for the worktree's script when the request asks for fresh or all
status: backlog
parent: S-0671
owner: alex
created: 2026-11-10T08:48:14Z
updated: 2026-11-10T08:48:14Z
transitions: []
stream: S-0671
tags: [flai, close-out, verify]
touches: [flai/internal/serve/verifyjob.go, flai/internal/serve/verifyjob_test.go]
---
# T-2920 flai serve's close-out job sets CLOSE_OUT_FRESH and CLOSE_OUT_ALL for the worktree's script when the request asks for fresh or all

## Work

- In `flai/internal/serve/verifyjob.go`, when the worktree has `scripts/close-out.sh` and the request asks for `fresh` or `all`, run the script with `CLOSE_OUT_FRESH=1` or `CLOSE_OUT_ALL=1` in its environment, beside what `closeOutEnv` gives. Pass no flags: a worktree whose script predates S-0671 would take a leading `--fresh` as the story and stop, while it ignores the variables.
- Leave the path with no script, `flai verify --record-issues` with `--all` and `--fresh`, and the container close-out in `verifyjob_container.go`, which passes `--fresh` already, as they are.
- Rewrite the comment on `command` that says "all and fresh do not reach it".
- In `flai/internal/serve/verifyjob_test.go`, add a test that reproduces I-0280: the stub close-out script records its environment, and a request with `Fresh` and `All` must reach it as the two variables. Before the change it sees neither.
- Waits for nothing: the variables' names are fixed in the story's `### Planning` note, so it runs beside T-2918.

## Done when

- The new test fails on the code before the change and passes after it.
- A request with neither `fresh` nor `all` runs the script with neither variable set.
- `flai test flai/internal/serve/verifyjob.go flai/internal/serve/verifyjob_test.go` passes.

## Notes
