---
id: S-0149
type: story
nature: improvement
title: Provide briefs in cases where large files are part of context
status: done
parent: E-0010
owner: alex
created: 2026-09-29T04:57:00Z
updated: 2026-09-29T05:42:53Z
transitions:
  - to: ready
    at: 2026-09-29T04:58:25Z
    by: alex
  - to: in-progress
    at: 2026-09-29T05:24:04Z
    by: agent-S-0149
  - to: review
    at: 2026-09-29T05:35:43Z
    by: agent-S-0149
  - to: done
    at: 2026-09-29T05:42:53Z
    by: alex
tags: [cli]
topics: [context]
touches: [flai/cmd, flai/internal/context, flai/internal/mcpserver/folder.go, design/adrs, design/system/agent-context.md, design/system/flai-cli.md, docs/users, CLAUDE.md, design/conventions/session-start.md, template]
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
- [x] Large referenced files get plain text and get a brief instead of full load by default
- [x] Agents should determine when the full file is needed in context based on its read of the brief and read them

## Tasks
- T-0533 An ADR refines ADR-0049: a large document named by a path written out is briefed, not loaded whole
- T-0534 flai prime --story briefs a large document named by a path written out, and tells the agent to read it before changing it
- T-0535 Users' documentation and the prime help describe the rule, and the packs are measured again

## Notes

- Decision: [ADR-0050](../../../design/adrs/0050-a-document-a-story-names-only-by-its-path-written-out-is-briefed-when-it-is.md), refining ADR-0049 on the designer's word in TH-0032. A document named only by its path written out is briefed when it is larger than an eighth of the budget (10 KB at 80 KB); a link, an ADR ID, and a small file still load whole, and a `#fragment` link still loads its section.
- Criterion 1 is verified by the behaviour tests in `flai/internal/context` (small and large written-out path, link and ID winning over a path, fragment link, ADR by path, no step from a size brief, the printed pack) and by the measurements below.
- Criterion 2: flai's part is what the pack tells the agent. Each such brief is followed by a line saying the story names the file, that it is briefed for its size, and that the agent decides from the brief whether it needs the body and reads the file, or the sections it will change, before relying on it or changing it; the `briefs` heading says the same. Verified in the printed pack of S-0141. The agent-facing prompt, `CLAUDE.md`, and `session-start.md` wording is S-0148's.
- Packs at 80 KB, before and after (the same numbers are in `agent-context.md` § Paths written out):

| Story | Before | After | Briefed for size |
|-------|--------|-------|------------------|
| S-0138 | 123 KB | 113 KB | `conventions.md` |
| S-0141 | 156 KB | 87 KB | `flaiover-dashboard.md` |
| S-0146 | 230 KB | 113 KB | `agent-context.md`, `conventions.md`, `flai-cli.md` |
| S-0147 | 188 KB | 102 KB | `flai-cli.md` |
