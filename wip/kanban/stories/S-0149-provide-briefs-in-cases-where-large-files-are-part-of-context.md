---
id: S-0149
type: story
nature: improvement
title: Provide briefs in cases where large files are part of context
status: backlog
parent: E-0010
owner: alex
created: 2026-09-29T04:57:00Z
updated: 2026-09-29T04:57:00Z
transitions: []
tags: [cli]
topics: [context]
touches: [flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0149 Provide briefs in cases where large files are part of context

## Goal

Write a follow-up story under E-0010. A path written out in plain text (not a markdown link, not an ADR ID) would be briefed, not loaded whole. A link, a #fragment link, and an ADR ID would still load. That refines ADR-0049's "named loads whole", so it probably needs a superseding ADR.

## Acceptance criteria
- [ ] Large referenced files get plain text and get a brief instead of full load by default
- [ ] Agents should determine when the full file is needed in context based on its read of the brief and read them

## Tasks

## Notes
