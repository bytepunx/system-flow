---
id: S-0271
type: story
nature: improvement
title: "Criteria and narrative state are commands: flai story tick checks a criterion and flai stream state writes Current state and Next steps"
status: backlog
parent: E-0017
owner: alex
created: 2026-10-05T01:35:30Z
updated: 2026-10-05T01:35:30Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0271 Criteria and narrative state are commands: flai story tick checks a criterion and flai stream state writes Current state and Next steps

## Goal

Agents tick acceptance criteria with `sed -i` on the story file (119 turns across 85 stories) and rewrite a narrative's `## Current state` and `## Next steps` with Edit or Write (111 turns across 44 stories), each a model turn that reads the file first. `flai story tick S-nnnn <n>` (and `--untick`) checks the nth criterion and refuses one that does not exist; `flai stream state S-nnnn --current "<text>" --next "<text>"` replaces the two sections, leaving the rest of the narrative alone, and appends nothing. Both are `item_tick` and `stream_state` over MCP and `item.tick` and `stream.state` on the host channel, so the dashboard's story page can tick a criterion too. `flai check` then knows the sections' shape.

## Acceptance criteria
- [ ] `flai story tick` and `flai stream state` exist with the behaviour above, as text and `--json`, and refuse a story that is not in progress or review
- [ ] The same operations exist over MCP and on the host channel, and the dashboard's story page ticks a criterion through them
- [ ] The conventions, the template's copies, the harness prompt, `design/system/flai-cli.md`, `agent-narrative.md`, and the user guide send the agent to them instead of editing the files

## Tasks

## Notes

From the epic's log classification. `flai task done` (its sibling story) may take `--tick <n>` once this exists.
