---
id: S-0704
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: ready
owner: alex
created: 2026-11-10T08:56:47Z
updated: 2026-11-10T08:56:57Z
transitions:
  - to: ready
    at: 2026-11-10T08:56:57Z
    by: system-flow
tags: []
touches: [wip/kanban/stories/S-0704-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
formula: publish
---
# S-0704 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-2952 Fetch the remote's branch and tags
- T-2953 Evaluate the release policy, stopping unless it is met
- T-2954 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
