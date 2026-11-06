---
id: T-2481
type: task
nature: feature
title: flai serve's kluster pass deploys a story entering review and again at each new head, one deploy at a time per project, and opens one thread on a failed deploy, only while the kluster host action is on
status: backlog
parent: S-0579
owner: alex
created: 2026-11-06T21:44:12Z
updated: 2026-11-06T21:47:05Z
transitions: []
stream: S-0579
tags: [kluster, serve]
touches: [flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go, flai/internal/serve/serve.go, flai/cmd/serve_actions.go]
---
# T-2481 flai serve's kluster pass deploys a story entering review and again at each new head, one deploy at a time per project, and opens one thread on a failed deploy, only while the kluster host action is on

## Work

Add `flai/internal/serve/kluster.go`, a pass flai serve runs on each look, as `forgePass` does in `forge.go`. Wire it in `serve.go`, and in `flai/cmd/serve_actions.go` where flai serve builds its supervisors, as it does `serve.NewComponents`. For each served project with the `kluster` host action on and a `kluster:` block:

- each story in review with no deployment record at its branch's head gets a deploy, through S-0578's package, run in the background so the look does not wait;
- one deploy runs at a time per project, and the rest wait in review order;
- flai serve's journal logs each start and outcome;
- a failed deploy opens one thread on the story naming the step, its exit code, and the log path, once per head, and moves nothing.

With the action off, or no block, the pass runs nothing.

Test in `kluster_test.go` with S-0578's package faked: a story entering review, a new head, two stories waiting their turn, a failure and its one thread, a second look after a failure at the same head opening no second thread, and the action off.

This is the first layer: it waits for no task.

## Done when

- The pass deploys and redeploys as above, one at a time, and logs each.
- A failure opens one thread per head and moves nothing.
- Every case above is tested, and `flai test` on the four paths passes.

## Notes
