---
id: T-1701
type: task
nature: feature
title: flai translates each definition in .claude/agents into an OpenCode agent through a tool-name table, with the role's model and read-only permissions where the definition grants no edit
status: cancelled
parent: S-0409
owner: alex
created: 2026-10-16T05:17:25Z
updated: 2026-11-07T21:53:43Z
transitions:
  - to: ready
    at: 2026-11-07T21:49:35Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T21:49:36Z
    by: system-flow
  - to: cancelled
    at: 2026-11-07T21:53:43Z
    by: alex
stream: S-0409
tags: [cli]
touches: [flai/internal/harness/opencodeagents.go, flai/internal/harness/opencodeagents_test.go]
usage:
  source: log
  seconds: 247
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 29
      output: 10856
      cache_read: 3604241
      cache_write: 239713
      cost: 2.3597
---
# T-1701 flai translates each definition in .claude/agents into an OpenCode agent through a tool-name table, with the role's model and read-only permissions where the definition grants no edit

## Work

- In `flai/internal/harness/opencodeagents.go`, read each definition with `readDefinition` and translate it: name, description, prompt, the tool allowlist through a table from Claude Code's names (`Read`, `Grep`, `Glob`, `Bash`, `Edit`, `Write`, `mcp__flai__*`) to OpenCode's, and `mode: subagent` for `explorer` and `verifier`, `mode: primary` for the planner, the orchestrator, and the analyzer.
- Set each sub-agent's `model` from the story's `agent.roles` for its role as `provider/model`, else leave it to inherit the story's.
- For a definition that grants neither `Edit` nor `Write`, set `edit: deny` and shell rules that refuse writes outside the temporary folder (ADR-0141).
- Test the five template definitions, a role with and without a model, and a tool the table does not know, which is refused with its name.
- Waits for nothing in this story.

## Done when

The five definitions translate, read-only ones carry the denials, and the tests pass.

## Notes
- 2026-11-07T21:53:43Z: moved to cancelled: E-0023 cancelled: OpenCode isn't going to receive first-class support given its aversion to supporting before hooks for tool calls (amongst other things)
