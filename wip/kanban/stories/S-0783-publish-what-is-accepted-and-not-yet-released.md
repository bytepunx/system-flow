---
id: S-0783
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-13T22:47:23Z
updated: 2026-11-13T22:47:32Z
transitions:
  - to: ready
    at: 2026-11-13T22:47:24Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T22:47:32Z
    by: agent-S-0783
tags: []
touches: [wip/kanban/stories/S-0783-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
formula: publish
---
# S-0783 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3340 Fetch the remote's branch and tags
- T-3341 Evaluate the release policy, stopping unless it is met
- T-3342 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
