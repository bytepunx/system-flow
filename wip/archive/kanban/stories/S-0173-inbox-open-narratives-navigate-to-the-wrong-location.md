---
id: S-0173
type: story
nature: remediation
title: Inbox Open Narratives Navigate To The Wrong Location
status: done
owner: alex
created: 2026-09-30T01:13:29Z
updated: 2026-10-01T07:48:04Z
transitions:
  - to: ready
    at: 2026-09-30T01:13:29Z
    by: alex
  - to: in-progress
    at: 2026-09-30T01:25:47Z
    by: agent-S-0173
  - to: ready
    at: 2026-10-01T07:31:51Z
    by: agent
  - to: in-progress
    at: 2026-10-01T07:34:04Z
    by: agent-S-0173
  - to: review
    at: 2026-10-01T07:47:10Z
    by: agent-S-0173
  - to: done
    at: 2026-10-01T07:48:04Z
    by: alex
tags: [dashboard]
topics: [client-side]
touches: [flaiover/src, flai/internal/check, design/system, docs, flai/internal/workitem/rules.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 836
  models:
    - model: claude-opus-5-5
      input: 206
      output: 54266
      cache_read: 13057027
      cache_write: 187002
      cost: 5.1936
---
# S-0173 Inbox Open Narratives Navigate To The Wrong Location

## Goal

Clicking an open narrative should navigate to the related story and the open thread with the question. It should never take the operator to the docs/wip/agents/[agent] markdown file since the operator can't see or answer threads from there.

## Acceptance criteria
- [x] Clicking any item in the items list should always navigate to the related item (task, story, or epic) and take the operator to the correct thread or part of the story
- [x] There are no item types in the items list that can ever take the operator to documentation

## Tasks
- T-0610 An open question in a narrative leads to its story's page, shown and answerable there
- T-0611 A thread in the inbox never leads to a document it could be answered from elsewhere
- T-0612 The design and the user guide say where each inbox entry leads
- T-0613 A Threads page under Workflow lists the open threads and opens on the one a link names
- T-0614 flai check accepts the back moves ADR-0055 allows

## Notes
- 2026-10-01: verified by flaiover's unit tests (hrefFor for every entry kind, the story page's open questions, the Threads page), not in a running dashboard. T-0613 adds the Threads page TH-0041 asked for; T-0614 fixes flai check rejecting the in-progress to ready move flai made when this story was restarted here.
- 2026-10-01T07:31:51Z: moved to ready: begun on another host whose branch never reached this one; restarting here (S-0177)
