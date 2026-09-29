---
id: S-0148
type: story
nature: feature
title: Agents flai starts, and any agent with a story, prime with the budgeted pack and fetch what it briefs
status: ready
parent: E-0010
owner: alex
created: 2026-09-29T03:08:25Z
updated: 2026-09-29T03:16:07Z
transitions:
  - to: ready
    at: 2026-09-29T03:16:07Z
    by: alex
tags: [cli, template]
touches: [flai/internal/harness, flai/internal/mcpserver, CLAUDE.md, template/root/CLAUDE.md.tmpl, template/root/design/conventions, design/conventions, design/system/conventions.md, design/system/agent-context.md, template/template.yaml, template/CHANGELOG.md, docs/users]
after: [S-0146, S-0147]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0148 Agents flai starts, and any agent with a story, prime with the budgeted pack and fetch what it briefs

## Goal

Agents use the budgeted context pack once S-0146 and S-0147 have landed, as [ADR-0049](../../../design/adrs/0049-a-story-s-context-pack-fits-a-size-budget-what-the-story-names-loads-whole-what.md) decides. This is the switch S-0138 was written for, deferred by TH-0030 because the pack was 537 KB. The prompt flai serve gives a story's agent, `CLAUDE.md`, the template's `CLAUDE.md.tmpl`, and the baseline `session-start.md` say to prime with `flai prime --story <id>` (or the MCP `prime` tool) when there is a story and `flai prime --cat` otherwise. They say that the pack is a brief, and that the agent reads a briefed document's body with `doc_get` and a heading before changing what it describes.

## Acceptance criteria
- [ ] `harness.Prompt` names `flai prime --story <id>` and says to fetch a briefed body with `doc_get` before changing what it describes; its test pins the text.
- [ ] `CLAUDE.md`, `template/root/CLAUDE.md.tmpl`, the conventions `README.md`, and `session-start.md` in the template and here say the same, `--cat` when there is no story; the template's version and changelog record it.
- [ ] The MCP servers' instructions tell an agent to call `prime` when it starts a story and `doc_search` or `doc_get` for what the pack briefs.
- [ ] `design/system/conventions.md` "Priming" and `agent-context.md` say agents prime with the pack; `docs/users/flai.md` says the same; all three test tiers, `make smoke` on a rendered template, and `flai check --strict` pass.

## Tasks

## Notes

Last of the three stories that build ADR-0049, from TH-0030. S-0138's commits "the prompt flai serve gives an agent primes with flai prime --story" and "CLAUDE.md and session-start.md prime with --story when there is a story", reverted there, hold a first draft of the wording, without ADR-0049's brief and fetch.

From S-0146 (TH-0032, 2026-09-29): the designer kept 80 KB and every brief, so packs of code stories are over budget until the conventions narrow, provided the agent reads a full document when its brief shows it bears on the story. The prompt, `CLAUDE.md`, and `session-start.md` must say that, as the pack's `briefs` and `decisions` headings do.
