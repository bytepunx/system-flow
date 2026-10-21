---
id: T-1776
type: task
nature: feature
title: "A failed step: the prompt and the convention have the agent block the task with its exit and output, ask on a thread, run nothing after, and run the step again when answered"
status: backlog
parent: S-0422
owner: alex
created: 2026-10-16T05:47:50Z
updated: 2026-10-21T12:59:28Z
transitions: []
stream: S-0422
tags: [cli]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/work-management.md, template/root/design/conventions/work-management.md]
after: [T-1774]
---
# T-1776 A failed step: the prompt and the convention have the agent block the task with its exit and output, ask on a thread, run nothing after, and run the step again when answered

## Work

- The instance paragraph of `harness.Prompt` (T-1774) gains the failure path: a command or script that exits non-zero, or a role run that ends `failed` or `asked`, blocks its task with `flai block` and the exit, writes the exit and the last 40 lines of its output in the task's Notes through `item_edit`, opens a thread on the story naming the step and the exit, runs nothing after it, writes `## Current state` and `## Next steps`, and ends the session, as for any question.
- When `flai serve` starts the agent again with the thread answered (`Answered`), the instance's prompt tells it to unblock the task, run that step again, and go on with the layers.
- The rule in `work-management.md`, here and in the template, says the same.
- Tests in `harness_test.go`: the failure path in the prompt; an answered instance is told to run the blocked step again.

Second layer: it writes the same prompt and rule as T-1774.

## Done when

`go test ./internal/harness/...` passes, and the markdown lint passes on both conventions.

## Notes

Rewritten by the planner on 2026-10-21, on TH-0471's answer 5 and TH-0582's approval. It was the serve-side failure path with `flai formula resume` before; the answer to the thread takes resume's place.
