---
id: S-0594
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: ready
owner: alex
created: 2026-11-07T21:47:12Z
updated: 2026-11-07T21:47:37Z
transitions:
  - to: ready
    at: 2026-11-07T21:47:37Z
    by: system-flow
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
formula: publish
---
# S-0594 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-2589 Fetch the remote's branch and tags
- T-2590 Evaluate the release policy, stopping unless it is met
- T-2591 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with no inputs.
