---
id: S-0137
type: story
nature: feature
title: flai prime --story adds the design, tech, and ADRs a story's topics, links, and ranking select, with a catalog of the rest
status: backlog
parent: E-0010
owner: alex
created: 2026-09-26T17:46:25Z
updated: 2026-09-26T17:46:25Z
transitions: []
tags: [cli]
touches: [flai/cmd/prime.go, flai/internal/context, flai/internal/search]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0137 flai prime --story adds the design, tech, and ADRs a story's topics, links, and ranking select, with a catalog of the rest

## Goal

An agent starts its story with the design and decisions it builds on, and knows what else exists, as [ADR-0047](../../../design/adrs/0047-an-agent-is-primed-with-what-its-story-s-topics-claim-and-links-select.md) decides. After the conventions, `flai prime --story` prints, each headed with its path, heading path, and reason: the `design/system`, `design/tech`, and ADR sections whose topics match; every document or ADR the story, its epic, and its tasks link or name by ID, and one step further (ADRs a selected section links, ADRs a selected ADR refines), with a superseded ADR replaced by what supersedes it; the five ADRs and five design sections ranked highest by BM25 against the title, goal, and criteria among those not selected; then a catalog of every document not loaded and the heading outline of each loaded in part.

## Acceptance criteria
- [ ] Each selection step is a function in the context package with behavior tests over a fixture repository: topics, links from story, epic, and tasks, one step of links, supersession, ranking, catalog, and no document printed twice (the first reason wins, the rest are listed).
- [ ] Reasons read `topics: cli`, `linked from S-nnnn`, `linked from design/system/x.md § Heading`, `refined by ADR-nnnn`, `supersedes ADR-nnnn`, `rank n`.
- [ ] `--json` returns every item with path, heading path, reason, and size; the header's size includes them.
- [ ] Ranking reuses `flai/internal/search`, indexing sections rather than whole files, without a new dependency.
- [ ] Run on this repository for S-0125 and for two archived stories, the output is attached to the story notes with its size, and every ADR those stories named that existed when they were created is either printed or in the catalog.
- [ ] `design/system/flai-cli.md`, `design/system/agent-context.md`, and `docs/users/flai.md` describe the pack; all three test tiers and `flai check --strict` pass.

## Tasks

## Notes

After S-0136. No size budget (TH-0021 4b); the ranked step's fixed count of five and five is from the S-0125 replay. Sections are cut at headings only.
