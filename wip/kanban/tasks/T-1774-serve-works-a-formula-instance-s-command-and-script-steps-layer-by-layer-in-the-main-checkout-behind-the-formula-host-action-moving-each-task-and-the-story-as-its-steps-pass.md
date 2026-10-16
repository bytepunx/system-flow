---
id: T-1774
type: task
nature: feature
title: serve works a formula instance's command and script steps layer by layer in the main checkout, behind the formula host action, moving each task and the story as its steps pass
status: backlog
parent: S-0422
owner: alex
created: 2026-10-16T05:47:38Z
updated: 2026-10-16T05:50:22Z
transitions: []
stream: S-0422
tags: [cli]
touches: [flai/internal/serve/formula.go, flai/internal/serve/formula_test.go, flai/internal/serve/serve.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go]
---
# T-1774 serve works a formula instance's command and script steps layer by layer in the main checkout, behind the formula host action, moving each task and the story as its steps pass

## Work

- Add the `formula` host action.
- At each look of the launcher, while the action is on, take each instance story (`formula:` set) in the backlog or in progress.
- Move the story to in-progress.
- Run its next layer's command steps, and its scripts under `scripts/`, in the main checkout with the inputs in the environment as `FLAI_INPUT_<NAME>`.
- Move each task to done when its step exits 0, and the story to review when every task is done.
- Journal each step under `formula`.
- Test a fixture formula of two command steps and a script step with a fake runner.

First layer: the role steps, the failure path, the acceptance, and the starters build on it.

## Done when

`go test ./internal/serve/... ./internal/hostapi/...` passes.

## Notes

Written on TH-0471's recommended answer to question 5. With the other answer, this task becomes the convention and prompt text that tell a story's agent how to work an instance.
