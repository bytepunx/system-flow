---
id: S-0138
type: story
nature: feature
title: The MCP prime tool returns a story's context pack, built by the same code as flai prime --story
status: done
parent: E-0010
owner: alex
created: 2026-09-26T17:46:25Z
updated: 2026-09-29T03:17:17Z
transitions:
  - to: ready
    at: 2026-09-26T22:40:39Z
    by: alex
  - to: in-progress
    at: 2026-09-29T01:07:54Z
    by: agent-S-0138
  - to: review
    at: 2026-09-29T03:15:33Z
    by: agent-S-0138
  - to: done
    at: 2026-09-29T03:17:17Z
    by: alex
tags: [cli]
touches: [flai/internal/mcpserver, flai/internal/context/story.go, flai/cmd/prime.go, design/system/conventions.md, design/system/agent-context.md, design/system/flai-cli.md, design/adrs, docs/users]
after: [S-0137]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0138 The MCP prime tool returns a story's context pack, built by the same code as flai prime --story

## Goal

Agents use the context pack as soon as it exists (TH-0021 5b), as [ADR-0047](../../../design/adrs/0047-an-agent-is-primed-with-what-its-story-s-topics-claim-and-links-select.md) decides. The MCP server gains a `prime` tool that returns the pack for a story, built by the same code as `flai prime --story`. Amended on TH-0030: at 537 KB the pack is too large to switch agents to, so [ADR-0049](../../../design/adrs/0049-a-story-s-context-pack-fits-a-size-budget-what-the-story-names-loads-whole-what.md) puts it on a budget, and the switch of the prompt, `CLAUDE.md`, the template, and `session-start.md` moves to S-0148, after S-0146 (the budget) and S-0147 (`doc_search` and `heading`).

## Acceptance criteria
- [x] `harness.Prompt`, `CLAUDE.md`, `template/root/CLAUDE.md.tmpl`, and `session-start.md` still prime with `flai prime --cat`; the switch is S-0148's, and the template does not change.
- [x] The MCP `prime` tool takes a story ID and returns the pack as `flai prime --story --json` does; the tool list tests are updated; the server's instructions do not yet tell agents to call it.
- [x] `design/system/conventions.md` "Priming" says why agents still prime with `--cat` and which stories switch them; ADR-0049 is accepted.
- [x] `docs/users/flai.md` documents the MCP tool; all three test tiers, `make smoke` on a rendered template, and `flai check --strict` pass.

## Tasks
- T-0510 The MCP prime tool returns a story's context pack
- T-0511 The prompt flai serve gives an agent primes with flai prime --story
- T-0512 CLAUDE.md and session-start.md prime with --story when there is a story
- T-0513 Document the MCP prime tool and run every check

## Notes

Last of the E-0010 stories from S-0125, after S-0137. A project whose installed flai is older than this release keeps priming with `--cat`; the MCP tool appears only when `flai mcp` is upgraded.
