---
id: S-0146
type: story
nature: feature
title: A story's context pack fits a size budget, with briefs of what its topics and links select
status: ready
parent: E-0010
owner: alex
created: 2026-09-29T03:08:15Z
updated: 2026-09-29T03:16:03Z
transitions:
  - to: ready
    at: 2026-09-29T03:16:03Z
    by: alex
tags: [cli]
touches: [flai/internal/context, flai/internal/search, flai/cmd/prime.go, design/system/agent-context.md, design/system/flai-cli.md, docs/users]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0146 A story's context pack fits a size budget, with briefs of what its topics and links select

## Goal

The context pack fits a size budget, as [ADR-0049](../../../design/adrs/0049-a-story-s-context-pack-fits-a-size-budget-what-the-story-names-loads-whole-what.md) decides, so that an agent can load it in one call. Conventions load by topic and are never cut. What the story, its epic, and its tasks name loads whole. Design and tech selected by topics load as briefs: title, first paragraph, and heading outline. An ADR reached by one step loads as its decision sentence. Ranked sections, each cut at its own heading, fill what is left.

## Acceptance criteria
- [ ] `flai prime --story` takes `--budget` (default 80 KB); a project sets its default under `prime.budget` in the config; the header gives the budget and each item's size, and says when conventions or named documents alone exceed it.
- [ ] Topic-selected design and tech files load as title, first paragraph, and heading outline; ADRs reached one step away load as ID, title, and the first sentence under `## Decision`; each brief says how to fetch the body.
- [ ] Ranked sections are cut at their own heading and fill the budget in rank order.
- [ ] `flai check` warns about an ADR whose `## Decision` does not open with a sentence.
- [ ] `flai prime --story S-0138` on this repository measures under 80 KB; the measurement is in the story notes; `--json` and the MCP `prime` tool carry the same pack.
- [ ] `design/system/agent-context.md`, `flai-cli.md`, `docs/users/flai.md`, and the reference say what changed; all three test tiers and `flai check --strict` pass.

## Tasks

## Notes

First of the three stories that build ADR-0049, from TH-0030. The story after this one and the `doc_search` story switch agents to the pack.
