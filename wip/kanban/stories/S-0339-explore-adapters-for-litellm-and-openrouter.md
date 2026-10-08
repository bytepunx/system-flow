---
id: S-0339
type: story
nature: research
title: Explore Adapters for LiteLLM and OpenRouter
status: backlog
owner: alex
created: 2026-10-08T07:17:16Z
updated: 2026-10-08T07:17:16Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-fable-5-1
  config:
    effort: medium
cost_of_delay:
  inputs:
    penalty_per_week: 500
    by: alex
    at: 2026-10-08T07:17:16Z
---
# S-0339 Explore Adapters for LiteLLM and OpenRouter

## Goal

Perform the research necessary to write a new epic to implement agent adapters for both LiteLLM and OpenRouter while exploring the necessary changes and abstractions to capture intent free from vendor specifics (Anthropic SDK API vs. OpenAI's).

## Acceptance criteria
- [ ] A new epic to create agent adapters for LiteLLM and OpenRouter
- [ ] New ADRs where necessary to capture changes necessary to move flai away from being Anthropic specific
- [ ] Documentation capturing the abstractions and changes needed to expand support for other vendor's agents that can inform the planner that will write the epic's stories

## Tasks

## Notes
