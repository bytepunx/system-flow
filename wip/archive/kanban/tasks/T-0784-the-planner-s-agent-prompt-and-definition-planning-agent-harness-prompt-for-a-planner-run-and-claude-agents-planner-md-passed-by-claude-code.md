---
id: T-0784
type: task
nature: feature
title: "The planner's agent, prompt, and definition: planning.agent, harness.Prompt for a planner run, and .claude/agents/planner.md passed by claude-code"
status: done
parent: S-0208
owner: alex
created: 2026-10-04T00:43:38Z
updated: 2026-10-04T00:52:30Z
transitions:
  - to: ready
    at: 2026-10-04T00:44:31Z
    by: agent-S-0208
  - to: in-progress
    at: 2026-10-04T00:44:33Z
    by: agent-S-0208
  - to: done
    at: 2026-10-04T00:52:30Z
    by: agent-S-0208
stream: S-0208
tags: []
touches: [flai/internal/manifest, flai/internal/harness, template/root/.claude/agents, ".claude/agents"]
usage:
  source: log
  seconds: 477
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 100
      output: 41938
      cache_read: 6055798
      cache_write: 161382
      cost: 2.9879
---
# T-0784 The planner's agent, prompt, and definition: planning.agent, harness.Prompt for a planner run, and .claude/agents/planner.md passed by claude-code

## Work

- `planning.agent` in the manifest: the same shape as `agent` (harness, model, config, roles), checked as `agent` is, with a function giving the planner's agent: `planning.agent` merged over the project's `agent` as a story's agent is merged over the default (ADR-0037, ADR-0065).
- `harness.Request` carries what a planner run needs: the item (an epic or a story) and the role `plan`, with no story. `harness.Prompt` gives a planner prompt for it: prime with `--role plan` and the item, call `inbox`, read the item and what it links, do the planning its state calls for (drafting for an epic with no stories; enriching for a story; both, revisiting every child not done or cancelled, for an epic with stories), write through flai only, end with a one-line summary (flai serve logs the activity from it, ADR-0079), and ask the operator with `thread_open` on the item when inputs are missing, then end.
- The environment of a planner session: `FLAI_ROLE=plan` and `FLAI_ITEM=<id>`, no `FLAI_STORY`, so that `flai guard` can tell the planner's calls apart.
- The claude-code adapter passes `.claude/agents/planner.md` as it passes the explorer and the verifier (ADR-0065), with the planning agent's model, and runs the session as it (`--agent planner`), alongside the explorer and verifier definitions.
- `template/root/.claude/agents/planner.md`, copied to `.claude/agents/planner.md`: the planner's definition, tools for reading and flai's MCP tools, no Edit or Write, pointing at `strategic-agents.md`.

Waits for nothing: the first layer, beside the guard task, whose paths it shares none of.

## Done when

- [ ] `planning.agent` parses, is checked, and merges over `agent`, with tests
- [ ] `harness.Prompt` for a planner request says each step above, and the claude-code argv for it carries `--agents` with planner, explorer, and verifier and `--agent planner`, with tests
- [ ] The template and the repository ship the same `planner.md`

## Notes
