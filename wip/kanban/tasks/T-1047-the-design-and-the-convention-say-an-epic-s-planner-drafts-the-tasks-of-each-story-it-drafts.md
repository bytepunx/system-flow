---
id: T-1047
type: task
nature: improvement
title: The design and the convention say an epic's planner drafts the tasks of each story it drafts
status: backlog
parent: S-0300
owner: alex
created: 2026-10-06T21:46:45Z
updated: 2026-10-06T22:32:06Z
transitions: []
stream: S-0300
tags: [planner, conventions]
touches: [design/system/strategic-agents.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md]
---
# T-1047 The design and the convention say an epic's planner drafts the tasks of each story it drafts

## Work

Criterion 3: the planner, asked to plan an epic, writes draft stories and their tasks. Today it drafts only the stories.

- In `design/system/strategic-agents.md`, under The planner › What it is told, change the table rows for an epic.
  - An epic with no stories: draft its stories, enrich each one as a story is enriched (touches, forecast, cost of delay value), and draft its tasks with `## Work`, `## Done when`, a nature, tags, `touches`, and `after`, in layers.
  - An epic with stories: draft the tasks of each story it adds, and of each story it revisits that is a draft with no tasks.
  - The thread names each story's tasks and layers.
  - The summary names the stories and tasks created.
- Say the tasks carry no draft flag. They are drafts because their story is one, and the story's agent reviews them when it pulls the story.
- In `design/conventions/strategic-agents.md`, under "As the planner", add the same rule for an epic. Copy it to `template/root/design/conventions/strategic-agents.md` and add a line to `template/CHANGELOG.md`, in the same story, as CLAUDE.md says.
- This task waits for nothing and runs alongside the prompt task: they share no path.

## Done when

- [ ] The design table and the convention say that an epic's planner drafts the tasks of each story it drafts, and what its thread and summary name.
- [ ] The two copies of the convention's "As the planner" section are identical.
- [ ] `template/CHANGELOG.md` records the change.
- [ ] `flai check --strict` and the markdown lint pass.

## Notes

Drafted by the planner. alex confirmed the approach in TH-0201 on 2026-10-06: the tasks get no draft flag, and no ADR is needed.
