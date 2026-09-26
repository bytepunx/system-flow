---
id: S-0136
type: story
nature: feature
title: flai prime --story prints the conventions a story's topics select, section by section, and lists what it left out
status: backlog
parent: E-0010
owner: alex
created: 2026-09-26T17:46:25Z
updated: 2026-09-26T17:46:25Z
transitions: []
tags: [cli]
touches: [flai/cmd/prime.go, flai/internal/conventions, flai/internal/context]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0136 flai prime --story prints the conventions a story's topics select, section by section, and lists what it left out

## Goal

An agent with a story reads the conventions that apply to it, as [ADR-0047](../../../design/adrs/0047-an-agent-is-primed-with-what-its-story-s-topics-claim-and-links-select.md) decides. `flai prime --story S-nnnn` prints `README.md` and each convention in read order with the sections whose topics do not match the story's topics left out, keeps the baseline marker and the `## Project additions` heading, follows with the open-issues table, and ends with one line per section left out (`code-quality.md § Go (cli)`). Plain `flai prime` is unchanged.

## Acceptance criteria
- [ ] `flai prime --story S-nnnn` prints what the goal says; a section is kept when its effective topics include `all` or any of the story's topics.
- [ ] `--json` returns, per convention, the sections kept and left out with their topics, and the story's topics with their sources.
- [ ] A header states the story, its topics, and the size of what was printed.
- [ ] With every convention at `[all]` the output is the same convention text as `flai prime --cat`; a behavior test with narrowed fixture conventions shows sections dropped and listed.
- [ ] An unknown or archived story ID fails with a message naming it.
- [ ] `design/system/flai-cli.md` and `docs/users/flai.md` describe the flag; the reference is regenerated; all three test tiers and `flai check --strict` pass.

## Tasks

## Notes

After S-0134 and S-0135. The design, tech, ADR, linked, and ranked parts of the pack are the next story; this one delivers the conventions filter on its own so it can be used as soon as the designer narrows topics.
