---
id: S-0818
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: ready
owner: alex
created: 2026-11-17T02:27:29Z
updated: 2026-11-17T02:27:30Z
transitions:
  - to: ready
    at: 2026-11-17T02:27:30Z
    by: system-flow
tags: []
touches: [wip/kanban/stories/S-0818-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
formula: publish
---
# S-0818 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3492 Fetch the remote's branch and tags
- T-3493 Evaluate the release policy, stopping unless it is met
- T-3494 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
