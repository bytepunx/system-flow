---
id: S-0146
type: story
nature: feature
title: A story's context pack fits a size budget, with briefs of what its topics and links select
status: done
parent: E-0010
owner: alex
created: 2026-09-29T03:08:15Z
updated: 2026-09-29T04:58:02Z
transitions:
  - to: ready
    at: 2026-09-29T03:16:03Z
    by: alex
  - to: in-progress
    at: 2026-09-29T03:17:47Z
    by: agent-S-0146
  - to: review
    at: 2026-09-29T03:56:21Z
    by: agent-S-0146
  - to: done
    at: 2026-09-29T04:58:02Z
    by: alex
tags: [cli]
touches: [flai/internal/context, flai/internal/search, flai/cmd/prime.go, design/system/agent-context.md, design/system/flai-cli.md, docs/users, flai/cmd/adr_test.go, flai/cmd/prime_test.go, flai/internal/check, flai/internal/manifest, flai/internal/mcpserver, docs/operators/settings.md, design/system/project-manifest.md]
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
- [x] `flai prime --story` takes `--budget` (default 80 KB); a project sets its default under `prime.budget` in the config; the header gives the budget and each item's size, and says when conventions or named documents alone exceed it.
- [x] Topic-selected design and tech files load as title, first paragraph, and heading outline; ADRs reached one step away load as ID, title, and the first sentence under `## Decision`; each brief says how to fetch the body.
- [x] Ranked sections are cut at their own heading and fill the budget in rank order.
- [x] `flai check` warns about an ADR whose `## Decision` does not open with a sentence.
- [x] `flai prime --story S-0138` on this repository is measured and the measurement is in the story notes (over 80 KB: the designer chose on TH-0032 to keep every brief and let the header say so; the criterion first read "measures under 80 KB"); `--json` and the MCP `prime` tool carry the same pack.
- [x] `design/system/agent-context.md`, `flai-cli.md`, `docs/users/flai.md`, and the reference say what changed; all three test tiers and `flai check --strict` pass.

## Tasks
- T-0522 A size budget on the pack, with sizes in the header
- T-0523 Briefs for topic-selected design and stepped ADRs, and the decision-sentence check
- T-0524 Ranked sections fill the budget in rank order
- T-0525 Measure S-0138, and update the MCP tool and the documents

## Notes

First of the three stories that build ADR-0049, from TH-0030. The story after this one and the `doc_search` story switch agents to the pack.

Measured on 2026-09-29 on `story/S-0146` (3dba494), `flai prime --story S-0138`, header included:

| Part | Size |
|------|------|
| Conventions, every section `[all]` | 48.0 KB |
| Open-issues table | 9.1 KB |
| Named: ADR-0047, ADR-0049, `design/system/conventions.md` | 36.3 KB of text |
| Briefs: 22 design and tech files, 33 ADR decision lines | 14.1 KB of text, 23 KB printed |
| Catalog and header | about 5 KB |
| Pack | 125,230 bytes text, 133,019 bytes `--json`; `exceeded: named`; nothing ranked |

ADR-0049's 80 KB estimate left out the issues table and the header, and S-0138 was re-scoped after it to name 36 KB instead of 9 KB. Conventions, issues, and briefs alone come to about 85 KB for a cli story. TH-0032 asked the designer: A (briefs yield), B (a 128 KB default), or C (keep every brief, over budget until the conventions narrow). The designer chose C, provided the agent reads a full document when a brief shows it matters. The pack's `briefs` and `decisions` headings say so; S-0148 carries the same into the prompt, `CLAUDE.md`, and `session-start.md`.

The same measure on other stories: S-0141 160 KB (names `flaiover-dashboard.md`, 72 KB, through T-0504), S-0146 207 KB (names `flai-cli.md`, 86 KB, through T-0525), and S-0147 189 KB. Each brief carries its path; how to fetch a body is said once, at the head of the `briefs` and `decisions` groups, rather than on every brief, which kept the S-0138 briefs to 23 KB from 28 KB.

A path written out as a file to update loads the whole file. TH-0032 proposes a follow-up that briefs a path written out in plain text; it awaits the designer.
