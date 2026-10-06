---
id: T-1048
type: task
nature: improvement
title: The planner's epic prompt tells it to enrich each story it drafts and draft that story's tasks
status: backlog
parent: S-0300
owner: alex
created: 2026-10-06T21:46:54Z
updated: 2026-10-06T21:46:54Z
transitions: []
stream: S-0300
tags: [planner, cli]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go]
---
# T-1048 The planner's epic prompt tells it to enrich each story it drafts and draft that story's tasks

## Work

Criterion 3, in code. `planPrompt` in `flai/internal/harness/harness.go` builds the planner's prompt. For an epic (`kind == workitem.Epic`), its `work` and `summary` text says only to draft or revisit stories. Change that text so the epic's planner does four things:

- enriches each story it drafts as it would a story: `flai touches suggest`, `flai forecast`, and `flai cod` on its ID, written through flai, with a `### Planning` heading in its Notes;
- drafts that story's tasks with `## Work`, `## Done when`, a nature, tags, file `touches`, and `after`, so that they form layers. It creates them in the backlog with `item_new` (type task, the story as parent) and does the same for each story it revisits that is a draft with no tasks;
- names each story's tasks and their layers in its one thread on the epic;
- ends with a summary that names by ID the stories and tasks it created and the stories it revisited.

Reuse the story prompt's task sentences, so that the two prompts say the same thing in the same words. Extend the epic case of the prompt test in `harness_test.go` (the table near `E-0016`), and the test that an epic's planner is not told a story's work, to check the new sentences.

This task waits for nothing and runs alongside the design task: they share no path.

## Done when

- [ ] The epic prompt tells the planner to enrich each story it drafts and draft its tasks, and its thread and summary name the tasks.
- [ ] The story prompt is unchanged.
- [ ] The harness tests check the new sentences, and `scripts/flai-test.sh` passes.

## Notes

Drafted by the planner on the recommendation in TH-0201. If the operator chooses alternative (c), a story planner run queued for each drafted story, this task changes `flai/internal/serve` instead.
