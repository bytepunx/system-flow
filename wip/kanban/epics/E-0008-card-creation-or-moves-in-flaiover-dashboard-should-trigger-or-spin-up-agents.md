---
id: E-0008
type: epic
nature: feature
title: Card creation or moves in flaiover dashboard should trigger or spin up agents
status: backlog
owner: alex
created: 2026-09-23T16:43:01Z
updated: 2026-09-23T16:46:55Z
transitions: []
tags: [dashboard, cli]
touches: [flaiover/src, flai/cmd]
---
# E-0008 Card creation or moves in flaiover dashboard should trigger or spin up agents

## Outcome

Instead of tell an agent like Claude to use the MCP endpoint to poll for ready items, moving or creating cards in to "ready" should spin up the correct configured agent. Between system flow yaml and front-matter on cards, the user should be able to specify the "harness" (the command line or API interface) and model (sonnet 5, opus 5.5, fable 5.1, etc.) that should take on the story and work on it.

These agents need to "prime" context via ADRs and conventions for the story and make use of the activity log, inbox, flai via MCP to submit questions for the user, create and work tasks, and monitor the inbox for responses so that they can successfully move the story to the review stage.

## Stories
- S-0103 Extend existing data structures, APIs, and UI to track agent configuration

## Notes
