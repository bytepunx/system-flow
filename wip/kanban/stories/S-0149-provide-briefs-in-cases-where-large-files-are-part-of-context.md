---
id: S-0149
type: story
nature: improvement
title: Provide briefs in cases where large files are part of context
status: in-progress
parent: E-0010
owner: alex
created: 2026-09-29T04:57:00Z
updated: 2026-09-29T05:31:10Z
transitions:
  - to: ready
    at: 2026-09-29T04:58:25Z
    by: alex
  - to: in-progress
    at: 2026-09-29T05:24:04Z
    by: agent-S-0149
tags: [cli]
topics: [context]
touches: [flai/cmd, flai/internal/context, flai/internal/mcpserver/folder.go, design/adrs, design/system/agent-context.md, design/system/flai-cli.md, docs/users]
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
- T-0533 An ADR refines ADR-0049: a large document named by a path written out is briefed, not loaded whole
- T-0534 flai prime --story briefs a large document named by a path written out, and tells the agent to read it before changing it
- T-0535 Users' documentation and the prime help describe the rule, and the packs are measured again

## Notes
