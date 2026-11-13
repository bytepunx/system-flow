---
id: S-0770
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: ready
owner: alex
created: 2026-11-13T13:16:45Z
updated: 2026-11-13T13:16:46Z
transitions:
  - to: ready
    at: 2026-11-13T13:16:46Z
    by: system-flow
tags: []
touches: [wip/kanban/stories/S-0770-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
formula: publish
---
# S-0770 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3297 Fetch the remote's branch and tags
- T-3298 Evaluate the release policy, stopping unless it is met
- T-3299 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
