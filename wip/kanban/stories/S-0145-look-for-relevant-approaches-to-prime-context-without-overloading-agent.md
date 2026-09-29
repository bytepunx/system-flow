---
id: S-0145
type: story
nature: research
title: Look for relevant approaches to prime context without overloading agent
status: ready
parent: E-0010
owner: alex
created: 2026-09-29T01:57:32Z
updated: 2026-09-29T01:57:54Z
transitions:
  - to: ready
    at: 2026-09-29T01:57:54Z
    by: alex
tags: [cli]
touches: [flai/cmd]
agent:
  harness: claude-code
  model: claude-fable-5-1
  config:
    effort: high
---
# S-0145 Look for relevant approaches to prime context without overloading agent

## Goal

In order to prime agents with the relevant context for the story, we have to pull a lot of text from disparate sources. Even with some improvements to filtering context can be well over 537 KB (several hundred tokens). The agents can't load that much context in a single run.

Let's do some additional research on different ways to make sure only relevant documentation is fed to the agent. Make sure to include all approaches (RAG, MCP, Memory, etc.) with details about their trade-offs. The output should be an ADR draft.

## Acceptance criteria
- [ ] Multiple potential approaches identified and analyzed
- [ ] Recommendation for approach made as apart of the ADR draft
- [ ] The approach has the potential to significantly reduce the number of tokens fed to the agent

## Tasks

## Notes
