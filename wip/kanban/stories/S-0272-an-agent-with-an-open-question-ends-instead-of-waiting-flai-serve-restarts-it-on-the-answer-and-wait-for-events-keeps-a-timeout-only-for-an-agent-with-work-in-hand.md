---
id: S-0272
type: story
nature: improvement
title: "An agent with an open question ends instead of waiting: flai serve restarts it on the answer, and wait_for_events keeps a timeout only for an agent with work in hand"
status: backlog
parent: E-0017
owner: alex
created: 2026-10-05T01:35:31Z
updated: 2026-10-05T01:35:31Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0272 An agent with an open question ends instead of waiting: flai serve restarts it on the answer, and wait_for_events keeps a timeout only for an agent with work in hand

## Goal

650 turns across 38 story runs woke from `wait_for_events` to find nothing, each a full-context model call. flai serve already restarts a story agent when its thread is answered, so an agent whose only pending work is a question has nothing to wait for. The convention and the harness prompt say to end the session on an open question after writing the narrative's state; `wait_for_events` answers at once with `end: true` when the caller's story has an open thread and no task in progress, and keeps its timeout for an agent that has work to go on with. The MCP server's instructions say the same.

## Acceptance criteria
- [ ] `wait_for_events` answers `end: true` and why when the calling agent's story has an open thread to the designer and no task in progress, and the harness prompt and conventions tell the agent to end then
- [ ] flai serve's restart on the answer is tested end to end: an agent that ended on a question is started again when the thread is answered, with the answer in its first inbox
- [ ] `design/system/flai-serve.md` (or where serve is described), the conventions, the template's copies, and the user guide describe the loop
- [ ] `flai stats` counts the empty wakes so that the saving is measured

## Tasks

## Notes

From the epic's log classification: MCP wait_for_events 650 turns, 1,235 minutes, 38 stories.
