---
id: S-0138
type: story
nature: feature
title: Agents flai starts, and any agent with a story, prime with flai prime --story
status: backlog
parent: E-0010
owner: alex
created: 2026-09-26T17:46:25Z
updated: 2026-09-26T17:46:25Z
transitions: []
tags: [cli, template]
touches: [flai/internal/harness, flai/internal/mcpserver, CLAUDE.md, template/root/CLAUDE.md.tmpl, template/root/design/conventions/session-start.md, design/conventions/session-start.md, docs/users]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0138 Agents flai starts, and any agent with a story, prime with flai prime --story

## Goal

Agents use the context pack as soon as it exists (TH-0021 5b), as [ADR-0047](../../../design/adrs/0047-an-agent-is-primed-with-what-its-story-s-topics-claim-and-links-select.md) decides. The prompt flai serve gives a story's agent says to prime with `flai prime --story <id>`; `CLAUDE.md`, the template's `CLAUDE.md.tmpl`, and the baseline `session-start.md` tell any agent with a story to do the same and to read what the catalog lists when it needs it; the MCP server gains a `prime` tool that returns the pack for a story.

## Acceptance criteria
- [ ] `harness.Prompt` names `flai prime --story <id>`; its test pins the new text.
- [ ] The MCP `prime` tool takes a story ID and returns the pack as `flai prime --story --json` does; the tool list tests are updated; the server's instructions mention it.
- [ ] `CLAUDE.md`, `template/root/CLAUDE.md.tmpl`, and `session-start.md` in the template and here say to prime with `--story` when there is a story and `flai prime --cat` otherwise; `design/system/conventions.md` "Priming" says the same and drops the not-yet-built note; the template's version and changelog record it.
- [ ] `docs/users/flai.md` documents the MCP tool; all three test tiers, `make smoke` on a rendered template, and `flai check --strict` pass.

## Tasks

## Notes

Last of the E-0010 stories from S-0125, after S-0137. A project whose installed flai is older than this release keeps priming with `--cat`; the MCP tool appears only when `flai mcp` is upgraded.
