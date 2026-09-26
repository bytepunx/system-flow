---
id: S-0135
type: story
nature: feature
title: Stories and epics carry topics, and flai works out a story's topics from them, its epic, and the projects its tags and claim reach
status: backlog
parent: E-0010
owner: alex
created: 2026-09-26T17:46:02Z
updated: 2026-09-26T17:46:02Z
transitions: []
tags: [cli, dashboard]
touches: [flai/internal/workitem, flai/cmd, flai/internal/mcpserver, flai/internal/hostapi, flaiover/src, design/system/work-hierarchy.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0135 Stories and epics carry topics, and flai works out a story's topics from them, its epic, and the projects its tags and claim reach

## Goal

A story or epic can name topics that no component names (`logging`, `release`), and flai can say what a story's topics are, as [ADR-0047](../../../design/adrs/0047-an-agent-is-primed-with-what-its-story-s-topics-claim-and-links-select.md) decides: its own and its epic's `topics`, the name, tags, and kind of every sub-project its `tags` or its claim (`touches` and its open tasks', as ADR-0046 defines it) reach, `code` when one of those is not the template, and `all`. `tags` keep meaning the component a story delivers to.

## Acceptance criteria
- [ ] Stories and epics accept an optional `topics:` front matter list, set with `--topics` on `flai story new` and `flai epic new`, with `flai edit`, and with the MCP `item_new` and `item_edit`; strict parsing accepts it and the round-trip test stays green.
- [ ] The dashboard's story and epic pages show and edit topics next to tags, through the host as other edits do.
- [ ] A function returns a story's topics with where each came from (own, epic, sub-project via tag, sub-project via claim, code, all), with behavior tests covering each source, a story with no tags or touches, and a touch outside every sub-project.
- [ ] `flai show S-nnnn` and `--json` print the story's topics and their sources.
- [ ] `design/system/work-hierarchy.md` documents the key; `docs/users/flai.md` and `docs/users/flaiover.md` document setting it; all three test tiers and `flai check --strict` pass.

## Tasks

## Notes

After S-0134, since both define the vocabulary `flai check` compares. Older flai refuses items with a key it does not know; say so in the release notes.
