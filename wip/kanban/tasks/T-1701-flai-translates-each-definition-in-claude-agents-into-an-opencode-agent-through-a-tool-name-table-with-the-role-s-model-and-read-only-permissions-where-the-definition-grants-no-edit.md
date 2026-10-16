---
id: T-1701
type: task
nature: feature
title: flai translates each definition in .claude/agents into an OpenCode agent through a tool-name table, with the role's model and read-only permissions where the definition grants no edit
status: backlog
parent: S-0409
owner: alex
created: 2026-10-16T05:17:25Z
updated: 2026-10-16T05:17:25Z
transitions: []
stream: S-0409
tags: [cli]
touches: [flai/internal/harness/opencodeagents.go, flai/internal/harness/opencodeagents_test.go]
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
