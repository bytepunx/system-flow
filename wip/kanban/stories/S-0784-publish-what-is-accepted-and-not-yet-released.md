---
id: S-0784
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: ready
owner: alex
created: 2026-11-14T07:05:35Z
updated: 2026-11-14T07:05:37Z
transitions:
  - to: ready
    at: 2026-11-14T07:05:37Z
    by: system-flow
tags: []
touches: [wip/kanban/stories/S-0784-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
formula: publish
---
# S-0784 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3351 Fetch the remote's branch and tags
- T-3352 Evaluate the release policy, stopping unless it is met
- T-3353 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
