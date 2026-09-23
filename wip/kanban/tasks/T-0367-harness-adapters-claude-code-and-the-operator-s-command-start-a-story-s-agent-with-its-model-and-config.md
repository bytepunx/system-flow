---
id: T-0367
type: task
nature: feature
title: "Harness adapters: claude-code and the operator's command start a story's agent with its model and config"
status: done
parent: S-0104
owner: alex
created: 2026-09-23T17:25:01Z
updated: 2026-09-23T17:29:25Z
transitions:
  - to: ready
    at: 2026-09-23T17:25:26Z
    by: system-flow
  - to: in-progress
    at: 2026-09-23T17:25:26Z
    by: system-flow
  - to: done
    at: 2026-09-23T17:29:25Z
    by: system-flow
stream: S-0104
tags: []
---

# T-0367 Harness adapters: claude-code and the operator's command start a story's agent with its model and config

## Work
A new package `internal/harness`. An adapter turns a story's agent (harness, model, config), the project, and the prompt into the argument list and environment to start. `claude-code` runs `claude -p` headless: `--model` is the story's model, `--mcp-config` points at this flai's `mcp`, and a few story config keys become flags (`effort`, `max_budget_usd`, `fallback_model`), each checked against the values it takes. What the agent may do (`--permission-mode`, `--allowedTools`) and which program runs are the operator's, set on the host with `flai serve agent harness`, not in the story. `command` is S-0079's operator command, now also with `{model}` and `{harness}`. A story with no harness uses `command` when one is set. The prompt tells the agent to work only its story through CLAUDE.md to review, to ask through `thread_open`, and to hold `wait_for_events` until answered.

## Done when
- Unit tests show the argument lists for each adapter.
- A config key or value that is not allowed is refused, with the reason.
- An unknown harness is refused.

## Notes
