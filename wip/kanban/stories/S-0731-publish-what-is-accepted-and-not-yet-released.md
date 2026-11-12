---
id: S-0731
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-12T03:44:58Z
updated: 2026-11-12T03:45:07Z
transitions:
  - to: ready
    at: 2026-11-12T03:44:59Z
    by: system-flow
  - to: in-progress
    at: 2026-11-12T03:45:07Z
    by: agent-S-0731
tags: []
touches: [wip/kanban/stories/S-0731-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
formula: publish
---
# S-0731 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3135 Fetch the remote's branch and tags
- T-3136 Evaluate the release policy, stopping unless it is met
- T-3137 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
