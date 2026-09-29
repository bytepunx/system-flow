---
id: S-0148
type: story
nature: feature
title: Agents flai starts, and any agent with a story, prime with the budgeted pack and fetch what it briefs
status: done
parent: E-0010
owner: alex
created: 2026-09-29T03:08:25Z
updated: 2026-09-29T05:31:49Z
transitions:
  - to: ready
    at: 2026-09-29T03:16:07Z
    by: alex
  - to: in-progress
    at: 2026-09-29T05:23:58Z
    by: agent-S-0148
  - to: review
    at: 2026-09-29T05:30:27Z
    by: agent-S-0148
  - to: done
    at: 2026-09-29T05:31:49Z
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
- [x] `harness.Prompt` names `flai prime --story <id>` and says to fetch a briefed body with `doc_get` before changing what it describes; its test pins the text.
- [x] `CLAUDE.md`, `template/root/CLAUDE.md.tmpl`, the conventions `README.md`, and `session-start.md` in the template and here say the same, `--cat` when there is no story; the template's version and changelog record it.
- [x] The MCP servers' instructions tell an agent to call `prime` when it starts a story and `doc_search` or `doc_get` for what the pack briefs.
- [x] `design/system/conventions.md` "Priming" and `agent-context.md` say agents prime with the pack; `docs/users/flai.md` says the same; all three test tiers, `make smoke` on a rendered template, and `flai check --strict` pass.

## Tasks
- T-0529 The prompt flai serve gives a story's agent primes with flai prime --story and fetches what the pack briefs
- T-0530 The MCP servers' instructions tell an agent to call prime when it starts a story and doc_search or doc_get for what it briefs
- T-0531 CLAUDE.md, the template, and session-start.md prime with the pack when there is a story and --cat when there is not
- T-0532 The living design and user docs say agents prime with the pack, and every tier passes

## Notes

Last of the three stories that build ADR-0049, from TH-0030. S-0138's commits "the prompt flai serve gives an agent primes with flai prime --story" and "CLAUDE.md and session-start.md prime with --story when there is a story", reverted there, hold a first draft of the wording, without ADR-0049's brief and fetch.

From S-0146 (TH-0032, 2026-09-29): the designer kept 80 KB and every brief, so packs of code stories are over budget until the conventions narrow, provided the agent reads a full document when its brief shows it bears on the story. The prompt, `CLAUDE.md`, and `session-start.md` must say that, as the pack's `briefs` and `decisions` headings do.

From the work (2026-09-29): the MCP servers' instructions share one sentence, `primeInstructions` in `flai/internal/mcpserver/folder.go`. S-0148's own pack is 116 KB against the 80 KB budget, over by what it names. `flai check --strict` passes except for four advisory `wip.overlap` warnings: S-0149, started first and in progress beside this story, widened its touches after S-0148 was opened to include `design/system/agent-context.md` and `docs/users`, which S-0148 had already claimed. S-0148 appends a section at the end of `agent-context.md` and changes one paragraph of `docs/users/flai.md` § Prime a session. A trial merge conflicted at the end of `agent-context.md` (TH-0034), where both appended; S-0148's paragraph moved up beside the one it follows, and the two merge cleanly.
