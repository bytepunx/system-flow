---
id: S-0208
type: story
nature: feature
title: The planner is an agent flai serve starts for an epic or a story, behind the plan host action
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:13Z
updated: 2026-10-02T11:54:40Z
transitions: []
tags: [flai, dashboard]
touches: [flai/internal/serve, flai/internal/harness, flai/internal/hostapi, flai/cmd, flai/internal/mcpserver, flaiover/src, ".claude/agents", template/]
after: [S-0199, S-0206, S-0207]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0208 The planner is an agent flai serve starts for an epic or a story, behind the plan host action

## Goal

The planner is an agent flai starts on demand for an epic or a story: it drafts stories from an epic, enriches stories with touches, forecasts, and cost of delay, and revisits an epic's children when it plans the epic. Like story agents it has a harness, model, and roles, and runs in the main checkout, writing work items only.

## Acceptance criteria
- [ ] A `plan` host action, off by default, gates it; the project's `agent` default (and a `planning.agent` override) gives its harness, model, and config
- [ ] `flai plan <E-nnnn|S-nnnn>`, the hostapi `plan.run`, and the MCP tool `plan` start a planner run for the item; `flai serve` records the run as it records story runs (log, outcome, session), and refuses a second run on the same item while one runs
- [ ] `harness.Prompt` has a planner prompt: prime with `--role plan`, read the item and what it links, do the planning the item's state calls for (drafting for an epic with no stories; enriching for a story; both, revisiting every open child, for an epic with stories), write through flai, log the activity, and end; ask the operator with `thread_open` on the item when inputs are missing
- [ ] The template ships `.claude/agents/planner.md`; the claude-code adapter passes it as it passes the explorer and verifier (ADR-0065)
- [ ] The planner's flai calls are guarded (`flai guard`): it may create and edit items and open threads, and may not move an item past backlog, accept, publish, or edit code
- [ ] `design/system/flai-cli.md` and a new `design/system/strategic-agents.md` describe it; the user guide says how to run it; tests cover the start, the refusal, and the guard

## Tasks

## Notes
