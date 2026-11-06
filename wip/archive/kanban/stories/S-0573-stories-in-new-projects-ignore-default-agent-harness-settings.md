---
id: S-0573
type: story
nature: remediation
title: Stories in new projects ignore default agent/harness settings
status: done
owner: alex
created: 2026-11-03T17:24:00Z
updated: 2026-11-03T17:55:48Z
transitions:
  - to: ready
    at: 2026-11-03T17:24:05Z
    by: alex
  - to: in-progress
    at: 2026-11-03T17:24:17Z
    by: agent-S-0573
  - to: review
    at: 2026-11-03T17:54:27Z
    by: system-flow
  - to: done
    at: 2026-11-03T17:55:48Z
    by: alex
tags: [cli]
touches: [flai/cmd, design/adrs/0251-a-story-whose-agent-names-no-harness-is-started-with-the-project-s-default.md, design/adrs/README.md, flai/internal/workitem/agent.go, flai/internal/workitem/agent_test.go, design/system/flai-cli.md, design/system/project-manifest.md, design/system/work-hierarchy.md, docs/users/flai-reference.md, docs/users/flai.md, flaiover/src/lib/components/SettingsPanel.svelte, flai/internal/serve/agents.go, flai/internal/serve/commit.go, flai/internal/serve/default_agent_test.go, flai/internal/serve/restart.go, flai/internal/serve/start.go, flai/internal/serve/usage.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1833
  priced_by: harness
  turns:
    - day: 2026-11-03
      ceremony: 2
      hand_edits: 2
      work: 55
  models:
    - model: claude-opus-5-5
      input: 238
      output: 57554
      cache_read: 12794321
      cache_write: 393201
      cost: 6.3476
  strategic:
    - kind: orchestrator
      seconds: 86
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 26
          output: 183
          cache_read: 839261
          cache_write: 84122
          cost: 0.2194
        - model: claude-sonnet-5-5
          input: 4
          output: 40
          cache_read: 12553
          cache_write: 29095
          cost: 0.0294
cost_of_delay:
  inputs:
    penalty_per_week: 50
    by: alex
    at: 2026-11-03T17:24:00Z
---
# S-0573 Stories in new projects ignore default agent/harness settings

## Goal

In new projects, after an agent has been set in system-flow.yaml and the harness and default agent are set in the configuration screens, new stories fail to start with the error that there's missing agent/harness configuration.

Diagnose the cause and correct this behavior. All configuration should be set via the dashboard without additional steps in the CLI required.

## Acceptance criteria
- [x] In new projects, where the harness and agent settings are configured, there are no additional steps required (i.e. editing the story's agent directly)
- [x] Stories in projects where the harness and agent configuration are provided start successfully without additional CLI commands or file edits required

## Tasks
- T-2454 workitem.Repo.AgentOf resolves a story's agent as the project's default under the story's own when the story names no harness
- T-2455 flai serve starts, restarts, resumes, and commits a story's agent with the agent Repo.AgentOf resolves, so a story made before the default was set starts
- T-2456 The design, the flai agent help, the user docs, and the dashboard's default agent panel say a story with no agent of its own starts with the default

## Notes
