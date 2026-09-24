---
id: S-0018
type: story
nature: feature
title: Template guide
status: done
parent: E-0004
owner: agent
created: 2026-09-15T16:09:00Z
updated: 2026-09-24T08:20:08Z
transitions:
  - to: ready
    at: 2026-09-24T07:58:34Z
    by: alex
  - to: in-progress
    at: 2026-09-24T08:12:04Z
    by: agent-S-0018
  - to: review
    at: 2026-09-24T08:19:48Z
    by: agent-S-0018
  - to: done
    at: 2026-09-24T08:20:08Z
    by: alex
tags: []
touches: [docs/contributors, design/system/template.md, docs/README.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0018 Template guide

## Goal
docs/contributors explains how to fork, edit, and version the template.

## Acceptance criteria
- [x] Covers template.yaml, variables, render rules, and testing

## Tasks
- T-0403 Write the template guide in docs/contributors
- T-0404 Link the guide and correct stale template text in the contributors index and template design
- T-0405 Verify the guide's examples against the tree's flai and lint

## Notes
The guide is `docs/contributors/template.md`. Every example in it was run against a scratch fork of `template/` with the tree's flai (T-0405). Running them found I-0040 and I-0041, linked from the guide as known limits.
