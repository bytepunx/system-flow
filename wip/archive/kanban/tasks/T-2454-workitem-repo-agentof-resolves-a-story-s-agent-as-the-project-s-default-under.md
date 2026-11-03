---
id: T-2454
type: task
nature: feature
title: workitem.Repo.AgentOf resolves a story's agent as the project's default under the story's own when the story names no harness
status: done
parent: S-0573
owner: alex
created: 2026-11-03T17:28:57Z
updated: 2026-11-03T17:36:23Z
transitions:
  - to: ready
    at: 2026-11-03T17:30:00Z
    by: system-flow
  - to: in-progress
    at: 2026-11-03T17:30:01Z
    by: system-flow
  - to: done
    at: 2026-11-03T17:36:23Z
    by: agent-S-0573
stream: S-0573
tags: []
touches: [flai/internal/workitem/agent.go, flai/internal/workitem/agent_test.go]
usage:
  source: log
  seconds: 382
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 11
      output: 2621
      cache_read: 582615
      cache_write: 17905
      cost: 0.2891
---
# T-2454 workitem.Repo.AgentOf resolves a story's agent as the project's default under the story's own when the story names no harness

## Work

Add `func (r *Repo) AgentOf(it *Item) *manifest.Agent` in `flai/internal/workitem/agent.go` (ADR-0251). When `it.Agent` names a harness, it is returned as it is. Otherwise the result is `r.Manifest.Agent.With(it.Agent)`, nil when both are empty. Behaviour tests in `agent_test.go`, beside the tests already there. Waits for nothing: first layer.

## Done when

- A story with no agent in a project with a default resolves to the default.
- A story whose agent sets only a model resolves to the default's harness with the story's model.
- A story whose agent names a harness, `none` included, resolves to its own agent unchanged, with none of the default's config or roles added.
- No story and no default give nil.
- `flai test flai/internal/workitem` passes.

## Notes
