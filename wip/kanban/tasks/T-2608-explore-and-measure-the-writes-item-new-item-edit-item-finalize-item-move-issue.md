---
id: T-2608
type: task
nature: research
title: "Explore and measure the writes: item.new, item.edit, item.finalize, item.move, issue.story, accept.run, and the close-out job"
status: in-progress
parent: S-0596
owner: alex
created: 2026-11-07T22:10:12Z
updated: 2026-11-07T22:12:52Z
transitions:
  - to: ready
    at: 2026-11-07T22:12:50Z
    by: agent-S-0596
  - to: in-progress
    at: 2026-11-07T22:12:52Z
    by: agent-S-0596
stream: S-0596
tags: [cli]
---
# T-2608 Explore and measure the writes: item.new, item.edit, item.finalize, item.move, issue.story, accept.run, and the close-out job

## Work

The serve log since 2026-11-04 shows `item.new` at 8 to 32 s, `item.edit` 8.5 s, `item.finalize` 6 to 19 s, `issue.story` 10 s, and `accept.run` 17 to 42 s, nearly all inside the `flai` process each starts (`exec.flai.*`). Read the code those commands run (the items lock, the checked creation and edit, the markdown lint, the replans, git, the acceptance's steps) and time them on a scratch clone under `/tmp`, never in the main checkout or a worktree. Read the close-out job's durations and steps from the serve journal and the verify records. Say where the time goes, which steps run one after another where they could run at once, and what is computed over the whole repository where the change needs only one item. Waits for no task: it writes nothing in the repository and returns its findings to the story's agent.

## Done when

- [ ] The findings name each write's cost by step, its I/O pattern by path and line, and each serial step that could run in parallel or be scoped, with a measurement.

## Notes
