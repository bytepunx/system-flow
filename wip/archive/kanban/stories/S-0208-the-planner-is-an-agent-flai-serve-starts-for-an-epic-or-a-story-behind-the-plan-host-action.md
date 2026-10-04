---
id: S-0208
type: story
nature: feature
title: The planner is an agent flai serve starts for an epic or a story, behind the plan host action
status: done
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:13Z
updated: 2026-10-04T03:59:15Z
transitions:
  - to: ready
    at: 2026-10-03T05:34:39Z
    by: alex
  - to: in-progress
    at: 2026-10-04T00:39:45Z
    by: system-flow
  - to: review
    at: 2026-10-04T03:52:47Z
    by: agent-S-0208
  - to: done
    at: 2026-10-04T03:59:15Z
    by: alex
tags: [flai, dashboard]
touches: [flai/internal/serve, flai/internal/harness, flai/internal/hostapi, flai/cmd, flai/internal/mcpserver, flai/internal/manifest, flai/internal/guard, flaiover/src, ".claude/agents", template, design/issues, design/system, design/adrs, docs/users, docs/operators, ".claude/settings.json"]
after: [S-0199, S-0206, S-0207]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3912
  models:
    - model: claude-haiku-4-5-20251001
      input: 306
      output: 12478
      cache_read: 2519495
      cache_write: 96166
      cost: 0.4349
    - model: claude-opus-5-5
      input: 936
      output: 325143
      cache_read: 63277997
      cache_write: 1327421
      cost: 27.0434
    - model: claude-sonnet-5
      input: 128
      output: 31471
      cache_read: 7818769
      cache_write: 203950
      cost: 2.3886
---
# S-0208 The planner is an agent flai serve starts for an epic or a story, behind the plan host action

## Goal

The planner is an agent flai starts on demand for an epic or a story: it drafts stories from an epic, enriches stories with touches, forecasts, and cost of delay, and revisits an epic's children when it plans the epic. Like story agents it has a harness, model, and roles, and runs in the main checkout, writing work items only.

## Acceptance criteria
- [x] A `plan` host action, off by default, gates it; the project's `agent` default (and a `planning.agent` override) gives its harness, model, and config
- [x] `flai plan <E-nnnn|S-nnnn>`, the hostapi `plan.run`, and the MCP tool `plan` start a planner run for the item; `flai serve` records the run as it records story runs (log, outcome, session), and refuses a second run on the same item while one runs
- [x] `harness.Prompt` has a planner prompt: prime with `--role plan`, read the item and what it links, do the planning the item's state calls for (drafting for an epic with no stories; enriching for a story; both, revisiting every open child, for an epic with stories), write through flai, log the activity, and end; ask the operator with `thread_open` on the item when inputs are missing
- [x] The template ships `.claude/agents/planner.md`; the claude-code adapter passes it as it passes the explorer and verifier (ADR-0065)
- [x] The planner's flai calls are guarded (`flai guard`): it may create and edit items and open threads, and may not move an item past backlog, accept, publish, or edit code
- [x] `design/system/flai-cli.md` and a new `design/system/strategic-agents.md` describe it; the user guide says how to run it; tests cover the start, the refusal, and the guard

## Tasks
- T-0784 The planner's agent, prompt, and definition: planning.agent, harness.Prompt for a planner run, and .claude/agents/planner.md passed by claude-code
- T-0785 flai guard holds a planner session to planning: items and threads through flai, no move past backlog, no accept, publish, or edit
- T-0786 flai serve runs the planner for an item behind the plan host action: flai plan, hostapi plan.run, the run recorded, a second run refused
- T-0787 The MCP tool plan starts a planner run for an epic or a story
- T-0788 The dashboard's epic and story pages offer Plan while the plan host action is on
- T-0789 Design, ADR, and user and operator docs for the planner: strategic-agents.md, flai-cli.md, project-manifest.md, the user guide, and settings
- T-0790 The template's settings run the guard on a planner session's Edit, Write, and NotebookEdit, in place of the adapter's --settings

## Notes

- The operator's `planner.md` (TH-0097) also describes planning a story's tasks and a task, which S-0255 delivers; until then `flai plan` refuses a task and the guard refuses the planner `flai task new`, so a planner that follows that text on a story with tasks meets a refusal.
- Close-out: `scripts/close-out.sh` passed lint and every test tier and stopped at `flai check --strict` on two findings outside this story (S-0173's unmerged branch, TH-0094 on the archived S-0206), I-0057's pattern, bumped; the verifier ran the remaining steps by hand and they passed.
