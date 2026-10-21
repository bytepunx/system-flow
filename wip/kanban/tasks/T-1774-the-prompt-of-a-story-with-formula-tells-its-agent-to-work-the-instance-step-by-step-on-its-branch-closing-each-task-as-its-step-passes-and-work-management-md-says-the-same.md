---
id: T-1774
type: task
nature: feature
title: The prompt of a story with formula tells its agent to work the instance step by step on its branch, closing each task as its step passes, and work-management.md says the same
status: backlog
parent: S-0422
owner: alex
created: 2026-10-16T05:47:38Z
updated: 2026-10-21T12:58:53Z
transitions: []
stream: S-0422
tags: [cli]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/serve/agents.go, design/conventions/work-management.md, template/root/design/conventions/work-management.md]
---
# T-1774 The prompt of a story with formula tells its agent to work the instance step by step on its branch, closing each task as its step passes, and work-management.md says the same

## Work

- `harness.Request` gains `Formula`, the story's `formula:` name, which `launcher.start` in `serve/agents.go` fills from the story.
- `harness.Prompt` adds a paragraph for an instance: the story is an instance of formula `<name>`; its tasks are its steps, to be worked in their layers and in no other order, each as its `## Work` names it; a flai command or a script under `scripts/` runs from the worktree, with its output read; a role step is started through the MCP tool `role_run`, as T-1775 says; close each task with `task_done` when its step exits 0, with the exit in the task's Notes; move the story to review when every task is done; invent no step and change no step's words.
- `work-management.md`, here and in the template, gains a rule with the same words, under the story-pulling rules, saying what an instance is and how its agent works it.
- Tests in `harness_test.go`: an instance's prompt names the formula and the rules; a story with no `formula` gets the prompt it gets today, unchanged.

First layer. T-1776 adds the failure path to the same prompt and rule, so it waits for this task.

## Done when

`go test ./internal/harness/... ./internal/serve/...` passes, the markdown lint passes on both conventions, and `flai check --strict` reports nothing.

## Notes

Rewritten by the planner on 2026-10-21, on TH-0471's answer 5 and TH-0582's approval. It was the serve-side step runner before.
