---
id: T-0800
type: task
nature: feature
title: The conventions, the design, the docs, and the template say the planner drafts and revisits a story's tasks
status: done
parent: S-0255
owner: alex
created: 2026-10-04T04:03:41Z
updated: 2026-10-04T04:12:30Z
transitions:
  - to: ready
    at: 2026-10-04T04:07:53Z
    by: agent-S-0255
  - to: in-progress
    at: 2026-10-04T04:07:53Z
    by: agent-S-0255
  - to: done
    at: 2026-10-04T04:12:30Z
    by: agent-S-0255
stream: S-0255
tags: []
touches: [design/conventions/strategic-agents.md, design/conventions/work-management.md, template/root/design/conventions/strategic-agents.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, template/template.yaml, design/system/strategic-agents.md, docs/users/flai.md]
after: [T-0797, T-0798]
usage:
  source: log
  seconds: 277
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 70
      output: 13472
      cache_read: 5832005
      cache_write: 92437
      cost: 2.0832
---

# T-0800 The conventions, the design, the docs, and the template say the planner drafts and revisits a story's tasks

## Work

Say what the planner now does for a story everywhere it is described: the convention `strategic-agents.md` section "As the planner" (draft a story's tasks, revisit open ones, the thread, never cancel or rewrite a task's words it did not write without asking); `work-management.md`'s rule that tasks are written by the agent that starts the story (the planner may draft them when the operator asks it to plan the story, and the story's agent reviews them when it pulls the story); the template's copies, with a template version bump and changelog entry; `design/system/strategic-agents.md` (what it is told); `docs/users/flai.md` (what `flai plan` does for a story). Waits for T-0797 and T-0798, whose behaviour it describes.

## Done when

- each document says what the planner does for a story, and the template's copies match this repository's baseline
- `template/template.yaml` and `template/CHANGELOG.md` carry the bump
- the markdown lint passes on the changed files

## Notes
