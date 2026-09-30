---
id: S-0173
type: story
nature: remediation
title: Inbox Open Narratives Navigate To The Wrong Location
status: ready
owner: alex
created: 2026-09-30T01:13:29Z
updated: 2026-09-30T01:13:29Z
transitions:
  - to: ready
    at: 2026-09-30T01:13:29Z
    by: alex
tags: [dashboard]
topics: [client-side]
touches: [flaiover/src]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0173 Inbox Open Narratives Navigate To The Wrong Location

## Goal

Clicking an open narrative should navigate to the related story and the open thread with the question. It should never take the operator to the docs/wip/agents/[agent] markdown file since the operator can't see or answer threads from there.

## Acceptance criteria
- [ ] Clicking any item in the items list should always navigate to the related item (task, story, or epic) and take the operator to the correct thread or part of the story
- [ ] There are no item types in the items list that can ever take the operator to documentation

## Tasks

## Notes
