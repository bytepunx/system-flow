---
id: S-0173
type: story
nature: remediation
title: Inbox Open Narratives Navigate To The Wrong Location
status: in-progress
owner: alex
created: 2026-09-30T01:13:29Z
updated: 2026-09-30T01:26:59Z
transitions:
  - to: ready
    at: 2026-09-30T01:13:29Z
    by: alex
  - to: in-progress
    at: 2026-09-30T01:25:47Z
    by: agent-S-0173
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
- T-0610 An open question in a narrative leads to its story's page, shown and answerable there
- T-0611 A thread in the inbox never leads to a document it could be answered from elsewhere
- T-0612 The design and the user guide say where each inbox entry leads

## Notes
