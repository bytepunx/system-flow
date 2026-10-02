---
id: T-0678
type: task
nature: experiment
title: The conventions and the claude-code prompt tell the story's agent to plan its tasks and work independent ones with sub-agents
status: done
parent: S-0176
owner: arobson
created: 2026-10-01T11:42:10Z
updated: 2026-10-01T11:52:00Z
transitions:
  - to: ready
    at: 2026-10-01T11:42:32Z
    by: agent-S-0176
  - to: in-progress
    at: 2026-10-01T11:42:32Z
    by: agent-S-0176
  - to: done
    at: 2026-10-01T11:52:00Z
    by: agent-S-0176
stream: S-0176
tags: []
touches: [flai/internal/harness, template/root/design/conventions, design/conventions, template/CHANGELOG.md, template/template.yaml]
usage:
  source: log
  seconds: 568
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 165
      output: 53024
      cache_read: 11717851
      cache_write: 167207
      cost: 4.4401
---
# T-0678 The conventions and the claude-code prompt tell the story's agent to plan its tasks and work independent ones with sub-agents

## Work

- `work-management.md` and `delegation.md`, in `template/root/design/conventions/` first and then this repository above the marker: when the agent writes a story's tasks it plans them. It says what each touches, which waits for which and why (`after`), and which can run together (no `after` between them and no overlap in `touches`), and records the plan's reasoning in the narrative's `## Decisions`.
- The story's agent runs the ready, non-overlapping tasks of a layer at once, one sub-agent each (a fork where the harness has them), waits, reviews their work, and moves each task. Only it commits, syncs the stream, moves items, and talks to the designer; a task sub-agent's question comes back in its final message (ADR-0059).
- `harness.Prompt` for `claude-code` says the same in brief.

## Done when

- The template and this repository's conventions say it, the template's version and changelog are bumped, and `go test -race -short ./internal/harness/...` passes with a test that the prompt asks for the plan.

## Notes
