---
id: T-2116
type: task
nature: research
title: "The agent adapters design records the Docker Agent survey: what the harness, its providers, its hooks, its Claude Code harness, and its usage reporting give against flai's concerns"
status: done
parent: S-0496
owner: alex
created: 2026-10-25T21:07:34Z
updated: 2026-10-25T21:09:26Z
transitions:
  - to: ready
    at: 2026-10-25T21:08:02Z
    by: system-flow
  - to: in-progress
    at: 2026-10-25T21:08:03Z
    by: system-flow
  - to: done
    at: 2026-10-25T21:09:26Z
    by: agent-S-0496
stream: S-0496
tags: []
touches: [design/system/agent-adapters.md]
usage:
  source: log
  seconds: 83
  estimated: true
  priced_by: harness
  models:
    - model: claude-fable-5-1
      input: 110
      output: 4888
      cache_read: 617387
      cache_write: 31037
      cost: 1.0206
---
# T-2116 The agent adapters design records the Docker Agent survey: what the harness, its providers, its hooks, its Claude Code harness, and its usage reporting give against flai's concerns

## Work

Add a Docker Agent (formerly cagent) subsection under LiteLLM and OpenRouter in `design/system/agent-adapters.md`, beside the other harnesses, and its row in the table Harnesses against the concerns: how it starts and resumes headless (`--exec --json`, `--session`), its MCP toolset, its sub-agents and their own models, its `pre_tool_use` hook with `agent_name`, its `permission_request` hook, its usage and cost reporting, its instruction file, its providers (Anthropic, OpenAI, Docker Model Runner, Ollama, any OpenAI-compatible endpoint), and its `claude-code` harness that drives the `claude` CLI on the operator's subscription login. Say what the documentation and the installed release (v1.124.0 on this host) settle and what they do not. Waits for nothing: the first layer.

## Done when

- `design/system/agent-adapters.md` has a Docker Agent subsection and its row in the harnesses table, each statement naming its source
- `updated` is bumped and the markdown lint passes

## Notes

Sources read on 2026-10-25: docs.docker.com/ai/docker-agent (providers, local models, configuration, hooks, permissions, sandbox, secrets, sessions, harnesses, headless, telemetry), the installed `docker agent` v1.124.0 help text, and `docker agent doctor` on this host.
